package main

import (
	"bufio"
	"encoding/json"
	"errors"
	"fmt"
	"os"
	"time"

	"github.com/Akstrov/PokedexCli/internal/api"
	"github.com/Akstrov/PokedexCli/internal/pokecashe"
)

type cliCommands struct {
	name        string
	description string
	callback    func(c *api.Config) error
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

func commandMap(config *api.Config) error {
	if config.Next == "" {
		return errors.New("no next location area")
	}
	cashe := config.Cashe
	if data, ok := cashe.Get(config.Next); ok {
		fmt.Println("wee in")
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
func commandMapB(config *api.Config) error {
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

func commandExit(config *api.Config) error {
	os.Exit(0)
	return nil
}
func commandHelp(config *api.Config) error {
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
		value, ok := commands[input]
		if ok {
			err := value.callback(&config)
			if err != nil {
				fmt.Println(err)
			}
		} else {
			fmt.Println("invalid command")
		}
	}
}
