package main

import (
	"bufio"
	"flag"
	"fmt"
	//"io"
	"os"
	//"strings"
	//"sync"

	tea "github.com/charmbracelet/bubbletea"
)

func main() {
	message := flag.String("message", "Hello, I'm lil guy!", "Message to display")
	flag.StringVar(message, "m", "Hello, I'm lil guy!", "Message to display (shorthand)")
	characterName := flag.String("character", "default", "Character to use")
	flag.StringVar(characterName, "c", "default", "Character to use (shorthand)")
	debug := flag.Bool("debug", false, "Run in debug mode")
	flag.Parse()

	config, err := loadConfig()
	if err != nil {
		fmt.Printf("Error loading config: %v\n", err)
		logError(err)
		os.Exit(1)
	}

	character, ok := config.Characters[*characterName]
	if !ok {
		if *debug {
			fmt.Printf("Character '%s' not found. Using default.\n", *characterName)
		}
		character = config.Characters["default"]
	}

	// Create channels for communication
	inputChan := make(chan string)
	doneChan := make(chan struct{})

	// Start the Bubble Tea program
	m := initialModel(character, *message, inputChan)
	p := tea.NewProgram(m)

	// Handle input in a separate goroutine
	go func() {
		defer close(inputChan)
		scanner := bufio.NewScanner(os.Stdin)
		for scanner.Scan() {
			select {
			case <-doneChan:
				return
			default:
				inputChan <- scanner.Text()
			}
		}
	}()

	// Run the program
	if _, err := p.Run(); err != nil {
		fmt.Printf("Error running program: %v\n", err)
		logError(err)
	}

	// Signal the input goroutine to stop
	close(doneChan)
}
