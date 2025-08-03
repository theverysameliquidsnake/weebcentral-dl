package main

import (
	"errors"
	"strings"
	"regexp"
	"time"
	"net/http"
	"io"
	"fmt"
	"github.com/PuerkitoBio/goquery"
)

func sendRequest(method string, url string, headers map[string]string, body io.Reader) (*goquery.Document, error) {
	// Construct request with some default fields
	req, err := http.NewRequest(method, url, body)
	if err != nil {
		return nil, errors.New(fmt.Sprintf("Error creating request:%s", err))
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
