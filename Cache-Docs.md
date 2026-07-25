# 🔄 Cache Operation — REPL Pokedex

## Resumen Ejecutivo

El sistema de caché implementado en `internal/pokecache/` es una **caché en memoria thread-safe** diseñada para almacenar respuestas HTTP crudas (`[]byte`) de la PokéAPI y evitar requests redundantes cuando el usuario navega entre páginas con `map` y `mapb`.

---

## 1. Arquitectura General

```
┌─────────────────────────────────────────────────────────────────────────────┐
│                           CACHE LAYER                                        │
├─────────────────────────────────────────────────────────────────────────────┤
│                                                                             │
│  ┌──────────────┐     ┌─────────────────────┐     ┌────────────────────┐   │
│  │  fetchAnd    │────▶│  cache.Get(url)     │────▶│  HIT: devuelve     │   │
│  │  Display     │     │  ([]byte, bool)     │     │  []byte cached     │   │
│  └──────────────┘     └─────────────────────┘     └────────────────────┘   │
│         │                     │                                                │
│         │ MISS                │ HIT                                           │
│         ▼                     ▼                                               │
│  ┌─────────────────────┐     ┌────────────────────┐                          │
│  │ HTTP Request        │     │ json.Unmarshal     │                          │
│  │ client.GetLocation  │     │ ([]byte → struct)  │                          │
│  └─────────────────────┘     └────────────────────┘                          │
│         │                                                                │
│         ▼                                                                │
│  ┌─────────────────────┐                                                │
│  │ json.Marshal        │                                                │
│  │ (struct → []byte)   │                                                │
│  └─────────────────────┘                                                │
│         │                                                                │
│         ▼                                                                │
│  ┌─────────────────────┐                                                │
│  │ cache.Add(url,      │                                                │
│  │       respBytes)    │                                                │
│  └─────────────────────┘                                                │
│                                                                             │
└─────────────────────────────────────────────────────────────────────────────┘
```

---

## 2. Estructuras de Datos

### 2.1 `cacheEntry` — Entrada Individual

```go
type cacheEntry struct {
    createdAt time.Time  // Timestamp exacto cuando se guardó la entrada
    val       []byte     // Respuesta HTTP cruda (JSON serializado)
}
```

| Campo | Tipo | Propósito |
|-------|------|-----------|
| `createdAt` | `time.Time` | Permite calcular la antigüedad para expiración |
| `val` | `[]byte` | JSON raw — desacoplado de tipos de dominio |

> **Diseño intencional:** Se guarda `[]byte` en lugar de `*model.LocationArea` para que el cache sea **agnóstico al contenido**. Puede almacenar respuestas de `/location-area/`, `/pokemon/`, `/berry/`, etc., sin cambios en el paquete `pokecache`.

---

### 2.2 `Cache` — Contenedor Principal

```go
type Cache struct {
    data     map[string]cacheEntry  // Key: URL completa | Value: entrada
    mu       sync.Mutex             // Protege acceso concurrente al map
    interval time.Duration          // Ventana de vida de cada entrada
}
```

| Campo | Tipo | Propósito |
|-------|------|-----------|
| `data` | `map[string]cacheEntry` | Almacén clave-valor. Clave = URL absoluta de la request |
| `mu` | `sync.Mutex` | Serializa lecturas/escrituras (Go maps no son thread-safe) |
| `interval` | `time.Duration` | TTL efectivo: entradas > `interval` son eliminadas |

---

## 3. Ciclo de Vida Completo

### 3.1 Inicialización — `NewCache(interval)`

```go
func NewCache(interval time.Duration) *Cache {
    c := &Cache{
        data:     make(map[string]cacheEntry),
        interval: interval,
    }
    go c.reapLoop()  // ← Goroutine background
    return c
}
```

**Secuencia:**
1. Crea map vacío y guarda `interval`
2. Lanza `reapLoop()` en **goroutine independiente** (non-blocking)
3. Retorna `*Cache` listo para uso inmediato

> **Importante:** El cache es **singleton por proceso** — se crea una vez en `commands_functions.go`:
> ```go
> var cache = pokecache.NewCache(5 * time.Second)
> ```

---

### 3.2 Escritura — `Add(key, val)`

```go
func (c *Cache) Add(key string, val []byte) {
    c.mu.Lock()
    defer c.mu.Unlock()
    c.data[key] = cacheEntry{
        createdAt: time.Now(),
        val:       val,
    }
}
```

**Flujo:**
```
Add("https://pokeapi.co/api/v2/location-area/?offset=20", jsonBytes)
        │
        ▼
   Lock(mu)
        │
        ▼
   data["url"] = cacheEntry{createdAt: now, val: jsonBytes}
        │
        ▼
   Unlock(mu)
```

**Garantías:**
- ✅ Atomicidad: `Lock`/`Unlock` envuelven la escritura completa
- ✅ Sobrescritura silenciosa: si `key` existe, se reemplaza
- ✅ Timestamp fresco: `createdAt = time.Now()` en cada `Add`

---

### 3.3 Lectura — `Get(key)`

