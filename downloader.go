package main

import (
	"fmt"
	"strings"
	"log"
	"errors"

	"github.com/PuerkitoBio/goquery"
	"github.com/playwright-community/playwright-go"
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

func extractChapterImageLinks(chapterUrl string) ([]string, error) {
	// Init playwright
	pw, err := playwright.Run()
	if err != nil {
		return nil, errors.New(fmt.Sprintf("Could not launch playwright:%s", err))
	}

	browser, err := pw.Chromium.Launch(
		playwright.BrowserTypeLaunchOptions {
			Headless: playwright.Bool(true),
		},
	)
	if err != nil {
		return nil, errors.New(fmt.Sprintf("Could not launch browser:%s", err))
	}

	page, err := browser.NewPage()
	if err != nil {
		return nil, errors.New(fmt.Sprintf("Could not create page:%s", err))
	}

	err = page.SetViewportSize(1920, 1080)
	if err != nil {
		return nil, errors.New(fmt.Sprintf("Could not set page height and width:%s", err))
	}

	// Navigate to chapter page and collect image links
	_, err = page.Goto(chapterUrl, playwright.PageGotoOptions {
		WaitUntil: playwright.WaitUntilStateNetworkidle,
		Timeout: playwright.Float(60000),
	})
	if err != nil {
		return nil, errors.New(fmt.Sprintf("Could not navigate to provided URL:%s", err))
	}

	_, err = page.WaitForSelector("img", playwright.PageWaitForSelectorOptions {
		State: playwright.WaitForSelectorStateVisible,
		Timeout: playwright.Float(30000),
	})
	if err != nil {
		return nil, errors.New(fmt.Sprintf("Could not load page images:%s", err))
	}

	// Extract all images
	images, err := page.Locator("img").All()
	if err != nil {
		return nil, errors.New(fmt.Sprintf("Could not extract page images:%s", err))
	}
	imageLinks := []string{}
	for index := range images {
		src, err := images[index].GetAttribute("src")
		if err != nil {
			return nil, errors.New(fmt.Sprintf("Could not get image src:%s", err))
		}
		
		srcLower := strings.ToLower(src)
		if strings.HasSuffix(srcLower, ".png") && !strings.HasSuffix(srcLower, "/static/images/brand.png") {
			imageLinks = append(imageLinks, src)
		}
	}

	// Close browser and playwright
	err = browser.Close()
	if err != nil {
		return nil, errors.New(fmt.Sprintf("Could not close browser:%s", err))
	}

	err = pw.Stop()
	if err != nil {
		return nil, errors.New(fmt.Sprintf("Could not stop playwright:%s", err))
	}

	return imageLinks, nil
}

func main() {
	// Search for URL for provided title
	mangaUrl, err := searchManga("Dandadan")
	if err != nil {
		log.Fatalln(err)
	}

	// Initiate start struct
	manga, err := constructManga(mangaUrl)
	if err != nil {
		log.Fatalln(err)
	}

	// Retrieve manga chapters with URL
	chapters, err := getChaptersFromList(manga.chapterListUrl)
	if err != nil {
		log.Fatalln(err)
	}
	manga.chapters = chapters

	// Get images for random chapter test
	images, err := extractChapterImageLinks(manga.chapters["Chapter 1"])
	if err != nil {
		log.Fatalln(err)
	}
	
	// Download random image test
	err = downloadImage(images[0], strings.Split(images[0], "/")[0])
}
