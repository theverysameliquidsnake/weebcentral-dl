package main

import (
//	"errors"
	"fmt"
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
		// add handler
		fmt.Println(err)
	}
	defer resp.Body.Close()
	
	// Create file
	out, err := os.Create(filePath)
	if err != nil {
		// add handler
		fmt.Println(err)
	}
	defer out.Close()

	// Write
	_, err = io.Copy(out, resp.Body)
	if err != nil {
		// add handler
		fmt.Println(err)
	}
	waitGroup.Done()

	return nil
}

func downloadChapter(downloadPath string, chapterTitle string, chapterUrl string) error {
	// Download all images from chapter
	images, err := extractChapterImageLinks(chapterUrl)
	if err != nil {
		// add handler
		fmt.Println(err)
	}

	//	Create chapter folder
	chapterFolderPath := filepath.Join(downloadPath, chapterTitle)
	err = createDirectory(chapterFolderPath)
	if err != nil {
		// add handler
		fmt.Println(err)
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
