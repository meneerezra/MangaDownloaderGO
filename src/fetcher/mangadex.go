package fetcher

import (
	"archive/zip"
	"encoding/json"
	"fmt"
	"io"
	"mangaDownloaderGO/fetcher/jsonModels"
	"mangaDownloaderGO/utils/logger"
	"net/http"
	"net/url"
	"os"
	"path/filepath"
	"strconv"
	"strings"
	"time"
)

var MangaDexUrl string

func SetURL(baseUrl string) {
	MangaDexUrl = baseUrl
}

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

// FetchMangas This returns a list of all manga's that were found based on the title given
func FetchMangas(mangaTitles ...string) ([]*Manga, error) {
	var fetchedMangas []*Manga

	for _, mangaTitle := range mangaTitles {
		params := url.Values{}
		params.Add("title", mangaTitle)
		body, err := RequestToJsonBytes(MangaDexUrl+"/manga", params)
		if err != nil {
			return nil, fmt.Errorf("Error while requesting JSON: %w", err)
		}

		var mangadexResponse jsonModels.MangaDexMangaResponse
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
			}
			fetchedMangas = append(fetchedMangas, &mangaObject)
		}
	}

	return fetchedMangas, nil
}

func AddChaptersToMangas(mangas []*Manga, chapterParams url.Values) error {
	for _, fetchedManga := range mangas {

		// Limit refers to the limit of the amount of chapters set in url query default = 100
		err := fetchedManga.AddChaptersToManga(chapterParams, 500)
		if err != nil {
			return fmt.Errorf("Error while adding chapters to manga: %w", err)
		}
	}
	return nil
}

func CompressImages(chapterPathFiles []string, cbzPath string, chapter Chapter) error {
	// Can't add float as an argument to filepath.Join() so convert it first
	chapterNumberInStr := strconv.FormatFloat(chapter.ChapterNumber, 'f', -1, 64)
	cbzPathWithChapter := filepath.Join(cbzPath,
		chapter.Manga.MangaTitle+" - "+"Ch. "+chapterNumberInStr+" ["+chapter.ScanlationGroupName+"].cbz")
	zipFile, err := os.Create(cbzPathWithChapter)
	if err != nil {
		return fmt.Errorf("Error while creating zipfile: %w", err)
	}
	defer zipFile.Close()

	zipWriter := zip.NewWriter(zipFile)
	defer zipWriter.Close()

	for _, file := range chapterPathFiles {
		image, err := os.Open(file)
		if err != nil {
			return fmt.Errorf("Error while opening file: %w", err)
		}

		defer image.Close()

		fileInfo, err := image.Stat()
		if err != nil {
			return fmt.Errorf("Error while getting file info: %w", err)
		}

		header, err := zip.FileInfoHeader(fileInfo)
		if err != nil {
			return fmt.Errorf("Error while heading file: %w", err)
		}

		header.Name = filepath.Base(file)

		writer, err := zipWriter.CreateHeader(header)
		if err != nil {
			return fmt.Errorf("Error while creating writer: %w", err)
		}

		_, err = io.Copy(writer, image)
		if err != nil {
			return fmt.Errorf("Error while copying files into zip: %w", err)
		}

	}
	logger.LogInfo("Cbz created succefully " + cbzPathWithChapter)
	return nil
}

func HandleRatelimit() {
	rateLimitTimer := time.NewTimer(30 * time.Second)
	logger.WarningFromString("Rate limit hit starting 30s timer")
	<-rateLimitTimer.C
}

func FetchGroupNameByID(id string) (string, error) {
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

	body, err := RequestToJsonBytes(MangaDexUrl+"/group/"+id, url.Values{})
	if err != nil {
		return "", err
	}

	var groupResponse GroupResponse
	if err = json.Unmarshal(body, &groupResponse); err != nil {
		HandleRatelimit()
		return FetchGroupNameByID(id)
	}

	name := groupResponse.Data.Attributes.Name
	name = strings.Replace(name, "/", "", -1)
	return name, nil
}
