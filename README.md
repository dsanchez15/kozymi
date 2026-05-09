# Komyzi

AI Agent Configuration Manager - Save, copy and manage configurations for your AI agents across projects.

## Current Status

This project is in active development. The CLI is currently a work in progress.

## Features (Planned)

- **Save Configurations**: Store configurations from any supported AI agent
- **Copy Configurations**: Apply saved configurations to new projects
- **Multi-Agent Support**: OpenCode, Kiro, Claude Code, and more
- **MCP Support**: Manage MCP configurations and skills
- **Cross-Platform**: Windows, macOS, Linux

## Installation

> Not yet available. Build from source to try the current version.

### Build from Source

```bash
git clone https://github.com/komyzi/komyzi.git
cd komyzi
go build -o komyzi.exe ./cmd/cli
```

## Usage

### Current Commands

```bash
# Show version
komyzi.exe --version
komyzi.exe version

# Save configuration
komyzi.exe save                           # Save current project config (uses folder name)
komyzi.exe save -n <name>                 # Save with custom name
komyzi.exe save --global                  # Save global config
komyzi.exe save -n <name> --global        # Save global config with custom name
komyzi.exe save --from ./path             # Save config from specific path
komyzi.exe save -n <name> --from ./path   # Save from path with custom name
komyzi.exe save --agent opencode          # Force specific agent

# Apply configuration
komyzi.exe apply <name>                   # Apply config to current directory
komyzi.exe apply <name> --to ./path       # Apply config to specific path
komyzi.exe apply <name> --agent opencode  # Apply config from specific agent
```

> Additional commands (`list`, `detect`, etc.) are planned but not yet implemented.

## Supported Agents

| Agent | Platform | Status |
|-------|----------|--------|
| OpenCode | Windows | Planned |
| Kiro | macOS/Linux | Planned |

## License

MIT - see [LICENSE](LICENSE) for details.
