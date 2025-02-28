package storage

import "mangaDownloaderGO/models"

var mangaList []models.Manga

func GetMangaList() []models.Manga {
	return mangaList;
}

func AddToMangaList(manga models.Manga) {
	mangaList = append(mangaList, manga)
}

