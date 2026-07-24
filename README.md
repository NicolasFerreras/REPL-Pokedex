# REPL Pokedex

Pokedex interactiva en terminal construida en Go, consumiendo la [PokéAPI](https://pokeapi.co/).

## Características

- REPL interactivo con loop de comandos
- Navegación por location areas con paginación (`map` / `mapb`)
- Cliente HTTP con timeout y manejo de errores
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

## Stack

- **Lenguaje:** Go
- **API:** PokéAPI
- **Arquitectura:** REPL con patrón de comandos
- **HTTP:** `net/http` (stdlib)
- **Tests:** Table-driven tests

## Licencia

Proyecto privado de aprendizaje
