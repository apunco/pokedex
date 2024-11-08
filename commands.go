package main

import (
	"bufio"
	"os"
)

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
	}
}

func commandHelp() error {
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

func commandExit() error {
	os.Exit(0)
	return nil
}
