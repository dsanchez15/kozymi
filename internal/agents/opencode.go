package agents

import (
	"os"
	"path/filepath"
	"runtime"

	"github.com/komyzi/komyzi/pkg/models"
)

// OpenCodeDetector implementa la detección de OpenCode
type OpenCodeDetector struct{}

// Name retorna el nombre del agente
func (d *OpenCodeDetector) Name() string {
	return "opencode"
}

// IsInstalled verifica si OpenCode está instalado
func (d *OpenCodeDetector) IsInstalled() bool {
	// Verificar si existe el directorio global de config
	globalDir := d.getGlobalConfigDir()
	if _, err := os.Stat(globalDir); err == nil {
		return true
	}
	// Verificar si existe el binario (en PATH)
	_, err := os.Stat(filepath.Join(os.Getenv("LOCALAPPDATA"), "opencode", "opencode.exe"))
	return err == nil
}

// DetectProject busca configuración de OpenCode en un directorio
func (d *OpenCodeDetector) DetectProject(projectPath string) (*models.AgentPaths, error) {
	configPath := filepath.Join(projectPath, "opencode.json")
	configPathC := filepath.Join(projectPath, "opencode.jsonc")
	rulesPath := filepath.Join(projectPath, "AGENTS.md")
	tuiPath := filepath.Join(projectPath, "tui.json")
	tuiPathC := filepath.Join(projectPath, "tui.jsonc")
	opencodeDir := filepath.Join(projectPath, ".opencode")

	// Verificar si existe alguno de los archivos principales o el directorio .opencode/
	hasConfig := fileExists(configPath) || fileExists(configPathC)
	hasRules := fileExists(rulesPath)
	hasTUI := fileExists(tuiPath) || fileExists(tuiPathC)
	hasDir := dirExists(opencodeDir)

	if !hasConfig && !hasRules && !hasTUI && !hasDir {
		return nil, nil // No se detectó OpenCode en este proyecto
	}

	// Usar la ruta de config que exista
	actualConfigPath := ""
	if fileExists(configPath) {
		actualConfigPath = configPath
	} else if fileExists(configPathC) {
		actualConfigPath = configPathC
	}

	// Usar la ruta de TUI que exista
	actualTUIPath := ""
	if fileExists(tuiPath) {
		actualTUIPath = tuiPath
	} else if fileExists(tuiPathC) {
		actualTUIPath = tuiPathC
	}

	return &models.AgentPaths{
		AgentType:   models.AgentOpenCode,
		ConfigPath:  actualConfigPath,
		RulesPath:   rulesPath,
		TUIPath:     actualTUIPath,
		OpenCodeDir: opencodeDir,
		IsGlobal:    false,
	}, nil
}

// DetectGlobal busca configuración global de OpenCode
func (d *OpenCodeDetector) DetectGlobal() (*models.AgentPaths, error) {
	globalDir := d.getGlobalConfigDir()

	configPath := filepath.Join(globalDir, "opencode.json")
	configPathC := filepath.Join(globalDir, "opencode.jsonc")
	rulesPath := filepath.Join(globalDir, "AGENTS.md")
	tuiPath := filepath.Join(globalDir, "tui.json")
	tuiPathC := filepath.Join(globalDir, "tui.jsonc")

	// Verificar si existe el directorio global
	if !dirExists(globalDir) {
		return nil, nil
	}

	// Usar la ruta de config que exista
	actualConfigPath := ""
	if fileExists(configPath) {
		actualConfigPath = configPath
	} else if fileExists(configPathC) {
		actualConfigPath = configPathC
	}

	// Usar la ruta de TUI que exista
	actualTUIPath := ""
	if fileExists(tuiPath) {
		actualTUIPath = tuiPath
	} else if fileExists(tuiPathC) {
		actualTUIPath = tuiPathC
	}

	return &models.AgentPaths{
		AgentType:   models.AgentOpenCode,
		ConfigPath:  actualConfigPath,
		RulesPath:   rulesPath,
		TUIPath:     actualTUIPath,
		OpenCodeDir: globalDir,
		IsGlobal:    true,
	}, nil
}

// getGlobalConfigDir retorna el directorio de configuración global de OpenCode
func (d *OpenCodeDetector) getGlobalConfigDir() string {
	// Primero intentar XDG_CONFIG_HOME
	if xdg := os.Getenv("XDG_CONFIG_HOME"); xdg != "" {
		return filepath.Join(xdg, "opencode")
	}

	// En Windows usar USERPROFILE\.config\opencode
	if runtime.GOOS == "windows" {
		if home := os.Getenv("USERPROFILE"); home != "" {
			return filepath.Join(home, ".config", "opencode")
		}
	}

	// Fallback a HOME/.config/opencode
	if home := os.Getenv("HOME"); home != "" {
		return filepath.Join(home, ".config", "opencode")
	}

	return ""
}

// fileExists verifica si un archivo existe
func fileExists(path string) bool {
	if path == "" {
		return false
	}
	info, err := os.Stat(path)
	if err != nil {
		return false
	}
	return !info.IsDir()
}

// dirExists verifica si un directorio existe
func dirExists(path string) bool {
	if path == "" {
		return false
	}
	info, err := os.Stat(path)
	if err != nil {
		return false
	}
	return info.IsDir()
}

// Registry contiene todos los detectores disponibles
var Registry = []Detector{
	&OpenCodeDetector{},
}
