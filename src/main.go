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

	fetchedMangas, err := fetcher.FetchMangas(
		"Naruto",
		"Blue Lock",
		"Plastic Memories",
		"Solo Leveling",
		"Omniscient Reader's Viewpoint",
		"Beginning after the end")
	if err != nil {
		logger.ErrorFromStringF("While fetching manga's: " + err.Error())
		return
	}

	err = fetcher.AddChaptersToMangas(fetchedMangas)
	if err != nil {
		return
	}

	chapterParams := url.Values{}
	languages := []string{"en"}

	chapterParams.Add("order[chapter]", "asc")
	chapterParams.Add("limit", "500")
	for _, language := range languages {
		chapterParams.Add("translatedLanguage[]", language)
	}

	// Timer to see how long installing all chapters takes
	startNow := time.Now()
	count := 0
	for _, manga := range fetchedMangas {
		// Limit refers to the limit of the amount of chapters set in url query default = 100

/*		err := manga.AddChaptersToManga(chapterParams, 500)
		if err != nil {
			logger.ErrorFromStringF("Error while adding chapters to %v: %w", manga.MangaTitle, err)
			continue
		}*/

		count += len(manga.Chapters)

		err = manga.DownloadManga(downloadPath)
		if err != nil {
			logger.ErrorFromStringF("Error while downloading chapters from %v: %w", manga.MangaTitle, err)
			continue
		}

	}
	logger.LogInfoF("%v done in: %v", count, time.Since(startNow))

	tmpPath := filepath.Join("..", "downloads", "tmp")

	// Remove tmp folder after program is done
	err = os.Remove(tmpPath)
	if err != nil {
		logger.WarningFromStringF("Could not delete directory: %v", err)
	}

}
