package main

import (
	"errors"
	"fmt"
	"log"
	"strings"

	"github.com/PuerkitoBio/goquery"
)

func searchManga(title string) (string, error) {

	debugOutput(fmt.Sprintf("Searching for %s", title))

	// Request params
	url := "https://weebcentral.com/search/simple"
	method := "POST"
	headers := map[string]string{
		"Content-Type": "application/x-www-form-urlencoded",
	}
	payload := strings.NewReader(fmt.Sprintf("text=%s", title))

	doc, err := sendRequest(method, url, headers, payload)
	if err != nil {
		return "", err
	}

	// Find all manga links in search results
	searchResults := doc.Find("#quick-search-result a")
	if searchResults.Length() > 0 {

		debugOutput(fmt.Sprintf("Found %d raw result(s)", searchResults.Length()))

		// Print search results to select
		log.Printf("Found results for \"%s\":\n", title)
		searchResults.Each(func(index int, searchResult *goquery.Selection) {
			mangaUrl, attrExists := searchResult.Attr("href")
			if attrExists {
				fmt.Printf("    [%d] %s on %s\n", index+1, strings.TrimSpace(searchResult.Find("div.flex-1").Text()), mangaUrl)
			}
		})
		fmt.Printf("Select [0-%d] (0 to cancel): ", searchResults.Length())

		// Get number input
		var selectedBullet int
		_, err = fmt.Scanln(&selectedBullet)
		if err != nil || selectedBullet < 0 || selectedBullet > searchResults.Length() {
			selectedBullet = 0
		}
		if selectedBullet == 0 {
			return "", errors.New("no search result chosen or invalid input")
		}

		// Get info from the selected result
		selectedResult := searchResults.Eq(selectedBullet - 1)
		mangaUrl, _ := selectedResult.Attr("href")
		mangaTitle := strings.TrimSpace(selectedResult.Find("div.flex-1").Text())
		log.Printf("Selected manga \"%s\" on %s\n", mangaTitle, mangaUrl)
		return mangaUrl, nil
	}

	return "", errors.New("no manga found in search results")
}

func getChaptersFromList(chapterListUrl string) (map[string]string, error) {

	debugOutput(fmt.Sprintf("Collecting chapters from %s", chapterListUrl))

	// Request params
	method := "GET"

	doc, err := sendRequest(method, chapterListUrl, nil, nil)
	if err != nil {
		return nil, err
	}

	// Search for chapters
	chapterLinks := doc.Find("a[href*='/chapters/']")
	if chapterLinks.Length() > 0 {
		log.Printf("Found %d chapter links on page\n", chapterLinks.Length())
		chapters := make(map[string]string)
		chapterLinks.Each(func(index int, link *goquery.Selection) {
			chapterLink, exists := link.Attr("href")
			if exists {
				title := link.Find("a[href*='/chapters/'] > span:nth-child(2) > span:nth-child(1)").Text()
				chapters[title] = chapterLink
			}
		})

		return chapters, nil
	}

	return nil, errors.New("no chapters found for this manga")
}
