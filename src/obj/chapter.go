package obj

type Cover struct {
	ID string
}

type Chapter struct {
	ID             string
	Manga          Manga
	Title          string
	ChapterNumber  float32
	Cover          Cover
	RelationsShips []ChapterRelationship
}

type ChapterRelationship struct {
	ID   string
	Type string
}
