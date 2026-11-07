package pokeapi

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
)

const baseURL = "https://pokeapi.co/api/v2"

type LocationAreasResp struct {
	Count    int     `json:"count"`
	Next     *string `json:"next"`
	Previous *string `json:"previous"`
	Results  []struct {
		Name string `json:"name"`
		URL  string `json:"url"`
	} `json:"results"`
}

func (c *Client) ListLocationAreas(pageURL *string) (LocationAreasResp, error) {
	url := baseURL + "/location-area"
	if pageURL != nil {
		url = *pageURL
	}

	if data, ok := c.cache.Get(url); ok {
		fmt.Println("[cache hit]", url)
		var cached LocationAreasResp
		if err := json.Unmarshal(data, &cached); err == nil {
			return cached, nil
		}
	}

	fmt.Println("[network request]", url)
	req, err := http.NewRequest(http.MethodGet, url, nil)
	if err != nil {
		return LocationAreasResp{}, err
	}

	res, err := c.httpClient.Do(req)
	if err != nil {
		return LocationAreasResp{}, err
	}
	defer res.Body.Close()

	if res.StatusCode != http.StatusOK {
		return LocationAreasResp{}, fmt.Errorf("unexpected status: %s", res.Status)
	}

	data, err := io.ReadAll(res.Body)
	if err != nil {
		return LocationAreasResp{}, err
	}

	var areas LocationAreasResp
	if err := json.Unmarshal(data, &areas); err != nil {
		return LocationAreasResp{}, err
	}

	c.cache.Add(url, data)
	return areas, nil
}
