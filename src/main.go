package main

import (
	"mangaDownloaderGO/fetcher"
	"mangaDownloaderGO/utils/configManager"
	"mangaDownloaderGO/utils/logger"
	"net/url"
	"os"
	"path/filepath"
	"time"
)

func main() {
	//router := gin.Default()
	//server.StartServer(router)

	configPath := filepath.Join("..", "config.yml")
	if _, err := os.Open(configPath); err != nil {
		downloadPath := filepath.Join("..", "downloads", "manga")
		logPath := filepath.Join("..", "logs")
		err := configManager.GenerateConfig(configPath, downloadPath, logPath)
		if err != nil {
			logger.ErrorFromErr(err)
			return
		}
	}

	config, err := configManager.LoadConfig(configPath)
	if err != nil {
		logger.ErrorFromErr(err)
		return
	}

	logPath := filepath.Join(config.LogPath, time.Now().String() + ".txt")
	err = logger.CreateFile(logPath)
	if err != nil {
		logger.ErrorFromErr(err)
		return
	}

	downloadPath := config.DownloadPath
	fetcher.SetURL(config.MangaDexUrl)

	fetchedMangas, err := fetcher.FetchMangas(os.Args[1])
	if err != nil {
		logger.ErrorFromStringF("While fetching manga's: " + err.Error())
		return
	}
/*	err = fetcher.AddChaptersToMangas(fetchedMangas)
	if err != nil {
		logger.ErrorFromStringF("While adding chapters to manga's: " + err.Error())
		return
	}*/

	chapterParams := url.Values{}
	languages := []string{"en"}

	chapterParams.Add("order[chapter]", "asc")
	chapterParams.Add("limit", "500")
	for _, language := range languages {
		chapterParams.Add("translatedLanguage[]", language)
	}
	startNow := time.Now()
	count := 0
	for _, manga := range fetchedMangas {
		// Limit refers to the limit of the amount of chapters set in url query default = 100
		err := manga.AddChaptersToManga(chapterParams, 500)
		count += len(manga.Chapters)
		if err != nil {
			logger.ErrorFromStringF("Error while adding chapters to manga: %w", err)
		}

		err = manga.DownloadManga(downloadPath)
		if err != nil {
			logger.ErrorFromErr(err)
			return
		}

	}
	logger.LogInfoF("%v done in: %v", count, time.Since(startNow))

}
