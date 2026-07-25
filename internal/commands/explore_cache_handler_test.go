package commands

import (
	"encoding/json"
	"testing"
	"time"

	"github.com/NicolasFerreras/REPL-Pokedex/internal/model"
	"github.com/NicolasFerreras/REPL-Pokedex/internal/pokecache"
)

func TestExploreCacheHit(t *testing.T) {
	c := pokecache.NewCache(15 * time.Second)

	pokemon := model.PokemonEncountersList{
		PokemonEncounters: []model.PokemonEncounter{
			{Pokemon: model.Pokemon{Name: "squirtle"}},
			{Pokemon: model.Pokemon{Name: "blastoise"}},
		},
	}
	data, err := json.Marshal(pokemon)
	if err != nil {
		t.Fatalf("failed to marshal test data: %v", err)
	}
	c.Add("test-url-hit", data)

	oldCache := cache
	cache = c
	defer func() { cache = oldCache }()

	err = exploreCache("test-url-hit")
	if err != nil {
		t.Errorf("exploreCache() with cache hit returned error: %v", err)
	}
}

func TestExploreCacheMiss(t *testing.T) {
	c := pokecache.NewCache(15 * time.Second)

	oldCache := cache
	cache = c
	defer func() { cache = oldCache }()

	err := exploreCache("https://pokeapi.co/api/v2/location-area/nonexistent-area-xyz")
	if err == nil {
		t.Error("exploreCache() with cache miss and bad endpoint should return error")
	}
}

func TestDisplayResultPokemon(t *testing.T) {
	result := model.PokemonEncountersList{
		PokemonEncounters: []model.PokemonEncounter{
			{Pokemon: model.Pokemon{Name: "bulbasaur"}},
			{Pokemon: model.Pokemon{Name: "ivysaur"}},
		},
	}

	err := displayResultPokemon(result)
	if err != nil {
		t.Errorf("displayResultPokemon() returned error: %v", err)
	}
}

func TestDisplayResultPokemonEmpty(t *testing.T) {
	result := model.PokemonEncountersList{
		PokemonEncounters: []model.PokemonEncounter{},
	}

	err := displayResultPokemon(result)
	if err != nil {
		t.Errorf("displayResultPokemon() with empty list returned error: %v", err)
	}
}
