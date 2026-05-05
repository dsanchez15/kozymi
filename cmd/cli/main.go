package main

import (
	"flag"
	"fmt"
	"os"
	"path/filepath"

	"github.com/komyzi/komyzi/internal/agents"
	"github.com/komyzi/komyzi/internal/storage"
	"github.com/komyzi/komyzi/pkg/models"
)

var version = "0.1.0"

func main() {
	if len(os.Args) < 2 {
		printUsage()
		os.Exit(1)
	}

	command := os.Args[1]

	// Handle --version flag
	if command == "--version" || command == "-v" {
		fmt.Println("komyzi version", version)
		return
	}

	switch command {
	case "save":
		handleSave(os.Args[2:])
	case "version":
		fmt.Println("komyzi version", version)
	default:
		fmt.Printf("Unknown command: %s\n", command)
		printUsage()
		os.Exit(1)
	}
}

func printUsage() {
	fmt.Println("Komyzi - AI Agent Configuration Manager")
	fmt.Println()
	fmt.Println("Usage:")
	fmt.Println("  komyzi save                           Save current project config (uses folder name)")
	fmt.Println("  komyzi save -n <name>                 Save with custom name")
	fmt.Println("  komyzi save --global                  Save global config")
	fmt.Println("  komyzi save -n <name> --global        Save global config with custom name")
	fmt.Println("  komyzi save --from <path>             Save config from specific path")
	fmt.Println("  komyzi save -n <name> --from <path>   Save from path with custom name")
	fmt.Println("  komyzi save --agent <type>            Force specific agent type")
	fmt.Println("  komyzi version                        Show version")
	fmt.Println("  komyzi --version                      Show version")
}

func handleSave(args []string) {
	// Parse flags
	fs := flag.NewFlagSet("save", flag.ExitOnError)
	nameFlag := fs.String("n", "", "Configuration name (optional, defaults to project folder name)")
	globalFlag := fs.Bool("global", false, "Save global configuration")
	fromFlag := fs.String("from", "", "Path to project directory")
	agentFlag := fs.String("agent", "", "Force specific agent type (opencode)")
	
	fs.Parse(args)

	configName := *nameFlag

	// For now, only support opencode
	if *agentFlag != "" && *agentFlag != "opencode" {
		fmt.Printf("Error: unsupported agent type: %s\n", *agentFlag)
		fmt.Println("Currently supported: opencode")
		os.Exit(1)
	}

	detector := &agents.OpenCodeDetector{}

	// Detect configuration
	var paths *models.AgentPaths
	var err error
	var sourceType string

	if *globalFlag {
		// Save global configuration
		fmt.Println("Detecting global OpenCode configuration...")
		paths, err = detector.DetectGlobal()
		if err != nil {
			fmt.Printf("Error detecting global config: %v\n", err)
			os.Exit(1)
		}
		if paths == nil {
			fmt.Println("No global OpenCode configuration found.")
			fmt.Printf("Expected location: %s\n", getGlobalConfigDirHint())
			os.Exit(1)
		}
		sourceType = "global"
	} else if *fromFlag != "" {
		// Save from specific path
		absPath, err := filepath.Abs(*fromFlag)
		if err != nil {
			fmt.Printf("Error resolving path: %v\n", err)
			os.Exit(1)
		}
		fmt.Printf("Detecting OpenCode configuration in %s...\n", absPath)
		paths, err = detector.DetectProject(absPath)
		if err != nil {
			fmt.Printf("Error detecting project config: %v\n", err)
			os.Exit(1)
		}
		if paths == nil {
			fmt.Printf("No OpenCode configuration found in %s\n", absPath)
			os.Exit(1)
		}
		sourceType = "project"
	} else {
		// Save current directory
		currentDir, err := os.Getwd()
		if err != nil {
			fmt.Printf("Error getting current directory: %v\n", err)
			os.Exit(1)
		}
		fmt.Printf("Detecting OpenCode configuration in current directory...\n")
		paths, err = detector.DetectProject(currentDir)
		if err != nil {
			fmt.Printf("Error detecting project config: %v\n", err)
			os.Exit(1)
		}
		if paths == nil {
			fmt.Printf("No OpenCode configuration found in current directory (%s)\n", currentDir)
			fmt.Println("Use --global to save global configuration or --from <path> to specify a project.")
			os.Exit(1)
		}
		sourceType = "project"
	}

	// Determine config name if not provided
	if configName == "" {
		if *globalFlag {
			configName = "global"
		} else {
			// Use the folder name of the project path
			configName = filepath.Base(paths.OpenCodeDir)
			if configName == "." || configName == "/" || configName == "\\" {
				configName = "default"
			}
			// Remove .opencode if it's the directory name
			if configName == ".opencode" {
				configName = filepath.Base(filepath.Dir(paths.OpenCodeDir))
			}
		}
	}

	// Save configuration
	repo, err := storage.NewRepository()
	if err != nil {
		fmt.Printf("Error initializing storage: %v\n", err)
		os.Exit(1)
	}

	if err := repo.SaveConfig(paths, configName); err != nil {
		fmt.Printf("Error saving configuration: %v\n", err)
		os.Exit(1)
	}

	fmt.Printf("✓ Configuration '%s' saved successfully from %s\n", configName, sourceType)
}

// getGlobalConfigDirHint retorna la ruta esperada para la config global
func getGlobalConfigDirHint() string {
	home, err := os.UserHomeDir()
	if err != nil {
		return "%USERPROFILE%\\.config\\opencode (Windows) or ~/.config/opencode"
	}
	return filepath.Join(home, ".config", "opencode")
}
