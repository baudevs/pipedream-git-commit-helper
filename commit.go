package main

import (
	"bufio"
	"fmt"
	"os"
	"os/exec"
	"strings"

	"github.com/fatih/color"
)

func buildCommitMessage(commitType string, workflowSteps []string, message string) string {
	wfMessage := ""
	for _, wfStep := range workflowSteps {
		parts := strings.Split(wfStep, ":")
		wfMessage += fmt.Sprintf("[[%s][%s]]", parts[0], parts[1])
	}
	return fmt.Sprintf("%s%s %s", commitType, wfMessage, message)
}

func captureMultiLineInput(promptMessage string) string {
	fmt.Println(promptMessage)

	reader := bufio.NewReader(os.Stdin)
	var lines []string

	for {
		input, _ := reader.ReadString('\n')
		input = strings.TrimSpace(input)

		if input == "" {
			break
		}

		lines = append(lines, input)
	}

	return strings.Join(lines, "\n")
}

func finalizeCommit(commitType string, workflowSteps []string) {
	commitMessage := captureMultiLineInput("Enter commit message (press Enter twice to finish):")
	finalCommitMessage := buildCommitMessage(commitType, workflowSteps, commitMessage)

	color.Yellow("The following command is ready for review:")
	fmt.Printf("git commit -m \"%s\"\n", finalCommitMessage)

	if confirmAction("Do you want to proceed with the commit?") {
		cmd := exec.Command("git", "commit", "-m", finalCommitMessage)
		cmd.Stdout = os.Stdout
		cmd.Stderr = os.Stderr
		if err := cmd.Run(); err != nil {
			color.Red("Error executing commit: %v", err)
		} else {
			color.Green("Commit successful.")
		}
	} else {
		color.Red("Commit aborted.")
	}
}
