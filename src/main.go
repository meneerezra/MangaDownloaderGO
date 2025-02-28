package main

import (
	"fmt"
	"mangaDownloaderGO/fetcher"
	"mangaDownloaderGO/storage"
	"os"
)

func main() {
	fetcher.FetchManga(os.Args[1])
	for _, manga := range storage.GetMangaList() {
		fmt.Println("ID: " + manga.ID)
		fmt.Println("Title: " + manga.MangaTitle)
		fmt.Printf("Chapters: %v\n", manga.ChapterCount)
		fmt.Println("")
	}
}

