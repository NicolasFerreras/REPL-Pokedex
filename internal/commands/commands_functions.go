package commands

import (
	"errors"
	"fmt"
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
