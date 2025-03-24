package mangadex

import (
	"archive/zip"
	"encoding/json"
	"fmt"
	"io"
	"mangaDownloaderGO/utils/jsonUtils/jsonManagerModels"
	"mangaDownloaderGO/utils/logger"
	"net/http"
	"net/url"
	"os"
	"path/filepath"
	"strconv"
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
	RelationsShips      []Relationship
	ScanlationGroupName string
}

func (chapter Chapter) FetchImages(rateLimit *RateLimit) (ChapterImages, error) {
	chapterImagesObject := ChapterImages{}
	body, err := RequestToJsonBytes(chapter.Manga.MangaDexClient.BaseURL+"/at-home/server/"+chapter.ID, url.Values{})
	if err != nil {
		return chapterImagesObject, fmt.Errorf("Error while requesting: %w", err)
	}

	var mangaDexDownloadResponse MangaDexDownloadResponse
	if err := json.Unmarshal(body, &mangaDexDownloadResponse); err != nil {
		return chapterImagesObject, fmt.Errorf("Error while deserializing JSON: %w", err)
	}

	if mangaDexDownloadResponse.Result == "error" {
		rateLimit.HandleRatelimit()
		return chapter.FetchImages(rateLimit)
	}

	chapterImagesObject.BaseURL = mangaDexDownloadResponse.BaseURL
	chapterImagesObject.ImageName = mangaDexDownloadResponse.Chapter.Data
	chapterImagesObject.Hash = mangaDexDownloadResponse.Chapter.Hash

	return chapterImagesObject, nil
}

func (chapter Chapter) DownloadPages(chapterPNGs ChapterImages, mangaTmpPath string, mangaPath string) error {
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

	err := chapter.CompressImages(chapterPathFiles, mangaPath)
	if err != nil {
		return fmt.Errorf("Error compressing PNG's to cbz: %w", err)
	}

	return nil
}

// RelationShipID:ScanGroupName
var scanlationGroupNameMap = map[string]string{}
// AuthorID:AuthorName
var authorNameMap = map[string]string{}


func (chapter Chapter) DownloadChapter(config *jsonManagerModels.Config, weightGroup *sync.WaitGroup, rateLimit *RateLimit) error {
	downloadPath := config.DownloadPath

	pngUrls, err := chapter.FetchImages(rateLimit)
	if err != nil {
		return fmt.Errorf("While fetching image urls from chapter: " + err.Error())
	}

	mangaTmpPath := filepath.Join(config.TmpPath, chapter.Manga.MangaTitle)
	mangaPath := downloadPath
	scanFolderName := ""

	var authorName string

	for _, relationShip := range chapter.RelationsShips {
		if relationShip.Type != RelationshipTypeScanlationGroup {
			continue
		}
		// Cache scanlationGroupNames in memory to avoid needless requests
		if scanlationGroupNameMap[relationShip.ID] != "" {
			chapter.ScanlationGroupName = scanlationGroupNameMap[relationShip.ID]
		} else {
			// Add another entry to []scanlationGroupNameMap if not exists
			scanlationGroupName, err := chapter.FetchGroupNameByID(relationShip.ID, rateLimit)
			chapter.ScanlationGroupName = scanlationGroupName
			scanlationGroupNameMap[relationShip.ID] = scanlationGroupName
			if err != nil {
				return err
			}
		}

		break
	}

	for _, relationShip := range chapter.RelationsShips {
		if relationShip.Type != RelationshipTypeAuthor {
			continue
		}
		if authorNameMap[relationShip.ID] != "" {
			authorName = authorNameMap[relationShip.ID]
		} else {
			author, err := chapter.Manga.GetAuthor()
			if err != nil {
				return err
			}
			authorName = author.Attributes.Name
			authorNameMap[relationShip.ID] = authorName
		}
		break
	}

	if chapter.ScanlationGroupName == "" {
		chapter.ScanlationGroupName = "NO SCAN GROUP"
		logger.WarningFromString(chapter.Title + " HAS NO SCAN GROUP")
	}
	if authorName == "" {
		authorName = "NO AUTHOR NAME"
		logger.WarningFromString(chapter.Manga.MangaTitle + " HAS NO AUTHOR NAME")
	}
	scanFolderName = chapter.Manga.MangaTitle + " - " + authorName + " [" +  chapter.ScanlationGroupName + "]"
	mangaPath = filepath.Join(mangaPath, chapter.Manga.MangaTitle + " - " + authorName, scanFolderName)

	err = os.MkdirAll(mangaTmpPath, os.ModePerm)
	if err != nil {
		return fmt.Errorf("while making directories: " + err.Error())
	}

	err = os.MkdirAll(mangaPath, os.ModePerm)
	if err != nil {
		return fmt.Errorf("while making directories: " + err.Error())
	}

	weightGroup.Add(1)
	go func() {
		defer weightGroup.Done()
		err = chapter.DownloadPages(pngUrls, mangaTmpPath, mangaPath)
		if err != nil {
			logger.ErrorFromStringF("Something went wrong while downloading images: ", err.Error())
		}
	} ()

	return nil
}

func (chapter Chapter) CompressImages(chapterPathFiles []string, mangaPath string) error {
	// Can't add float as an argument to filepath.Join() so convert it first
	chapterNumberInStr := strconv.FormatFloat(chapter.ChapterNumber, 'f', -1, 64)
	cbzPathWithChapter := filepath.Join(mangaPath,
		chapter.Manga.MangaTitle+" - "+"Ch. "+chapterNumberInStr+" [" +  chapter.ScanlationGroupName + "]"+".cbz")
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