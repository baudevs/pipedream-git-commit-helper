package main

import (
	"bufio"
	"fmt"
	"os"
	"os/exec"
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

type Config struct {
	Schema    string            `yaml:"schema"`    // The schema version
	Workflows map[string]string `yaml:"workflows"` // Mapping of workflow directories to names
	Steps     map[string]string `yaml:"steps"`     // Mapping of step directories to namespaces
}

const configFileName = "pipedream-config.yaml"

/* const lastCommitFile = ".last_commit_info" */
// checkGitChanges checks for both staged and unstaged changes
func checkGitChanges() ([]string, []string) {
	repo, err := git.PlainOpen(".")
	if err != nil {
		fmt.Println("Error opening Git repository:", err)
		return nil, nil
	}

	// Get git worktree
	w, err := repo.Worktree()
	if err != nil {
		fmt.Println("Error getting git worktree:", err)
		return nil, nil
	}

	// Get git status
	status, err := w.Status()
	if err != nil {
		fmt.Println("Error getting git status:", err)
		return nil, nil
	}

	// Slices to hold staged and unstaged files
	stagedFiles := []string{}
	unstagedFiles := []string{}

	// Loop through status and differentiate between staged and unstaged changes
	for file, s := range status {
		if s.Staging != git.Unmodified { // This means the file is staged for commit
			stagedFiles = append(stagedFiles, file)
		}
		if s.Worktree != git.Unmodified { // This means the file is modified but not staged
			unstagedFiles = append(unstagedFiles, file)
		}
	}

	return stagedFiles, unstagedFiles
}

func main() {
	// Check if user provided the init command
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
		}
	}

	// Ensure that the project is initialized
	if !isInitialized() {
		fmt.Println("Error: This is not an initialized Pipedream project.")
		fmt.Println("Please run 'pdcommit init' to initialize the project.")
		return
	}

	// Load project config
	config, err := loadConfig()
	if err != nil {
		fmt.Println("Error loading configuration:", err)
		return
	}

	// Get git changes (both staged and unstaged)
	stagedFiles, unstagedFiles := checkGitChanges()

	if len(stagedFiles) == 0 && len(unstagedFiles) == 0 {
		fmt.Println("No changes detected.")
		return
	}

	// Warn about unstaged changes
	if len(unstagedFiles) > 0 {
		fmt.Println("Warning: You have unstaged changes. These changes will not be included in the commit.")
		// Prompt to stage everything
		stage := prompt("Do you want to stage all changes? (y/n): ", "y")
		if stage == "y" || stage == "Y" {
			// Stage all changes
			cmd := exec.Command("git", "add", ".")
			cmd.Stdout = os.Stdout
			cmd.Stderr = os.Stderr
			if err := cmd.Run(); err != nil {
				fmt.Println("Error staging files:", err)
				return
			}
			fmt.Println("All changes have been staged.")
			// Re-check for staged files after staging
			stagedFiles, _ = checkGitChanges()
		}
	}

	// If no staged files after attempting to stage, abort
	if len(stagedFiles) == 0 {
		fmt.Println("No files are staged for commit. Commit aborted.")
		return
	}

	workflowSteps := proposeWorkflowFromGit(stagedFiles, config)
	if len(workflowSteps) == 0 {
		fmt.Println("No matching workflow or step found.")
		return
	}

	// Prompt for commit type and message
	commitType := prompt("Commit type (add, fix, change, remove): ", "fix")

	finalizeCommit(commitType, workflowSteps)
}

func initializeProject() {
	if _, err := os.Stat(configFileName); err == nil {
		fmt.Println("This project is already initialized.")
		return
	}

	// Initialize the config struct with the schema version
	config := Config{
		Schema:    "baudevs/2024-09-29",
		Workflows: make(map[string]string),
		Steps:     make(map[string]string),
	}

	// Traverse the project directory to find workflow.yaml files and steps
	err := filepath.Walk(".", func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return err
		}

		// If a workflow.yaml file is found, add the workflow
		if info.Name() == "workflow.yaml" {
			dir := filepath.Dir(path)
			workflowName := filepath.Base(dir)

			// Prompt user to rename workflow
			promptWorkflow := prompt(fmt.Sprintf("Detected workflow in %s. Enter workflow name (default: %s): ", dir, workflowName), workflowName)
			config.Workflows[workflowName] = promptWorkflow

			// Now scan the directory for steps (each subdirectory)
			subdirs, err := os.ReadDir(dir)
			if err != nil {
				return err
			}
			for _, subdir := range subdirs {
				if subdir.IsDir() {
					stepName := subdir.Name()
					fullPath := filepath.Join(dir, stepName)

					// Prompt user to rename step
					promptStep := prompt(fmt.Sprintf("Detected step %s in %s. Enter step namespace (default: %s): ", stepName, fullPath, stepName), stepName)
					config.Steps[filepath.Join(workflowName, stepName)] = promptStep
				}
			}
		}

		return nil
	})

	if err != nil {
		fmt.Printf("Error scanning the project directory: %v\n", err)
		return
	}

	// Save the config to the config file
	saveConfig(config)
	fmt.Println("Pipedream project initialized successfully with schema version baudevs/2024-09-29.")
}

func isInitialized() bool {
	_, err := os.Stat(configFileName)
	return err == nil
}

