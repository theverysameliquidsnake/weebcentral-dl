package main

import (
	"errors"
	"fmt"
	"log"
	"sync"
	"strings"
	"os"
	"io"
	"net/http"
	"path/filepath"
)

func downloadImage(imageUrl string, filePath string, waitGroup *sync.WaitGroup) error {
	// Get image from request
	resp, err := http.Get(imageUrl)
	if err != nil {
		return errors.New(concatErrorString("Could not send request to get image: %s", err))
	}
	defer resp.Body.Close()
	
	// Create file
	out, err := os.Create(filePath)
	if err != nil {
		return errors.New(concatErrorString("Could not create file for downloaded image: %s", err))
	}
	defer out.Close()

	// Write
	_, err = io.Copy(out, resp.Body)
	if err != nil {
		return errors.New(concatErrorString("Could not copy downloaded image to file: %s", err))
	}
	waitGroup.Done()
	log.Println(fmt.Sprintf("Successfully downloaded: %s", imageUrl))

	return nil
}

func downloadChapter(downloadPath string, chapterTitle string, chapterUrl string) error {
	// Download all images from chapter
	images, err := extractChapterImageLinks(chapterUrl)
	if err != nil {
		return err
	}
	log.Println(fmt.Sprintf("Found %d images in %s", len(images), chapterTitle))

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
		imageFilePath := filepath.Join(chapterFolderPath, parts[len(parts) - 1])
		go downloadImage(image, imageFilePath, &wg)
	}
	wg.Wait()

	return nil
}
