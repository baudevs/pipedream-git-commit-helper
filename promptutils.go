package main

import (
	"fmt"
	"os"

	"github.com/manifoldco/promptui"
)

func prompt(message, defaultVal string) string {
	prompt := promptui.Prompt{
		Label:   message,
		Default: defaultVal,
	}
	result, err := prompt.Run()
	if err != nil {
		fmt.Printf("Prompt failed %v\n", err)
		return defaultVal
	}
	return result
}

func selectCommitType() string {
	prompt := promptui.Select{
		Label: "Select Commit Type",
		Items: []string{"add", "fix", "change", "remove"},
	}

	_, commitType, err := prompt.Run()
	if err != nil {
		fmt.Printf("Prompt failed %v\n", err)
		os.Exit(1)
	}

	return commitType
}

func confirmAction(label string) bool {
	prompt := promptui.Prompt{
		Label:     label,
		IsConfirm: true,
	}

	_, err := prompt.Run()
	return err == nil
}
