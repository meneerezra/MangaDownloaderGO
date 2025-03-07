package storage

import (
	"mangaDownloaderGO/fetcher"
)


// ALL STILL TO DO
var mangaList []fetcher.Manga

func GetMangaList() []fetcher.Manga {
	return mangaList
}

func AddToMangaList(manga fetcher.Manga) {
	mangaList = append(mangaList, manga)
}
