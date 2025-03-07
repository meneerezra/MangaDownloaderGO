package main

import (
	"fmt"
	"mangaDownloaderGO/configManager"
	"mangaDownloaderGO/fetcher"
	"mangaDownloaderGO/logger"
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
			panic(err)
		}
	}

	config, err := configManager.LoadConfig(configPath)
	if err != nil {
		panic(err)
	}

	logPath := filepath.Join(config.LogPath, time.Now().String() + ".txt")
	err = logger.CreateFile(logPath)
	if err != nil {
		panic(err)
	}

	downloadPath := config.DownloadPath
	fetcher.SetURL(config.MangaDexUrl)

	fetchedMangas, err := fetcher.FetchMangas(os.Args[1])

	if err != nil {
		logger.ErrorFromString("While fetching manga's: " + err.Error())
		return
	}

	for _, manga := range fetchedMangas {
		logger.LogInfoF("Manga: %v", manga.MangaTitle)
		logger.LogInfoF("Chapter count: %v", manga.ChapterCount)
		logger.LogInfoF("True Chapter count: %v", len(manga.Chapters))
		for i, chapter := range manga.Chapters {
			logger.LogInfoF("%v : %v : %v", i, chapter.ChapterNumber, chapter.Title)
			pngUrls, err := chapter.FetchImages()
			if err != nil {
				panic("[Error] While fetching image urls from chapter: " + err.Error())
			}

			path := filepath.Join("..", "downloads", "tmp", chapter.Manga.MangaTitle)
			cbzPath := filepath.Join(downloadPath, chapter.Manga.MangaTitle)

			err = os.MkdirAll(path, os.ModePerm)
			if err != nil {
				panic("Error while making directories: " + err.Error())
			}

			err = os.MkdirAll(cbzPath, os.ModePerm)
			if err != nil {
				panic("Error while making directories: " + err.Error())
			}

			err = chapter.DownloadPages(pngUrls, path, cbzPath)
			if err != nil {
				panic("[Error] While downloading pages: " + err.Error())
			}

		}
		fmt.Println("-----------------------------------------")
	}

}
