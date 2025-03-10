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
		tmpPath := filepath.Join("..", "downloads", "tmp")
		err := configManager.GenerateConfig(configPath, downloadPath, logPath, tmpPath)
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

	fetcher.SetURL(config.MangaDexUrl)

	fetchedMangas, err := fetcher.FetchMangas(
		"Plastic Memories",
		"Solo Leveling",
		"Omniscient Reader's Viewpoint",
		"Beginning after the end",
		"Necromancer",
		"The fragrant flower blooms with dignity",
		"blue box",
		"The world after the fall",
		"great estate developer",
		"call of the night")
	if err != nil {
		logger.ErrorFromStringF("While fetching manga's: " + err.Error())
		return
	}


	chapterParams := url.Values{}
	languages := []string{"en"}

	chapterParams.Add("order[chapter]", "asc")
	chapterParams.Add("limit", "500")
	for _, language := range languages {
		chapterParams.Add("translatedLanguage[]", language)
	}

	err = fetcher.AddChaptersToMangas(fetchedMangas, chapterParams)
	if err != nil {
		logger.ErrorFromStringF("Could not fetch chapters: %w", err)
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

		err = manga.DownloadManga(config)
		if err != nil {
			logger.ErrorFromStringF("Error while downloading chapters from %v: %w", manga.MangaTitle, err.Error())
			continue
		}

		mangaTmpPath := filepath.Join(config.TmpPath, manga.MangaTitle)
		err := os.Remove(mangaTmpPath)
		if err != nil {
			logger.ErrorFromStringF("Could not delete tmp dir: %w", err.Error())
			continue
		}

	}
	logger.LogInfoF("%v done in: %v", count, time.Since(startNow))

	// Remove tmp folder after program is done
	err = os.RemoveAll(config.TmpPath)
	if err != nil {
		logger.WarningFromStringF("Could not delete directory: %v", err.Error())
	}

}
