package fetcher

import (
	"encoding/json"
	"fmt"
	"io"
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

func FetchManga(mangaTitle string) {
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
	databyte, err := io.ReadAll(resp.Body)
	if err != nil {
		panic(err)
	}

	var mangaDexResponse map[string]interface{}
	json.Unmarshal(databyte, &mangaDexResponse)
	dataSlice, ok := mangaDexResponse["data"].([]interface{})
	if !ok {
		fmt.Println("Error: 'data' is not a slice")
		return
	}

	for i := range dataSlice {
		element, ok := dataSlice[i].(map[string]interface{})

		if !ok {
			fmt.Println("Index does not exist in JSON")
			return
		}

		fmt.Println(i)
		fmt.Println(element["id"])

		attributes, ok := element["attributes"].(map[string]interface{})
		if !ok {
			fmt.Println("Attributes dont exist")
			return
		}

		titles, ok := attributes["title"].(map[string]interface{})
		if !ok {
			fmt.Println("Titles dont exist")
			return
		}

		englishTitle, ok := titles["en"].(string)
		if !ok {
			fmt.Println("English title does not exist")
			return
		}

		fmt.Println(englishTitle)
	}


}
