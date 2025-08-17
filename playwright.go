package main

import (
	"fmt"
	"strings"

	"github.com/playwright-community/playwright-go"
)

func installPlaywright() error {
	// Install playwright dependencies if it is missing
	installOptions := playwright.RunOptions{
		Browsers: []string{"chromium"},
	}
	err := playwright.Install(&installOptions)
	if err != nil {
		return fmt.Errorf("could not install Playwright: %w", err)
	}

	return nil
}

func extractChapterImageLinks(chapterUrl string) ([]string, error) {
	// Init playwright
	pw, err := playwright.Run()
	if err != nil {
		return nil, fmt.Errorf("could not launch playwright: %w", err)
	}

	browser, err := pw.Chromium.Launch(
		playwright.BrowserTypeLaunchOptions{
			Headless: playwright.Bool(true),
		},
	)
	if err != nil {
		return nil, fmt.Errorf("could not launch browser: %w", err)
	}

	page, err := browser.NewPage()
	if err != nil {
		return nil, fmt.Errorf("could not create page: %w", err)
	}

	err = page.SetViewportSize(1920, 1080)
	if err != nil {
		return nil, fmt.Errorf("could not set page height and width: %w", err)
	}

	// Navigate to chapter page and collect image links
	_, err = page.Goto(chapterUrl, playwright.PageGotoOptions{
		WaitUntil: playwright.WaitUntilStateNetworkidle,
		Timeout:   playwright.Float(60000),
	})
	if err != nil {
		return nil, fmt.Errorf("could not navigate to provided URL: %w", err)
	}

	_, err = page.WaitForSelector("img", playwright.PageWaitForSelectorOptions{
		State:   playwright.WaitForSelectorStateVisible,
		Timeout: playwright.Float(30000),
	})
	if err != nil {
		return nil, fmt.Errorf("could not load page images: %w", err)
	}

	// Extract all images
	images, err := page.Locator("img").All()
	if err != nil {
		return nil, fmt.Errorf("could not extract page images: %w", err)
	}
	imageLinks := []string{}
	for index := range images {
		src, err := images[index].GetAttribute("src")
		if err != nil {
			return nil, fmt.Errorf("could not get image src: %w", err)
		}

		srcLower := strings.ToLower(src)
		if strings.HasSuffix(srcLower, ".png") && !strings.HasSuffix(srcLower, "/static/images/brand.png") {
			imageLinks = append(imageLinks, src)
		}
	}

	// Close browser and playwright
	err = browser.Close()
	if err != nil {
		return nil, fmt.Errorf("could not close browser: %w", err)
	}

	err = pw.Stop()
	if err != nil {
		return nil, fmt.Errorf("could not stop playwright: %w", err)
	}

	return imageLinks, nil
}
