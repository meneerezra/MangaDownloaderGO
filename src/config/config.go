package config

import (
	"fmt"
	"gopkg.in/yaml.v3"
	"os"
)

type ConfigFile struct {
	MangaDexUrl string `yaml:"mangadex_url"`
	DownloadPath string `yaml:"download_path"`
}

func LoadConfig(path string) (*ConfigFile, error) {
	configFile := &ConfigFile{}
	file, err := os.Open(path)
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

func GenerateConfig(path string, downloadPath string) error {
	configFile, err := os.Create(path)
	if err != nil {
		return err
	}

	basicConfig := "mangadex_url: 'https://api.mangadex.org'" +
		"\ndownload_path: '" + downloadPath + "'"

	configFile.Write([]byte(basicConfig))
	return nil
}
