package main

import (
//	"errors"
	"fmt"
	"sync"
	"strings"
	"os"
	"io"
	"net/http"
)

func downloadImage(imageUrl string, filePath string, waitGroup *sync.WaitGroup) error {
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
	waitGroup.Done()

	return nil
}

func downloadChapter(chapterUrl string) error {
	// Download all images from chapter
	images, err := extractChapterImageLinks(chapterUrl)
	if err != nil {
		// add handler
	}
	
	// Download images async
	wg := sync.WaitGroup{}
	for index := range images {
		wg.Add(1)
		parts := strings.Split(images[index], "/")
		filename := parts[len(parts) - 1]
		go downloadImage(images[index], fmt.Sprintf("dan/%s", filename), &wg)
	}
	wg.Wait()

	return nil
}
