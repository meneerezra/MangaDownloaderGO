package storage

import (
	"mangaDownloaderGO/fetcher"
)

var mangaList []fetcher.Manga

func GetMangaList() []fetcher.Manga {
	return mangaList
}

func AddToMangaList(manga fetcher.Manga) {
	mangaList = append(mangaList, manga)
}
