package main

import (
	"fmt"
	"os"
)

const (
	ColorReset  = "\033[0m"
	ColorRed    = "\033[31m"
	ColorGreen  = "\033[32m"
	ColorYellow = "\033[33m"
	ColorBlue   = "\033[34m"
	ColorCyan   = "\033[36m"
	ColorWhite  = "\033[97m"
)

func main() {
	fmt.Printf("%sWelcome to Pipedream Git Commit Helper!%s\n", ColorBlue, ColorReset)

	if len(os.Args) > 1 {
		switch os.Args[1] {
		case "--help":
			displayHelp()
			return
		case "man":
			displayMan()
			return
		case "init":
			initializeProject()
			return
		case "sync":
			syncConfig()
			return
		}
	}

	if !isInitialized() {
		fmt.Println(ColorRed, "Error: This is not an initialized Pipedream project.", ColorReset)
		fmt.Println("Please run 'pdcommit init' to initialize the project.")
		return
	}

	config, err := loadConfig()
	if err != nil {
		fmt.Println(ColorRed, "Error loading configuration:", err, ColorReset)
		return
	}

	stagedFiles, unstagedFiles := checkGitChanges()
	if len(stagedFiles) == 0 && len(unstagedFiles) == 0 {
		fmt.Println(ColorYellow, "No changes detected.", ColorReset)
		return
	}

	if len(unstagedFiles) > 0 {
		fmt.Printf("%sWarning: You have unstaged changes.%s\n", ColorYellow, ColorReset)
		if confirmAction("Do you want to stage all changes?") {
			stageAllChanges()
			stagedFiles, _ = checkGitChanges()
		}
	}

	if len(stagedFiles) == 0 {
		fmt.Println(ColorRed, "No files are staged for commit. Commit aborted.", ColorReset)
		return
	}

	workflowSteps := proposeWorkflowFromGit(stagedFiles, config)
	if len(workflowSteps) == 0 {
		fmt.Println(ColorRed, "No matching workflow or step found.", ColorReset)
		return
	}

	commitType := selectCommitType()
	finalizeCommit(commitType, workflowSteps)
}
