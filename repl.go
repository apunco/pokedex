package main

import (
	"bufio"
	"errors"
	"fmt"
	"os"
	"strings"

	pokeapi "github.com/apunco/go/pokedex/internal/pokeapi"
)

type config struct {
	pokeApiClient pokeapi.Client
	prevUrl       *string
	nextUrl       *string
	parameters    []string
}

func startRepl(cfg *config) {
	scanner := bufio.NewScanner(os.Stdin)
	writer := bufio.NewWriter(os.Stdout)

	for {
		writer.WriteString("pokedex > ")
		writer.Flush()

		if scanner.Scan() {
			input := cleanInput(scanner.Text())
			commandName := input[0]

			command, exists := getCommands()[commandName]
			if exists {
				err := validateParameters(input, cfg)
				if err != nil {
					fmt.Println()
					getCommands()["help"].callback(cfg)
					continue
				}

				err = command.callback(cfg)
				if err != nil {
					fmt.Println(err)
				}
				continue
			} else {
				fmt.Println("Unknown command")
				getCommands()["help"].callback(cfg)
				continue
			}
		}

		if err := scanner.Err(); err != nil {
			fmt.Fprintln(os.Stderr, "reading standard input: ", err)
		}
	}
}

func cleanInput(input string) []string {
	inputLower := strings.ToLower(input)
	words := strings.Fields(inputLower)
	return words
}

func validateParameters(input []string, cfg *config) error {
	cfg.parameters = nil
	commandName := input[0]
	parametersCount := getCommands()[commandName].parameterCount

	if parametersCount > 0 {
		if len(input) == 1+parametersCount {
			cfg.parameters = input[1:]
		} else {
			return errors.New("missing parameters")
		}
	}
	return nil
}
