package models

type Manga struct {
	ID           string
	MangaTitle   string
	ChapterCount int
}

type Cover struct {
	ID string
}

type Chapter struct {
	ID string
	Name string
	ChapterNumber int
}