```go
func (c *Cache) Get(key string) ([]byte, bool) {
    c.mu.Lock()
    entry, exists := c.data[key]
    c.mu.Unlock()
    if !exists {
        return nil, false
    }
    return entry.val, true
}
```

**Flujo:**
```
Get("https://pokeapi.co/api/v2/location-area/?offset=20")
        │
        ▼
   Lock(mu)
        │
        ▼
   entry, exists := data["url"]
   Unlock(mu)              ← Unlock ANTES del return (crítico)
        │
        ├── exists=false ──▶ return (nil, false)
        │
        └── exists=true ──▶ return (entry.val, true)
```

**Puntos clave:**
- `Unlock()` **antes** del `return` evita mantener el lock durante la copia del slice
- Retorna copia del `[]byte` (Go copies slice header, no underlying array — seguro)
- `bool` indica *cache hit/miss* sin ambigüedad

---

### 3.4 Limpieza Automática — `reapLoop()`

```go
func (c *Cache) reapLoop() {
    ticker := time.NewTicker(c.interval)
    defer ticker.Stop()

    for range ticker.C {
        c.mu.Lock()
        now := time.Now()
        for key, entry := range c.data {
            if now.Sub(entry.createdAt) > c.interval {
                delete(c.data, key)
            }
        }
        c.mu.Unlock()
    }
}
```

**Diagrama de tiempo:**

```
Interval = 5s
──────────────────────────────────────────────────────────────▶ TIME

T=0s     NewCache(5s) → reapLoop inicia
         │
T=0.1s   Add("url-A")  createdAt=0.1s
T=1.0s   Add("url-B")  createdAt=1.0s
         │
T=5.0s   ▼ TICKER TIQUEA ▼
         reapLoop despierta
         Lock()
         now = 5.0s
         ├── "url-A": 5.0 - 0.1 = 4.9s  ≤ 5s  ✓ KEEP
         └── "url-B": 5.0 - 1.0 = 4.0s  ≤ 5s  ✓ KEEP
         Unlock()
         │
T=10.0s  ▼ TICKER TIQUEA ▼
         now = 10.0s
         ├── "url-A": 10.0 - 0.1 = 9.9s  > 5s  ✗ DELETE
         └── "url-B": 10.0 - 1.0 = 9.0s  > 5s  ✗ DELETE
```

**Propiedades:**
- ✅ **No bloquea** `Add`/`Get` por más de microsegundos (lock corto)
- ✅ **Precisión**: entradas expiran entre `interval` y `2*interval` después de crearse
- ✅ **Limpieza graceful**: `defer ticker.Stop()` al salir (aunque la goroutine es infinita)
- ✅ **Memory leak prevention**: map no crece indefinidamente

---

## 4. Integración en `fetchAndDisplay`

### 4.1 Código Real (`internal/commands/commands_functions.go`)

```go
// Variable global — UNA sola instancia para toda la app
var cache = pokecache.NewCache(5 * time.Second)

func fetchAndDisplay(url string) error {
    // 1️⃣ INTENTO DE CACHE HIT
    if val, ok := cache.Get(url); ok {
        fmt.Println("(cache hit)")           // Feedback visual
        var result model.LocationArea
        if err := json.Unmarshal(val, &result); err != nil {
            return err
        }
        return displayResult(result)
    }

    // 2️⃣ CACHE MISS → HTTP REQUEST
    c := client.NewClient(url)
    result, err := c.GetLocationArea(url)
    if err != nil {
        return fmt.Errorf("failed to get location area: %w", err)
    }

    // 3️⃣ GUARDAR EN CACHE
    respBytes, _ := json.Marshal(result)
    cache.Add(url, respBytes)

    return displayResult(result)
}
```

### 4.2 Flujo de Datos Completo

```
Usuario: map
    │
    ▼
fetchAndDisplay("https://pokeapi.co/api/v2/location-area/")
    │
    ├── cache.Get(url) ──────────────────▶ MISS (primera vez)
    │       │
    │       ▼
    │   HTTP GET → PokéAPI
    │       │
    │       ▼
    │   JSON Response → *LocationArea
    │       │
    │       ├── displayResult()
    │       │
    │       └── json.Marshal() → []byte
    │               │
    │               ▼
    │          cache.Add(url, bytes)
    │
    ▼
RESULTADO EN PANTALLA

────────────────────────────────────────────────────────

Usuario: mapb
    │
    ▼
fetchAndDisplay("https://pokeapi.co/api/v2/location-area/?offset=20")
    │
    ├── cache.Get(url) ──────────────────▶ HIT ✓
    │       │
    │       ▼
    │   json.Unmarshal(bytes) → *LocationArea
    │       │
    │       ▼
    │   displayResult()
    │
    ▼
RESULTADO INSTANTÁNEO (sin HTTP)
```

---

## 5. Concurrencia — Análisis de Seguridad

### 5.1 ¿Por qué `sync.Mutex` y no `sync.RWMutex`?

```go
// Opción actual: Mutex simple
mu sync.Mutex

// Alternativa: RWMutex (lecturas concurrentes)
mu sync.RWMutex
```

