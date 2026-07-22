package commands

import (
	"errors"
	"fmt"
)

// ErrExit is a sentinel error used to signal the REPL should exit
var ErrExit = errors.New("exit")

// CliCommand represents a CLI command with name, description and callback
type CliCommand struct {
	Name        string
	Description string
	Callback    func() error
}

// Commands is defined at package level to avoid recreation on every call
var Commands = map[string]CliCommand{
	"exit": {
		Name:        "exit",
		Description: "Exit the Pokedex",
		Callback:    CommandExit,
	},
	"help": {
		Name:        "help",
		Description: "Display help information",
		Callback:    CommandHelp,
	},
}

// CommandExit exits the Pokedex
func CommandExit() error {
	fmt.Println("Closing the Pokedex... Goodbye!")
	return ErrExit
}

// CommandHelp displays available commands
func CommandHelp() error {
	fmt.Println("Available commands:")
	fmt.Println("  help - Display help information")
	fmt.Println("  exit - Exit the Pokedex")
	return nil
}

// UserInput processes user input and executes the corresponding command
func UserInput(input string) error {
	if len(input) == 0 {
		return fmt.Errorf("no command entered. Please enter a command")
	}

	cmd, exists := Commands[input]
	if !exists {
		return fmt.Errorf("unknown command: %s. Type 'help' for available commands", input)
	}

	return cmd.Callback()
}
