package api

import (
	"encoding/json"
	"io"
	"net/http"

	"github.com/Akstrov/PokedexCli/internal/pokecashe"
)

type LocationAreas struct {
	Count    int    `json:"count"`
	Next     string `json:"next"`
	Previous string `json:"previous"`
	Results  []struct {
		Name string `json:"name"`
		URL  string `json:"url"`
	} `json:"results"`
}

type Config struct {
	Next     string
	Previous string
	Cashe    *pokecashe.Cashe
}

func GetLocationAreas(config *Config, direction string) (LocationAreas, error) {
	var url string
	if direction == "next" {
		url = config.Next
	} else {
		url = config.Previous
	}
	res, err := http.Get(url)
	if err != nil {
		return LocationAreas{}, err
	}
	body, err := io.ReadAll(res.Body)
	res.Body.Close()
	if err != nil {
		return LocationAreas{}, err
	}
	cashe := config.Cashe
	cashe.Add(url, body)
	locations := LocationAreas{}
	err = json.Unmarshal(body, &locations)
	if err != nil {
		return LocationAreas{}, err
	}

	return locations, nil
}
