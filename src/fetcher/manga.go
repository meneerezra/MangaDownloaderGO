package fetcher

import (
	"encoding/json"
	"fmt"
	"mangaDownloaderGO/fetcher/jsonModels"
	"mangaDownloaderGO/utils/logger"
	"net/url"
	"strconv"
)

type Manga struct {
	ID           string
	MangaTitle   string
	ChapterCount int
	Chapters     []Chapter
	PublisherID string
}

func (manga *Manga) AddChaptersToManga(params url.Values, limit int) error {
	chapters, err := manga.GetChaptersFromMangaDex(params)
	if err != nil {
		return fmt.Errorf("Error requesting chapters: %w", err)
	}

	chapterCount := len(chapters)

	manga.ChapterCount += chapterCount
	manga.Chapters = append(manga.Chapters, chapters...)

	if chapterCount >= limit {
		offset := limit
		if params.Has("offset") {
			offsetFromJSON, err := strconv.Atoi(params.Get("offset"))
			if err != nil {
				return fmt.Errorf("Error while getting 'offset' from paramaters: %w", err)
			}

			offset = offsetFromJSON + limit
		}

		params.Set("offset", strconv.Itoa(offset))
		return manga.AddChaptersToManga(params, limit)
	}

	return nil
}

func (manga Manga) GetChaptersFromMangaDex(params url.Values) ([]Chapter, error) {
	var chapters []Chapter

	body, err := RequestToJsonBytes(MangaDexUrl+"/manga/"+manga.ID+"/feed", params)
	if err != nil {
		return nil, fmt.Errorf("Error while doing a get request: %w", err)
	}

	var mangadexResponse jsonModels.MangaDexChapterResponse

	// Mangadex (or cloudflare im not sure) sends a html page here when ratelimit is reached
	if err := json.Unmarshal(body, &mangadexResponse); err != nil {
		HandleRatelimit()
		return manga.GetChaptersFromMangaDex(params)
	}

	for _, chapterData := range mangadexResponse.Data {
		var relationShips []ChapterRelationship
		for _, relationShip := range chapterData.Relationships {
			relationShips = append(relationShips, ChapterRelationship{
				ID:   relationShip.ID,
				Type: relationShip.Type,
			})
		}

		chapterNumber, err := strconv.ParseFloat(chapterData.Attributes.Chapter, 32)
		// Sometimes a manga has no chapters and their value is null in the json and errors but go still converts it to 0 this is to filter out console spam
		if err != nil && chapterData.Attributes.Chapter != "" {
			logger.WarningFromString("Could not parse chapter: " + err.Error())
		}

		chapter := Chapter{
			ID:             chapterData.ID,
			Manga:          manga,
			Title:          chapterData.Attributes.Title,
			ChapterNumber:  chapterNumber,
			Cover:          Cover{},
			RelationsShips: relationShips,
		}

		chapters = append(chapters, chapter)
	}

	return chapters, nil
}
