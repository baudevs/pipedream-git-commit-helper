package main

import (
	"fmt"
	"os"
	"path/filepath"

	"github.com/fatih/color" // Import for colorization
	"gopkg.in/yaml.v2"
)

// Declare the config file name
const configFileName = "pipedream-config.yaml"

// Declare the Config struct type
type Step struct {
	Namespace string `yaml:"namespace"`
}

type Workflow struct {
	Name  string `yaml:"name"`
	Steps []Step `yaml:"steps"`
}

type Config struct {
	Schema    string            `yaml:"schema"`
	Workflows map[string]string `yaml:"workflows"`
	Steps     map[string]string `yaml:"steps"`
}

// Check if the project is initialized (i.e., config file exists)
func isInitialized() bool {
	_, err := os.Stat(configFileName)
	return err == nil
}

// Load configuration from the config file
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

	if config.Schema != "baudevs/2024-09-29" {
		color.Yellow("Warning: This project is using an unsupported schema version: %s", config.Schema)
	}

	return config, nil
}

// Save the configuration to the config file
func saveConfig(config Config) {
	data, err := yaml.Marshal(&config)
	if err != nil {
		color.Red("Error saving configuration: %v", err)
		return
	}
	err = os.WriteFile(configFileName, data, 0644)
	if err != nil {
		color.Red("Error writing configuration file: %v", err)
	}
}

// Initialize the Pipedream project by detecting workflows and steps
func initializeProject() {
	if _, err := os.Stat(configFileName); err == nil {
		color.Yellow("This project is already initialized.")
		return
	}

	config := Config{
		Schema:    "baudevs/2024-09-29",
		Workflows: make(map[string]string),
		Steps:     make(map[string]string),
	}

	err := filepath.Walk(".", func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return err
		}

		if info.Name() == "workflow.yaml" {
			dir := filepath.Dir(path)
			workflowName := filepath.Base(dir)
			promptWorkflow := prompt(fmt.Sprintf("Detected workflow in %s. Enter workflow name (default: %s): ", dir, workflowName), workflowName)
			config.Workflows[workflowName] = promptWorkflow

			subdirs, err := os.ReadDir(dir)
			if err != nil {
				return err
			}
			for _, subdir := range subdirs {
				if subdir.IsDir() {
					stepName := subdir.Name()
					fullPath := filepath.Join(dir, stepName)
					promptStep := prompt(fmt.Sprintf("Detected step %s in %s. Enter step namespace (default: %s): ", stepName, fullPath, stepName), stepName)
					config.Steps[filepath.Join(workflowName, stepName)] = promptStep
				}
			}
		}
		return nil
	})

	if err != nil {
		color.Red("Error scanning the project directory: %v", err)
		return
	}

	saveConfig(config)
	color.Green("Pipedream project initialized successfully with schema version baudevs/2024-09-29.")
}
