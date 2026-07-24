package commands

import (
	"errors"
	"fmt"

	"github.com/NicolasFerreras/REPL-Pokedex/internal/client"
	"github.com/NicolasFerreras/REPL-Pokedex/internal/model"
)

// ErrExit is a sentinel error used to signal the REPL should exit
var ErrExit = errors.New("exit")

// CommandExit exits the Pokedex
func CommandExit() error {
	fmt.Println("Closing the Pokedex... Goodbye!")
	return ErrExit
}

// CommandHelp displays available commands
func CommandHelp() error {
	for _, cmd := range GetCommands() {
		fmt.Printf("%s: %s\n", cmd.Name, cmd.Description)
	}
	return nil
}

// UserInput processes user input and executes the corresponding command
func UserInput(input string) error {
	if len(input) == 0 {
		return fmt.Errorf("no command entered. Please enter a command")
	}

	cmd, exists := GetCommands()[input]
	if !exists {
		return fmt.Errorf("unknown command: %s. Type 'help' for available commands", input)
	}

	return cmd.Callback()
}
func CommandMap() error {
	url := model.ConfigData.BaseURL
	if model.ConfigData.NextURL != nil {
		url = *model.ConfigData.NextURL
	}
	return fetchAndDisplay(url)
}

func CommandMapBack() error {
	if model.ConfigData.PreviousURL == nil {
		fmt.Println("you're on the first page")
		return nil
	}
	return fetchAndDisplay(*model.ConfigData.PreviousURL)
}

func fetchAndDisplay(url string) error {
	c := client.NewClient(url)
	result, err := c.GetLocationArea(url)
	if err != nil {
		return fmt.Errorf("failed to get location area: %w", err)
	}

	// Actualizar paginación
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

	fmt.Println("Location Areas:")
	for _, r := range result.Results {
		fmt.Printf("- %s\n", r.Name)
	}
	return nil
}
