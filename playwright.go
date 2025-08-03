package main

import (
	"fmt"
	"strings"
	"errors"
	"github.com/playwright-community/playwright-go"
)

func installPlaywright() error {
	// Install playwright dependencies if it is missing
	installOptions := playwright.RunOptions {
		Browsers: []string{"chromium"},
	}
	err := playwright.Install(&installOptions)
	if err != nil {
		return errors.New(fmt.Sprintf("Could not install Playwright:%s", err))
	}

	return nil
}

func extractChapterImageLinks(chapterUrl string) ([]string, error) {
	// Init playwright
	pw, err := playwright.Run()
	if err != nil {
		return nil, errors.New(fmt.Sprintf("Could not launch playwright:%s", err))
	}

	browser, err := pw.Chromium.Launch(
		playwright.BrowserTypeLaunchOptions {
			Headless: playwright.Bool(true),
		},
	)
	if err != nil {
		return nil, errors.New(fmt.Sprintf("Could not launch browser:%s", err))
	}

	page, err := browser.NewPage()
	if err != nil {
		return nil, errors.New(fmt.Sprintf("Could not create page:%s", err))
	}

	err = page.SetViewportSize(1920, 1080)
	if err != nil {
		return nil, errors.New(fmt.Sprintf("Could not set page height and width:%s", err))
	}

	// Navigate to chapter page and collect image links
	_, err = page.Goto(chapterUrl, playwright.PageGotoOptions {
		WaitUntil: playwright.WaitUntilStateNetworkidle,
		Timeout: playwright.Float(60000),
	})
	if err != nil {
		return nil, errors.New(fmt.Sprintf("Could not navigate to provided URL:%s", err))
	}

	_, err = page.WaitForSelector("img", playwright.PageWaitForSelectorOptions {
		State: playwright.WaitForSelectorStateVisible,
		Timeout: playwright.Float(30000),
	})
	if err != nil {
		return nil, errors.New(fmt.Sprintf("Could not load page images:%s", err))
	}

	// Extract all images
	images, err := page.Locator("img").All()
	if err != nil {
		return nil, errors.New(fmt.Sprintf("Could not extract page images:%s", err))
	}
	imageLinks := []string{}
	for index := range images {
		src, err := images[index].GetAttribute("src")
		if err != nil {
			return nil, errors.New(fmt.Sprintf("Could not get image src:%s", err))
		}
		
		srcLower := strings.ToLower(src)
		if strings.HasSuffix(srcLower, ".png") && !strings.HasSuffix(srcLower, "/static/images/brand.png") {
			imageLinks = append(imageLinks, src)
		}
	}

	// Close browser and playwright
	err = browser.Close()
	if err != nil {
		return nil, errors.New(fmt.Sprintf("Could not close browser:%s", err))
	}

	err = pw.Stop()
	if err != nil {
		return nil, errors.New(fmt.Sprintf("Could not stop playwright:%s", err))
	}

	return imageLinks, nil
}
