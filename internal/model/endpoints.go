package model

type Config struct {
	PreviousURL *string
	NextURL     *string
	BaseURL     string
	Arg         string
}

var ConfigData = Config{
	PreviousURL: nil,
	NextURL:     nil,
	BaseURL:     "https://pokeapi.co/api/v2/location-area/",
}
