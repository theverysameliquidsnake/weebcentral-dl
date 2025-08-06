package main

import (
	"archive/zip"
	"path/filepath"
	"os"
	"io"
)

func compressChapter(downloadPath string, chapterTitle string, compressMethod string) error {
	// Create archive file
	archive, err := os.Create(filepath.Join(downloadPath, chapterTitle + "." + compressMethod))
	if err != nil {
		// add
	}

	writter := zip.NewWritter(archive)



}
