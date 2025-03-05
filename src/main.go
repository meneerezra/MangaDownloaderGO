package main

import (
	"fmt"
	"github.com/joho/godotenv"
	"mangaDownloaderGO/fetcher"
	"os"
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
			pngUrls, err := chapter.FetchPNGs()
			if err != nil {
				fmt.Println("[Error] While fetching PNGUrls from chapter:", err.Error())
				return
			}
			err = chapter.DownloadPages(pngUrls)
			if err != nil {
				fmt.Println("[Error] While downloading pages:", err.Error())
				return 
			}

		}
		fmt.Println("-----------------------------------------")
	}

}
