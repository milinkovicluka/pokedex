package pokeapi

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
)

type LocationArea struct {
	Name              string `json:"name"`
	PokemonEncounters []struct {
		Pokemon struct {
			Name string `json:"name"`
			URL  string `json:"url"`
		} `json:"pokemon"`
	} `json:"pokemon_encounters"`
}

func (c *Client) GetLocationArea(name string) (LocationArea, error) {
	url := baseURL + "/location-area/" + name

	if data, ok := c.cache.Get(url); ok {
		var cached LocationArea
		if err := json.Unmarshal(data, &cached); err == nil {
			return cached, nil
		}
	}

	req, err := http.NewRequest(http.MethodGet, url, nil)
	if err != nil {
		return LocationArea{}, err
	}

	res, err := c.httpClient.Do(req)
	if err != nil {
		return LocationArea{}, err
	}
	defer res.Body.Close()

	if res.StatusCode == http.StatusNotFound {
		return LocationArea{}, fmt.Errorf("location area %q not found", name)
	}
	if res.StatusCode != http.StatusOK {
		return LocationArea{}, fmt.Errorf("unexpected status: %s", res.Status)
	}

	data, err := io.ReadAll(res.Body)
	if err != nil {
		return LocationArea{}, err
	}

	var area LocationArea
	if err := json.Unmarshal(data, &area); err != nil {
		return LocationArea{}, err
	}

	c.cache.Add(url, data)
	return area, nil
}
