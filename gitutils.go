package main

import (
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"strings"

	"github.com/fatih/color" // Import to handle colors
	"github.com/go-git/go-git/v5"
)

func checkGitChanges() ([]string, []string) {
	repo, err := git.PlainOpen(".")
	if err != nil {
		fmt.Println(ColorRed, "Error opening Git repository:", err, ColorReset)
		return nil, nil
	}

	w, err := repo.Worktree()
	if err != nil {
		fmt.Println(ColorRed, "Error getting git worktree:", err, ColorReset)
		return nil, nil
	}

	status, err := w.Status()
	if err != nil {
		fmt.Println(ColorRed, "Error getting git status:", err, ColorReset)
		return nil, nil
	}

	stagedFiles := []string{}
	unstagedFiles := []string{}

	for file, s := range status {
		if s.Staging != git.Unmodified {
			stagedFiles = append(stagedFiles, file)
		}
		if s.Worktree != git.Unmodified {
			unstagedFiles = append(unstagedFiles, file)
		}
	}

	color.Cyan("Checking for changes...")
	for _, file := range stagedFiles {
		color.Green("Staged: %s", file)
	}
	for _, file := range unstagedFiles {
		color.Yellow("Unstaged: %s", file)
	}

	return stagedFiles, unstagedFiles
}

func proposeWorkflowFromGit(files []string, config Config) []string {
	var workflowSteps []string

	for _, file := range files {
		dir := filepath.Dir(file)
		color.Cyan("Checking file: %s", file)
		color.Cyan("Directory: %s", dir)

		for wfDir, wfName := range config.Workflows {
			if strings.Contains(dir, wfDir) {
				color.Green("Matched workflow: %s -> %s", wfDir, wfName)
				for stepDir, stepName := range config.Steps {
					if strings.Contains(dir, stepDir) {
						color.Blue("Matched step: %s -> %s", stepDir, stepName)
						workflowSteps = append(workflowSteps, fmt.Sprintf("%s:%s", wfName, stepName))
					}
				}
			}
		}
	}

	return workflowSteps
}

// stageAllChanges stages all changes in the Git repository using `git add .`
func stageAllChanges() {
	cmd := exec.Command("git", "add", ".")
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	err := cmd.Run()
	if err != nil {
		color.Red("Error staging files: %v", err)
	} else {
		color.Green("All changes have been staged.")
	}
}
