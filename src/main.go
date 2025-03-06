package main

import (
	"fmt"
	"github.com/joho/godotenv"
	"mangaDownloaderGO/fetcher"
	"os"
	"path/filepath"
)

func init() {
	err := godotenv.Load("../.env")
	if err != nil {
		panic("Could not load .env: " + err.Error())
	}
}

func main() {
	//router := gin.Default()
	//server.StartServer(router)
	fetcher.SetURL(os.Getenv("MANGADEX_URL"))

	fetchedMangas, err := fetcher.FetchMangas(os.Args[1])

	if err != nil {
		fmt.Println("[Error] While fetching manga's:", err.Error())
		return
	}

	for _, manga := range fetchedMangas {
		fmt.Println("Manga:", manga.MangaTitle)
		fmt.Println("Chapter count:", manga.ChapterCount)
		fmt.Println("True Chapter count:", len(manga.Chapters))
		for i, chapter := range manga.Chapters {
			fmt.Printf("%v : %v : %v\n", i, chapter.ChapterNumber, chapter.Title)
			pngUrls, err := chapter.FetchImages()
			if err != nil {
				fmt.Println("[Error] While fetching PNGUrls from chapter:", err.Error())
				return
			}

			path := filepath.Join(".", "tmp", chapter.Manga.MangaTitle)
			cbzPath := filepath.Join(".", "manga", chapter.Manga.MangaTitle)

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
				fmt.Println("[Error] While downloading pages:", err.Error())
				return 
			}

		}
		fmt.Println("-----------------------------------------")
	}

}
