package pokemonVault

type MyPokemonList struct {
	Pokemons []PokemonVault
}

type PokemonVault struct {
	PokemonName   string
	HasBeenCaught bool
}

type PokemonVaultMethods interface {
	AddPokemon(pokemonName string)
	GetPokemonCaught() []string
}

var DefaultPokemonVault PokemonVaultMethods = &MyPokemonList{
	Pokemons: []PokemonVault{},
}

func (p *MyPokemonList) AddPokemon(pokemonName string) {
	p.Pokemons = append(p.Pokemons, PokemonVault{
		PokemonName:   pokemonName,
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
