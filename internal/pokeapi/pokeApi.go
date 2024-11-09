package pokeapi

import (
	"net/http"
)

const (
	baseUrl = "https://pokeapi.co/api/v2"
)

type pokeapi struct {
	client http.Client
}
