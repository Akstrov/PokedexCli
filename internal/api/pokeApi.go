package api

import (
	"encoding/json"
	"io"
	"math/rand"
	"net/http"
	"net/url"

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
	Next           string
	Previous       string
	Cashe          *pokecashe.Cashe
	CaughtPokemons map[string]Pokemon
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
	defer res.Body.Close()
	body, err := io.ReadAll(res.Body)
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

type Pokemons struct {
	PokemonEncounters []struct {
		Pokemon struct {
			Name string `json:"name"`
			//URL  string `json:"url"`
		} `json:"pokemon"`
	} `json:"pokemon_encounters"`
}

func GetPokemonsInLocation(config *Config, location string) (Pokemons, error) {
	url := "https://pokeapi.co/api/v2/location-area/" + url.QueryEscape(location)
	res, err := http.Get(url)
	if err != nil {
		return Pokemons{}, err
	}
	defer res.Body.Close()
	body, err := io.ReadAll(res.Body)
	if err != nil {
		return Pokemons{}, err
	}
	cashe := config.Cashe
	pokemons := Pokemons{}
	cashe.Add(url, body)
	err = json.Unmarshal(body, &pokemons)
	if err != nil {
		return Pokemons{}, err
	}
	return pokemons, nil
}

type Pokemon struct {
	Name string `json:"name"`
	Exp  int    `json:"base_experience"`
}

func CatchPokemon(config *Config, pokemonName string) (bool, error) {
	url := "https://pokeapi.co/api/v2/pokemon/" + url.QueryEscape(pokemonName)
	res, err := http.Get(url)
	if err != nil {
		return false, err
	}
	defer res.Body.Close()
	body, err := io.ReadAll(res.Body)
	if err != nil {
		return false, err
	}
	pokemon := Pokemon{}
	err = json.Unmarshal(body, &pokemon)
	if err != nil {
		return false, err
	}
	//check if we caught the pokemon
	exp := pokemon.Exp
	//calculate probability at random
	prob := rand.Intn(exp)
	if prob > 50 {
		return false, nil
	}
	if config.CaughtPokemons == nil {
		config.CaughtPokemons = make(map[string]Pokemon)
	}
	config.CaughtPokemons[pokemonName] = pokemon
	return true, nil
}
