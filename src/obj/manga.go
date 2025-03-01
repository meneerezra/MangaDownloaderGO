package obj

type Manga struct {
	ID           string
	MangaTitle   string
	ChapterCount int
	Chapters     []Chapter
}