// Load the config and validate the schema version
func loadConfig() (Config, error) {
	var config Config
	data, err := os.ReadFile(configFileName)
	if err != nil {
		return config, err
	}
	err = yaml.Unmarshal(data, &config)
	if err != nil {
		return config, err
	}

	// Validate the schema version
	if config.Schema != "baudevs/2024-09-29" {
		fmt.Printf("Warning: This project is using an unsupported schema version: %s\n", config.Schema)
	}

	return config, nil
}

func saveConfig(config Config) {
	data, err := yaml.Marshal(&config)
	if err != nil {
		fmt.Println("Error saving configuration:", err)
		return
	}
	err = os.WriteFile(configFileName, data, 0644)
	if err != nil {
		fmt.Println("Error writing configuration file:", err)
	}
}

// proposeWorkflowFromGit traverses the directory structure and matches git changes with config
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

// buildCommitMessage constructs the commit message
func buildCommitMessage(commitType string, workflowSteps []string, message string) string {
	wfMessage := ""
	for _, wfStep := range workflowSteps {
		parts := strings.Split(wfStep, ":")
		wfMessage += fmt.Sprintf("[[%s][%s]]", parts[0], parts[1])
	}
	return fmt.Sprintf("%s%s %s", commitType, wfMessage, message)
}

// Capture a multi-line commit message from user input
func captureMultiLineInput(promptMessage string) string {
	fmt.Println(promptMessage)

	reader := bufio.NewReader(os.Stdin)
	var lines []string

	for {
		input, _ := reader.ReadString('\n')
		input = strings.TrimSpace(input)

		// End capturing if the user enters an empty line
		if input == "" {
			break
		}

		lines = append(lines, input)
	}

	return strings.Join(lines, "\n")
}

// checkChanges detects both staged and unstaged changes
func checkChanges() (bool, bool) {
	fmt.Println("Checking changes...")
	cmd := exec.Command("git", "status", "--porcelain")
	output, err := cmd.Output()
	fmt.Println("Output:", string(output))
	if err != nil {
		fmt.Println("Error running git status:", err)
		return false, false
	}

	stagedFiles := false
	unstagedFiles := false
	scanner := bufio.NewScanner(strings.NewReader(string(output)))

	for scanner.Scan() {
		line := scanner.Text()
		fmt.Println("Line:", line)

		// Handle staged changes (Modified, Added, Deleted)
		if strings.HasPrefix(line, "M ") || strings.HasPrefix(line, "A ") || strings.HasPrefix(line, "D ") {
			stagedFiles = true
		}

		// Handle unstaged changes (but tracked, meaning the file has been modified but not staged)
		if strings.HasPrefix(line, " M") || strings.HasPrefix(line, " A") || strings.HasPrefix(line, " D") {
			unstagedFiles = true
		}
	}

	return stagedFiles, unstagedFiles
}

// finalizeCommit reviews and confirms the commit process
func finalizeCommit(commitType string, workflowSteps []string) {
	// Check if there are staged or unstaged changes
	stagedFiles, unstagedFiles := checkChanges()

	if !stagedFiles && !unstagedFiles {
		fmt.Println("No changes detected. Nothing to commit.")
		return
	}

	// Warn about unstaged changes
	if unstagedFiles {
		fmt.Println("Warning: You have unstaged changes. These changes will not be included in the commit.")

		// Prompt to stage everything
		stage := prompt("Do you want to stage all changes? (y/n): ", "y")
		if strings.ToLower(stage) == "y" {
			cmd := exec.Command("git", "add", ".")
			cmd.Stdout = os.Stdout
			cmd.Stderr = os.Stderr

			if err := cmd.Run(); err != nil {
				fmt.Println("Error staging files:", err)
				return
			}

			fmt.Println("All changes have been staged.")
			// Re-check for staged files after adding
			stagedFiles, _ = checkChanges()
		}
	}

	// If no staged files after prompt, abort
	if !stagedFiles {
		fmt.Println("No files are staged for commit. Commit aborted.")
		return
	}

	// Capture the multi-line commit message
	commitMessage := captureMultiLineInput("Enter commit message (press Enter twice to finish):")

	// Build the final commit message
	finalCommitMessage := buildCommitMessage(commitType, workflowSteps, commitMessage)

	// Show the command to the user
	fmt.Println("The following command is ready for review:")
	fmt.Printf("git commit -m \"%s\"\n", finalCommitMessage)

	// Ask if they want to proceed
	proceed := prompt("Do you want to proceed with the commit? (y/n): ", "y")

	if strings.ToLower(proceed) == "y" {
		// Execute the git commit command
		cmd := exec.Command("git", "commit", "-m", finalCommitMessage)
		cmd.Stdout = os.Stdout
		cmd.Stderr = os.Stderr

		if err := cmd.Run(); err != nil {
			fmt.Println("Error executing commit:", err)
		} else {
			fmt.Println("Commit successful.")
		}
	} else {
		fmt.Println("Commit aborted.")
	}
}

// prompt gets user input with a default value
func prompt(message, defaultVal string) string {
	fmt.Print(message)
	var input string
	fmt.Scanln(&input)
	if input == "" {
		return defaultVal
	}
	return input
}

/* func saveLastCommit(commitType string, workflowSteps []string) {
	f, _ := os.Create(lastCommitFile)
	defer f.Close()
	f.WriteString(fmt.Sprintf("last_type=%s\n", commitType))
	f.WriteString(fmt.Sprintf("last_workflow=%v\n", workflowSteps))
} */

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
