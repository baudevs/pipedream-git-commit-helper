package main

import "fmt"

// DisplayHelp shows the available commands and options
func displayHelp() {
	fmt.Print(`Pipedream Git Commit Helper

Usage:
  pdcommit init               Initialize a Pipedream project
  pdcommit                    Commit changes detected in the project

Options:
  --help                      Show this help message
  man                         Display the manual with detailed descriptions of each command
`)
}

// DisplayMan shows the manual with detailed descriptions
func displayMan() {
	fmt.Print(`Pipedream Git Commit Helper Manual

COMMANDS

pdcommit init
  - Initializes a Pipedream project, scans for workflows and steps, and stores the mappings in pipedream-config.yaml.
  
pdcommit
  - Analyzes git changes, proposes commit messages, and helps automate the commit process.
  
--help
  - Displays a list of available commands and options.

man
  - Shows the full manual with detailed explanations of each command.
`)
}
