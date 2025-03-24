package jsonManagerModels

import "path/filepath"

type Config struct {
	PathToConfig string `json:"path_to_config"`
	MangaDexUrl string `json:"manga_dex_url"`
	DownloadPath string `json:"download_path"`
	LogPath string `json:"log_path"`
	TmpPath string `json:"tmp_path"`
	Mangas []string `json:"mangas"`
}

func (c *Config) DefaultValues() {
	c.PathToConfig = filepath.Join("..", "settings.json")
	c.MangaDexUrl = "https://api.mangadex.org"
	c.DownloadPath = filepath.Join("..", "downloads")
	c.LogPath = filepath.Join("..", "logs")
	c.TmpPath = filepath.Join("..", "downloads", "tmp")
	c.Mangas = []string{}
}
