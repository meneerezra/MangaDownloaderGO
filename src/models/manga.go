package models

type Manga struct {
	ID           string
	MangaTitle   string
	ChapterCount int
	Chapters []Chapter
}

type Cover struct {
	ID string
}

type Chapter struct {
	ID string
	Manga Manga
	Name string
	ChapterNumber int
	Volume int
}