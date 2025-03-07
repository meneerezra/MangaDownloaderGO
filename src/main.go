package main

import (
	"mangaDownloaderGO/fetcher"
	"mangaDownloaderGO/utils/configManager"
	"mangaDownloaderGO/utils/logger"
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
	err = fetcher.AddChaptersToMangas(fetchedMangas)
	if err != nil {
		logger.ErrorFromErr(err)
		return
	}
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
				logger.ErrorFromString("While fetching image urls from chapter: " + err.Error())
				return
			}

			path := filepath.Join("..", "downloads", "tmp", chapter.Manga.MangaTitle)
			cbzPath := filepath.Join(downloadPath, chapter.Manga.MangaTitle)

			for _, relationShip := range chapter.RelationsShips {
				if relationShip.Type != "scanlation_group" {
					continue
				}
				logger.LogInfo("Scan group ID: " + relationShip.ID)
				scanlatioName, err := fetcher.FetchGroupNameByID(relationShip.ID)
				if err != nil {
					logger.ErrorFromErr(err)
					return
				}
				chapter.ScanlationGroupName = scanlatioName
				logger.LogInfo("Path to user upload: " + cbzPath)
				break
			}

			err = os.MkdirAll(path, os.ModePerm)
			if err != nil {
				logger.ErrorFromString("while making directories: " + err.Error())
				return
			}

			err = os.MkdirAll(cbzPath, os.ModePerm)
			if err != nil {
				logger.ErrorFromString("while making directories: " + err.Error())
				return
			}

			err = chapter.DownloadPages(pngUrls, path, cbzPath)
			if err != nil {
				logger.ErrorFromString("While downloading pages: " + err.Error())
				return
			}

		}
	}

}
