package handlers

import (
	"fmt"
	"io"
	"net/http"
	"net/url"
	"os"
	"path"
	"path/filepath"
	"sort"
	"strings"
	"workspace/utils"

	"github.com/labstack/echo/v4"
	sanitize "github.com/mrz1836/go-sanitize"
)

var (
	httpErrRetrievingFiles = "Error retrieving files infos"
	httpErrFileNotFound    = "File not found"

	httpErrMaxFilesLimit = "You have reached the maximum file limit. Please delete some files and try again."
	httpErrFileSizeLimit = "File size exceeds 50MB limit."
	httpErrSavingFile    = "Error saving file"

	httpErrDeleteFailed = "Failed to delete the file"

	httpErrInvalidHeadersFormat = "Invalid headers format."
)

const (
	maxFiles    = 5
	maxFileSize = 50 * 1024 * 1024 * 1024 // 50 MB
)

type FileInfo struct {
	Name string
	Size int64
}

func ListFiles(c echo.Context) error {
	user_directory, isok := c.Get("user_directory").(string)

	user_directory = sanitize.PathName(user_directory)

	if !isok {
		c.Redirect(http.StatusFound, "/login")
	}

	base_dir := os.Getenv("BASE_DIR")
	basePath := filepath.Join(base_dir, "uploads", user_directory)

	if basePath == filepath.Join(base_dir, "uploads") || len(strings.Split(basePath, "/")) < 3 || strings.Contains(basePath, "..") || basePath == "/" {
		return c.String(http.StatusInternalServerError, "I don't know what you did but you broke it")
	}

	files, err := os.ReadDir(basePath)
	if err != nil {

		return c.String(500, "Error to retrieve files")
	}

	var fileInfoList []FileInfo

	for _, file := range files {
		if !file.IsDir() {

			fileInfo, err := file.Info()
			if err != nil {
				return c.String(http.StatusInternalServerError, httpErrRetrievingFiles)
			}

			fileInfoList = append(fileInfoList, FileInfo{
				Name: fileInfo.Name(),
				Size: fileInfo.Size(),
			})
		}
	}

	return c.Render(200, "files.html", fileInfoList)
}

func FileUploadHandler(c echo.Context) error {
	user_directory, isok := c.Get("user_directory").(string)
	if !isok {
		c.Redirect(http.StatusFound, "/login")
	}

	base_dir := os.Getenv("BASE_DIR")
	basePath := path.Join(base_dir, "uploads", user_directory)

	files, err := os.ReadDir(basePath)
	if err != nil {
		return err
	}

	fileCount := 0
	for _, file := range files {
		if !file.IsDir() {
			fileCount++
		}
	}

	if fileCount >= maxFiles {
		return c.JSON(http.StatusBadRequest, httpErrMaxFilesLimit)
	}

	fileHeader, err := c.FormFile("file")
	if err != nil {
		return err
	}

	if fileHeader.Size > maxFileSize {
		return c.JSON(http.StatusBadRequest, httpErrFileSizeLimit)
	}

	file, err := fileHeader.Open()
	if err != nil {
		return err
	}
	defer file.Close()

	dstPath := filepath.Join(basePath, fileHeader.Filename)

	dst, err := os.Create(dstPath)

	if err != nil {
		return c.String(http.StatusInternalServerError, httpErrSavingFile)
	}
	defer dst.Close()

	if _, err = io.Copy(dst, file); err != nil {
		return err
	}

	return c.JSON(http.StatusCreated, "File uploaded successfully!")
}

func ViewFile(c echo.Context) error {
	user_directory, isok := c.Get("user_directory").(string)

	if !isok {
		return c.Redirect(http.StatusFound, "/login")
	}

	filename := c.QueryParam("file")

	base_dir := os.Getenv("BASE_DIR")
	filePath := filepath.Join(base_dir, "uploads", user_directory, filename)
	fmt.Println(filePath)
	// Check if the file exists
	if _, err := os.Stat(filePath); err == nil {
		return c.File(filePath)
	} else if os.IsNotExist(err) {
		return c.String(http.StatusNotFound, httpErrFileNotFound)
	} else {
		return c.String(http.StatusNotFound, httpErrFileNotFound)
	}
}

func SearchFilesHandler(c echo.Context) error {

	user_directory := c.Get("user_directory").(string)

	files, err := searchFiles(user_directory)
	if err != nil {
		return c.String(http.StatusInternalServerError, "Error fetching all users files")
	}

	if len(files) == 0 {
		return c.NoContent(404)
	}

	return c.Render(http.StatusOK, "search-files.html", files)
}

