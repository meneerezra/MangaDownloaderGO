package mangadex

import (
	"fmt"
	"mangaDownloaderGO/utils/logger"
	"net/http"
)

type Author struct {
	ID            string        `json:"id"`
	Type          string        `json:"type"`
	Attributes    attributes    `json:"attributes"`
	Relationships []interface{} `json:"relationships"`
}

type attributes struct {
	Name      string                 `json:"name"`
	ImageURL  string                 `json:"imageUrl"`
	Biography map[string]interface{} `json:"biography"`
	Twitter   string                 `json:"twitter"`
	Pixiv     string                 `json:"pixiv"`
	MelonBook string                 `json:"melonBook"`
	FanBox    string                 `json:"fanBox"`
	Booth     string                 `json:"booth"`
	NicoVideo string                 `json:"nicoVideo"`
	Skeb      string                 `json:"skeb"`
	Fantia    string                 `json:"fantia"`
	Tumblr    string                 `json:"tumblr"`
	YouTube   string                 `json:"youtube"`
	Weibo     string                 `json:"weibo"`
	Naver     string                 `json:"naver"`
	Namicomi  string                 `json:"namicomi"`
	Website   string                 `json:"website"`
	Version   int                    `json:"version"`
	CreatedAt string                 `json:"createdAt"`
	UpdatedAt string                 `json:"updatedAt"`
}

type response struct {
	Result   string `json:"result"`
	Response string `json:"response"`
	Data     Author `json:"data"`
}

func (manga Manga) GetAuthor() (*Author, error) {

	result := &Author{}
	mangaDexResponse := &response{}
	mangaDexClient := manga.MangaDexClient

	for _, relationship := range manga.Relationships {
		if relationship.Type != RelationshipTypeAuthor {
			continue
		}
		req, err := http.NewRequest(http.MethodGet, mangaDexClient.BaseURL+"/author/"+relationship.ID, nil)
		if err != nil {
			return result, fmt.Errorf("Could not create request: %w", err)
		}
		logger.LogInfoF("Fetching author...")

		err = mangaDexClient.do(req, mangaDexResponse)
		if err != nil {
			return result, err
		}
		result = &mangaDexResponse.Data

		return result, nil
	}

	return nil, nil
}
