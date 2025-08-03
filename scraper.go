package main

import (
	"strings"
	"fmt"
	"errors"
	"log"
	"github.com/PuerkitoBio/goquery"	
)

func searchManga(title string) (string, error) {
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
		// Get info from the first result
		firstResult := searchResults.First();
		mangaUrl, attrExists := firstResult.Attr("href")
		if attrExists {
			mangaTitle := strings.TrimSpace(firstResult.Find("div.flex-1").First().Text())
			log.Println("Found manga:", mangaTitle)
			log.Println("URL:", mangaUrl)
			return mangaUrl, nil
		}
	}
	
	return "", errors.New("No manga found in search results")
}

func getChaptersFromList(chapterListUrl string) (map[string]string, error) {
	// Request params
	method := "GET"
	
	doc, err := sendRequest(method, chapterListUrl, nil, nil)
	if err != nil {
		return nil, err
	}
	
	// Search for chapters
	chapterLinks := doc.Find("a[href*='/chapters/']")
	if chapterLinks.Length() > 0 {
		log.Println(fmt.Sprintf("Found %d raw chapter links on page", chapterLinks.Length()))
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
