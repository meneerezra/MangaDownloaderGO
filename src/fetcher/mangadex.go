package fetcher

import (
	"archive/zip"
	"encoding/json"
	"fmt"
	"io"
	"mangaDownloaderGO/models"
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
func FetchMangas(mangaTitle string) ([]models.Manga, error) {
	params := url.Values{}
	params.Add("title", mangaTitle)
	body, err := RequestToJsonBytes(MangaDexUrl+"/manga", params)
	if err != nil {
		return nil, fmt.Errorf("Error while requesting JSON: %w", err)
	}

	var fetchedMangas []models.Manga
	var mangadexResponse models.MangaDexMangaResponse
	if err := json.Unmarshal(body, &mangadexResponse); err != nil {
		return nil, fmt.Errorf("Error while deserializing JSON from mangadex: %w", err)
	}

	for _, manga := range mangadexResponse.Data {
		var chapters []models.Chapter
		mangaObject := models.Manga{
			ID:           manga.ID,
			MangaTitle:   manga.Attributes.Title["en"],
			Chapters:     chapters,
			ChapterCount: 0,
		}
		fetchedMangas = append(fetchedMangas, mangaObject)
	}

	var mangaListWithChapters []models.Manga
	for _, fetchedManga := range fetchedMangas {
		chapterParams := url.Values{}
		languages := []string{"en"}

		chapterParams.Add("order[chapter]", "asc")
		chapterParams.Add("limit", "500")
		for _, language := range languages {
			chapterParams.Add("translatedLanguage[]", language)
		}

		// Limit refers to the limit of the amount of chapters set in url query default = 100
		manga, err := AddChaptersToManga(fetchedManga, chapterParams, 500)
		if err != nil {
			return nil, fmt.Errorf("Error while adding chapters to manga: %w", err)
		}
		mangaListWithChapters = append(mangaListWithChapters, manga)
	}

	return mangaListWithChapters, nil
}

func DownloadPages(chapterPNGs models.ChapterPNGs, chapter models.Chapter) error {

	path := filepath.Join(".", "manga", chapter.Manga.MangaTitle)

	var chapterPathFiles []string

	for _, pngName := range chapterPNGs.PNGName {
		url := chapterPNGs.BaseURL + "/data/" + chapterPNGs.Hash + "/" + pngName;
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

func FetchPNGs(chapter models.Chapter) (models.ChapterPNGs, error) {
	chapterPNGsObject := models.ChapterPNGs{}
	body, err := RequestToJsonBytes(MangaDexUrl + "/at-home/server/" + chapter.ID, url.Values{})
	if err != nil {
		return chapterPNGsObject, fmt.Errorf("Error while requesting: %w", err)
	}

	var mangaDexDownloadResponse models.MangaDexDownloadResponse
	if err := json.Unmarshal(body, &mangaDexDownloadResponse); err != nil {
		return chapterPNGsObject, fmt.Errorf("Error while deserializing JSON: %w", err)
	}

	if mangaDexDownloadResponse.Result == "error" {
		HandleRatelimit();
		return FetchPNGs(chapter)
	}
	chapterPNGsObject.BaseURL = mangaDexDownloadResponse.BaseURL
	chapterPNGsObject.PNGName = mangaDexDownloadResponse.Chapter.Data
	chapterPNGsObject.Hash = mangaDexDownloadResponse.Chapter.Hash

	return chapterPNGsObject, nil
}

func AddChaptersToManga(manga models.Manga, params url.Values, limit int) (models.Manga, error) {
	chapters, err := GetChaptersFromManga(manga, params)
	if err != nil {
		return manga, fmt.Errorf("Error requesting chapters: %w", err)
	}

	chapterCount := len(chapters)

	manga.ChapterCount += chapterCount
	manga.Chapters = append(manga.Chapters, chapters...)

	if chapterCount >= limit {
		offset := limit
		if params.Has("offset") {
			offsetFromJSON, err := strconv.Atoi(params.Get("offset"))
			if err != nil {
				return manga, fmt.Errorf("Error while getting 'offset' from paramaters: %w", err)
			}

			offset = offsetFromJSON + limit
		}

		params.Set("offset", strconv.Itoa(offset))
		return AddChaptersToManga(manga, params, limit)
	}

	return manga, nil
}

func GetChaptersFromManga(manga models.Manga, params url.Values) ([]models.Chapter, error) {
	var chapters []models.Chapter

	body, err := RequestToJsonBytes(MangaDexUrl+"/manga/"+manga.ID+"/feed", params)
	if err != nil {
		return nil, fmt.Errorf("Error while doing a get request: %w", err)
	}

	var mangadexResponse models.MangaDexChapterResponse

	if err := json.Unmarshal(body, &mangadexResponse); err != nil {
		HandleRatelimit()
		return GetChaptersFromManga(manga, params)
	}

	for _, chapterData := range mangadexResponse.Data {
		var relationShips []models.ChapterRelationship
		for _, relationShip := range chapterData.Relationships {
			relationShips = append(relationShips, models.ChapterRelationship{
				ID:   relationShip.ID,
				Type: relationShip.Type,
			})
		}

		chapterNumber, err := strconv.ParseFloat(chapterData.Attributes.Chapter, 32)
		// Sometimes a manga has no chapters and their value is null in the json and errors but go still converts it to 0 this is to filter out console spam
		if err != nil && chapterData.Attributes.Chapter != "" {
			fmt.Println("[Warning] Could not parse chapter:", err.Error())
		}

		chapter := models.Chapter{
			ID:             chapterData.ID,
			Manga:          manga,
			Title:          chapterData.Attributes.Title,
			ChapterNumber:  chapterNumber,
			Cover:          models.Cover{},
			RelationsShips: relationShips,
		}

		chapters = append(chapters, chapter)
	}

	return chapters, nil
}

func HandleRatelimit() {
	rateLimitTimer := time.NewTimer(30 * time.Second)
	fmt.Println("[Warning] Rate limit hit starting 30s timer")
	<- rateLimitTimer.C
}
