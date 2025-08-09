package main

import (
	"fmt"
	"log"
)

var isDebugOutputEnabled bool

func main() {
	log.SetFlags(0)
	log.SetOutput(new(customWriter))

	// Verify provided args
	args, err := getArgs()
	if err != nil {
		log.Fatalln(err)
	}
	if !args.hasEnoughArgs || args.help {
		printHelp()
		return
	}

	// Install Playwright dependencies
	if args.install {
		err := installPlaywright()
		if err != nil {
			log.Fatalln(err)
		}
	}

	// Check if debug output enabled
	isDebugOutputEnabled = args.verbose	

	// Search manga
	if len(args.title) > 0 {
		// Begin searching process
		mangaUrl, err := searchManga(args.title)
		if err != nil {
			log.Println(err)
			return
		}

		// Get manga id and slug from URL
		id, slug, err := extractAttrFromUrl(mangaUrl)
		if err != nil {
			log.Println(err)
			return
		}

		// Retrieve manga chapters with URL
		chapters, err := getChaptersFromList(fmt.Sprintf("https://weebcentral.com/series/%s/full-chapter-list", id))
		if err != nil {
			log.Println(err)
			return
		}

		// Filter chapters if prefix, first or last is set
		for chapterTitle := range chapters {
			isToDownload, err := isChapterToDownload(args.prefix, args.first, args.isFirstSet, args.last, args.isLastSet, chapterTitle)
			if err != nil {
				log.Println(err)
				return
			}
			if !isToDownload {
				delete(chapters, chapterTitle)
			}
		}
		if len(chapters) == 0 {
			log.Println("No chapters to download")
			return
		}
		
		// Download chapters to manga title folder
		downloadFolderPath, err := resolveDownloadFolderPath(slug, args.output)
		if err != nil {
			log.Println(err)
			return
		}

		for chapterTitle, chapterUrl := range chapters {
			err = downloadChapter(downloadFolderPath, chapterTitle, chapterUrl)
			if err != nil {
				log.Println(err)
				return
			}
		}

		// Compress chapters if needed
		if len(args.compress) > 0 {
			for chapterTitle, _ := range chapters {
				err = compressChapter(downloadFolderPath, chapterTitle, args.compress)
				if err != nil {
					log.Println(err)
					return
				}
			}
		}
	}
}
