package main

import (
	"bufio"
	"errors"
	"fmt"
	"os"
)

const welcomePrompt = "Welcome to the Pokedex! Please enter a command (type 'help' for assistance): "

func main() {
	scanner := bufio.NewScanner(os.Stdin)
	for {
		fmt.Print(welcomePrompt)

		scanner.Scan()
		input := scanner.Text()

		cleanedInput := cleanInput(input)

		if len(cleanedInput) == 0 {
			continue
		}

		err := userInput(cleanedInput[0])
		if errors.Is(err, errExit) {
			os.Exit(0)
		}
		if err != nil {
			fmt.Println(err)
		}

	}
}
