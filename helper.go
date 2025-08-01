package main

import (
	"errors"
	"strings"
	"regexp"
	"net/http"
)

func sendRequest(method string, url string, headers []string, payload string)

func extractSeriesId(mangaUrl string) (string, error) {
	// Extract the series ID from the URL
    // Example: from https://weebcentral.com/series/01J76XYFM1TWGNNQ2Y2T8V7E8Y/Wistoria-Wand-and-Sword get 01J76XYFM1TWGNNQ2Y2T8V7E8Y
	if len(mangaUrl) > 0 {
		parts := strings.Split(mangaUrl, "/")
		for index := range parts {
			// Series IDs are 26 characters long
			if len(parts[index]) == 26 && regexp.MustCompile(`^[a-zA-Z0-9]*$`).MatchString(parts[index]) {
				return parts[index], nil
			}
		}
	}

	return "", errors.New("Url is empty")
}
