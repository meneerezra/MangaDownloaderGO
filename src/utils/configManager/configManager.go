package configManager

import (
	"fmt"
	"gopkg.in/yaml.v3"
	"os"
)

type Config struct {
	MangaDexUrl string `yaml:"mangadex_url"`
	DownloadPath string `yaml:"download_path"`
	LogPath string `yaml:"log_path"`
}

func LoadConfig(path string) (*Config, error) {
	configFile := &Config{}
	file, err := os.Open(path)
	if err != nil {
		return nil, fmt.Errorf("[Error] Could not load configManager.yml: %w", err.Error())
	}
	defer file.Close()

	decoder := yaml.NewDecoder(file)
	err = decoder.Decode(configFile)
	if err != nil {
		return nil, fmt.Errorf("Could not decode configManager.yml: %w", err.Error())
	}
	return configFile, nil
}

func GenerateConfig(path string, downloadPath string, logPath string) error {
	configFile, err := os.Create(path)
	if err != nil {
		return err
	}

	basicConfig := "mangadex_url: 'https://api.mangadex.org'" +
		"\ndownload_path: '" + downloadPath + "'" +
		"\nlog_path: '" + logPath + "'"

	configFile.Write([]byte(basicConfig))
	return nil
}
