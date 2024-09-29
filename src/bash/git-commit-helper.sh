#!/bin/bash

# File to store the last used commit information
LAST_COMMIT_FILE=".last_commit_info"

# Colors for the prompts
RED='\033[0;31m'
GREEN='\033[0;32m'
YELLOW='\033[1;33m'
BLUE='\033[1;34m'
NC='\033[0m' # No color

# Function to extract workflows and steps from workflow.yaml files
function get_workflow_and_steps() {
    workflow_step_combinations=()

    # Find all workflow.yaml files
    yaml_files=$(find . -name "workflow.yaml")
    
    # Iterate over each workflow.yaml file
    for yaml_file in $yaml_files; do
        # Extract workflow name (top-level name property)
        workflow_name=$(yq e '.name' "$yaml_file")
        
        # Extract step namespaces (steps under 'steps' property)
        step_names=$(yq e '.steps[].namespace' "$yaml_file")

        # For each step, create a workflow:step combination
        for step_name in $step_names; do
            workflow_step_combinations+=("$workflow_name:$step_name")
        done
    done
}

# Function to load the last commit information
function load_last_commit() {
    if [[ -f $LAST_COMMIT_FILE ]]; then
        source $LAST_COMMIT_FILE
        echo "Last commit: $last_type[$last_workflow][$last_step]"
    else
        echo "No previous commit found."
        last_type=""
        last_workflow=""
        last_step=""
    fi
}

# Function to save the last commit information
function save_last_commit() {
    echo "last_type=$type" > $LAST_COMMIT_FILE
    echo "last_workflow=$workflow" >> $LAST_COMMIT_FILE
    echo "last_step=$step" >> $LAST_COMMIT_FILE
}

