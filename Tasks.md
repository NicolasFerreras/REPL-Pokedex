# 📋 Tasks

## 🟡 Caché en memoria para respuestas HTTP

**Estado:** In Progress
**Prioridad:** Media
**Fecha:** 24 de Julio 2026

### Descripción

Como usuario de la CLI Pokedex, quiero que las respuestas de PokéAPI se almacenen en caché en memoria, para que las consultas repetidas a la misma zona sean instantáneas sin realizar requests HTTP redundantes.

### Criterios de Aceptación

- [ ] DADO que el usuario consulta una URL de PokéAPI por segunda vez, CUANDO se ejecuta `map` o `mapb`, ENTONCES la respuesta se obtiene del caché (sin request HTTP).
- [ ] DADO que el caché está vacío, CUANDO el usuario consulta una URL nueva, ENTONCES se realiza el request HTTP y se almacena el resultado.
- [ ] DADO que hay múltiples goroutines accediendo al caché, CUANDO se leen o escriben datos, ENTONCES no ocurren race conditions (protegido por `sync.RWMutex`).
- [ ] DADO que el usuario cierra la CLI, CUANDO se termina la ejecución, ENTONCES el caché se libera correctamente (sin memory leaks).

### Subtareas

1. [ ] Definir estructura `Cache` con `map[string]model.LocationArea` + `sync.RWMutex`
2. [ ] Implementar métodos `Get(url) (model.LocationArea, bool)` y `Set(url string, data model.LocationArea)`
3. [ ] Integrar caché en `fetchAndDisplay` (`internal/commands/commands_functions.go`)
4. [ ] Verificar que no hay race conditions con `go test -race ./...`
5. [ ] Documentar comportamiento del caché en README

### Contexto Técnico

- **Ubicación sugerida:** `internal/cache/cache.go`
- **Estructura:** `map[string]model.LocationArea` + `sync.RWMutex`
- **Caché volátil:** se pierde al cerrar la CLI (sin persistencia)
- **Sin dependencias externas**

### Notas

- El usuario quiere implementar el código él mismo — esta story es solo documentación
- No hay manejo de expiración de caché (válido para uso interactivo)
