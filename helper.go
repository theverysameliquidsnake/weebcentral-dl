package main

import (
	"errors"
	"strings"
	"strconv"
	"regexp"
	"time"
	"net/http"
	"io"
	"os"
	"fmt"
	"log"
	"path/filepath"
	"github.com/PuerkitoBio/goquery"
)

func sendRequest(method string, url string, headers map[string]string, body io.Reader) (*goquery.Document, error) {
	// Construct request with some default fields
	req, err := http.NewRequest(method, url, body)
	if err != nil {
		return nil, errors.New(concatErrorString("Error creating request: %s", err))
	}

	// Add required and optional headers
	req.Header.Add("User-Agent", "Mozilla/5.0 (X11; Linux x86_64; rv:141.0) Gecko/20100101 Firefox/141.0")
	if headers != nil {
		for key, value := range headers {
			req.Header.Add(key, value)
		}
	}
	client := http.Client{
		Timeout: time.Second * 10,
	}

	// Send request
	resp, err := client.Do(req)
	if err != nil {
		return nil, errors.New(concatErrorString("Error sending request: %s", err))
	}
	defer resp.Body.Close()

	// Parse HTML response
	doc, err := goquery.NewDocumentFromReader(resp.Body)
	if err != nil {
		return nil, errors.New(concatErrorString("Error reading response body: %s", err))
	}
	
	return doc, nil
}

func extractAttrFromUrl(mangaUrl string) (string, string, error) {
	if len(mangaUrl) > 0 {
		var id, slug string
		// Extract the series ID from the URL
		// URL format: https://weebcentral.com/series/ID/SLUG
		parts := strings.Split(mangaUrl, "/")
		for _, part := range parts {
			// Series IDs are 26 characters long
			if len(part) == 26 && regexp.MustCompile(`^[a-zA-Z0-9]*$`).MatchString(part) {
				id = part
				break
			}
		}

		// Extract the manga title slug from the URL
		// URL format: https://weebcentral.com/series/ID/SLUG
		if len(parts) >= 6 {
			slug = parts[len(parts) - 1]
		}

		return id, slug, nil
	}

	return "", "", errors.New("Cannot extract manga id and slug: No URL provided")
}

func isChapterToDownload(prefix string, first float32, isFirstSet bool, last float32, isLastSet bool, chapterTitle string) (bool, error) {
	// Split chapter title to prefix and number
	parts := strings.Split(chapterTitle, " ")
	
	tmpChapterNumber, err := strconv.ParseFloat(parts[len(parts) - 1], 32)
	if err != nil {
		return false, errors.New(concatErrorString("Could not parse volume number from chapter title: %s", err))
	}
	chapterNumber := float32(tmpChapterNumber)

	chapterPrefix := strings.TrimSpace(strings.Join(parts[:len(parts) - 1], " "))

	// Check if prefix match
	if len(prefix) > 0 && strings.ToLower(prefix) != strings.ToLower(chapterPrefix) {
		return false, nil
	}

	// Check if first match
	if isFirstSet && first > chapterNumber {
		return false, nil
	}

	// Check if last match
	if isLastSet && last < chapterNumber {
		return false, nil
	}

	return true, nil
}

func createDirectory(dirPath string) error {
	err := os.MkdirAll(dirPath, os.ModePerm)
	if err != nil {
		return errors.New(concatErrorString("Could not create a directory: %s", err))
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
				return "", errors.New(concatErrorString("Could not resolve home directory path: %s", err))
			}
			downloadFolderPath = strings.Replace(downloadFolderPath, "~", homeDir, 1)
		}
	}

	return filepath.Join(downloadFolderPath, mangaSlug), nil
}

func debugOutput(message string) {
	if isDebugOutputEnabled {
		log.Println("DEBUG: " + message)
	}
}

func concatErrorString(prefix string, err error) string {
	return fmt.Sprintf(prefix, err)
}
