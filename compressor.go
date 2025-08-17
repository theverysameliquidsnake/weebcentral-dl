package main

import (
	"archive/zip"
	"fmt"
	"io"
	"io/ioutil"
	"log"
	"os"
	"path/filepath"
)

func compressChapter(downloadPath string, chapterTitle string, compressMethod string) error {

	debugOutput(fmt.Sprintf("Compressing %s", chapterTitle))

	// Create archive file
	archive, err := os.Create(filepath.Join(downloadPath, chapterTitle+"."+compressMethod))
	if err != nil {
		return fmt.Errorf("could not init archive file: %w", err)
	}

	writter := zip.NewWriter(archive)

	// Add chapter images to archive
	files, err := ioutil.ReadDir(filepath.Join(downloadPath, chapterTitle))
	if err != nil {
		return fmt.Errorf("could not read images in directory: %w", err)
	}
	for _, file := range files {
		f, err := os.Open(filepath.Join(downloadPath, chapterTitle, file.Name()))
		if err != nil {
			return fmt.Errorf("could not open file: %w", err)
		}
		w, err := writter.Create(f.Name())
		if err != nil {
			return fmt.Errorf("could not create file: %w", err)
		}
		_, err = io.Copy(w, f)
		if err != nil {
			return fmt.Errorf("could not copy image to archive: %w", err)
		}
		err = f.Close()
		if err != nil {
			return fmt.Errorf("could not close file: %w", err)
		}
	}

	// Close archive
	writter.Close()
	err = archive.Close()
	if err != nil {
		return fmt.Errorf("could not close archive: %w", err)
	}

	log.Printf("Compressed %s to %s\n", chapterTitle, compressMethod)

	return nil
}
