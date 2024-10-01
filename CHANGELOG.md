# Changelog

All notable changes to this project will be documented in this file.

The format is based on [Keep a Changelog](https://keepachangelog.com/en/1.0.0/),
and this project adheres to [Semantic Versioning](https://semver.org/spec/v2.0.0.html).

## [Unreleased]

### Added

- New `pdcommit sync` command to synchronize the configuration with the current project structure
- Functionality to add new workflows and steps during sync
- Option to remove non-existent workflows and steps from config during sync

### Changed

- Moved original bash script to its own folder
- Updated help text to include information about the new `sync` command

### Fixed

- Improved error handling and user feedback during configuration operations

## [0.2.0] - 2024-10-01

### Added

- Initial Go implementation of the git commit helper
- Project initialization with `pdcommit init` command
- Configuration management using `pipedream-config.yaml`
- Workflow and step detection during initialization

### Changed

- Migrated core functionality from bash script to Go

## [0.1.0] - 2024-09-28

### Added

- Initial commit: Added git commit helper script in bash

[Unreleased]: https://github.com/yourusername/yourrepository/compare/v0.2.0...HEAD
[0.2.0]: https://github.com/yourusername/yourrepository/compare/v0.1.0...v0.2.0
[0.1.0]: https://github.com/yourusername/yourrepository/releases/tag/v0.1.0
