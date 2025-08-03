package main

import (
	"os"
	"fmt"
	"strings"
)

type Args struct {
	hasValidArg bool
	help bool
	title string
}

func printHelp() {
	// Print cli usage help
	fmt.Println("usage: weebcentral-dl <operation> [...]")
	fmt.Println("operations:")
	// TO DO
	// continue usage help text
}

func getArgs() *Args {
	// Parse args
	var parsedArgs Args
	if len(os.Args) > 1 {
		args := os.Args[1:]
		// Check for help arg
		for _, arg := range args {
			if arg == "-h" || arg == "--help" {
				parsedArgs.help = true
				parsedArgs.hasValidArg = true
				break
			}
		}

		// Check for title arg
		for index, arg := range args {
			if arg == "-t" || arg == "--title" {
				// Check if it has string entry
				if len(args) > (index + 1) && !strings.HasPrefix(args[index + 1], "-") {
					parsedArgs.title = args[index + 1]
					parsedArgs.hasValidArg = true
					break
				}
			}
		}
	}

	return &parsedArgs
}
