package main

import (
	"bufio"
	"flag"
	"fmt"
	"log"
	"os"
	"os/signal"
	"path/filepath"
	"strings"
	"syscall"
	"time"

	"github.com/BurntSushi/toml"
	"github.com/fatih/color"
)

type Character struct {
	Faces []string `toml:"faces"`
}

type Config struct {
	Characters map[string]Character
}

const outputLines = 10 // Number of output lines to display

func main() {
	message := flag.String("message", "Hello, I'm lil guy!", "Message to display")
	characterName := flag.String("character", "default", "Character to use")
	debug := flag.Bool("debug", false, "Run in debug mode")
	flag.Parse()

	config, err := loadConfig()
	if err != nil {
		fmt.Printf("Error loading config: %v\n", err)
		os.Exit(1)
	}

	if *debug {
		fmt.Printf("Loaded characters: %v\n", config.Characters)
		fmt.Printf("Requested character: %s\n", *characterName)
	}

	character, ok := config.Characters[*characterName]
	if !ok {
		if *debug {
			fmt.Printf("Character '%s' not found. Using default.\n", *characterName)
		}
		character = config.Characters["default"]
	}

	if *debug {
		fmt.Printf("Selected character: %v\n", character)
		fmt.Println("Press Enter to continue...")
		bufio.NewReader(os.Stdin).ReadBytes('\n')
	}

	if len(character.Faces) == 0 {
		fmt.Printf("No faces defined for character '%s'. Using default.\n", *characterName)
		character = config.Characters["default"]
	}

	log.Printf("Selected character: %v\n", character)

	// If default character is also empty, use a fallback
	if len(character.Faces) == 0 {
		character.Faces = []string{"(o_o)"}
	}

	frames := []string{"<", "-", ">", " "}

	// Clear screen and hide cursor
	fmt.Print("\033[2J\033[H\033[?25l")
	defer fmt.Print("\033[?25h") // Show cursor when done

	c := color.New(color.FgCyan)

	sigChan := make(chan os.Signal, 1)
	signal.Notify(sigChan, os.Interrupt, syscall.SIGTERM)

	outputChan := make(chan string)
	go readStdin(outputChan)

	go func() {
		<-sigChan
		fmt.Print("\033[?25h") // Show cursor before exiting
		os.Exit(0)
	}()

	characterIndex := 0
	outputBuffer := make([]string, outputLines)
	for i := range outputBuffer {
		outputBuffer[i] = "" // Initialize with empty strings
	}

	for {
		for _, frame := range frames {
			// Move cursor to top-left and print animation
			fmt.Print("\033[H")

			lines := strings.Split(character.Faces[characterIndex], "\n")
			isMultiLine := len(lines) > 1

			for i, line := range lines {
				if i == 0 && !isMultiLine {
					// Only apply arms to single-line characters
					leftArm, rightArm := getArms(frame)
					c.Printf("  %s %s %s\n", leftArm, line, rightArm)
				} else {
					// For multi-line characters or subsequent lines, don't add arms
					c.Printf("    %s\n", line)
				}
			}

			fmt.Printf("\n  %s\n\n", *message)

			// Display the last few lines of output
			for _, line := range outputBuffer {
				if line != "" {
					fmt.Printf("  %s\n", line)
				}
			}

			// Add padding to cover any previous longer messages
			padding := strings.Repeat(" ", 50)
			fmt.Printf("%s\n", padding)

			time.Sleep(250 * time.Millisecond)

			// Check for new output
			select {
			case newOutput := <-outputChan:
				// Shift the buffer and add the new output
				outputBuffer = append(outputBuffer[1:], newOutput)
			default:
				// No new output, continue with the current buffer
			}
		}

		// Cycle to the next face for the character
		characterIndex = (characterIndex + 1) % len(character.Faces)
	}
}

func loadConfig() (*Config, error) {
	home, err := os.UserHomeDir()
	if err != nil {
		return nil, err
	}

	configPath := filepath.Join(home, ".config", "lil-guy", "characters.toml")
	fmt.Printf("Loading config from: %s\n", configPath)

	// Read and print raw file contents
	rawContent, err := os.ReadFile(configPath)
	if err != nil {
		return nil, fmt.Errorf("error reading config file: %v", err)
	}
	fmt.Println("Raw config file contents:")
	fmt.Println(string(rawContent))

	var config Config
	config.Characters = make(map[string]Character)

	_, err = toml.DecodeFile(configPath, &config.Characters)
	if err != nil {
		return nil, fmt.Errorf("error decoding TOML: %v", err)
	}

	fmt.Printf("Loaded config: %+v\n", config)

	// Ensure default character exists
	if _, ok := config.Characters["default"]; !ok {
		config.Characters["default"] = Character{Faces: []string{"(o_o)"}}
	}

	return &config, nil
}

func readStdin(outputChan chan<- string) {
	scanner := bufio.NewScanner(os.Stdin)
	for scanner.Scan() {
		outputChan <- scanner.Text()
	}
}

func getArms(frame string) (string, string) {
	switch frame {
	case "<":
		return "<", "<"
	case ">":
		return ">", ">"
	default:
		return frame, frame
	}
}
