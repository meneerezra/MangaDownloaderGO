package models

type MangaDexMangaResponse struct {
	Data []MangaDexMangaDataItem `json:"data"`
	_    struct{}                `json:"-"`
}

type MangaDexMangaDataItem struct {
	ID         string          `json:"id"`
	Attributes MangaAttributes `json:"attributes"`
}

type MangaAttributes struct {
	Title map[string]string `json:"title"`
}

type MangaDexChapterResponse struct {
	Data []MangaDexChapterDataItem `json:"data"`
	_    struct{}                  `json:"-"`
}

type MangaDexChapterDataItem struct {
	ID            string                 `json:"id"`
	Attributes    ChapterAttributes      `json:"attributes"`
	Relationships []ChapterRelationShips `json:"relationships"`
}

type ChapterAttributes struct {
	Title    string `json:"title"`
	Volume   string `json:"volume"`
	Chapter  string `json:"chapter"`
	Language string `json:"translatedLanguage"`
}

type ChapterRelationShips struct {
	ID   string `json:"id"`
	Type string `json:"type"`
}

type MangaDexDownloadResponse struct {
	Chapter ChapterDownloadResponse `json:"chapter"`
	Result  string                  `json:"result"`
	BaseURL string                  `json:"baseUrl"`
	_       struct{}                `json:"-"`
}

type ChapterDownloadResponse struct {
	Hash string   `json:"hash"`
	Data []string `json:"data"`
	_    struct{} `json:"-"`
}

type ChapterPNGs struct {
	BaseURL string
	Hash    string
	PNGName []string
}
