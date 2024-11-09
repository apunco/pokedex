package main

import (
	"bufio"
	"fmt"
	"os"
	"strings"

	pokeapi "github.com/apunco/go/pokedex/internal/pokeapi"
)

type config struct {
	pokeApiClient pokeapi.Client
	prevUrl       *string
	nextUrl       *string
	location      string
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

			switch commandName {
			case "explore":
				if len(input) > 0 {
					cfg.location = input[1]
				} else {
					fmt.Println("missing location parameter")
					continue
				}
			}

			command, exists := getCommands()[commandName]
			if exists {
				err := command.callback(cfg)
				if err != nil {
					fmt.Println(err)
				}
				continue
			} else {
				fmt.Println("Unknown command")
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
