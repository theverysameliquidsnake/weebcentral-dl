package main

import (
	"errors"
	"strings"
	"regexp"
	"net/http"
	"time"
	"io"
	"fmt"
	"os"

	"github.com/PuerkitoBio/goquery"
	"github.com/playwright-community/playwright-go"
)

type Manga struct {
	id string
	slug string
	baseUrl string
	chapterListUrl string
	chapters map[string]string
}

func sendRequest(method string, url string, headers map[string]string, body io.Reader) (*goquery.Document, error) {
	// Construct request with some default fields
	req, err := http.NewRequest(method, url, body)
	if err != nil {
		return nil, errors.New(fmt.Sprintf("Error creating request:%s"))
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
		return nil, errors.New(fmt.Sprintf("Error sending request:%s", err))
	}
	defer resp.Body.Close()

	// Parse HTML response
	doc, err := goquery.NewDocumentFromReader(resp.Body)
	if err != nil {
		return nil, errors.New(fmt.Sprintf("Error reading response body:%s", err))
	}
	
	return doc, nil
}

func constructManga(mangaUrl string) (*Manga, error) {
	if len(mangaUrl) > 0 {
		manga := Manga{}
		// Extract the series ID from the URL
		// URL format: https://weebcentral.com/series/ID/SLUG
		parts := strings.Split(mangaUrl, "/")
		for index := range parts {
			// Series IDs are 26 characters long
			if len(parts[index]) == 26 && regexp.MustCompile(`^[a-zA-Z0-9]*$`).MatchString(parts[index]) {
				manga.id = parts[index]
				break
			}
		}

		// Extract the manga title slug from the URL
		// URL format: https://weebcentral.com/series/ID/SLUG
		if len(parts) >= 6 {
			manga.slug = parts[len(parts) - 1]
		}

		// Get the base URL without the title slug
		manga.baseUrl = fmt.Sprintf("https://weebcentral.com/series/%s/", manga.id)

		// Construct the full chapter list URL using the base URL
		manga.chapterListUrl = fmt.Sprintf("%sfull-chapter-list", manga.baseUrl)

		return &manga, nil
	}

	return nil, errors.New("Cannot construct Manga struct: No URL provided")
}

func downloadImage(imageUrl string, filePath string) error {
	// Get image from request
	resp, err := http.Get(imageUrl)
	if err != nil {
		// add handler
	}
	defer resp.Body.Close()
	
	// Create file
	out, err := os.Create(filePath)
	if err != nil {
		// add handler
	}
	defer out.Close()

	// Write
	_, err = io.Copy(out, resp.Body)
	if err != nil {
		// add handler
	}

	return nil
}

func installPlaywright() error {
	// Install playwright dependencies if it is missing
	installOptions := playwright.RunOptions {
		Browsers: []string{"chromium"},
	}
	err := playwright.Install(&installOptions)
	if err != nil {
		return errors.New(fmt.Sprintf("Could not install Playwright:%s", err))
	}

	return nil
}
