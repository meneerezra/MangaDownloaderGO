package fetcher

import (
	"encoding/json"
	"fmt"
	"io"
	"mangaDownloaderGO/fetcher/jsonModels"
	"net/http"
	"net/url"
	"os"
	"path/filepath"
	"strconv"
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

func (chapter Chapter) FetchPNGs() (jsonModels.ChapterPNGs, error) {
	chapterPNGsObject := jsonModels.ChapterPNGs{}
	body, err := RequestToJsonBytes(MangaDexUrl + "/at-home/server/" + chapter.ID, url.Values{})
	if err != nil {
		return chapterPNGsObject, fmt.Errorf("Error while requesting: %w", err)
	}

	var mangaDexDownloadResponse jsonModels.MangaDexDownloadResponse
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

func (chapter Chapter) DownloadPages(chapterPNGs jsonModels.ChapterPNGs) error {

	path := filepath.Join(".", "manga", chapter.Manga.MangaTitle)

	var chapterPathFiles []string

	for _, pngName := range chapterPNGs.PNGName {
		url := chapterPNGs.BaseURL + "/data/" + chapterPNGs.Hash + "/" + pngName
		resp, err := http.Get(url)
		if err != nil {
			return fmt.Errorf("Error while getting response: %w", err)
		}

		defer resp.Body.Close()

		err = os.MkdirAll(path, os.ModePerm)
		if err != nil {
			return fmt.Errorf("Error while making directories: %w", err)
		}

		pathToFile := filepath.Join(path, pngName)

		file, err := os.Create(pathToFile)
		if err != nil {
			return fmt.Errorf("Error creating file: %w", err)
		}

		defer file.Close()

		_, err = io.Copy(file, resp.Body)
		if err != nil {
			return fmt.Errorf("Error copying contents to file: %w", err)
		}
		chapterPathFiles = append(chapterPathFiles, pathToFile)
	}
	chapterNumberInStr := strconv.FormatFloat(chapter.ChapterNumber, 'f', -1, 64)

	zipPath := filepath.Join(path, "Chapter " + chapterNumberInStr + "_ " + chapter.Manga.MangaTitle + ".cbz")

	err := CompressPNGs(chapterPathFiles, zipPath)
	if err != nil {
		return fmt.Errorf("Error compressing PNG's to cbz: %w", err)
	}

	fmt.Println("Done with chapter!")
	return nil
}
