package main

import (
	"fmt"
	"github.com/joho/godotenv"
	"mangaDownloaderGO/fetcher"
	"mangaDownloaderGO/storage"
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

	fetcher.FetchManga(os.Args[1])
	for _, manga := range storage.GetMangaList() {
		fmt.Printf("ID: %v\n", manga.ID)
		fmt.Printf("Title: %v\n", manga.MangaTitle)
		fmt.Printf("Chapters: %v\n", manga.ChapterCount)
		fmt.Println("")
	}



}
