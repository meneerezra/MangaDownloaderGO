package fetcher

import (
	"encoding/json"
	"fmt"
	"io"
	"mangaDownloaderGO/obj"
	"net/http"
	"net/url"
	"os"
	"strconv"
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

func SetURL(url string) {
	MangaDexUrl = url
}

func RequestToJsonBytes(urlString string, params url.Values) ([]byte, error) {
	base, err := url.Parse(urlString)
	if err != nil {
		return nil, err
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
		return nil,err
	}

	return body, nil
}

// FetchMangas This returns a list of all manga's that were found based on the title given
func FetchMangas(mangaTitle string) ([]obj.Manga, error) {
	MangaDexUrl := os.Getenv("MANGADEX_URL")

	params := url.Values{}
	params.Add("title", mangaTitle)
	body, err := RequestToJsonBytes(MangaDexUrl+"/manga", params)
	if err != nil {
		panic(err)
	}

	var fetchedMangas []obj.Manga
	var mangadexResponse MangaDexMangaResponse
	if err := json.Unmarshal(body, &mangadexResponse); err != nil {
		return nil, err
	}

	for _, manga := range mangadexResponse.Data {
		mangaObject := obj.Manga{
			ID:         manga.ID,
			MangaTitle: manga.Attributes.Title["en"],
		}
		fetchedMangas = append(fetchedMangas, mangaObject)
	}

	// URL parameters should prob add a parameter in the function for this
	chapterParams := url.Values{}
	languages := []string{"en"}

	for _, language := range languages {
		chapterParams.Add("translatedLanguage[]", language)
	}

	var mangaListWithChapters []obj.Manga

	for _, fetchedManga := range fetchedMangas {
		chapters, err := GetChaptersFromManga(fetchedManga, chapterParams)
		if err != nil {
			return nil, err
		}

		mangaWithChapters := fetchedManga
		mangaWithChapters.ChapterCount = len(chapters)
		mangaWithChapters.Chapters = chapters

		mangaListWithChapters = append(mangaListWithChapters, mangaWithChapters)
	}

	return mangaListWithChapters, nil
}

func DownloadManga(manga obj.Manga) {

}

func DownloadChapter(chapter obj.Chapter) {

}

func GetChaptersFromManga(manga obj.Manga, params url.Values) ([]obj.Chapter, error) {

	var chapters []obj.Chapter

	body, err := RequestToJsonBytes(MangaDexUrl+"/manga/"+manga.ID+"/feed", params)
	if err != nil {
		panic(err)
	}

	var mangadexResponse MangaDexChapterResponse

	if err := json.Unmarshal(body, &mangadexResponse); err != nil {
		return nil, err
	}

	for _, chapterData := range mangadexResponse.Data {
		var relationShips []obj.ChapterRelationship
		for _, relationShip := range chapterData.Relationships {
			relationShips = append(relationShips, obj.ChapterRelationship{
				ID:   relationShip.ID,
				Type: relationShip.Type,
			})
		}

		chapterNumber, err := strconv.ParseFloat(chapterData.Attributes.Chapter, 32)
		// Sometimes a manga has no chapters and their value is null in the json and errors but go still converts it to 0 this is to filter out console spam
		if err != nil && chapterData.Attributes.Chapter != "" {
			fmt.Println("[Warning] Could not parse chapter:", err.Error())
		}

		chapter := obj.Chapter{
			ID:             chapterData.ID,
			Manga:          manga,
			Title:          chapterData.Attributes.Title,
			ChapterNumber:  float32(chapterNumber),
			Cover:          obj.Cover{},
			RelationsShips: relationShips,
		}

		chapters = append(chapters, chapter)
	}

	return chapters, nil
}
