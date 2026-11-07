package main

import (
	"time"

	"github.com/milinkovicluka/pokedex/internal/pokeapi"
)

func main() {
	client := pokeapi.NewClient(10*time.Second, 5*time.Minute)
	cfg := newConfig(client)
	startRepl(cfg)
}
