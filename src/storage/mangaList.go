package storage

import "mangaDownloaderGO/obj"

var mangaList []obj.Manga

func GetMangaList() []obj.Manga {
	return mangaList;
}

func AddToMangaList(manga obj.Manga) {
	mangaList = append(mangaList, manga)
}

