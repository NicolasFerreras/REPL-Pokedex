package commands

import (
	"encoding/json"
	"fmt"

	"github.com/NicolasFerreras/REPL-Pokedex/internal/catchlogic"
	"github.com/NicolasFerreras/REPL-Pokedex/internal/client"
	"github.com/NicolasFerreras/REPL-Pokedex/internal/errors"
	"github.com/NicolasFerreras/REPL-Pokedex/internal/model"
	"github.com/NicolasFerreras/REPL-Pokedex/internal/pokemonVault"
)

func catchPokemon(url, pokemonName string) error {
	// Check cache first

	for _, p := range pokemonVault.DefaultPokemonVault.GetPokemonCaught() {
		if p == pokemonName {
			fmt.Printf("%s has already been caught!\n", pokemonName)
			return nil
		}
	}

	if data, found := cache.Get(url); found {
		fmt.Println("Data fetched from cache:")
		var result model.PokemonDetails

		if err := json.Unmarshal(data, &result); err != nil {
			return fmt.Errorf(errors.ErrUnmarshalData, err)
		}

		fmt.Println(catchlogic.Catch(result))
		return nil
	}

	// If not in cache, fetch from API
	c := client.NewClient(url)
	result, err := c.GetPokemon(url)
	if err != nil {
		return fmt.Errorf(errors.ErrFetchData, err)
	}

	response, err := json.Marshal(result)
	if err != nil {
		return fmt.Errorf(errors.ErrMarshalData, err)
	}
	cache.Add(url, response)

	fmt.Println(catchlogic.Catch(*result))
	return nil

}
