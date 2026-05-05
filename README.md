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
```

> Additional commands (`save`, `list`, `apply`, `detect`, etc.) are planned but not yet implemented.

## Supported Agents

| Agent | Platform | Status |
|-------|----------|--------|
| OpenCode | Windows | Planned |
| Kiro | macOS/Linux | Planned |

## License

MIT - see [LICENSE](LICENSE) for details.
