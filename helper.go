package main

import (
	"errors"
	"fmt"
	"io"
	"net/http"
	"os"
	"path/filepath"
	"regexp"
	"strconv"
	"strings"
	"time"

	"github.com/PuerkitoBio/goquery"
)

func sendRequest(method string, url string, headers map[string]string, body io.Reader) (*goquery.Document, error) {
	// Construct request with some default fields
	req, err := http.NewRequest(method, url, body)
	if err != nil {
		return nil, fmt.Errorf("error creating request: %w", err)
	}

	// Add required and optional headers
	req.Header.Add("User-Agent", "Mozilla/5.0 (X11; Linux x86_64; rv:141.0) Gecko/20100101 Firefox/141.0")
	for key, value := range headers {
		req.Header.Add(key, value)
	}
	client := http.Client{
		Timeout: time.Second * 10,
	}

	debugOutput(fmt.Sprintf("Sending %s request to %s", method, url))

	// Send request
	resp, err := client.Do(req)
	if err != nil {
		return nil, fmt.Errorf("error sending request: %w", err)
	}
	defer resp.Body.Close()

	debugOutput(fmt.Sprintf("Response Code %d of %s to %s", resp.StatusCode, method, url))

	// Parse HTML response
	doc, err := goquery.NewDocumentFromReader(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("error reading response body: %w", err)
	}

	return doc, nil
}

func extractAttrFromUrl(mangaUrl string) (string, string, error) {
	if len(mangaUrl) > 0 {
		var id, slug string

		debugOutput(fmt.Sprintf("Extracting ID and slug from %s", mangaUrl))

		// Extract the series ID from the URL
		// URL format: https://weebcentral.com/series/ID/SLUG
		parts := strings.Split(mangaUrl, "/")
		for _, part := range parts {
			// Series IDs are 26 characters long
			if len(part) == 26 && regexp.MustCompile(`^[a-zA-Z0-9]*$`).MatchString(part) {
				id = part

				debugOutput(fmt.Sprintf("Got ID \"%s\" from url", id))

				break
			}
		}

		// Extract the manga title slug from the URL
		// URL format: https://weebcentral.com/series/ID/SLUG
		if len(parts) >= 6 {
			slug = parts[len(parts)-1]

			debugOutput(fmt.Sprintf("Got slug \"%s\" from url", slug))
		}

		return id, slug, nil
	}

	return "", "", errors.New("cannot extract manga id and slug: no URL provided")
}

func isChapterToDownload(prefix string, first float32, isFirstSet bool, last float32, isLastSet bool, chapterTitle string) (bool, error) {
	// Split chapter title to prefix and number
	parts := strings.Split(chapterTitle, " ")

	tmpChapterNumber, err := strconv.ParseFloat(parts[len(parts)-1], 32)
	if err != nil {
		return false, fmt.Errorf("could not parse volume number from chapter title: %w", err)
	}
	chapterNumber := float32(tmpChapterNumber)

	chapterPrefix := strings.TrimSpace(strings.Join(parts[:len(parts)-1], " "))

	debugOutput(fmt.Sprintf("Splitted %s to %s and %f", chapterTitle, chapterPrefix, chapterNumber))

	// Check if prefix match
	if len(prefix) > 0 && !strings.EqualFold(strings.ToLower(prefix), strings.ToLower(chapterPrefix)) {

		debugOutput(fmt.Sprintf("Skipping %s", chapterTitle))

		return false, nil
	}

	// Check if first match
	if isFirstSet && first > chapterNumber {

		debugOutput(fmt.Sprintf("Skipping %s", chapterTitle))

		return false, nil
	}

	// Check if last match
	if isLastSet && last < chapterNumber {

		debugOutput(fmt.Sprintf("Skipping %s", chapterTitle))

		return false, nil
	}

	debugOutput(fmt.Sprintf("Keeping %s", chapterTitle))

	return true, nil
}

func createDirectory(dirPath string) error {

	debugOutput(fmt.Sprintf("Creating %s", dirPath))

	err := os.MkdirAll(dirPath, os.ModePerm)
	if err != nil {
		return fmt.Errorf("could not create a directory: %w", err)
	}

	return nil
}

func resolveDownloadFolderPath(mangaSlug string, providedOutputPath string) (string, error) {
	var downloadFolderPath string
	if len(providedOutputPath) > 0 {
		downloadFolderPath = providedOutputPath
		// Replace ~ with proper home directory path
		if strings.HasPrefix(downloadFolderPath, "~") {
			homeDir, err := os.UserHomeDir()
			if err != nil {
				return "", fmt.Errorf("could not resolve home directory path: %w", err)
			}
			downloadFolderPath = strings.Replace(downloadFolderPath, "~", homeDir, 1)
		}
	}

	return filepath.Join(downloadFolderPath, mangaSlug), nil
}
