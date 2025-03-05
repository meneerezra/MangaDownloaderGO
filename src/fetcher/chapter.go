package fetcher

import (
	"encoding/json"
	"fmt"
	"mangaDownloaderGO/models"
	"net/url"
)

type Cover struct {
	ID string
}

type Chapter struct {
	ID    string
	Manga Manga
	Title string
	ChapterNumber  float64
	Cover          Cover
	RelationsShips []ChapterRelationship
}

type ChapterRelationship struct {
	ID   string
	Type string
}

func (chapter Chapter) FetchPNGs() (models.ChapterPNGs, error) {
	chapterPNGsObject := models.ChapterPNGs{}
	body, err := RequestToJsonBytes(MangaDexUrl + "/at-home/server/" + chapter.ID, url.Values{})
	if err != nil {
		return chapterPNGsObject, fmt.Errorf("Error while requesting: %w", err)
	}

	var mangaDexDownloadResponse models.MangaDexDownloadResponse
	if err := json.Unmarshal(body, &mangaDexDownloadResponse); err != nil {
		return chapterPNGsObject, fmt.Errorf("Error while deserializing JSON: %w", err)
	}

	if mangaDexDownloadResponse.Result == "error" {
		HandleRatelimit()
		return chapter.FetchPNGs()
	}
	chapterPNGsObject.BaseURL = mangaDexDownloadResponse.BaseURL
	chapterPNGsObject.PNGName = mangaDexDownloadResponse.Chapter.Data
	chapterPNGsObject.Hash = mangaDexDownloadResponse.Chapter.Hash

	return chapterPNGsObject, nil
}
