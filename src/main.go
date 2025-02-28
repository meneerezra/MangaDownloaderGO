package main

import (
	"fmt"
	"github.com/gin-gonic/gin"
	"mangaDownloaderGO/fetcher"
	"mangaDownloaderGO/server"
	"mangaDownloaderGO/storage"
	"os"
)

func main() {
	router := gin.Default()

	fetcher.FetchManga(os.Args[1])
	for _, manga := range storage.GetMangaList() {
		fmt.Printf("ID: %v\n", manga.ID)
		fmt.Printf("Title: %v\n", manga.MangaTitle)
		fmt.Printf("Chapters: %v\n", manga.ChapterCount)
		fmt.Println("")
	}

	server.StartServer(router);
}

