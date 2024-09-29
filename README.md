# Pipedream Git Commit Helper

Pipedream Git Commit Helper is a CLI tool designed to simplify the process of committing changes in projects that use [Pipedream](https://pipedream.com) workflows (BauDevs is not in any way affiliated with Pipedream. We did this tool to make our lives easier). 

This tool detects Git changes, helps you interactively stage files, and proposes commit messages based on detected workflows and steps. It's highly interactive, color-coded, and supports custom configurations for your Pipedream project.

## Features

- **Automatic Workflow and Step Detection**: The tool automatically detects Pipedream workflows and steps from the project directory.
- **Interactive Prompts**: Choose commit types, confirm staging, and finalize your commit with easy-to-use interactive prompts.
- **Color-Coded Output**: The tool uses color to differentiate important information, making the CLI experience more user-friendly.
- **Multiline Commit Messages**: You can enter detailed commit messages with multiple lines.
- **Configuration**: The tool allows for custom workflow and step mappings through a `pipedream-config.yaml` file.

## Installation

To install the tool, follow these steps:

1. Clone the repository:

   ```bash
   git clone https://github.com/bau-devs/pipedream-git-commit-helper.git
   ```

2. Navigate to the project directory:

   ```bash
   cd pipedream-git-commit-helper
   ```

3. Install the tool (you need to have git, go, installed and in your path, for instructions on how to install go, see [here](https://go.dev/doc/install) and for git, see [here](https://git-scm.com/downloads)):

> Warning: If the install script fails, make sure you have the correct permissions to write to your home directory, or try doing

```bash
   chmod +x install.sh
    ./install.sh
```

- On macOS, run the install script:

```bash
 ./install.sh
```

- On Linux, run the install script:

```bash
    ./install.sh
```

- On Windows, run these commands (make sure to have git, git-bash, and go installed):

```bash
    go build -o pdcommit.exe
    mv pdcommit.exe $env:Path
```

4. Reload your shell:

```bash
   source ~/.bashrc
```

*or*

```bash
  source ~/.zshrc
```

5. Test the tool:

```bash
   pdcommit --help
```

## Usage

Here’s how you can use the Pipedream Git Commit Helper:

***Initialize a Pipedream Project***

If your project has not been initialized, you must run the init command to detect workflows and steps:

```bash
pdcommit init
```

This will scan your project for workflows (workflow.yaml files) and steps (directories containing the code for each step). The detected workflows and steps will be stored in a pipedream-config.yaml file.

It will by default choose the name in the workflow yaml and the namespace in the steps property of the workflow yaml as names to use for when committing. You can change as it asks you one by one and the mapping is stored in a pipedream-config.yaml file. You can also at any point you want change the mapping by editing the pipedream-config.yaml file manually.

***Committing Changes***

1. Detect Git Changes: The tool automatically detects changes in your Git repository.
2. Stage Changes: If there are unstaged changes, it prompts you to stage all changes:

```bash
pdcommit
```

***If uncommitted changes are detected, it will prompt you to stage them:***

```bash
Warning: You have unstaged changes.
Do you want to stage all changes? (y/N) y      
```

***If you choose to stage all changes, it will stage all changes and then proceed with the commit process.***

```bash
Staging all changes...
```

3.	Propose Workflow and Step: Based on the Git changes, the tool will propose a matching workflow and step:

```bash
Checking file: services/workflow_name/step_namespace/entry.js
Matched workflow: workflow_name -> step_namespace
```

4.	Choose Commit Type: It will ask you to choose a commit type (feat, fix, chore, etc.) and enter a detailed commit message.

```bash
Select Commit Type: 
- add
- fix
- change
- remove
```

5.	Enter Commit Message: Enter a commit message. You can also add multiline messages by pressing Enter twice to finish

```bash
Enter Commit Message:
```

6.	Review Commit: The tool displays the full commit command for review, and asks for confirmation before executing the commit: 

```bash
The following command is ready for review:
git commit -m "add[[workflow_name][step_namespace]] New feature added"
Do you want to proceed with the commit? (y/N)
```

now you can do 

```bash
git fetch && git merge
git push
```

Go to pipedream and pull the branch and you should see be able to see the changes in the workflow and test it.

## Help & Manual

For a list of all the commands and how to use them, run:

```bash
pdcommit --help
```

To view the manual, run:

```bash
pdcommit man
```

or

```bash
man pdcommit
```

## Configuration

The tool uses a `pipedream-config.yaml` file to store the mappings between workflows and steps. The file is stored in the user's home directory by default.

```yaml
schema: baudevs/2024-09-29
workflows:
  workflow_1_name-p_ABC123: workflow_1_name
  workflow_2_name-p_XYZ456: workflow_2_name
steps:
  workflow_1_name-p_ABC123/step_1_namespace: step_1_namespace
  workflow_1_name-p_ABC123/step_2_namespace: step_2_namespace
  workflow_2_name-p_XYZ456/step_3_namespace: step_3_namespace
  workflow_2_name-p_XYZ456/step_4_namespace : step_4_namespace
```

You can edit the file to change the mappings between workflows and steps.

## Contributing

We welcome contributions to this project! Feel free to open issues or submit pull requests on the GitHub repository.

To contribute:

1. Fork the repository.
2. Create a new branch for your feature or bugfix.
3. Submit a pull request when ready.

## Licensing

This project is licensed under a **dual license model**:

1. **MIT License**: The default license for open-source usage. This allows free use of the software for personal or internal projects under the [MIT License](LICENSE).

2. **Commercial License**: For commercial usage, such as deployment in production environments, enterprise applications, or providing services to customers, a commercial license is required. See the [Commercial License](LICENSE_COMMERCIAL) for more details.

   - The commercial license includes support, warranties, and custom agreements tailored for business needs.

### License Summary

- If you are using this tool for **non-commercial, open-source, or personal use**, you are free to use it under the terms of the **MIT License**.
- If you intend to use this tool for **commercial purposes** or require **enterprise support**, please contact us for a commercial license.

## Installation

Follow the standard installation steps outlined above.

## Support

- **Open-source users**: Feel free to submit issues or contribute via pull requests.
- **Commercial users**: You will receive direct support as part of the commercial license agreement.

For more information about licensing or support, please contact [licenses++at++baudevs.com].

## Credits

This tool was developed by BauDevs to streamline the Git commit process for Pipedream-based projects.

*Visit our Github page at [BauDevs On Github](https://github.com/BauDevs)*

*Visit our website at [BauDevs Website](https://baudevs.com)*

