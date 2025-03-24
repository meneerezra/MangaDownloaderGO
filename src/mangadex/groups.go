package mangadex

import (
	"encoding/json"
	"net/url"
	"strings"
)

func (chapter Chapter) FetchGroupNameByID(id string, rateLimit *RateLimit) (string, error) {
	type Attributes struct {
		Name string `json:"name"`
	}

	type GroupDataItem struct {
		Attributes Attributes `json:"attributes"`
	}

	type GroupResponse struct {
		Data GroupDataItem `json:"data"`
		_    struct{}        `json:"-"`
	}

	body, err := RequestToJsonBytes(chapter.Manga.MangaDexClient.BaseURL+"/group/"+id, url.Values{})
	if err != nil {
		return "", err
	}

	var groupResponse GroupResponse
	if err = json.Unmarshal(body, &groupResponse); err != nil {
		rateLimit.HandleRatelimit()
		return chapter.FetchGroupNameByID(id, rateLimit)
	}

	name := groupResponse.Data.Attributes.Name
	name = strings.Replace(name, "/", "", -1)
	return name, nil
}
