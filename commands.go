package main

import (
	"errors"
	"fmt"
)

// errExit is a sentinel error used to signal the REPL should exit
var errExit = errors.New("exit")

type cliCommand struct {
	name        string
	description string
	callback    func() error
}

// commands is defined at package level to avoid recreation on every call
var commands = map[string]cliCommand{
	"exit": {
		name:        "exit",
		description: "Exit the Pokedex",
		callback:    commandExit,
	},
	"help": {
		name:        "help",
		description: "Display help information",
		callback:    commandHelp,
	},
}

func commandExit() error {
	fmt.Println("Closing the Pokedex... Goodbye!")
	return errExit
}

func commandHelp() error {
	fmt.Println("Available commands:")
	fmt.Println("  help - Display help information")
	fmt.Println("  exit - Exit the Pokedex")
	return nil
}

func userInput(input string) error {
	if len(input) == 0 {
		return fmt.Errorf("no command entered. Please enter a command")
	}

	cmd, exists := commands[input]
	if !exists {
		return fmt.Errorf("unknown command: %s. Type 'help' for available commands", input)
	}

	return cmd.callback()
}
