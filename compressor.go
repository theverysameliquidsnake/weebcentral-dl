package main

import (
	"archive/zip"
	"path/filepath"
	"os"
	"io"
	"io/ioutil"
)

func compressChapter(downloadPath string, chapterTitle string, compressMethod string) error {
	// Create archive file
	archive, err := os.Create(filepath.Join(downloadPath, chapterTitle + "." + compressMethod))
	if err != nil {
		// add
	}

	writter := zip.NewWriter(archive)

	// Add chapter images to archive
	files, err := ioutil.ReadDir(filepath.Join(downloadPath, chapterTitle))
	if err != nil {
		// add
	}
	for _, file := range files {
		f, err := os.Open(filepath.Join(downloadPath, chapterTitle, file.Name()))
		if err != nil {
			// add
		}
		w, err := writter.Create(f.Name())
		if err != nil {
			// add
		}
		_, err = io.Copy(w, f)
		if err != nil {
			// add
		}
		err = f.Close()
		if err != nil {
			// add
		}
	}

	// Close archive
	writter.Close()
	err = archive.Close()
	if err != nil {
		// add
	}

	return nil
}
