package config

import (
	"encoding/json"
	"os"
)

const configFileName = "/.gatorconfig.json"

func getConfigFilePath() (string, error) {
	home, err := os.UserHomeDir()
	jsonLoc := home + configFileName
	if err != nil {
		return "", err
	}
	return jsonLoc, nil
}

func Read() (Config, error) {
	jsonLoc, err := getConfigFilePath()
	if err != nil {
		return Config{}, err
	}
	jsonBytes, err := os.ReadFile(jsonLoc)
	if err != nil {
		return Config{}, err
	}
	cfg := Config{}
	err = json.Unmarshal(jsonBytes, &cfg)
	if err != nil {
		return Config{}, err
	}
	return cfg, nil
}
