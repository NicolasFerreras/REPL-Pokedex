# REPL Pokedex

Pokedex interactiva en terminal construida en Go, consumiendo la [PokéAPI](https://pokeapi.co/).

## Características

- REPL interactivo con loop de comandos
- Navegación por location areas con paginación (`map` / `mapb`)
- Cliente HTTP con timeout y manejo de errores
- **Caché en memoria thread-safe** con expiración automática (TTL)
- Sistema de comandos extensible con patrón callback
- Tests unitarios con tabla-driven tests (idiomático Go)

## Comandos

| Comando | Descripción |
|---------|-------------|
| `help`  | Muestra los comandos disponibles |
| `map`   | Muestra las siguiente página de location areas |
| `mapb`  | Muestra la página anterior de location areas |
| `exit`  | Cierra la Pokedex |

## Estructura del proyecto

```
REPL-Pokedex/
├── cmd/
│   └── pokedex/
│       └── main.go              # Entry point
├── internal/
│   ├── client/
│   │   └── client.go            # Cliente HTTP para PokéAPI
│   ├── commands/
│   │   ├── commands.go          # Registro de comandos
│   │   ├── commands_functions.go # Implementación de comandos
│   │   └── commands_test.go     # Tests de comandos
│   ├── model/
│   │   ├── endpoints.go         # Config y datos de paginación
│   │   └── location_area.go     # Modelos de LocationArea
│   ├── pokecache/
│   │   ├── cache.go             # Structs + NewCache + Add + Get
│   │   ├── cache_handler.go     # reapLoop (limpieza automática)
│   │   └── cache_test.go        # Tests del cache
│   └── repl/
│       ├── repl.go              # Lógica del REPL
│       └── repl_test.go         # Tests del REPL
├── go.mod
└── README.md
```

## Requisitos

- Go 1.26.5 o superior

## Instalación y ejecución

```bash
# Clonar el repositorio
git clone https://github.com/NicolasFerreras/REPL-Pokedex.git
cd REPL-Pokedex

# Ejecutar
go run ./cmd/pokedex
```

## Tests

```bash
go test ./...
```

## API

Este proyecto consume la [PokéAPI](https://pokeapi.co/) — una API gratuita y open source de Pokémon.

Endpoint utilizado:
- `GET /api/v2/location-area/` — Lista de areas de ubicación con paginación

## Caché (pokecache)

Implementación propia de caché en memoria **thread-safe** con expiración automática:

- **Almacenamiento:** `map[string]cacheEntry` protegido por `sync.Mutex`
- **Entrada:** `cacheEntry{ createdAt time.Time, val []byte }` — JSON crudo
- **Limpieza:** `reapLoop()` en goroutine con `time.Ticker` — elimina entradas > TTL
- **TTL configurable:** 5 segundos por defecto
- **Integración:** `fetchAndDisplay` usa `cache.Get/Add` — cache hit = instantáneo

📖 **Documentación completa:** [`Cache-Docs.md`](Cache-Docs.md)

## Stack

- **Lenguaje:** Go
- **API:** PokéAPI
- **Arquitectura:** REPL con patrón de comandos
- **HTTP:** `net/http` (stdlib)
- **Tests:** Table-driven tests

## Licencia

Proyecto privado de aprendizaje
