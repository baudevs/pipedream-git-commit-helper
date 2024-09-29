package main

import (
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"strings"

	"github.com/go-git/go-git/v5"
)

// proposeWorkflowFromGit traverses the directory structure and matches Git changes with config
func proposeWorkflowFromGit(files []string, config Config) []string {
	var workflowSteps []string

	for _, file := range files {
		dir := filepath.Dir(file) // Get the directory of the changed file
		fmt.Println("Checking file:", file)
		fmt.Println("Directory:", dir)

		// Iterate over workflows in the config
		for wfDir, wfName := range config.Workflows {
			// Check if the directory contains the workflow directory
			if strings.Contains(dir, wfDir) {
				fmt.Printf("Matched workflow: %s -> %s\n", wfDir, wfName)
				// Now, iterate over the steps
				for stepDir, stepName := range config.Steps {
					if strings.Contains(dir, stepDir) {
						fmt.Printf("Matched step: %s -> %s\n", stepDir, stepName)
						workflowSteps = append(workflowSteps, fmt.Sprintf("%s:%s", wfName, stepName))
					}
				}
			}
		}
	}

	return workflowSteps
}

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

	return stagedFiles, unstagedFiles
}

func stageAllChanges() {
	cmd := exec.Command("git", "add", ".")
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	if err := cmd.Run(); err != nil {
		fmt.Println(ColorRed, "Error staging files:", err, ColorReset)
	}
}
