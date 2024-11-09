package main

import (
	"bufio"
	"errors"
	"fmt"
	"os"
)

type cliCommand struct {
	name        string
	description string
	callback    func(*config) error
}

func getCommands() map[string]cliCommand {
	return map[string]cliCommand{
		"help": {
			name:        "help",
			description: "Displays a help message",
			callback:    commandHelp,
		},
		"exit": {
			name:        "exit",
			description: "Exit the Pokedex",
			callback:    commandExit,
		},
		"map": {
			name:        "map",
			description: "Displays the names of the next 20 areas in the Pokemon world",
			callback:    getNextArea,
		},
		"mapb": {
			name:        "map",
			description: "Displays the names of the previous 20 areas in the Pokemon world",
			callback:    getPreviousArea,
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
