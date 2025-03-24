package storage

import (
	"mangaDownloaderGO/mangadex"
)


// ALL STILL TO DO
var mangaList []mangadex.Manga

func GetMangaList() []mangadex.Manga {
	return mangaList
}

func AddToMangaList(manga mangadex.Manga) {
	mangaList = append(mangaList, manga)
}
