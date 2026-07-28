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
    data map[string]cacheEntry  // Key: URL completa | Value: entrada
    mu   sync.Mutex             // Protege acceso concurrente al map
}
```

| Campo | Tipo | Propósito |
|-------|------|-----------|
| `data` | `map[string]cacheEntry` | Almacén clave-valor. Clave = URL absoluta de la request |
| `mu` | `sync.Mutex` | Serializa lecturas/escrituras (Go maps no son thread-safe) |

> **Nota:** El `interval` **NO se guarda** en el struct. Se pasa como parámetro a `reapLoop()` al crear el cache.

---

## 3. Ciclo de Vida Completo

### 3.1 Inicialización — `NewCache(interval)`

```go
func NewCache(interval time.Duration) *Cache {
    c := &Cache{
        data: make(map[string]cacheEntry),
    }
    go c.reapLoop(interval)  // ← Goroutine background con intervalo por parámetro
    return c
}
```

**Secuencia:**
1. Crea map vacío
2. Lanza `reapLoop(interval)` en **goroutine independiente** (non-blocking)
3. Retorna `*Cache` listo para uso inmediato

> **Uso real en `pagination_handler.go`:**
> ```go
> var cache = pokecache.NewCache(15 * time.Second)  // TTL = 15 segundos
> ```

---

### 3.2 Escritura — `Add(key, val)`

```go
func (c *Cache) Add(key string, val []byte) {
    c.mu.Lock()
    c.data[key] = cacheEntry{
        createdAt: time.Now(),
        val:       val,
    }
    c.mu.Unlock()
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

### 3.4 Limpieza Automática — `reapLoop(interval)`

```go
func (c *Cache) reapLoop(interval time.Duration) {
    ticker := time.NewTicker(interval)
    for range ticker.C {
        c.mu.Lock()
        for key, entry := range c.data {
            if time.Since(entry.createdAt) > interval {
                delete(c.data, key)
            }
        }
        c.mu.Unlock()
    }
}
```

**Diagrama de tiempo (interval = 15s):**

```
Interval = 15s
──────────────────────────────────────────────────────────────▶ TIME

T=0s     NewCache(15s) → reapLoop inicia con interval=15s
         │
T=0.1s   Add("url-A")  createdAt=0.1s
T=1.0s   Add("url-B")  createdAt=1.0s
         │
T=15s    ▼ TICKER TIQUEA ▼
         reapLoop despierta
         Lock()
         now = 15s
         ├── "url-A": 15 - 0.1 = 14.9s  ≤ 15s  ✓ KEEP
         └── "url-B": 15 - 1.0 = 14.0s  ≤ 15s  ✓ KEEP
         Unlock()
         │
T=30s    ▼ TICKER TIQUEA ▼
         now = 30s
         ├── "url-A": 30 - 0.1 = 29.9s  > 15s  ✗ DELETE
         └── "url-B": 30 - 1.0 = 29.0s  > 15s  ✗ DELETE
```

**Propiedades:**
- ✅ **No bloquea** `Add`/`Get` por más de microsegundos (lock corto)
- ✅ **Precisión**: entradas expiran entre `interval` y `2*interval` después de crearse
- ✅ **Memory leak prevention**: map no crece indefinidamente

---

## 4. Integración en `fetchAndDisplay`

### 4.1 Código Real (`internal/commands/pagination_handler.go`)

```go
// Variable global — UNA sola instancia para toda la app
// TTL = 15 segundos (configurado aquí, no hardcodeado en pokecache)
var cache = pokecache.NewCache(15 * time.Second)

func fetchAndDisplay(url string) error {
    // 1️⃣ INTENTO DE CACHE HIT
    if data, found := cache.Get(url); found {
        fmt.Println("Data fetched from cache:")
        var result model.LocationArea
        json.Unmarshal(data, &result)
        updatePagination(result)
        return displayResult(result)
    }

    // 2️⃣ CACHE MISS → HTTP REQUEST
    c := client.NewClient(url)
    result, err := c.GetLocationArea(url)
    if err != nil {
        return fmt.Errorf("failed to fetch data: %w", err)
    }

    // 3️⃣ GUARDAR EN CACHE
    response, err := json.Marshal(result)
    if err != nil {
        return fmt.Errorf("failed to marshal data: %w", err)
    }
    cache.Add(url, response)

    updatePagination(*result)
    return displayResult(*result)
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

### 6.1 Intervalo Real Usado

```go
// pagination_handler.go línea 13
var cache = pokecache.NewCache(15 * time.Second)
```

| Intervalo | Pros | Contras |
|-----------|------|---------|
| **15s** (actual) | Balance hit-rate / freshness | — |
| 5s | Datos más frescos | Más requests a API |
| 60s | Menos requests | Datos stale, más memoria |

**15 segundos** es el sweet spot para CLI interactivo: el usuario navega `map` → `mapb` en <15s typical.

---

### 6.2 Clave de Cache (Key)

```go
// En fetchAndDisplay, la URL COMPLETA es la key
cache.Get("https://pokeapi.co/api/v2/location-area/?offset=20&limit=20")
```

**Por qué URL completa:**
- Distingue `/?offset=0` vs `/?offset=20` vs `/?offset=40`
- Incluye query params → cache por página exacta
- Simple, sin normalización necesaria

---

## 7. Limitaciones Conocidas

| Limitación | Impacto | Mitigación |
|------------|---------|------------|
| **Volátil** — se pierde al cerrar CLI | Cold start = requests fríos | Acceptable para CLI |
| **Sin validación de contenido** | JSON inválido se cachea | `Unmarshal` falla al leer |
| **Sin compresión** | `[]byte` crudo = más memoria | `gzip` si payloads >100KB |
| **Single-process** | No compartido entre instancias | Redis si se escala |

---

## 8. Testing — Cobertura

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

## 9. Referencia Rápida — API Pública

```go
package pokecache

// Cache — thread-safe in-memory cache con TTL automático
type Cache struct { ... }

// NewCache — crea cache con intervalo de limpieza
// Lanza reapLoop(interval) en background
func NewCache(interval time.Duration) *Cache

// Add — guarda valor con timestamp actual
// Sobrescribe si key existe
func (c *Cache) Add(key string, val []byte)

// Get — recupera valor si existe
// Retorna (val, true) en hit; (nil, false) en miss
func (c *Cache) Get(key string) ([]byte, bool)

// reapLoop — interno, corre cada interval
// Elimina entries donde now - createdAt > interval
func (c *Cache) reapLoop(interval time.Duration)
```

---

## 10. Archivos Relacionados

| Archivo | Responsabilidad |
|---------|-----------------|
| `internal/pokecache/cache.go` | Structs + `NewCache` + `Add` + `Get` |
| `internal/pokecache/cache_handler.go` | `reapLoop` |
| `internal/pokecache/cache_test.go` | Tests table-driven + concurrent |
| `internal/commands/pagination_handler.go` | Integración en `fetchAndDisplay` + TTL 15s |

---

## 11. Caché en `catch` (Pokémon)

El comando `catch` reutiliza la misma caché global compartida desde `pagination_handler.go` (`var cache = pokecache.NewCache(15 * time.Second)`), con el mismo patrón que `fetchAndDisplay`:

```
catchPokemon(url, name)
  ├─ 1️⃣  Vault check: ¿ya fue capturado? → "has already been caught!"
  ├─ 2️⃣  cache.Get(url) → HIT: usar datos almacenados
  ├─ 3️⃣  cache.Get(url) → MISS: HTTP GET → cache.Add(url, jsonBytes)
  └─ 4️⃣  catchlogic.Catch(pokemonDetails) → probabilidad de captura
```

**Diferencia con otros comandos:** Antes de consultar el cache o hacer el fetch, el vault se verifica primero. Si el Pokémon ya fue capturado, no se hace request HTTP ni se accede al cache.

### Invalidación de caché en catch

Cuando `Catch()` retorna `"was caught"`, la entrada del cache del Pokémon se elimina automáticamente para mantener coherencia:

```go
cache.Delete(url)  // Invalidar caché del Pokémon recién capturado
```

---

## 12. PokemonVault — Inventario

El vault (`internal/pokemonVault/`) almacena los datos completos de los Pokémon capturados:

| Método | Retorna | Descripción |
|--------|---------|-------------|
| `AddPokemon(details)` | — | Registra un Pokémon con todos sus datos |
| `GetPokemonCaught()` | `[]string` | Lista de nombres capturados |
| `GetPokemonDetails(name)` | `PokemonVault` | Datos completos por nombre |
| `DisplayPokemonDetails(pokemon)` | `string` | String formateado con stats, tipos, height, weight |

**Estructuras:**
```go
type PokemonVault struct {
    PokemonName   string
    PokemonHeight int
    PokemonWeight int
    PokemonStats  []PokemonStat  // HP, Attack, Defense, Sp.Atk, Sp.Def, Speed
    PokemonTypes  []PokemonType  // Tipos (electric, water, fire...)
    HasBeenCaught bool
}
```

**Instancia global:** `DefaultPokemonVault PokemonVaultMethods` — accesible desde cualquier comando vía `pokemonVault.DefaultPokemonVault`.

---

## 13. Lógica de Captura

Implementada en `internal/catchlogic/catch.go`:

```go
func possiblyCatchPokemon(baseExperience int) bool {
    if baseExperience <= 0 {
        baseExperience = 1 // mínimo para dar alguna chance de captura
    }
    maxCaptureChance := baseExperience + 50
    captureChance := rand.Intn(maxCaptureChance)
    // Random number between 0 and maxCaptureChance
    return captureChance >= baseExperience
}
```

Se toma la `baseExperience` del Pokémon y se le suma 50, obteniendo `maxCaptureChance`. Se genera un número aleatorio entre 0 y `maxCaptureChance`. Si el número generado es **mayor o igual** a `baseExperience`, el Pokémon se captura. Si es menor, escapa.

- `baseExperience ≤ 0` → mínimo forzado a `1` (nunca 0% de chance)
- Mayor `baseExperience` → mayor umbral de captura → menos probable captura
- Menor `baseExperience` → menor umbral de captura → más probable captura

Función clave: `possiblyCatchPokemon(baseExperience int) bool` en `internal/catchlogic/catch.go`.