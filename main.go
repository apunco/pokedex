package main

import (
	"time"

	"github.com/apunco/go/pokedex/internal/pokeapi"
)

func main() {

	pokeClient := pokeapi.NewClient(5 * time.Second)
	cfg := &config{
		pokeApiClient: pokeClient,
	}
	startRepl(cfg)
}
