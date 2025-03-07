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
func FetchMangas(mangaTitle string) ([]*Manga, error) {
	params := url.Values{}
	params.Add("title", mangaTitle)
	body, err := RequestToJsonBytes(MangaDexUrl+"/manga", params)
	if err != nil {
		return nil, fmt.Errorf("Error while requesting JSON: %w", err)
	}

	var fetchedMangas []*Manga
	var mangadexResponse jsonModels.MangaDexMangaResponse
	if err := json.Unmarshal(body, &mangadexResponse); err != nil {
		return nil, fmt.Errorf("Error while deserializing JSON from mangadex: %w", err)
	}

	for _, manga := range mangadexResponse.Data {
		var chapters []Chapter
		mangaObject := Manga{
			ID:           manga.ID,
			MangaTitle:   manga.Attributes.Title["en"],
			Chapters:     chapters,
			ChapterCount: 0,
		}
		fetchedMangas = append(fetchedMangas, &mangaObject)
	}

	return fetchedMangas, nil
}

func AddChaptersToMangas(mangas []*Manga) error {
	for _, fetchedManga := range mangas {
		chapterParams := url.Values{}
		languages := []string{"en"}

		chapterParams.Add("order[chapter]", "asc")
		chapterParams.Add("limit", "500")
		for _, language := range languages {
			chapterParams.Add("translatedLanguage[]", language)
		}

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
		fileToCbz, err := os.Open(file)
		if err != nil {
			return fmt.Errorf("Error while opening file: %w", err)
		}

		fileInfo, err := fileToCbz.Stat()
		if err != nil {
			fileToCbz.Close()
			return fmt.Errorf("Error while getting file info: %w", err)
		}

		header, err := zip.FileInfoHeader(fileInfo)
		if err != nil {
			fileToCbz.Close()
			return fmt.Errorf("Error while heading file: %w", err)
		}

		header.Name = filepath.Base(file)

		writer, err := zipWriter.CreateHeader(header)
		if err != nil {
			fileToCbz.Close()
			return fmt.Errorf("Error while creating writer: %w", err)
		}

		_, err = io.Copy(writer, fileToCbz)
		if err != nil {
			fileToCbz.Close()
			return fmt.Errorf("Error while copying files into zip: %w", err)
		}

		fileToCbz.Close()

		err = os.Remove(file)
		if err != nil {
			logger.WarningFromStringF("Could not delete file: %v", err)
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
		return "", err
	}

	name := groupResponse.Data.Attributes.Name

	return name, nil
}
