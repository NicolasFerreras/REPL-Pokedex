# REPL Pokedex

Pokedex interactiva en terminal construida en Go, consumiendo la [PokéAPI](https://pokeapi.co/).

## Características

- REPL interactivo con loop de comandos
- Navegación por location areas con paginación (`map` / `mapb`)
- Exploración de Pokémon en un área específica (`explore`)
- Captura de Pokémon con probabilidad basada en `baseExperience` (`catch`)
- Inspección detallada de Pokémon capturados (`inspect`)
- Pokedex con listado de todos los Pokémon atrapados (`pokedex`)
- Cliente HTTP con timeout y manejo de errores
- **Caché en memoria thread-safe** con expiración automática (TTL)
- **PokemonVault** — inventario persistente de Pokémon capturados durante la sesión
- Sistema de comandos extensible con patrón callback
- Sistema de error handling centralizado
- Tests unitarios con tabla-driven tests (idiomático Go)

## Comandos

| Comando | Argumento | Descripción |
|---------|-----------|-------------|
| `help` | — | Muestra los comandos disponibles |
| `map` | — | Muestra la siguiente página de location areas |
| `mapb` | — | Muestra la página anterior de location areas |
| `explore` | `<área>` | Explora una location area específica y lista los Pokémon |
| `catch` | `<pokemon>` | Atrapa un Pokémon (probabilidad basada en baseExperience) |
| `inspect` | `<pokemon>` | Muestra detalles de un Pokémon capturado |
| `pokedex` | — | Lista todos los Pokémon capturados |
| `exit` | — | Cierra la Pokedex |

## Uso

```bash
# Iniciar la Pokedex
go run ./cmd/pokedex

# Ejemplo de flujo
> map
> explore santos-nurf
> catch pikachu
Throwing a Pokeball at pikachu...
pikachu was caught!
> inspect pikachu
Name: Pikachu
Height: 4
Weight: 60
Stats:
- HP: 35
- Attack: 55
- Defense: 40
- Special Attack: 50
- Special Defense: 50
- Speed: 90
Types:
- electric
> pokedex
Your Pokedex:
- pikachu
> exit
Closing the Pokedex... Goodbye!
```

## Estructura del proyecto

```
REPL-Pokedex/
├── cmd/
│   └── pokedex/
│       └── main.go                   # Entry point, REPL loop, command routing
├── internal/
│   ├── client/
│   │   └── client.go                 # Cliente HTTP para PokéAPI
│   ├── catchlogic/
│   │   └── catch.go                  # Lógica de captura (probabilidad)
│   ├── commands/
│   │   ├── commands.go               # Registro de comandos
│   │   ├── commands_functions.go     # Implementación de comandos
│   │   ├── commands_test.go          # Tests de comandos
│   │   ├── catch_cache_handler.go    # Handler del comando catch (vault + caché)
│   │   ├── explore_cache_handler.go  # Handler del comando explore (caché)
│   │   └── pagination_handler.go     # Paginación (map/mapb)
│   ├── errors/
│   │   └── errors.go                 # Constantes de error centralizadas
│   ├── model/
│   │   ├── endpoints.go              # Config de URLs y paginación
│   │   ├── location_area.go          # Modelo LocationArea
│   │   ├── pokemon.go                # Modelo PokemonDetails
│   │   └── pokemon_encounters.go     # Modelo PokemonEncountersList
│   ├── pokecache/
│   │   ├── cache.go                  # Struct Cache + NewCache
│   │   ├── cache_handler.go          # reapLoop (limpieza automática)
│   │   └── cache_test.go             # Tests del cache
│   ├── pokemonVault/
│   │   └── pokemonVault.go           # Inventario de Pokémon capturados
│   ├── repl/
│   │   ├── repl.go                   # Lógica del REPL
│   │   └── repl_test.go              # Tests del REPL
├── go.mod
├── go.sum
├── README.md
└── Cache-Docs.md
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

### Endpoints utilizados

| Endpoint | Uso |
|----------|-----|
| `GET /api/v2/location-area/` | Lista de location areas con paginación |
| `GET /api/v2/location-area/<nombre>` | Detalles de un área ( Pokémon disponibles) |
| `GET /api/v2/pokemon/<nombre>` | Detalles de un Pokémon (nombre, base_experience, height, weight, stats, types) |

## Caché (pokecache)

Implementación propia de caché en memoria **thread-safe** con expiración automática:

- **Almacenamiento:** `map[string]cacheEntry` protegido por `sync.Mutex`
- **Entrada:** `cacheEntry{ createdAt time.Time, val []byte }` — JSON crudo
- **Limpieza:** `reapLoop()` en goroutine con `time.Ticker` — elimina entradas expiradas
- **TTL configurable:** 15 segundos
- **Integración:** Usada por `fetchAndDisplay` (map/mapb), `exploreCache`, y `catchPokemon`
- **Clave de caché:** URL completa con query params

### Flujo de caché

```
Petición
  ├─ cache.Get(url) → HIT: usar datos almacenados
  └─ cache.Get(url) → MISS: fetch API → cache.Add(url, response) → usar datos
```

📖 **Documentación completa:** [`Cache-Docs.md`](Cache-Docs.md)

## PokemonVault

Inventario en memoria de los Pokémon capturados durante la sesión:

- **`AddPokemon(pokemonDetails)`** — registra un Pokémon capturado con todos sus datos
- **`GetPokemonCaught()`** — retorna lista de nombres de Pokémon atrapados
- **`GetPokemonDetails(name)`** — retorna los datos completos de un Pokémon capturado
- **`DisplayPokemonDetails(pokemon)`** — retorna string formateado con stats, tipos, height y weight
- **Instancia global** `DefaultPokemonVault` — accesible desde cualquier comando

## Lógica de captura

La probabilidad de captura se basa en el `baseExperience` del Pokémon:

```
probabilidad = baseExperience / (baseExperience + 50)
```

- `baseExperience = 0` → se usa mínimo de 1 (nunca 0% de chance)
- Mayor `baseExperience` → mayor probabilidad de captura

## Stack

| Componente | Tecnología |
|------------|------------|
| **Lenguaje** | Go |
| **API** | PokéAPI |
| **Arquitectura** | REPL con patrón de comandos |
| **HTTP** | `net/http` (stdlib) |
| **Tests** | Table-driven tests |
| **Caché** | In-memory, thread-safe, TTL |

## Error handling

Mensajes de error centralizados en `internal/errors/errors.go`:

| Constante | Mensaje |
|-----------|---------|
| `ErrNoCommandEntered` | `no command entered. Please enter a command` |
| `ErrUnknownCommand` | `unknown command: %s. Type 'help' for available commands` |
| `ErrCommandNeedsArg` | `command %s requires an argument` |
| `ErrFetchData` | `failed to fetch data: %w` |
| `ErrUnmarshalData` | `failed to unmarshal data: %w` |
| `ErrMarshalData` | `failed to marshal data: %w` |
| `ErrPokemonNotCaught` | `you haven't caught %s yet` |

Códigos HTTP mapeados a mensajes amigables:

| Código | Mensaje |
|--------|---------|
| 400 | `bad request — check your input` |
| 404 | `resource not found — the location or Pokemon does not exist` |
| 429 | `too many requests — slow down and try again` |
| 500 | `internal server error — the API is having issues` |

## Licencia

Proyecto privado de aprendizaje
