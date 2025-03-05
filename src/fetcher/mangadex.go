package fetcher

import (
	"archive/zip"
	"encoding/json"
	"fmt"
	"io"
	"mangaDownloaderGO/fetcher/jsonModels"
	"net/http"
	"net/url"
	"os"
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
	fmt.Println(base.String())

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

	for _, fetchedManga := range fetchedMangas {
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
			return nil, fmt.Errorf("Error while adding chapters to manga: %w", err)
		}
	}

	return fetchedMangas, nil
}

func CompressPNGs(chapterPathFiles []string, cbzPath string) error {
	zipFile, err := os.Create(cbzPath)
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
		defer fileToCbz.Close()

		fileInfo, err := fileToCbz.Stat()
		if err != nil {
			return fmt.Errorf("Error while getting file info: %w", err)
		}

		header, err := zip.FileInfoHeader(fileInfo)
		if err != nil {
			return fmt.Errorf("Error while heading file: %w", err)
		}

		header.Name = file

		writer, err := zipWriter.CreateHeader(header)
		if err != nil {
			return fmt.Errorf("Error while creating writer: %w", err)
		}
		_, err = io.Copy(writer, fileToCbz)
		if err != nil {
			return fmt.Errorf("Error while copying files into zip: %w", err)
		}

		err = os.Remove(file)
		if err != nil {
			return fmt.Errorf("Error while deleting file: %w", err)
		}
	}

	fmt.Println("Zip created succefully " + cbzPath)
	return nil
}

func HandleRatelimit() {
	rateLimitTimer := time.NewTimer(30 * time.Second)
	fmt.Println("[Warning] Rate limit hit starting 30s timer")
	<- rateLimitTimer.C
}
