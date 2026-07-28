package catchlogic

import (
	"fmt"
	"math/rand"

	"github.com/NicolasFerreras/REPL-Pokedex/internal/model"
	"github.com/NicolasFerreras/REPL-Pokedex/internal/pokemonVault"
)

func Catch(pokemon model.PokemonDetails) string {
	baseExperience := pokemon.BaseExperience
	if !possiblyCatchPokemon(baseExperience) {
		return fmt.Sprintf("%s escaped!", pokemon.Name)
	}
	pokemonVault.DefaultPokemonVault.AddPokemon(pokemon)
	return fmt.Sprintf("%s was caught!", pokemon.Name)
}

func possiblyCatchPokemon(baseExperience int) bool {
	if baseExperience <= 0 {
		baseExperience = 1 // mínimo para dar alguna chance de captura
	}
	maxCaptureChance := baseExperience + 50
	captureChance := rand.Intn(maxCaptureChance)
	// Random number between 0 and maxCaptureChance
	return captureChance >= baseExperience
}
