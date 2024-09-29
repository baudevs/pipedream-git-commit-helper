package main

import (
	"bufio"
	"fmt"
	"os"
	"os/exec"
	"strings"
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

	fmt.Println("The following command is ready for review:")
	fmt.Printf("git commit -m \"%s\"\n", finalCommitMessage)

	if confirmAction("Do you want to proceed with the commit?") {
		cmd := exec.Command("git", "commit", "-m", finalCommitMessage)
		cmd.Stdout = os.Stdout
		cmd.Stderr = os.Stderr
		if err := cmd.Run(); err != nil {
			fmt.Println(ColorRed, "Error executing commit:", err, ColorReset)
		} else {
			fmt.Println(ColorGreen, "Commit successful.", ColorReset)
		}
	} else {
		fmt.Println(ColorRed, "Commit aborted.", ColorReset)
	}
}
