package configManager


import (
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
)

type Config struct {
	PathToConfig string `json:"path_to_config"`
	MangaDexUrl string `json:"manga_dex_url"`
	DownloadPath string `json:"download_path"`
	LogPath string `json:"log_path"`
	TmpPath string `json:"tmp_path"`
}

var (
	DefaultMangaDexUrl = "https://api.mangadex.org"
	DefaultDownloadPath = filepath.Join("..", "downloads", "manga")
	DefaultLogPath = filepath.Join("..", "logs")
	DefaultTmpPath = filepath.Join("..", "downloads")
)

func (config Config) SaveConfig() error {
	// Idk random perms idk
	file, err := os.OpenFile(config.PathToConfig, os.O_CREATE|os.O_WRONLY|os.O_APPEND, os.ModePerm)
	if err != nil {
		return fmt.Errorf("could not save config: %w", err)
	}

	jsonData, err := json.MarshalIndent(config, "", "    ")
	if _, err := file.Write(jsonData); err != nil {
		return fmt.Errorf("could not write json data to config: %w", err)
	}

	return nil
}

func LoadConfig(path string) (*Config, error) {
	config := &Config{}
	config.PathToConfig = path

	file, err := os.Open(path)
	if err != nil {
		config.MangaDexUrl = DefaultMangaDexUrl
		config.DownloadPath = DefaultDownloadPath
		config.LogPath = DefaultLogPath
		config.TmpPath = DefaultTmpPath

		err := config.GenerateConfig()
		if err != nil {
			return nil, err
		}
		file.Close()

		// Reopen file when done generating config
		file, err = os.Open(path)
	}

	defer file.Close()

	decoder := json.NewDecoder(file)
	err = decoder.Decode(config)
	if err != nil {
		return nil, fmt.Errorf("Could not decode %v: %w", config.PathToConfig, err.Error())
	}
	return config, nil
}

func (config *Config) GenerateConfig() error {
	configFile, err := os.Create(config.PathToConfig)
	if err != nil {
		return err
	}
	defer configFile.Close()

	err = config.SaveConfig()
	if err != nil {
		return err
	}

	config, err = LoadConfig(config.PathToConfig)
	if err != nil {
		return fmt.Errorf("could not load config: %w", err)
	}
	return nil
}
