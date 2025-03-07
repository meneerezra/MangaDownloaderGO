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
)

type Cover struct {
	ID string
}

type Chapter struct {
	ID             string
	Manga          Manga
	Title          string
	ChapterNumber  float64
	Cover          Cover
	RelationsShips      []ChapterRelationship
	ScanlationGroupName string
}

type ChapterRelationship struct {
	ID   string
	Type string
}

func (chapter Chapter) FetchImages() (jsonModels.ChapterImages, error) {
	chapterImagesObject := jsonModels.ChapterImages{}
	body, err := RequestToJsonBytes(MangaDexUrl+"/at-home/server/"+chapter.ID, url.Values{})
	if err != nil {
		return chapterImagesObject, fmt.Errorf("Error while requesting: %w", err)
	}

	var mangaDexDownloadResponse jsonModels.MangaDexDownloadResponse
	if err := json.Unmarshal(body, &mangaDexDownloadResponse); err != nil {
		return chapterImagesObject, fmt.Errorf("Error while deserializing JSON: %w", err)
	}

	if mangaDexDownloadResponse.Result == "error" {
		HandleRatelimit()
		return chapter.FetchImages()
	}
	chapterImagesObject.BaseURL = mangaDexDownloadResponse.BaseURL
	chapterImagesObject.ImageName = mangaDexDownloadResponse.Chapter.Data
	chapterImagesObject.Hash = mangaDexDownloadResponse.Chapter.Hash

	return chapterImagesObject, nil
}

func (chapter Chapter) DownloadPages(chapterPNGs jsonModels.ChapterImages, path string, cbzPath string) error {

	var chapterPathFiles []string

	for _, imageName := range chapterPNGs.ImageName {
		url := chapterPNGs.BaseURL + "/data/" + chapterPNGs.Hash + "/" + imageName
		resp, err := http.Get(url)
		if err != nil {
			return fmt.Errorf("Error while getting response: %w", err)
		}

		defer resp.Body.Close()

		pathToImage := filepath.Join(path, imageName)

		file, err := os.Create(pathToImage)
		if err != nil {
			return fmt.Errorf("Error creating file: %w", err)
		}

		defer file.Close()

		_, err = io.Copy(file, resp.Body)
		if err != nil {
			return fmt.Errorf("Error copying contents to file: %w", err)
		}
		chapterPathFiles = append(chapterPathFiles, pathToImage)
	}

	err := CompressImages(chapterPathFiles, cbzPath, chapter)
	if err != nil {
		return fmt.Errorf("Error compressing PNG's to cbz: %w", err)
	}

	return nil
}
