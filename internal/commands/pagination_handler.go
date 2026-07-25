package commands

import (
	"encoding/json"
	"fmt"
	"time"

	"github.com/NicolasFerreras/REPL-Pokedex/internal/client"
	"github.com/NicolasFerreras/REPL-Pokedex/internal/model"
	"github.com/NicolasFerreras/REPL-Pokedex/internal/pokecache"
)

var cache = pokecache.NewCache(15 * time.Second)

func fetchAndDisplay(url string) error {
	// Check cache first
	if data, found := cache.Get(url); found {
		fmt.Println("Data fetched from cache:")
		var result model.LocationArea

		json.Unmarshal(data, &result)

		updatePagination(result)
		return displayResult(result)
	}

	// If not in cache, fetch from API
	c := client.NewClient(url)
	result, err := c.GetLocationArea(url)
	if err != nil {
		return fmt.Errorf("failed to fetch data: %w", err)
	}

	response, err := json.Marshal(result)
	if err != nil {
		return fmt.Errorf("failed to marshal data: %w", err)
	}
	cache.Add(url, response)

	updatePagination(*result)
	return displayResult(*result)

}

func displayResult(result model.LocationArea) error {
	fmt.Println("Location Areas:")
	for _, r := range result.Results {
		fmt.Printf("- %s\n", r.Name)
	}
	return nil
}

func updatePagination(result model.LocationArea) {
	if result.Next != "" {
		model.ConfigData.NextURL = &result.Next
	} else {
		model.ConfigData.NextURL = nil
	}
	if result.Previous != "" {
		model.ConfigData.PreviousURL = &result.Previous
	} else {
		model.ConfigData.PreviousURL = nil
	}
}
