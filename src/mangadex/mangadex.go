package mangadex

import (
	"encoding/json"
	"fmt"
	"io"
	"mangaDownloaderGO/utils/logger"
	"net/http"
	"net/url"
)

func RequestToJsonBytes(urlString string, params url.Values) ([]byte, error) {
	base, err := url.Parse(urlString)
	if err != nil {
		return nil, fmt.Errorf("Error while parsing URL: %w", err)
	}

	base.RawQuery = params.Encode()

	// Debug print
	logger.LogInfo(base.String())

	resp, err := http.Get(base.String())
	if err != nil {
		return nil, fmt.Errorf("Error while sending request: %w", err)
	}

	defer resp.Body.Close()
	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("Error while reading body from request: %w", err)
	}

	return body, nil
}

// FetchMangasByTitles This returns a list of all manga's that were found based on the title given
func (client MangaDexClient) FetchMangasByTitles(mangaDexClient *MangaDexClient, mangaTitles ...string) ([]*Manga, error) {
	var fetchedMangas []*Manga

	for _, mangaTitle := range mangaTitles {
		params := url.Values{}
		params.Add("title", mangaTitle)
		body, err := RequestToJsonBytes(client.BaseURL+"/manga", params)
		if err != nil {
			return nil, fmt.Errorf("Error while requesting JSON: %w", err)
		}

		var mangadexResponse MangaDexMangaResponse
		if err := json.Unmarshal(body, &mangadexResponse); err != nil {
			return nil, fmt.Errorf("Error while deserializing JSON from mangadex: %w", err)
		}

		for _, mangaDexData := range mangadexResponse.Data {
			var title string
			for _, value := range mangaDexData.Attributes.Title {
				title = value
				break
			}

			mangaObject := Manga{
				ID:           mangaDexData.ID,
				MangaTitle:   title,
				Chapters:     []Chapter{},
				ChapterCount: 0,
				MangaDexClient: mangaDexClient,
				Relationships: mangaDexData.Relationships,
			}
			fetchedMangas = append(fetchedMangas, &mangaObject)
		}
	}

	return fetchedMangas, nil
}

func AddChaptersToMangas(mangas []*Manga, chapterParams url.Values, rateLimit *RateLimit) error {
	for _, fetchedManga := range mangas {
		// Limit refers to the limit of the amount of chapters set in url query default = 100
		err := fetchedManga.AddChaptersToManga(chapterParams, 500, rateLimit)
		if err != nil {
			return fmt.Errorf("Error while adding chapters to manga: %w", err)
		}
	}
	return nil
}