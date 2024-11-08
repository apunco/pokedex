package main

import (
	"bufio"
	"fmt"
	"os"
	"strings"
)

type cliCommand struct {
	name        string
	description string
	callback    func() error
}

func startRepl() {
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
				err := command.callback()
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
