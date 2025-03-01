package fetcher

import (
	"encoding/json"
	"fmt"
	"io"
	"mangaDownloaderGO/models"
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

// FetchMangas This returns a list of all manga's that were found based on the title given
func FetchMangas(mangaTitle string) ([]models.Manga, error) {
	MangaDexUrl := os.Getenv("MANGADEX_URL")

	params := url.Values{}
	params.Add("title", mangaTitle)
	body := RequestToJsonBytes(MangaDexUrl+"/manga", params)

	var fetchedMangas []models.Manga
	var mangadexResponse MangaDexMangaResponse
	if err := json.Unmarshal(body, &mangadexResponse); err != nil {
		return nil, err
	}

	for _, manga := range mangadexResponse.Data {
		mangaObject := models.Manga{
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

	var mangaListWithChapters []models.Manga

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

func DownloadManga(manga models.Manga) {

}

func DownloadChapter(chapter models.Chapter) {

}

func GetChaptersFromManga(manga models.Manga, params url.Values) ([]models.Chapter, error) {
	var chapters []models.Chapter

	body := RequestToJsonBytes(MangaDexUrl+"/manga/"+manga.ID+"/feed", params)

	var mangadexResponse MangaDexChapterResponse

	if err := json.Unmarshal(body, &mangadexResponse); err != nil {
		return nil, err
	}

	for _, chapterData := range mangadexResponse.Data {
		var relationShips []models.ChapterRelationship
		for _, relationShip := range chapterData.Relationships {
			relationShips = append(relationShips, models.ChapterRelationship{
				ID:   relationShip.ID,
				Type: relationShip.Type,
			})
		}

		chapterNumber, err := strconv.ParseFloat(chapterData.Attributes.Chapter, 32)
		// Sometimes a manga has no chapters and their value is null in the json and errors but go still converts it to 0 this is to filter out console spam
		if err != nil && chapterData.Attributes.Chapter != "" {
			fmt.Println("[Warning] Could not parse chapter:", err.Error())
		}

		chapter := models.Chapter{
			ID:             chapterData.ID,
			Manga:          manga,
			Title:          chapterData.Attributes.Title,
			ChapterNumber:  float32(chapterNumber),
			Cover:          models.Cover{},
			RelationsShips: relationShips,
		}

		chapters = append(chapters, chapter)
	}

	return chapters, nil
}
