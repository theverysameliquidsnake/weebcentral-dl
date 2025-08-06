package main

import (
	"fmt"
	"log"
)

func main() {
	// Verify provided args
	args := getArgs()
	if !args.hasEnoughArgs || args.help {
		printHelp()
		return
	}

	// Search manga
	if len(args.title) > 0 {
		// Install Playwright dependencies
		err := installPlaywright()
		if err != nil {
			log.Fatalln(err)
		}
		
		// Begin searching process
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

		// Filter chapters if prefix, first or last is set
		for chapterTitle := range chapters {
			isToDownload, err := isChapterToDownload(args.prefix, args.first, args.last, chapterTitle)
			if err != nil {
				log.Fatalln(err)
			}
			if !isToDownload {
				delete(chapters, chapterTitle)
			}
		}
		
		// Download chapters to manga title folder
		downloadFolderPath, err := resolveDownloadFolderPath(slug, args.output)
		if err != nil {
			log.Fatalln(err)
		}

		for chapterTitle, chapterUrl := range chapters {
			err = downloadChapter(downloadFolderPath, chapterTitle, chapterUrl)
			if err != nil {
				log.Fatalln(err)
			}
		}	
	}
}
