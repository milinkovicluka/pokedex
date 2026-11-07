package main

import (
	"bufio"
	"fmt"
	"os"
	"strings"

	"github.com/milinkovicluka/pokedex/internal/pokeapi"
)

type config struct {
	commands      map[string]cliCommand
	order         []string
	pokeapiClient pokeapi.Client
	nextURL       *string
	prevURL       *string
	pokedex       map[string]pokeapi.Pokemon
}

type cliCommand struct {
	name        string
	description string
	callback    func(cfg *config, args []string) error
}

func newConfig(client pokeapi.Client) *config {
	list := []cliCommand{
		{
			name:        "help",
			description: "Displays a help message",
			callback:    commandHelp,
		},
		{
			name:        "map",
			description: "Displays the next 20 location areas",
			callback:    commandMap,
		},
		{
			name:        "mapb",
			description: "Displays the previous 20 location areas",
			callback:    commandMapb,
		},
		{
			name:        "explore",
			description: "Lists the Pokemon in a location area: explore <area-name>",
			callback:    commandExplore,
		},
		{
			name:        "catch",
			description: "Try to catch a Pokemon: catch <pokemon-name>",
			callback:    commandCatch,
		},
		{
			name:        "inspect",
			description: "Shows details of a caught Pokemon: inspect <pokemon-name>",
			callback:    commandInspect,
		},
		{
			name:        "pokedex",
			description: "Lists all the Pokemon you have caught",
			callback:    commandPokedex,
		},
		{
			name:        "exit",
			description: "Exit the Pokedex",
			callback:    commandExit,
		},
	}

	cfg := &config{
		commands:      make(map[string]cliCommand),
		pokeapiClient: client,
		pokedex:       make(map[string]pokeapi.Pokemon),
	}
	for _, cmd := range list {
		cfg.commands[cmd.name] = cmd
		cfg.order = append(cfg.order, cmd.name)
	}
	return cfg
}

func startRepl(cfg *config) {
	scanner := bufio.NewScanner(os.Stdin)
	for {
		fmt.Print("Pokedex > ")
		if !scanner.Scan() {
			break
		}
		words := cleanInput(scanner.Text())
		if len(words) == 0 {
			continue
		}
		command, ok := cfg.commands[words[0]]
		if !ok {
			fmt.Println("Unknown command")
			continue
		}
		if err := command.callback(cfg, words[1:]); err != nil {
			fmt.Println(err)
		}
	}
}

func cleanInput(text string) []string {
	return strings.Fields(strings.ToLower(text))
}
