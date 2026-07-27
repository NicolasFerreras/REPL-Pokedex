package model

type Config struct {
	PreviousURL *string
	NextURL     *string
	BaseURL     string
	Arg         string
	PokemonURL  string
}

var ConfigData = Config{
	PreviousURL: nil,
	NextURL:     nil,
	BaseURL:     "https://pokeapi.co/api/v2/location-area/",
	PokemonURL:  "https://pokeapi.co/api/v2/pokemon/",
}
