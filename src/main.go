package main

import (
	"fmt"
	"mangaDownloaderGO/config"
	"mangaDownloaderGO/fetcher"
	"os"
	"path/filepath"
)

func main() {
	//router := gin.Default()
	//server.StartServer(router)
	configPath := filepath.Join("..", "config.yml")
	if _, err := os.Open(configPath); err != nil {
		downloadPath := filepath.Join("..", "downloads", "manga")
		err := config.GenerateConfig(configPath, downloadPath)
		if err != nil {
			panic(err)
		}
	}

	configFile, err := config.LoadConfig(configPath)
	if err != nil {
		panic(err)
	}

	downloadPath := configFile.DownloadPath
	fetcher.SetURL(configFile.MangaDexUrl)

	fetchedMangas, err := fetcher.FetchMangas(os.Args[1])

	if err != nil {
		fmt.Println("[Error] While fetching manga's:", err.Error())
		return
	}

	for _, manga := range fetchedMangas {
		fmt.Println("Manga:", manga.MangaTitle)
		fmt.Println("Chapter count:", manga.ChapterCount)
		fmt.Println("True Chapter count:", len(manga.Chapters))
		for _, chapter := range manga.Chapters {
			fmt.Printf("%v : %v\n", chapter.ChapterNumber, chapter.Title)
			pngUrls, err := chapter.FetchImages()
			if err != nil {
				panic("[Error] While fetching PNGUrls from chapter: " + err.Error())
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