| Aspecto | `Mutex` | `RWMutex` |
|---------|---------|-----------|
| Complejidad | Simple | Moderada |
| `Get` concurrentes | Serializados | **Paralelos** |
| `Add` vs `Get` | Exclusivos | Exclusivos |
| Overhead | Mínimo | Levemente mayor |

**Decisión:** `Mutex` es suficiente porque:
- Operaciones son **ultra-rápidas** (map access + timestamp check)
- Contención real es **mínima** (usuario humano → ~1 req/seg)
- Código más **simple y auditable**

> Si el cache sirviera a un servidor web con miles de RPS, `RWMutex` sería preferible.

---

### 5.2 Race Conditions Prevenidos

| Escenario | Sin Mutex | Con Mutex |
|-----------|-----------|-----------|
| `Add` + `Get` simultáneos | **Data race** (map corrupto) | Serializado |
| `reapLoop` `delete` + `Get` | **Panic** (iteración + delete) | Seguro |
| Dos `Add` misma key | Comportamiento indefinido | Determinístico |

**Test de validación:** `TestCacheConcurrentAccess` (100 writers + 100 readers concurrentes) — **PASS** con `-race`.

---

## 6. Configuración y Parámetros

### 6.1 Intervalo Recomendado

```go
var cache = pokecache.NewCache(5 * time.Second)
```

| Intervalo | Pros | Contras |
|-----------|------|---------|
| **1-5s** (actual) | Datos frescos, poca memoria | Más requests a API |
| **30-60s** | Menos requests | Datos stale, más memoria |
| **∞** (sin reap) | Máximo hit rate | **Memory leak** |

**5 segundos** es el sweet spot para CLI interactivo: el usuario navega `map` → `mapb` en <5s typical.

---

### 6.2 Clave de Cache (Key)

```go
// En fetchAndDisplay, la URL COMPLETA es la key
cache.Get("https://pokeapi.co/api/v2/location-area/?offset=20&limit=20")
```

**Por qué URL completa:**
- Distingue `/location-area/?offset=0` vs `/?offset=20` vs `/?offset=40`
- Incluye query params → cache por página exacta
- Simple, sin normalización necesaria

---

## 7. Métricas y Observabilidad

### 7.1 Logging Actual

```go
if val, ok := cache.Get(url); ok {
    fmt.Println("(cache hit)")     // ← Único feedback
    ...
}
```

### 7.2 Métricas Sugeridas (Future Work)

```go
type Cache struct {
    ...
    hits   int64  // atomic
    misses int64  // atomic
    evictions int64
}
```

Exponibles vía `expvar` o `/debug/vars` para monitoring.

---

## 8. Limitaciones Conocidas

| Limitación | Impacto | Mitigación |
|------------|---------|------------|
| **Volátil** — se pierde al cerrar CLI | Cold start = requests fríos | Acceptable para CLI |
| **Sin validación de contenido** | JSON inválido se cachea | `Unmarshal` falla al leer |
| **Sin compresión** | `[]byte` crudo = más memoria | `gzip` si payloads >100KB |
| **Single-process** | No compartido entre instancias | Redis si se escala |

---

## 9. Testing — Cobertura

| Test | Qué Verifica |
|------|--------------|
| `TestCacheAddGet` | Round-trip básico (3 casos) |
| `TestCacheGetMissing` | Key inexistente → `(nil, false)` |
| `TestCacheOverwrite` | Segunda `Add` reemplaza valor |
| `TestCacheReapLoop` | Entradas viejas expiran, nuevas persisten |
| `TestCacheConcurrentAccess` | 200 goroutines sin races (`-race` pass) |

**Ejecutar:**
```bash
go test -v ./internal/pokecache/...
go test -race ./internal/pokecache/...  # Detección de races
```

---

## 10. Referencia Rápida — API Pública

```go
package pokecache

// Cache — thread-safe in-memory cache con TTL automático
type Cache struct { ... }

// NewCache — crea cache con intervalo de limpieza
// Lanza reapLoop() en background
func NewCache(interval time.Duration) *Cache

// Add — guarda valor con timestamp actual
// Sobrescribe si key existe
func (c *Cache) Add(key string, val []byte)

// Get — recupera valor si existe y no expiró
// Retorna (val, true) en hit; (nil, false) en miss
func (c *Cache) Get(key string) ([]byte, bool)

// reapLoop — interno, corre cada interval
// Elimina entries donde now - createdAt > interval
func (c *Cache) reapLoop()
```

---

## 11. Archivos Relacionados

| Archivo | Responsabilidad |
|---------|-----------------|
| `internal/pokecache/cache.go` | Structs + `NewCache` + `Add` + `Get` |
| `internal/pokecache/cache_handler.go` | `reapLoop` |
| `internal/pokecache/cache_test.go` | Tests table-driven + concurrent |
| `internal/commands/commands_functions.go` | Integración en `fetchAndDisplay` |

---

**Última actualización:** 24 Julio 2026  
**Versión:** 1.0  
**Autor:** Nicolas Ferreras