package main

import (
	"fmt"
	"io"
	"log"
	"net/http"
	"os"
	"path/filepath"
	"strings"
	"sync"
)

func downloadImage(imageUrl string, filePath string, waitGroup *sync.WaitGroup) error {
	// Get image from request
	resp, err := http.Get(imageUrl)
	if err != nil {
		return fmt.Errorf("could not send request to get image: %w", err)
	}
	defer resp.Body.Close()

	// Create file
	out, err := os.Create(filePath)
	if err != nil {
		return fmt.Errorf("could not create file for downloaded image: %w", err)
	}
	defer out.Close()

	// Write
	_, err = io.Copy(out, resp.Body)
	if err != nil {
		return fmt.Errorf("could not copy downloaded image to file: %w", err)
	}
	waitGroup.Done()

	debugOutput(fmt.Sprintf("Successfully downloaded: %s", imageUrl))

	return nil
}

func downloadChapter(downloadPath string, chapterTitle string, chapterUrl string) error {
	// Download all images from chapter
	images, err := extractChapterImageLinks(chapterUrl)
	if err != nil {
		return err
	}
	log.Printf("Found %d images in %s\n", len(images), chapterTitle)

	//	Create chapter folder
	chapterFolderPath := filepath.Join(downloadPath, chapterTitle)
	err = createDirectory(chapterFolderPath)
	if err != nil {
		return err
	}

	// Download images async
	wg := sync.WaitGroup{}
	for _, image := range images {
		wg.Add(1)
		parts := strings.Split(image, "/")
		imageFilePath := filepath.Join(chapterFolderPath, parts[len(parts)-1])
		go downloadImage(image, imageFilePath, &wg)
	}
	wg.Wait()
	log.Printf("Downloaded %d images in %s\n", len(images), chapterTitle)

	return nil
}
