package main

import (
	"fmt"
	"os"
	"path/filepath"

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
	Schema    string            `yaml:"schema"`
	Workflows map[string]string `yaml:"workflows"`
	Steps     map[string]string `yaml:"steps"`
}

const configFileName = "pipedream-config.yaml"

func isInitialized() bool {
	_, err := os.Stat(configFileName)
	return err == nil
}

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

func initializeProject() {
	if _, err := os.Stat(configFileName); err == nil {
		fmt.Println(ColorYellow, "This project is already initialized.", ColorReset)
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
		fmt.Printf("Error scanning the project directory: %v\n", err)
		return
	}

	saveConfig(config)
	fmt.Println(ColorGreen, "Pipedream project initialized successfully with schema version baudevs/2024-09-29.", ColorReset)
}
