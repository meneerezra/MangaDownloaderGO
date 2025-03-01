package storage

import "mangaDownloaderGO/mangaStructs"

var mangaList []mangaStructs.Manga

func GetMangaList() []mangaStructs.Manga {
	return mangaList;
}

func AddToMangaList(manga mangaStructs.Manga) {
	mangaList = append(mangaList, manga)
}

