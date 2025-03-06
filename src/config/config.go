package config

import (
	"fmt"
	"gopkg.in/yaml.v3"
	"os"
	"path/filepath"
)

type ConfigFile struct {
	MangaDexUrl string `yaml:"mangadex_url"`
}

func LoadConfig(path string) (*ConfigFile, error) {
	configFile := &ConfigFile{}
	pathToConfig := filepath.Join(path, "config.yml")
	file, err := os.Open(pathToConfig)
	if err != nil {
		return nil, fmt.Errorf("[Error] Could not load config.yml: %w", err.Error())
	}
	defer file.Close()

	decoder := yaml.NewDecoder(file)
	err = decoder.Decode(configFile)
	if err != nil {
		return nil, fmt.Errorf("Could not decode config.yml: %w", err.Error())
	}
	return configFile, nil
}
