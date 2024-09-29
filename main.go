package main

import (
	"bufio"
	"fmt"
	"os"
	"path/filepath"
	"strings"

	"github.com/go-git/go-git/v5"
	"gopkg.in/yaml.v2"
)

type Step struct {
	Namespace string `yaml:"namespace"`
}

type Workflow struct {
	Name  string `yaml:"name"`
	Steps []Step `yaml:"steps"`
}

const lastCommitFile = ".last_commit_info"

func main() {
	repo, err := git.PlainOpen(".")
	if err != nil {
		fmt.Println("Error opening Git repository:", err)
		return
	}

	// Get git changes
	w, err := repo.Worktree()
	if err != nil {
		fmt.Println("Error getting git worktree:", err)
		return
	}

	status, err := w.Status()
	if err != nil {
		fmt.Println("Error getting git status:", err)
		return
	}

	changedFiles := []string{}
	for file := range status {
		changedFiles = append(changedFiles, file)
	}

	if len(changedFiles) == 0 {
		fmt.Println("No changes detected.")
		return
	}

	workflowSteps := proposeWorkflowFromGit(changedFiles)
	if len(workflowSteps) == 0 {
		fmt.Println("No matching workflow or step found.")
		return
	}

	// Prompt for commit type and message
	commitType := prompt("Commit type (add, fix, change, remove): ", "fix")
	message := prompt("Enter commit message: ", "")

	commitMsg := buildCommitMessage(commitType, workflowSteps, message)
	fmt.Println("The following command is ready for review:")
	fmt.Printf("git commit -m \"%s\"\n", commitMsg)

	if prompt("Do you want to proceed with the commit? (y/n): ", "y") == "y" {
		_, err = w.Add(".")
		if err != nil {
			fmt.Println("Error adding files:", err)
			return
		}

		_, err = w.Commit(commitMsg, &git.CommitOptions{})
		if err != nil {
			fmt.Println("Error committing:", err)
			return
		}

		fmt.Println("Commit successful.")
		saveLastCommit(commitType, workflowSteps)
	} else {
		fmt.Println("Commit aborted.")
	}
}

func prompt(message, defaultVal string) string {
	fmt.Print(message)
	reader := bufio.NewReader(os.Stdin)
	input, _ := reader.ReadString('\n')
	input = strings.TrimSpace(input)
	if input == "" {
		return defaultVal
	}
	return input
}

func proposeWorkflowFromGit(files []string) []string {
	var workflowSteps []string

	for _, file := range files {
		dir := filepath.Dir(file)
		workflowFile := filepath.Join(dir, "workflow.yaml")

		// Traverse up directories to find workflow.yaml
		for dir != "/" && !fileExists(workflowFile) {
			dir = filepath.Dir(dir)
			workflowFile = filepath.Join(dir, "workflow.yaml")
		}

		if fileExists(workflowFile) {
			workflow := parseWorkflow(workflowFile)
			stepCandidate := filepath.Base(filepath.Dir(file))

			for _, step := range workflow.Steps {
				if step.Namespace == stepCandidate {
					workflowSteps = append(workflowSteps, fmt.Sprintf("%s:%s", workflow.Name, step.Namespace))
				}
			}
		}
	}

	return workflowSteps
}

func parseWorkflow(file string) Workflow {
	var workflow Workflow
	data, _ := os.ReadFile(file)
	yaml.Unmarshal(data, &workflow)
	return workflow
}

func fileExists(filename string) bool {
	_, err := os.Stat(filename)
	return !os.IsNotExist(err)
}

func buildCommitMessage(commitType string, workflowSteps []string, message string) string {
	wfMessage := ""
	for _, wfStep := range workflowSteps {
		parts := strings.Split(wfStep, ":")
		wfMessage += fmt.Sprintf("[[%s][%s]]", parts[0], parts[1])
	}
	return fmt.Sprintf("%s%s %s", commitType, wfMessage, message)
}

func saveLastCommit(commitType string, workflowSteps []string) {
	f, _ := os.Create(lastCommitFile)
	defer f.Close()
	f.WriteString(fmt.Sprintf("last_type=%s\n", commitType))
	f.WriteString(fmt.Sprintf("last_workflow=%v\n", workflowSteps))
}
