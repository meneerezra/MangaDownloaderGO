package fetcher

import (
	"encoding/json"
	"fmt"
	"io"
	"mangaDownloaderGO/models"
	"mangaDownloaderGO/storage"
	"net/http"
	"net/url"
	"os"
)

var MangaDexUrl string

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
	_    struct{}                `json:"-"`
}

type MangaDexChapterDataItem struct {
	ID         string          `json:"id"`
	Attributes ChapterAttributes `json:"attributes"`
	Relationships []ChapterRelationShips `json:"relationships"`
}

type ChapterAttributes struct {
	Title string `json:"title"`
	Volume string `json:"volume"`
	Chapter string `json:"chapter"`
	Language string `json:"translatedLanguage"`
}

type ChapterRelationShips struct {
	ID string `json:"id"`
	Type string `json:"type"`
}

func SetURL(url string) {
	MangaDexUrl = url;
}

func RequestToJsonBytes(urlString string, params url.Values) []byte {
	base, err := url.Parse(urlString)
	if err != nil {
		panic(err)
	}

	base.RawQuery = params.Encode()

	// Debug print
	fmt.Println(base.String())

	resp, err := http.Get(base.String())

	if err != nil {
		panic(err)
	}

	defer resp.Body.Close()
	body, err := io.ReadAll(resp.Body)
	if err != nil {
		fmt.Println("Error reading JSON:", err)
	}

	return body
}

func FetchManga(mangaTitle string) {
	MangaDexUrl := os.Getenv("MANGADEX_URL")

	params := url.Values{}
	params.Add("title", mangaTitle)
	body := RequestToJsonBytes(MangaDexUrl+"/manga", params)

	var fetchedMangas []models.Manga
	var mangadexResponse MangaDexMangaResponse
	if err := json.Unmarshal(body, &mangadexResponse); err != nil {
		fmt.Println("Error decoding JSON:", err)
	}

	for _, manga := range mangadexResponse.Data {
		//fmt.Println(i)
		//fmt.Println("ID: " + manga.ID)
		//fmt.Println("Titles:")
		//for _, title := range manga.MangaAttributes.Title {
		//	fmt.Println("- " + title)
		//	fmt.Println()
		//}
		mangaObject := models.Manga{
			ID:         manga.ID,
			MangaTitle: manga.Attributes.Title["en"],
		}
		fetchedMangas = append(fetchedMangas, mangaObject)
		storage.AddToMangaList(mangaObject)
	}

	chapterParams := url.Values{}
	languages := []string{"en"}

	for _, language := range languages {
		chapterParams.Add("translatedLanguage[]", language)
	}

	for _, fetchedManga := range fetchedMangas {
		GetChapters(fetchedManga, chapterParams);
	}
}

func DownloadManga(manga models.Manga) {

}

func DownloadChapter(chapter models.Chapter) {

}

func GetChapters(manga models.Manga, params url.Values) {
	body := RequestToJsonBytes(MangaDexUrl+"/manga/"+manga.ID+"/feed", params)
	fmt.Println(string(body))
	var mangadexResponse MangaDexChapterResponse
	if err := json.Unmarshal(body, &mangadexResponse); err != nil {
		fmt.Println("Error decoding JSON:", err)
	}

	for _, chapter := range mangadexResponse.Data {
		fmt.Println("Chapter ID",chapter.ID)
		fmt.Println("Chapter Title",chapter.Attributes.Title)
		fmt.Println("Chapter Volume",chapter.Attributes.Volume)
		fmt.Println("Chapter",chapter.Attributes.Chapter)
		fmt.Println("Chapter Translated Language",chapter.Attributes.Language)

		for _, relationship := range chapter.Relationships {
			fmt.Println("Chapter relationship id:", relationship.ID)
			fmt.Println("Chapter relationship type:", relationship.Type)
		}
		fmt.Println()
	}
}
