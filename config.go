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

func syncConfig() {
	if !isInitialized() {
		color.Red("Error: This is not an initialized Pipedream project.")
		fmt.Println("Please run 'pdcommit init' to initialize the project.")
		return
	}

	config, err := loadConfig()
	if err != nil {
		color.Red("Error loading configuration: %v", err)
		return
	}

	newConfig := Config{
		Schema:    config.Schema,
		Workflows: make(map[string]string),
		Steps:     make(map[string]string),
	}

	err = filepath.Walk(".", func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return err
		}

		if info.Name() == "workflow.yaml" {
			dir := filepath.Dir(path)
			workflowName := filepath.Base(dir)
			if existingName, ok := config.Workflows[workflowName]; ok {
				newConfig.Workflows[workflowName] = existingName
			} else {
				promptWorkflow := prompt(fmt.Sprintf("New workflow detected in %s. Enter workflow name (default: %s): ", dir, workflowName), workflowName)
				newConfig.Workflows[workflowName] = promptWorkflow
			}

			subdirs, err := os.ReadDir(dir)
			if err != nil {
				return err
			}
			for _, subdir := range subdirs {
				if subdir.IsDir() {
					stepName := subdir.Name()
					fullPath := filepath.Join(workflowName, stepName)
					if existingNamespace, ok := config.Steps[fullPath]; ok {
						newConfig.Steps[fullPath] = existingNamespace
					} else {
						promptStep := prompt(fmt.Sprintf("New step %s detected in %s. Enter step namespace (default: %s): ", stepName, fullPath, stepName), stepName)
						newConfig.Steps[fullPath] = promptStep
					}
				}
			}
		}
		return nil
	})

	if err != nil {
		color.Red("Error scanning the project directory: %v", err)
		return
	}

	// Check for removed workflows and steps
	for workflowName, workflowValue := range config.Workflows {
		if _, exists := newConfig.Workflows[workflowName]; !exists {
			if confirmAction(fmt.Sprintf("Workflow %s no longer exists. Remove it from config?", workflowName)) {
				color.Yellow("Removed workflow: %s", workflowName)
			} else {
				newConfig.Workflows[workflowName] = workflowValue
			}
		}
	}

	for stepPath, stepValue := range config.Steps {
		if _, exists := newConfig.Steps[stepPath]; !exists {
			if confirmAction(fmt.Sprintf("Step %s no longer exists. Remove it from config?", stepPath)) {
				color.Yellow("Removed step: %s", stepPath)
			} else {
				newConfig.Steps[stepPath] = stepValue
			}
		}
	}

	saveConfig(newConfig)
	color.Green("Configuration synced successfully.")
}
