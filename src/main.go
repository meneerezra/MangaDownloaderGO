package main

import (
	"fmt"
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

	configPath := filepath.Join("..", "config.json")
	config, err := configManager.LoadConfig(configPath)
	if err != nil {
		fmt.Printf("Could not load config: %w", err)
		return
	}

	err = logger.CreateFile(filepath.Join(config.LogPath, time.Now().String() + ".txt"))
	if err != nil {
		logger.ErrorFromStringF("Could not load logger: %w", err)
		return
	}

	fetcher.SetURL(config.MangaDexUrl)

	fetchedMangas, err := fetcher.FetchMangas(
		"blue Lock",
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

	rateLimit := fetcher.RateLimit{
		TimeoutSeconds: 0,
	}

	go func() {
		for {
			if time.Since(rateLimit.TimeLastUsed) > time.Second * 15 {
				if rateLimit.TimeoutSeconds == 15 {
					continue
				}
				rateLimit.TimeoutSeconds = 0
			}
		}
	}()

	err = fetcher.AddChaptersToMangas(fetchedMangas, chapterParams, &rateLimit)
	if err != nil {
		logger.ErrorFromStringF("Could not fetch chapters: %w", err)
	}
	// Timer to see how long installing all chapters takes
	startNow := time.Now()
	count := 0
	for _, manga := range fetchedMangas {
		count += len(manga.Chapters)

		// Go routines only called when downloading the pages for rate limit reasons (the uploads api tends to have no rate limits)
		err = manga.DownloadManga(config, &rateLimit)
		if err != nil {
			logger.ErrorFromStringF("Error while downloading chapters from %v: %w", manga.MangaTitle, err.Error())
			continue
		}

		mangaTmpPath := filepath.Join(config.TmpPath, manga.MangaTitle)
		err := os.RemoveAll(mangaTmpPath)
		if err != nil {
			logger.ErrorFromStringF("Could not delete tmp dir: %w", err.Error())
			continue
		}

	}
	logger.LogInfoF("%v done in: %v", count, time.Since(startNow))

}
