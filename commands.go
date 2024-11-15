package main

import (
	"bufio"
	"errors"
	"fmt"
	"math/rand/v2"
	"os"

	pokeapi "github.com/apunco/go/pokedex/internal/pokeapi"
)

type cliCommand struct {
	name           string
	description    string
	callback       func(*config) error
	parameterCount int
}

func getCommands() map[string]cliCommand {
	return map[string]cliCommand{
		"help": {
			name:           "help",
			description:    "Displays a help message",
			callback:       commandHelp,
			parameterCount: 0,
		},
		"exit": {
			name:           "exit",
			description:    "Exit the Pokedex",
			callback:       commandExit,
			parameterCount: 0,
		},
		"map": {
			name:           "map",
			description:    "Displays the names of the next 20 areas in the Pokemon world",
			callback:       getNextArea,
			parameterCount: 0,
		},
		"mapb": {
			name:           "mapb",
			description:    "Displays the names of the previous 20 areas in the Pokemon world",
			callback:       getPreviousArea,
			parameterCount: 0,
		},
		"explore": {
			name:           "explore",
			description:    "Displays pokemon found at the location sent as a parameter: explore {location}",
			callback:       exploreLocation,
			parameterCount: 1,
		},
		"catch": {
			name:           "catch",
			description:    "Attempts to catch a pokemon and add it to a pokedex",
			callback:       catchPokemon,
			parameterCount: 1,
		},
		"inspect": {
			name:           "inspect",
			description:    "Inspect a pokemon from your pokedex",
			callback:       inspectPokemon,
			parameterCount: 1,
		},
		"pokedex": {
			name:           "pokedex",
			description:    "Returns all pokemons currently in the pokedex",
			callback:       inspectPokedex,
			parameterCount: 0,
		},
	}
}

func commandHelp(cfg *config) error {
	commands := getCommands()
	writer := bufio.NewWriter(os.Stdout)

	writer.WriteString("\nWelcome to the Pokedex!\nUsage:\n\n")
	for key, value := range commands {
		writer.WriteString(key + ":" + value.description + "\n")
	}

	writer.WriteString("\n")
	writer.Flush()
	return nil
}

func commandExit(cfg *config) error {
	os.Exit(0)
	return nil
}

func getNextArea(cfg *config) error {

	locations, err := cfg.pokeApiClient.GetLocations(cfg.nextUrl)
	if err != nil {
		return err
	}

	for _, val := range locations.Results {
		fmt.Println(val.Name)
	}

	if locations.Next == "" {
		cfg.nextUrl = nil
	} else {
		cfg.nextUrl = &locations.Next
	}

	if locations.Previous == "" {
		cfg.prevUrl = nil
	} else {
		cfg.prevUrl = &locations.Previous
	}

	return nil
}

func getPreviousArea(cfg *config) error {

	if cfg.prevUrl == nil {
		return errors.New("you are on the first page")
	}

	locations, err := cfg.pokeApiClient.GetLocations(cfg.prevUrl)
	if err != nil {
		return err
	}

	for _, val := range locations.Results {
		fmt.Println(val.Name)
	}

	cfg.nextUrl = &locations.Next
	cfg.prevUrl = &locations.Previous

	return nil
}

func exploreLocation(cfg *config) error {
	loc, err := cfg.pokeApiClient.ExploreLocation(cfg.parameters[0])
	if err != nil {
		return err
	}

	fmt.Println("Exploring " + cfg.parameters[0])
	fmt.Println("Found Pokemon:")
	for _, val := range loc.PokemonEncounters {
		println("- " + val.Pokemon.Name)
	}

	return nil
}

func catchPokemon(cfg *config) error {
	pokemonName := cfg.parameters[0]
	fmt.Println("Throwing a Pokeball at " + pokemonName)

	pokemon, err := cfg.pokeApiClient.GetPokemon(pokemonName)
	if err != nil {
		return err
	}

	if pokemonCaught(pokemon.BaseExperience) {
		fmt.Println(pokemonName + " was caught!")
		cfg.pokedex[pokemonName] = pokemon
	} else {
		fmt.Println(pokemonName + " escaped!")
	}

	return nil
}

func inspectPokemon(cfg *config) error {
	pokemonName := cfg.parameters[0]

	if _, ok := cfg.pokedex[pokemonName]; !ok {
		return errors.New("Pokemon " + pokemonName + " has not been caught yet!")
	}

	pokemon := cfg.pokedex[pokemonName]

	fmt.Println("Name: " + pokemon.Name)
	fmt.Printf("Height: %d\n", pokemon.Height)
	fmt.Printf("Weight: %d\n", pokemon.Weight)
	fmt.Println("Stats:")
	printPokemonStats(pokemon.Stats)
	fmt.Println("Types:")
	printPokemonTypes(pokemon.Types)

	return nil
}

func pokemonCaught(baseExperience int) bool {
	return baseExperience < rand.IntN(700)
}

func printPokemonStats(pokemonStats []pokeapi.StatDetail) {
	for _, stat := range pokemonStats {
		fmt.Printf("  -%s: %d\n", stat.StatType.Name, stat.BaseStat)
	}
}

func printPokemonTypes(pokemonTypes []pokeapi.Types) {
	for _, pokemonType := range pokemonTypes {
		fmt.Printf("  -%s\n", pokemonType.Type.Name)
	}
}

func inspectPokedex(cfg *config) error {
	fmt.Println("Your Pokedex:")

	for key, _ := range cfg.pokedex {
		fmt.Printf(" - %s\n", key)
	}

	return nil
}
