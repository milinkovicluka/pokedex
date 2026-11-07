package main

import (
	"errors"
	"fmt"
)

func commandMap(cfg *config, args []string) error {
	areas, err := cfg.pokeapiClient.ListLocationAreas(cfg.nextURL)
	if err != nil {
		return err
	}
	cfg.nextURL = areas.Next
	cfg.prevURL = areas.Previous

	for _, area := range areas.Results {
		fmt.Println(area.Name)
	}
	return nil
}

func commandMapb(cfg *config, args []string) error {
	if cfg.prevURL == nil {
		return errors.New("you're on the first page")
	}
	areas, err := cfg.pokeapiClient.ListLocationAreas(cfg.prevURL)
	if err != nil {
		return err
	}
	cfg.nextURL = areas.Next
	cfg.prevURL = areas.Previous

	for _, area := range areas.Results {
		fmt.Println(area.Name)
	}
	return nil
}
