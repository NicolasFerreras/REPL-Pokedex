package pokemonVault

import (
	"fmt"
	"strings"

	"github.com/NicolasFerreras/REPL-Pokedex/internal/model"
)

type MyPokemonList struct {
	Pokemons []PokemonVault
}

type PokemonVault struct {
	PokemonName   string
	PokemonHeight int
	PokemonWeight int
	PokemonStats  []PokemonStat
	PokemonTypes  []PokemonType
	HasBeenCaught bool
}

type PokemonStat struct {
	HP                 int
	Attack             int
	Defense            int
	SpecialAttack      int
	SpecialDefense     int
	Speed              int
}

type PokemonType struct {
	TypeName string
}

type PokemonVaultMethods interface {
	AddPokemon(pokemon model.PokemonDetails)
	GetPokemonCaught() []string
	GetPokemonDetails(pokemonName string) PokemonVault
	DisplayPokemonDetails(pokemon PokemonVault) string
}

var DefaultPokemonVault PokemonVaultMethods = &MyPokemonList{
	Pokemons: []PokemonVault{},
}

func (p *MyPokemonList) AddPokemon(pokemon model.PokemonDetails) {
	p.Pokemons = append(p.Pokemons, PokemonVault{
		PokemonName:   pokemon.Name,
		PokemonHeight: pokemon.Height,
		PokemonWeight: pokemon.Weight,
		PokemonStats:  convertStats(pokemon.Stats),
		PokemonTypes:  convertTypes(pokemon.Types),
		HasBeenCaught: true,
	})
}

func (p MyPokemonList) GetPokemonCaught() []string {
	var caughtPokemons []string
	for _, pokemon := range p.Pokemons {
		if pokemon.HasBeenCaught {
			caughtPokemons = append(caughtPokemons, pokemon.PokemonName)
		}
	}
	return caughtPokemons
}

func (p MyPokemonList) GetPokemonDetails(pokemonName string) PokemonVault {
	for _, pokemon := range p.Pokemons {
		if strings.EqualFold(pokemon.PokemonName, pokemonName) {
			return pokemon
		}
	}
	return PokemonVault{}
}

func (p MyPokemonList) DisplayPokemonDetails(pokemon PokemonVault) string {
	var sb strings.Builder
	sb.WriteString(fmt.Sprintf("Name: %s\n", pokemon.PokemonName))
	sb.WriteString(fmt.Sprintf("Height: %d\n", pokemon.PokemonHeight))
	sb.WriteString(fmt.Sprintf("Weight: %d\n", pokemon.PokemonWeight))
	sb.WriteString("Stats:\n")
	for _, stat := range pokemon.PokemonStats {
		sb.WriteString(fmt.Sprintf("- HP: %d\n", stat.HP))
		sb.WriteString(fmt.Sprintf("- Attack: %d\n", stat.Attack))
		sb.WriteString(fmt.Sprintf("- Defense: %d\n", stat.Defense))
		sb.WriteString(fmt.Sprintf("- Special Attack: %d\n", stat.SpecialAttack))
		sb.WriteString(fmt.Sprintf("- Special Defense: %d\n", stat.SpecialDefense))
		sb.WriteString(fmt.Sprintf("- Speed: %d\n", stat.Speed))
	}
	sb.WriteString("Types:\n")
	for _, t := range pokemon.PokemonTypes {
		sb.WriteString(fmt.Sprintf("- %s\n", t.TypeName))
	}
	sb.WriteString("\n")
	return sb.String()
}

func convertStats(stats []model.Stats) []PokemonStat {
	pokemonStats := PokemonStat{}
	for _, stat := range stats {
		switch strings.ToLower(stat.Stat.Name) {
		case "hp":
			pokemonStats.HP = stat.BaseStat
		case "attack":
			pokemonStats.Attack = stat.BaseStat
		case "defense":
			pokemonStats.Defense = stat.BaseStat
		case "special-attack":
			pokemonStats.SpecialAttack = stat.BaseStat
		case "special-defense":
			pokemonStats.SpecialDefense = stat.BaseStat
		case "speed":
			pokemonStats.Speed = stat.BaseStat
		}
	}
	return []PokemonStat{pokemonStats}
}

func convertTypes(types []model.Types) []PokemonType {
	var pokemonTypes []PokemonType
	for _, t := range types {
		pokemonTypes = append(pokemonTypes, PokemonType{
			TypeName: t.Type.Name,
		})
	}
	return pokemonTypes
}
