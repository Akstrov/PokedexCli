package main

import (
	"bufio"
	"encoding/json"
	"errors"
	"fmt"
	"net/url"
	"os"
	"strings"
	"time"

	"github.com/Akstrov/PokedexCli/internal/api"
	"github.com/Akstrov/PokedexCli/internal/pokecashe"
)

type cliCommands struct {
	name        string
	description string
	callback    func(c *api.Config, param string) error
}

func getCommands() map[string]cliCommands {
	return map[string]cliCommands{
		"exit": {
			name:        "exit",
			description: "exit the program",
			callback:    commandExit,
		},
		"help": {
			name:        "help",
			description: "displays a help message",
			callback:    commandHelp,
		},
		"map": {
			name:        "map",
			description: "displays the next 20 Location areas",
			callback:    commandMap,
		},
		"mapb": {
			name:        "mapb",
			description: "displays the previous 20 Location areas",
			callback:    commandMapB,
		},
		"explore": {
			name:        "explore",
			description: "explore the current Location area",
			callback:    commandExplore,
		},
	}
}

func printLocationAreas(config *api.Config, locations api.LocationAreas) {
	if config.Previous == "" {
		config.Previous = config.Next
		config.Next = locations.Next
	} else if config.Next == "" {
		config.Next = config.Previous
		config.Previous = locations.Previous
	} else {
		config.Next = locations.Next
		config.Previous = locations.Previous
	}
	for _, location := range locations.Results {
		fmt.Println(location.Name)
	}
}

func commandExplore(config *api.Config, param string) error {
	if param == "" {
		return errors.New("explore requires a Location name")
	}
	fmt.Printf("Exploring %s...\n", param)
	cashe := config.Cashe
	url := "https://pokeapi.co/api/v2/location-area/" + url.QueryEscape(param)
	if data, ok := cashe.Get(url); ok {
		pokemons := api.Pokemons{}
		err := json.Unmarshal(data, &pokemons)
		if err != nil {
			return err
		}
		for _, pokemon := range pokemons.PokemonEncounters {
			fmt.Printf(" - %s\n", pokemon.Pokemon.Name)
		}
		return nil
	}
	pokemons, err := api.GetPokemonsInLocation(config, param)
	if err != nil {
		return err
	}
	for _, pokemon := range pokemons.PokemonEncounters {
		fmt.Printf(" - %s\n", pokemon.Pokemon.Name)
	}
	return nil
}

func commandMap(config *api.Config, param string) error {
	if param != "" {
		return errors.New("map requires no parameters")
	}
	if config.Next == "" {
		return errors.New("no next location area")
	}
	cashe := config.Cashe
	if data, ok := cashe.Get(config.Next); ok {
		locations := api.LocationAreas{}
		err := json.Unmarshal(data, &locations)
		if err != nil {
			return err
		}
		printLocationAreas(config, locations)
		return nil
	}
	locations, err := api.GetLocationAreas(config, "next")
	if err != nil {
		return err
	}
	printLocationAreas(config, locations)
	return nil
}
func commandMapB(config *api.Config, param string) error {
	if param != "" {
		return errors.New("mapb requires no parameters")
	}
	if config.Previous == "" {
		return errors.New("no previous location area")
	}
	cashe := config.Cashe
	if data, ok := cashe.Get(config.Previous); ok {
		locations := api.LocationAreas{}
		err := json.Unmarshal(data, &locations)
		if err != nil {
			return err
		}
		printLocationAreas(config, locations)
		return nil
	}
	locations, err := api.GetLocationAreas(config, "previous")
	if err != nil {
		return err
	}
	printLocationAreas(config, locations)
	return nil
}

func commandExit(config *api.Config, param string) error {
	if param != "" {
		return errors.New("exit requires no parameters")
	}
	os.Exit(0)
	return nil
}
func commandHelp(config *api.Config, param string) error {
	if param != "" {
		return errors.New("help requires no parameters")
	}
	fmt.Printf("\nWelcome to the Pokedex!\nUsage:\n\n")
	for _, command := range getCommands() {
		fmt.Printf("%s: %s\n", command.name, command.description)
	}
	fmt.Printf("\n")
	return nil
}

func main() {
	scanner := bufio.NewScanner(os.Stdin)
	commands := getCommands()
	config := api.Config{
		Next:     "https://pokeapi.co/api/v2/location-area/?offset=0&limit=20",
		Previous: "",
		Cashe:    pokecashe.NewCashe(5 * time.Minute),
	}
	for {
		fmt.Print("pokedex > ")
		scanner.Scan()
		input := scanner.Text()
		inputs := strings.Split(input, " ")
		if len(inputs) == 0 {
			fmt.Println("invalid command")
			continue
		}
		if len(inputs) == 1 {
			inputs = append(inputs, "")
		}
		command := inputs[0]
		value, ok := commands[command]
		if ok {
			err := value.callback(&config, inputs[1])
			if err != nil {
				fmt.Println(err)
			}
		} else {
			fmt.Println("invalid command")
		}
	}
}