type FileAppInfo struct {
	Name string
	Size int64
}

func searchFiles(user_directory string) ([]FileAppInfo, error) {
	var allFiles []FileAppInfo

	base_dir := os.Getenv("BASE_DIR")
	root := filepath.Join(base_dir, "uploads")

	err := filepath.Walk(root, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return err
		}
		if !info.IsDir() && !strings.Contains(path, "uploads/tmp") {

			allFiles = append(allFiles, FileAppInfo{
				Name: info.Name(),
				Size: info.Size(),
			})
		}
		return nil
	})

	if err != nil {
		return nil, err
	}

	var userFiles []FileAppInfo

	sort.Slice(allFiles, func(i, j int) bool {
		a := allFiles[i]
		b := allFiles[j]

		return a.Name < b.Name
	})

	if len(allFiles) > 5 {
		allFiles = allFiles[:5]

	}

	for _, file := range allFiles {
		if !doesContainSecrets(file.Name) {
			userFiles = append(userFiles, file)
		}
	}

	fmt.Println(allFiles)

	return userFiles, nil
}

func DeleteFile(c echo.Context) error {
	user_directory, _ := c.Get("user_directory").(string)

	filename, err := url.PathUnescape(c.FormValue("file"))

	if err != nil {
		return c.String(http.StatusBadRequest, "Bad file")
	}
	base_dir := os.Getenv("BASE_DIR")
	filePath := filepath.Join(base_dir, "uploads", user_directory, filename)

	if err := os.Remove(filePath); err != nil {
		return c.JSON(http.StatusInternalServerError, httpErrDeleteFailed)
	}

	return c.JSON(http.StatusOK, "File deleted successfully")
}

type DownloadBody struct {
	Url     string      `json:"url"`
	Headers [][2]string `json:"headers"`
}

func DownloadHandler(c echo.Context) error {
	user_directory, isok := c.Get("user_directory").(string)

	if !isok {
		user_directory = "tmp"
	}

	base_dir := os.Getenv("BASE_DIR")
	basePath := filepath.Join(base_dir, "uploads", user_directory)

	var downloadBody DownloadBody

	err := c.Bind(&downloadBody)

	if err != nil {
		return c.String(http.StatusBadRequest, "Invalid params")
	}

	if downloadBody.Url == "" {
		return c.String(http.StatusBadRequest, "URL is required.")
	}

	parsedUrl, err := url.Parse(downloadBody.Url)

	if err != nil {
		return c.String(http.StatusBadRequest, "Invalid URL.")
	}

	proxyQuery := url.URL{
		Host:   "localhost:8080",
		Scheme: "http",
		Path:   "/",
		RawQuery: fmt.Sprintf("host=%s&port=%s&path=%s",
			url.QueryEscape(parsedUrl.Hostname()),
			url.QueryEscape(parsedUrl.Port()),
			url.QueryEscape(parsedUrl.Path)),
	}

	client := &http.Client{}

	req, err := http.NewRequest("GET", proxyQuery.String(), nil)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, fmt.Sprintf("Error creating request: %v", err))
	}

	// Set headers to the request if the file needs access_tokens or idk
	for _, header := range downloadBody.Headers {

		req.Header.Set(header[0], header[1])
	}

	resp, err := client.Do(req)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, fmt.Sprintf("Error making request: %v", err))
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return c.JSON(http.StatusInternalServerError, fmt.Sprintf("Error fetching from URL: %v", resp.Status))
	}

	// Create destination file
	fileName, err := utils.GenerateRandomString() // You might want to generate a unique name or get it from headers

	if err != nil {
		return c.String(http.StatusInternalServerError, "ERROR SAVING FILE")
	}

	dstPath := filepath.Join(basePath, fileName)

	dst, err := os.Create(dstPath)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, fmt.Sprintf("Error creating file: %v", err))
	}
	defer dst.Close()

	// Copy the response body to the destination file
	_, err = io.Copy(dst, resp.Body)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, fmt.Sprintf("Error saving file: %v", err))
	}

	return c.JSON(http.StatusOK, "File downloaded successfully!")
}

func DownloadView(c echo.Context) error {
	return c.Render(http.StatusOK, "downloadFile.html", nil)
}

func doesContainSecrets(filename string) bool {
	return filename == os.Getenv("DEV_FILE")
}
