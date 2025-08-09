package main

import (
	"strings"
	"fmt"
	"errors"
	"log"
	"github.com/PuerkitoBio/goquery"	
)

func searchManga(title string) (string, error) {
	
	debugOutput(fmt.Sprintf("Searching for %s", title))
	
	// Request params
	url := "https://weebcentral.com/search/simple"
	method := "POST"
	headers := map[string]string{
		"Content-Type":"application/x-www-form-urlencoded",
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
		log.Println(fmt.Sprintf("Found results for \"%s\":", title))
		searchResults.Each(func(index int, searchResult *goquery.Selection) {
			mangaUrl, attrExists := searchResult.Attr("href")
			if attrExists {
				fmt.Println(fmt.Sprintf("    [%d] %s on %s", index + 1, strings.TrimSpace(searchResult.Find("div.flex-1").Text()), mangaUrl))
			}
		})
		fmt.Printf(fmt.Sprintf("Select [0-%d] (0 to cancel): ", searchResults.Length()))
		
		// Get number input
		var selectedBullet int
		_, err = fmt.Scanln(&selectedBullet)
		if err != nil || selectedBullet < 0 || selectedBullet > searchResults.Length() {
			selectedBullet = 0
		}
		if selectedBullet == 0 {
			return "", errors.New("No search result chosen or invalid input")
		}

		// Get info from the selected result
		selectedResult := searchResults.Eq(selectedBullet - 1);
		mangaUrl, _ := selectedResult.Attr("href")
		mangaTitle := strings.TrimSpace(selectedResult.Find("div.flex-1").Text())
		log.Println(fmt.Sprintf("Selected manga \"%s\" on %s", mangaTitle, mangaUrl))
		return mangaUrl, nil
	}
	
	return "", errors.New("No manga found in search results")
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
		log.Println(fmt.Sprintf("Found %d chapter links on page", chapterLinks.Length()))
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

	return nil, errors.New("No chapters found for this manga")
}
