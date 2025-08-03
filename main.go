package main

import (
	"fmt"
	"log"
)

func main() {
	// Verify provided args
	args := getArgs()
	if !args.hasValidArg || args.help {
		printHelp()
		return
	}

	// Search manga
	if len(args.title) > 0 {
		mangaUrl, err := searchManga(args.title)
		if err != nil {
			log.Fatalln(err)
		}
		
		// Get manga id and slug from URL
		id, slug, err := extractAttrFromUrl(mangaUrl)
		if err != nil {
			log.Fatalln(err)
		}

		// Retrieve manga chapters with URL
		chapters, err := getChaptersFromList(fmt.Sprintf("https://weebcentral.com/series/%s/full-chapter-list", id))
		if err != nil {
			log.Fatalln(err)
		}
		fmt.Println(slug, len(chapters), chapters["Chapter 1"])
	}

	/*// Retrieve manga chapters with URL
	chapters, err := getChaptersFromList(manga.chapterListUrl)
	if err != nil {
		log.Fatalln(err)
	}
	manga.chapters = chapters

	// Get images for random chapter test
	err = downloadChapter(manga.chapters["Chapter 1"])
	if err != nil {
		log.Fatalln(err)
	}
	
	// Download random image test err = downloadImage(images[0], "dan/dan.png")*/ }
