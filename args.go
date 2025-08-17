package main

import (
	"errors"
	"fmt"
	"os"
	"strconv"
	"strings"
)

type Args struct {
	hasEnoughArgs bool
	help          bool
	install       bool
	title         string
	output        string
	compress      string
	prefix        string
	first         float32
	isFirstSet    bool
	last          float32
	isLastSet     bool
	verbose       bool
}

func printHelp() {
	// Print cli usage help
	fmt.Println("cli downloader for weebcentral.com")
	fmt.Println("")
	fmt.Println("usage: weebcentral-dl -h | -v | -i")
	fmt.Println("usage: weebcentral-dl [-t title] [-f number] [-l number] [-p prefix]")
	fmt.Println("                      [-o directory] [-c format]")
	fmt.Println("")
	fmt.Println("options:")
	fmt.Println(concatHelpStringOption("-h", "--help", "display help message and exit"))
	fmt.Println(concatHelpStringOption("-t", "--title=title", "search manga by specified title"))
	fmt.Println(concatHelpStringOption("-f", "--first=number", "filter chapters equal or newer than specified number"))
	fmt.Println(concatHelpStringOption("-l", "--last=number", "filter chapters equal or older than specified number"))
	fmt.Println(concatHelpStringOption("-p", "--prefix=prefix", "filter chapters by specified chapter prefix"))
	fmt.Println(concatHelpStringOption("-o", "--output=directory", "download to specified directory"))
	fmt.Println(concatHelpStringOption("-c", "--compress=format", "compress to specified format (cbz or zip)"))
	fmt.Println(concatHelpStringOption("-i", "--install", "install required Playwright dependencies"))
	fmt.Println(concatHelpStringOption("-v", "--verbose", "enable detailed log output"))
}

func concatHelpStringOption(shortOpt string, longOpt string, desc string) string {
	helpOptionEntry := strings.Repeat(" ", 2)
	helpOptionEntry += shortOpt + ", " + longOpt
	helpOptionEntry += strings.Repeat(" ", 32-len(helpOptionEntry))
	helpOptionEntry += desc

	return helpOptionEntry
}

func getArgs() (*Args, error) {
	// Parse args
	var parsedArgs Args
	if len(os.Args) > 1 {
		args := os.Args[1:]

		for index, arg := range args {
			switch arg {
			// Check for help arg
			case "-h", "--help":
				parsedArgs.help = true
				parsedArgs.hasEnoughArgs = true

			// Check for Playwright installation arg
			case "-i", "--install":
				parsedArgs.install = true
				parsedArgs.hasEnoughArgs = true

			// Check for manga title arg
			case "-t", "--title":
				// Check if it has string entry following
				if len(args) > (index+1) && !strings.HasPrefix(args[index+1], "-") {
					parsedArgs.title = args[index+1]
					parsedArgs.hasEnoughArgs = true
				}

			// Check for output path arg
			case "-o", "--output":
				// Check if it has string entry following
				if len(args) > (index+1) && !strings.HasPrefix(args[index+1], "-") {
					parsedArgs.output = args[index+1]
				}

			// Check for compress option arg
			case "-c", "--compress":
				// Check if it has string entry following (optional with zip being default)
				var compressVariant string
				if len(args) > (index+1) && !strings.HasPrefix(args[index+1], "-") {
					compressVariant = args[index+1]
				}
				switch compressVariant {
				case "zip":
					parsedArgs.compress = "zip"
				case "cbz":
					parsedArgs.compress = "cbz"
				default:
					parsedArgs.compress = "zip"
				}

			// Check for prefix volume arg
			case "-p", "--prefix":
				if len(args) > (index+1) && !strings.HasPrefix(args[index+1], "-") {
					parsedArgs.prefix = args[index+1]
				}

			// Check for first num of volume arg
			case "-f", "--first":
				// Check if it has string entry following
				if len(args) > (index+1) && !strings.HasPrefix(args[index+1], "-") {
					first, err := strconv.ParseFloat(args[index+1], 32)
					if err != nil {
						return nil, errors.New("invalid value for \"-f\" or \"--first\" flag")
					}
					parsedArgs.isFirstSet = true
					parsedArgs.first = float32(first)
				}

			// Check for last num of volume arg
			case "-l", "--last":
				// Check if it has string entry following
				if len(args) > (index+1) && !strings.HasPrefix(args[index+1], "-") {
					last, err := strconv.ParseFloat(args[index+1], 32)
					if err != nil {
						return nil, errors.New("invalid value for \"-l\" or \"--last\" flag")
					}
					parsedArgs.last = float32(last)
					parsedArgs.isLastSet = true
				}

			// Check for enable verbose arg
			case "-v", "--verbose":
				parsedArgs.verbose = true
			}
		}
	}

	return &parsedArgs, nil
}
