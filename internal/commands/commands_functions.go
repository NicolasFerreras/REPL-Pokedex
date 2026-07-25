package commands

import (
	"errors"
	"fmt"

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

func CommandExplore() error {
	arg := model.ConfigData.Arg
	url := model.ConfigData.BaseURL + arg

	return exploreCache(url)
}
