package main

import (
	"os"
	"fmt"
	"strings"
	"strconv"
	"errors"
)

type Args struct {
	hasEnoughArgs bool
	help bool
	install bool
	title string
	output string
	compress string
	prefix string
	first float32
	isFirstSet bool
	last float32
	isLastSet bool
	verbose bool
}

func printHelp() {
	// Print cli usage help
	fmt.Println("usage: weebcentral-dl [--option] [argument]")
	fmt.Println("operations:")
	// TO DO
	// continue usage help text
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
				if len(args) > (index + 1) && !strings.HasPrefix(args[index + 1], "-") {
					parsedArgs.title = args[index + 1]
					parsedArgs.hasEnoughArgs = true
				}

			// Check for output path arg
			case "-o", "--output":
				// Check if it has string entry following
				if len(args) > (index + 1) && !strings.HasPrefix(args[index + 1], "-") {
					parsedArgs.output = args[index + 1]
				}
			
			// Check for compress option arg
			case "-c", "--compress":
				// Check if it has string entry following (optional with zip being default)	
				var compressVariant string
				if len(args) > (index + 1) && !strings.HasPrefix(args[index + 1], "-") {
					compressVariant = args[index + 1]
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
				if len(args) > (index + 1) && !strings.HasPrefix(args[index + 1], "-") {
					parsedArgs.prefix = args[index + 1]
				}
			
			// Check for first num of volume arg
			case "-f", "--first":
				// Check if it has string entry following
				if len(args) > (index + 1) && !strings.HasPrefix(args[index + 1], "-") {
					first, err := strconv.ParseFloat(args[index + 1], 32)
					if err != nil {
						return nil, errors.New("Invalid value for \"-f\" or \"--first\" flag")
					}
					parsedArgs.isFirstSet = true
					parsedArgs.first = float32(first)
				}

			// Check for last num of volume arg
			case "-l", "--last":
				// Check if it has string entry following
				if len(args) > (index + 1) && !strings.HasPrefix(args[index + 1], "-") {
					last, err := strconv.ParseFloat(args[index + 1], 32)
					if err != nil {
						return nil, errors.New("Invalid value for \"-l\" or \"--last\" flag")
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
