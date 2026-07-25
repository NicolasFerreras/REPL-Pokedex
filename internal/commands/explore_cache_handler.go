package commands

import (
	"encoding/json"
	"fmt"

	"github.com/NicolasFerreras/REPL-Pokedex/internal/client"
	"github.com/NicolasFerreras/REPL-Pokedex/internal/model"
)

func exploreCache(url string) error {
	// Check cache first
	if data, found := cache.Get(url); found {
		fmt.Println("Data fetched from cache:")
		var result model.PokemonEncountersList

		json.Unmarshal(data, &result)

		return displayResultPokemon(result)
	}

	// If not in cache, fetch from API
	c := client.NewClient(url)
	result, err := c.GetPokemon(url)
	if err != nil {
		return fmt.Errorf("failed to fetch data: %w", err)
	}

	response, err := json.Marshal(result)
	if err != nil {
		return fmt.Errorf("failed to marshal data: %w", err)
	}
	cache.Add(url, response)

	return displayResultPokemon(*result)

}

func displayResultPokemon(result model.PokemonEncountersList) error {
	fmt.Println("Pokemon Encounters:")
	for _, r := range result.PokemonEncounters {
		fmt.Printf("- %s\n", r.Pokemon.Name)
	}
	return nil
}
