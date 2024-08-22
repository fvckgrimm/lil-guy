package main

import (
	"fmt"
	"os"
	"path/filepath"

	"github.com/BurntSushi/toml"
)

type Config struct {
	Characters map[string]Character
}

func loadConfig() (*Config, error) {
	home, err := os.UserHomeDir()
	if err != nil {
		return nil, err
	}

	configPath := filepath.Join(home, ".config", "lil-guy", "characters.toml")
	fmt.Printf("Loading config from: %s\n", configPath)

	var config Config
	config.Characters = make(map[string]Character)

	_, err = toml.DecodeFile(configPath, &config.Characters)
	if err != nil {
		return nil, fmt.Errorf("error decoding TOML: %v", err)
	}

	// Ensure default character exists
	if _, ok := config.Characters["default"]; !ok {
		config.Characters["default"] = Character{Faces: []string{"(o_o)"}}
	}

	return &config, nil
}
