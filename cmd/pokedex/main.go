package main

import (
	"bufio"
	"errors"
	"fmt"
	"os"

	"github.com/NicolasFerreras/REPL-Pokedex/internal/commands"
	"github.com/NicolasFerreras/REPL-Pokedex/internal/repl"
)

const welcomePrompt = "Welcome to the Pokedex! Please enter a command (type 'help' for assistance): "

func main() {
	scanner := bufio.NewScanner(os.Stdin)
	for {
		fmt.Print(welcomePrompt)

		scanner.Scan()
		input := scanner.Text()

		cleanedInput := repl.CleanInput(input)

		if len(cleanedInput) == 0 {
			continue
		}

		err := commands.UserInput(cleanedInput[0])
		if errors.Is(err, commands.ErrExit) {
			os.Exit(0)
		}
		if err != nil {
			fmt.Println(err)
		}

	}
}
