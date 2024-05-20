package config

import (
	"encoding/json"
	"os"
	"path/filepath"
)

type CredentialsConfig struct {
	Spotify struct {
		ClientID     string `json:"client_id"`
		ClientSecret string `json:"client_secret"`
	} `json:"spotify"`
}

func NewCredentialsConfig() *CredentialsConfig {
	binary, err := os.Executable()
	if err != nil {
		panic(err)
	}
	binaryDirectory := filepath.Dir(binary)
	configurationPath := filepath.Join(binaryDirectory, "credentials.json")
	if _, fileExistsErr := os.Stat(configurationPath); fileExistsErr == nil {
		configurationFile, fileOpenErr := os.Open(configurationPath)
		if fileOpenErr != nil {
			panic(fileOpenErr)
		}
		defer configurationFile.Close()
		decoder := json.NewDecoder(configurationFile)
		configuration := CredentialsConfig{}
		fileOpenErr = decoder.Decode(&configuration)
		if fileOpenErr != nil {
			panic(fileOpenErr)
		}
		return &configuration
	} else {
		panic(fileExistsErr)
	}
}
