package main

import (
	"fmt"
	"time"
	"net/http"
	"strings"
	"log"
	
	"github.com/PuerkitoBio/goquery"
)

type manga struct {
	id string
	title string
	slug string
	baseUrl string
	rssUrl string
	chapterListUrl string
}	

func searchManga(title string) string {
	// Construct request
	const url = "https://weebcentral.com/search/simple"
	const method = "POST"
	payload := strings.NewReader(fmt.Sprintf("text=%s", title))
	
	req, err := http.NewRequest(method, url, payload)
	if err != nil {
		log.Fatalln("Error creating request:", err)
	}

	req.Header.Add("Content-Type", "application/x-www-form-urlencoded")
	req.Header.Add("User-Agent", "Mozilla/5.0 (X11; Linux x86_64; rv:141.0) Gecko/20100101 Firefox/141.0")

	client := http.Client{
		Timeout: time.Second * 10,
	}
	
	// Make request
	resp, err := client.Do(req)
	if err != nil {
		log.Fatalln("Error making request:", err)
	}
	defer resp.Body.Close()
	
	// Parse HTML response
	doc, err := goquery.NewDocumentFromReader(resp.Body)
	if err != nil {
		log.Fatalln("Error reading response body:", err)
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
			return mangaUrl
		}
	}
	
	log.Println("No manga found in search results")
	return ""
}

func getMangaSlug(mangaUrl string) string {
	// Extract the manga title slug from the URL
	parts := strings.Split(mangaUrl, "/")
	// URL format: https://weebcentral.com/series/ID/SLUG
	if len(parts) >= 6 {
		return parts[len(parts) - 1]
	}

	return ""
}

func getBaseUrl(mangaUrl string) string {
	// Get the base URL without the title slug
	seriesId, err := extractSeriesId(mangaUrl)
	if err == nil {
		return fmt.Sprintf("https://weebcentral.com/series/%s/", seriesId)
	}

	return ""
}

func getRssUrl(mangaUrl string) string {
	// Construct the RSS feed URL using the series ID
	seriesId, err := extractSeriesId(mangaUrl)
	if err == nil {
		return fmt.Sprintf("https://weebcentral.com/series/%s/rss", seriesId)
	}
	
	return ""
}

func getChapterListUrl(mangaUrl string) string {
	// Construct the full chapter list URL using the base URL
	baseUrl := getBaseUrl(mangaUrl)
	if len(baseUrl) > 0 {
		return fmt.Sprintf("%sfull-chapter-list", baseUrl)
	}

	return ""
}

func getChaptersFromList(chapterListUrl string) {
	// Get chaptr links from the full chapter list page
	// Construct request
	const method = "GET"
	
	req, err := http.NewRequest(method, chapterListUrl, nil)
	if err != nil {
		log.Fatalln("Error creating request:", err)
	}

	req.Header.Add("User-Agent", "Mozilla/5.0 (X11; Linux x86_64; rv:141.0) Gecko/20100101 Firefox/141.0")

	client := http.Client{
		Timeout: time.Second * 10,
	}

	// Make request
	resp, err := client.Do(req)
	if err != nil {
		log.Fatalln("Error making request:", err)
	}
	defer resp.Body.Close()

	// Parse HTML response
	doc, err := goquery.NewDocumentFromReader(resp.Body)
	if err != nil {
		log.Fatalln("Error reading response body:", err)
	}

	chapterLinks := doc.Find("a[href*='/chapters/']")
	if chapterLinks.Length() > 0 {
		log.Println(fmt.Sprintf("Found %d raw chapter links on page", chapterLinks.Length()))
		// Continue here
	}
}

func main() {
	mangaUrl := searchManga("Dandadan")
	chapterListUrl := getChapterListUrl(mangaUrl)
	getChaptersFromList(chapterListUrl)
}
