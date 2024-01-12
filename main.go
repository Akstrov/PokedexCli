package main

import (
	"bufio"
	"errors"
	"fmt"
	"os"
	"github.com/Akstrov/PokedexCli/internal"
)

type cliCommands struct {
	name        string
	description string
	callback    func(c *internal.Config) error
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

func commandMap(config *internal.Config) error {
	if config.Next == "" {
		return errors.New("No next location area")
	}
	locations, err := internal.GetLocationAreas(config, "next")
	if err != nil {
		return err
	}
	config.Next = locations.Next
	config.Previous = locations.Previous
	for _, location := range locations.Results {
		fmt.Println(location.Name)
	}
	return nil
}
func commandMapB(config *internal.Config) error {
	if config.Previous == "" {
		return errors.New("no previous location area")
	}
	locations, err := internal.GetLocationAreas(config, "previous")
	if err != nil {
		return err
	}
	config.Next = locations.Next
	config.Previous = locations.Previous
	for _, location := range locations.Results {
		fmt.Println(location.Name)
	}
	return nil
}

func commandExit(config *internal.Config) error {
	os.Exit(0)
	return nil
}
func commandHelp(config *internal.Config) error {
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
	config := internal.Config{
		Next:     "https://pokeapi.co/api/v2/location-area/",
		Previous: "",
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
