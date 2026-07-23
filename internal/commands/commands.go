package commands

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
