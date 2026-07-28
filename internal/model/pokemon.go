package model

type PokemonDetails struct {
	BaseExperience int     `json:"base_experience"`
	Name           string  `json:"name"`
	Height         int     `json:"height"`
	Weight         int     `json:"weight"`
	Stats          []Stats `json:"stats"`
	Types          []Types `json:"types"`
}

type Stat struct {
	Name string `json:"name"`
}
type Stats struct {
	BaseStat int  `json:"base_stat"`
	Stat     Stat `json:"stat"`
}

type Type struct {
	Name string `json:"name"`
}
type Types struct {
	Slot int  `json:"slot"`
	Type Type `json:"type"`
}
