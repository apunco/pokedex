package main

import (
	"bufio"
	"errors"
	"fmt"
	"math/rand/v2"
	"os"
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

	//print pokemons
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

func pokemonCaught(baseExperience int) bool {
	return baseExperience < rand.IntN(700)
}
