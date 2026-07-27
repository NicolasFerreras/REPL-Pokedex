package main

import (
	"bufio"
	"errors"
	"fmt"
	"os"

	"github.com/NicolasFerreras/REPL-Pokedex/internal/commands"
	errs "github.com/NicolasFerreras/REPL-Pokedex/internal/errors"
	"github.com/NicolasFerreras/REPL-Pokedex/internal/model"
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
		if cleanedInput[0] == "explore" {
			if len(cleanedInput) < 2 {
				fmt.Println(fmt.Errorf(errs.ErrCommandNeedsArg, "explore"))
				continue
			}
			model.ConfigData.Arg = cleanedInput[1]
		}

		if cleanedInput[0] == "catch" {
			if len(cleanedInput) < 2 {
				fmt.Println(fmt.Errorf(errs.ErrCommandNeedsArg, "catch"))
				continue
			}
			model.ConfigData.Arg = cleanedInput[1]
		}

		err := commands.UserInput(cleanedInput[0])
		if errors.Is(err, commands.ErrExit) {
			os.Exit(0)
		}
		if err != nil {
			fmt.Println(err)
		}

		if err := scanner.Err(); err != nil {
			fmt.Println(fmt.Errorf(errs.ErrScannerIO, err))
		}
	}
}
