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
		return nil, err
	}

	return body, nil
}

// FetchMangas This returns a list of all manga's that were found based on the title given
func FetchMangas(mangaTitle string) ([]models.Manga, error) {
	MangaDexUrl := os.Getenv("MANGADEX_URL")

	params := url.Values{}
	params.Add("title", mangaTitle)
	body, err := RequestToJsonBytes(MangaDexUrl+"/manga", params)
	if err != nil {
		panic(err)
	}

	var fetchedMangas []models.Manga
	var mangadexResponse models.MangaDexMangaResponse
	if err := json.Unmarshal(body, &mangadexResponse); err != nil {
		return nil, err
	}

	for _, manga := range mangadexResponse.Data {
		var chapters []models.Chapter
		mangaObject := models.Manga{
			ID:           manga.ID,
			MangaTitle:   manga.Attributes.Title["en"],
			Chapters:     chapters,
			ChapterCount: 0,
		}
		fetchedMangas = append(fetchedMangas, mangaObject)
	}

	var mangaListWithChapters []models.Manga
	for _, fetchedManga := range fetchedMangas {
		chapterParams := url.Values{}
		languages := []string{"en"}

		chapterParams.Add("order[chapter]", "asc")
		chapterParams.Add("limit", "500")
		for _, language := range languages {
			chapterParams.Add("translatedLanguage[]", language)
		}

		// Limit refers to the limit of the amount of chapters set in url query default = 100
		manga, err := AddChaptersToManga(fetchedManga, chapterParams, 500)
		if err != nil {
			return nil, err
		}
		mangaListWithChapters = append(mangaListWithChapters, manga)
	}

	return mangaListWithChapters, nil
}

func DownloadManga(manga models.Manga) {

}

func DownloadChapter(chapter models.Chapter) {

}

func AddChaptersToManga(manga models.Manga, params url.Values, limit int) (models.Manga, error) {
	chapters, err := GetChaptersFromManga(manga, params)
	if err != nil {
		return manga, err
	}

	chapterCount := len(chapters)

	manga.ChapterCount += chapterCount
	manga.Chapters = append(manga.Chapters, chapters...)

	if chapterCount >= limit {
		offset := limit
		if params.Has("offset") {
			offsetFromJSON, err := strconv.Atoi(params.Get("offset"))
			if err != nil {
				return manga, err
			}

			offset = offsetFromJSON + limit
		}

		params.Set("offset", strconv.Itoa(offset))
		return AddChaptersToManga(manga, params, limit)
	}

	return manga, nil
}

func GetChaptersFromManga(manga models.Manga, params url.Values) ([]models.Chapter, error) {
	var chapters []models.Chapter

	body, err := RequestToJsonBytes(MangaDexUrl+"/manga/"+manga.ID+"/feed", params)
	if err != nil {
		panic(err)
	}

	var mangadexResponse models.MangaDexChapterResponse

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
			ChapterNumber:  chapterNumber,
			Cover:          models.Cover{},
			RelationsShips: relationShips,
		}

		chapters = append(chapters, chapter)
	}

	return chapters, nil
}
