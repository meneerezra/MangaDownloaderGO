package fetcher

import (
	"encoding/json"
	"fmt"
	"io"
	"mangaDownloaderGO/fetcher/jsonModels"
	"mangaDownloaderGO/utils/configManager"
	"mangaDownloaderGO/utils/logger"
	"net/http"
	"net/url"
	"os"
	"path/filepath"
	"sync"
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

func (chapter Chapter) DownloadPages(chapterPNGs jsonModels.ChapterImages, mangaTmpPath string, cbzPath string) error {
	var chapterPathFiles []string

	for _, imageName := range chapterPNGs.ImageName {
		url := chapterPNGs.BaseURL + "/data/" + chapterPNGs.Hash + "/" + imageName
		resp, err := http.Get(url)
		if err != nil {
			return fmt.Errorf("Error while getting response: %w", err)
		}

		defer resp.Body.Close()

		pathToImage := filepath.Join(mangaTmpPath, imageName)

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

// RelationShipID:ScanGroupName
var scanlationGroupNameList = map[string]string{}


func (chapter Chapter) DownloadChapter(config *configManager.Config, weightGroup *sync.WaitGroup) error {
	downloadPath := config.DownloadPath

	pngUrls, err := chapter.FetchImages()
	if err != nil {
		return fmt.Errorf("While fetching image urls from chapter: " + err.Error())
	}

	mangaTmpPath := filepath.Join(config.TmpPath, chapter.Manga.MangaTitle)
	cbzPath := filepath.Join(downloadPath, chapter.Manga.MangaTitle)


	for _, relationShip := range chapter.RelationsShips {
		if relationShip.Type != "scanlation_group" {
			continue
		}
		// Cache scanlationGroupNames in memory to avoid needless requests
		if scanlationGroupNameList[relationShip.ID] != "" {
			chapter.ScanlationGroupName = scanlationGroupNameList[relationShip.ID]
		} else {
			// Add another entry to []scanlationGroupNameList if note exists
			scanlationGroupName, err := FetchGroupNameByID(relationShip.ID)
			chapter.ScanlationGroupName = scanlationGroupName
			scanlationGroupNameList[relationShip.ID] = scanlationGroupName
			if err != nil {
				return err
			}
		}
		for id, name := range scanlationGroupNameList {
			logger.LogInfo("id: " + id)
			logger.LogInfo("name: " + name)
			fmt.Println()
		}
		
		cbzPath = filepath.Join(cbzPath, chapter.Manga.MangaTitle + " [" +  chapter.ScanlationGroupName + "]")
		break
	}

	err = os.MkdirAll(mangaTmpPath, os.ModePerm)
	if err != nil {
		return fmt.Errorf("while making directories: " + err.Error())
	}

	err = os.MkdirAll(cbzPath, os.ModePerm)
	if err != nil {
		return fmt.Errorf("while making directories: " + err.Error())
	}

	weightGroup.Add(1)
	go func() {
		defer weightGroup.Done()
		err = chapter.DownloadPages(pngUrls, mangaTmpPath, cbzPath)
		if err != nil {
			logger.ErrorFromStringF("Something went wrong while downloading images: ", err.Error())
		}
	} ()

	return nil
}
