package commands

// CliCommand represents a CLI command with name, description and callback
type CliCommand struct {
	Name        string
	Description string
	Callback    func() error
}

// Commands is defined at package level to avoid recreation on every call

func GetCommands() map[string]CliCommand {
	return map[string]CliCommand{
		"help": {
			Name:        "help",
			Description: "Display help information",
			Callback:    CommandHelp,
		},
		"exit": {
			Name:        "exit",
			Description: "Exit the Pokedex",
			Callback:    CommandExit,
		},
		"map": {
			Name:        "map",
			Description: "Display location areas",
			Callback:    CommandMap,
		},
		"mapb": {
			Name:        "mapb",
			Description: "Display previous location areas",
			Callback:    CommandMapBack,
		},
		"explore": {
			Name:        "explore",
			Description: "Explore a specific location area",
			Callback:    CommandExplore,
		},
		"catch": {
			Name:        "catch",
			Description: "Catch a specific Pokemon",
			Callback:    CommandCatch,
		},
		"inspect": {
			Name:        "inspect",
			Description: "Inspect a specific Pokemon",
			Callback:    CommandInspect,
		},
	}
}