# Function to propose workflow and step based on git changes
function propose_workflow_from_git() {
    echo -e "${YELLOW}Checking git changes to infer workflow and step...${NC}"

    # Get the list of modified files from git status
    git_changes=$(git status --porcelain | awk '{print $2}')

    if [[ -z "$git_changes" ]]; then
        echo -e "${RED}No modified files found in git.${NC}"
        return 1
    fi

    # Variable to track proposed workflows and steps
    workflow_steps=()

    # Loop over modified files and check if they are related to any workflow.yaml
    for file in $git_changes; do
        echo -e "${BLUE}Checking file: $file${NC}"
        
        # Find the directory that contains the modified file
        dir=$(dirname "$file")
        echo "Starting directory: $dir"

        # Traverse parent directories to find workflow.yaml
        while [[ "$dir" != "/" && "$dir" != "." ]]; do
            echo "Looking for workflow.yaml in: $dir"
            workflow_file=$(find "$dir" -maxdepth 1 -name "workflow.yaml" -print -quit)

            if [[ -f "$workflow_file" ]]; then
                echo -e "${GREEN}Found workflow.yaml in $dir${NC}"

                # Extract workflow name
                workflow=$(yq e '.name' "$workflow_file")

                # Extract the folder name where the modified file is located
                step_candidate=$(basename "$(dirname "$file")")
                echo "Step candidate: $step_candidate"

                # Extract step namespaces from the workflow.yaml
                step_names=$(yq e '.steps[].namespace' "$workflow_file")

                # Check if the folder name matches any of the steps in the workflow.yaml
                for step_name in $step_names; do
                    if [[ "$step_candidate" == "$step_name" ]]; then
                        step=$step_name
                        echo -e "${YELLOW}Proposed: ${BLUE}[${workflow}]${NC} -> ${GREEN}[${step}]${NC}"
                        
                        # Append to the list of proposed workflow and steps
                        workflow_steps+=("$workflow:$step")
                    fi
                done
            fi

            # Move up one directory
            dir=$(dirname "$dir")
        done
    done

    if [[ ${#workflow_steps[@]} -eq 0 ]]; then
        echo -e "${RED}No matching workflow:step found from git changes.${NC}"
        return 1
    else
        echo "Proposed workflow:step combinations:"
        for wf in "${workflow_steps[@]}"; do
            echo -e "${GREEN}* $wf${NC}"
        done
        return 0
    fi
}

# Function to build the final commit message
function build_commit_message() {
    message_type=$1
    workflow_steps=("${@:2:$#-2}")
    message=${@: -1}

    # Build the workflow:step format
    wf_message=""

    for wf_step in "${workflow_steps[@]}"; do
        # Split workflow and step
        wf=$(echo "$wf_step" | cut -d':' -f1)
        step=$(echo "$wf_step" | cut -d':' -f2)

        # Append the workflow and step directly
        wf_message+="[[$wf][$step]]"
    done

    # Return the final commit message
    echo "$message_type$wf_message $message"
}

# Load the last commit if available
load_last_commit

# Ask if the user wants to run git add (default is 'y')
read -p "Do you want to run 'git add .' to stage all changes? (y/n) [y]: " run_git_add
run_git_add=${run_git_add:-y}  # Default to 'y' if no input

if [[ $run_git_add == "y" ]]; then
    git add .
    echo "All changes have been staged."
fi

# Check git for possible workflow and step suggestions
propose_workflow_from_git
if [[ $? -eq 0 ]]; then
    # Default is 'y' for using proposed workflow and step
    read -p "Use the proposed workflow and step from git changes? (y/n) [y]: " use_git_suggestion
    use_git_suggestion=${use_git_suggestion:-y}  # Default to 'y' if no input

    if [[ $use_git_suggestion == "y" ]]; then
        echo "Using proposed workflow:step from git changes."
        git_proposed=1
    else
        git_proposed=0
    fi
else
    git_proposed=0
fi

# If no git suggestion, ask if the user wants to use the last commit info (default is 'y')
if [[ $git_proposed -eq 0 && -n $last_type && -n $last_workflow && -n $last_step ]]; then
    read -p "No matching workflow found from git. Use the same commit type, workflow, and step as the last commit? (y/n) [y]: " reuse_last
    reuse_last=${reuse_last:-y}  # Default to 'y' if no input
else
    reuse_last="n"
fi

if [[ $git_proposed -eq 1 ]]; then
    # Workflow and step were proposed by git changes
    type="fix"  # Default type in case git proposes the workflow
else
    if [[ $reuse_last == "y" ]]; then
        type=$last_type
        workflow=$last_workflow
        step=$last_step
        echo "Reusing: $type[$workflow][$step]"
    else
        # Prompt user for the commit type using fzf
        type=$(printf "add\nfix\nchange\nremove" | fzf --prompt "Commit type: ")

        if [[ -z $type ]]; then
            echo "No commit type selected. Exiting."
            exit 1
        fi

        echo "Selected commit type: $type"

        # Fetch workflows and steps from workflow.yaml files
        get_workflow_and_steps

        # Prompt user to select workflow:step combination using fzf
        workflow_step=$(printf "%s\n" "${workflow_step_combinations[@]}" | fzf --prompt "Workflow:Step: ")

        if [[ -z $workflow_step ]]; then
            echo "No workflow:step selected. Exiting."
            exit 1
        fi

        echo "Selected workflow:step combination: $workflow_step"

        # Split workflow:step into individual variables
        workflow=$(echo $workflow_step | cut -d':' -f1)
        step=$(echo $workflow_step | cut -d':' -f2)
    fi
fi

# Prompt user for commit message
read -p "Enter the commit message: " message

# Build the final commit message
final_commit_message=$(build_commit_message "$type" "${workflow_steps[@]}" "$message")

# Show the command to the user and prompt for confirmation or modification
echo -e "${YELLOW}The following command is ready for review:${NC}"
echo "git commit -m \"$final_commit_message\""

# Prompt user to manually copy or edit the command and execute
read -p "Press Enter to run the command, or Ctrl+C to cancel. You can modify the command before executing if needed."

# Execute the command
git commit -m "$final_commit_message"

# Save the last commit info for next time
save_last_commit