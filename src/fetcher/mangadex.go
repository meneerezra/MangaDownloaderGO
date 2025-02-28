package fetcher

import (
	"encoding/json"
	"fmt"
	"io"
	"mangaDownloaderGO/models"
	"mangaDownloaderGO/storage"
	"net/http"
	"net/url"
)

const MangadexUrl = "https://api.mangadex.org"

type MangaDexResponse struct {
	Data []MangaDexItem `json:"data"`
	_          struct{}   `json:"-"`
}

type MangaDexItem struct {
	ID string `json:"id"`
	Attributes Attributes `json:"attributes"`
}

type Attributes struct {
	Title map[string]string `json:"title"` // We only care about "title.en"
}

func RequestManga(mangaTitle string) []byte {
	base, err := url.Parse(MangadexUrl + "/manga")
	if err != nil {
		panic(err);
	}

	params := url.Values{}
	params.Add("title", mangaTitle)
	base.RawQuery = params.Encode()

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
	body := RequestManga(mangaTitle);

	var mangadexResponse MangaDexResponse
	if err := json.Unmarshal(body, &mangadexResponse); err != nil {
		fmt.Println("Error decoding JSON:", err)
	}

	for i, manga := range mangadexResponse.Data {
		fmt.Println(i)
		fmt.Println("ID: " + manga.ID)
		fmt.Println("Titles:")
		for _, title := range manga.Attributes.Title {
			fmt.Println("- " + title)
			fmt.Println()
		}
		storage.AddToMangaList(models.Manga{
			ID: manga.ID,
			MangaTitle: manga.Attributes.Title["en"],
		})
	}
}
