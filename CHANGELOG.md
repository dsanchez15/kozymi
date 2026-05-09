# Changelog

All notable changes to this project will be documented in this file.

The format is based on [Keep a Changelog](https://keepachangelog.com/en/1.0.0/),
and this project adheres to [Semantic Versioning](https://semver.org/spec/v2.0.0.html).

## [Unreleased]

## [0.2.0] - 2026-05-05

### Added
- OpenCode configuration detector for Windows (project and global levels)
- `komyzi save` command with support for:
  - Saving current project configuration (auto-detects folder name)
  - Saving with custom name (`-n` flag)
  - Saving global configuration (`--global` flag)
  - Saving from specific path (`--from` flag)
  - Force agent type (`--agent` flag)
- `komyzi apply` command with support for:
  - Applying saved configuration to current directory
  - Applying to specific path (`--to` flag)
  - Agent type selection (`--agent` flag)
- Storage architecture with portable/non-portable separation:
  - `~/.komyzi/agents/<agent>/configs/<name>/` for non-portables
  - `~/.komyzi/agents/<agent>/shared/` for portables (skills, agents, themes, commands)
- Metadata tracking (`agent.json`) with source, date, and shared references
- AGENTS.md with absolute rules for AI agent behavior (gitignored)
- Command documentation synchronization skill
- PR template for GitHub

## [0.1.0] - 2026-05-04

### Added
- Initial project setup
- Go module with Bubble Tea for TUI
- Basic CLI entry point
- CI/CD workflows for GitHub Actions
- README, LICENSE, CONTRIBUTING files
- GitHub Actions for CI and Release automation