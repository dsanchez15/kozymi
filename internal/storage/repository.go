package storage

import (
	"encoding/json"
	"fmt"
	"io"
	"os"
	"path/filepath"
	"time"

	"github.com/komyzi/komyzi/pkg/models"
)

// Repository maneja el almacenamiento de configuraciones
type Repository struct {
	baseDir string
}

// NewRepository crea un nuevo repositorio de configuraciones
func NewRepository() (*Repository, error) {
	homeDir, err := os.UserHomeDir()
	if err != nil {
		return nil, fmt.Errorf("error getting home directory: %w", err)
	}

	baseDir := filepath.Join(homeDir, ".komyzi")
	return &Repository{baseDir: baseDir}, nil
}

// SaveConfig guarda una configuración de agente
func (r *Repository) SaveConfig(agentPaths *models.AgentPaths, name string) error {
	// Determinar directorio base según si es global o proyecto
	configsDir := filepath.Join(r.baseDir, "agents", string(agentPaths.AgentType), "configs")
	if agentPaths.IsGlobal {
		configsDir = filepath.Join(configsDir, "global")
	} else {
		configsDir = filepath.Join(configsDir, name)
	}

	// Crear directorios necesarios
	if err := os.MkdirAll(configsDir, 0755); err != nil {
		return fmt.Errorf("error creating config directory: %w", err)
	}

	// Crear directorio shared si no existe
	sharedDir := filepath.Join(r.baseDir, "agents", string(agentPaths.AgentType), "shared")
	if err := os.MkdirAll(sharedDir, 0755); err != nil {
		return fmt.Errorf("error creating shared directory: %w", err)
	}

	// Guardar archivos no-portables
	var sharedRefs []string

	// Config principal (opencode.json)
	if agentPaths.ConfigPath != "" {
		src, err := os.Open(agentPaths.ConfigPath)
		if err == nil {
			defer func() { _ = src.Close() }()
			dst := filepath.Join(configsDir, "config.json")
			if err := copyFile(src, dst); err != nil {
				return fmt.Errorf("error copying config: %w", err)
			}
		}
	}

	// Reglas (AGENTS.md)
	if agentPaths.RulesPath != "" && fileExists(agentPaths.RulesPath) {
		src, err := os.Open(agentPaths.RulesPath)
		if err == nil {
			defer func() { _ = src.Close() }()
			dst := filepath.Join(configsDir, "rules.md")
			if err := copyFile(src, dst); err != nil {
				return fmt.Errorf("error copying rules: %w", err)
			}
		}
	}

	// TUI config (tui.json)
	if agentPaths.TUIPath != "" {
		src, err := os.Open(agentPaths.TUIPath)
		if err == nil {
			defer func() { _ = src.Close() }()
			dst := filepath.Join(configsDir, "tui.json")
			if err := copyFile(src, dst); err != nil {
				return fmt.Errorf("error copying tui config: %w", err)
			}
		}
	}

	// Guardar elementos portables desde .opencode/
	if agentPaths.OpenCodeDir != "" && dirExists(agentPaths.OpenCodeDir) {
		refs, err := r.savePortables(agentPaths)
		if err != nil {
			return fmt.Errorf("error saving portables: %w", err)
		}
		sharedRefs = append(sharedRefs, refs...)
	}

	// Crear metadata
	config := models.AgentConfig{
		Agent:      agentPaths.AgentType,
		Name:       name,
		Source:     models.SourceProject,
		SourcePath: agentPaths.OpenCodeDir,
		DateSaved:  time.Now(),
		SharedRefs: sharedRefs,
	}

	if agentPaths.IsGlobal {
		config.Source = models.SourceGlobal
		config.SourcePath = agentPaths.OpenCodeDir
	}

	metadataPath := filepath.Join(configsDir, "agent.json")
	metadataJSON, err := json.MarshalIndent(config, "", "  ")
	if err != nil {
		return fmt.Errorf("error marshaling metadata: %w", err)
	}

	if err := os.WriteFile(metadataPath, metadataJSON, 0644); err != nil {
		return fmt.Errorf("error writing metadata: %w", err)
	}

	return nil
}

// savePortables guarda los elementos portables y retorna las referencias
func (r *Repository) savePortables(paths *models.AgentPaths) ([]string, error) {
	var refs []string

	// Directorios portables a copiar
	portableDirs := []string{"skills", "agents", "themes", "commands"}

	for _, dirName := range portableDirs {
		srcDir := filepath.Join(paths.OpenCodeDir, dirName)
		if !dirExists(srcDir) {
			continue
		}

		dstDir := filepath.Join(r.baseDir, "agents", string(paths.AgentType), "shared", dirName)
		
		// Copiar todo el contenido
		if err := copyDir(srcDir, dstDir); err != nil {
			return refs, fmt.Errorf("error copying %s: %w", dirName, err)
		}

		refs = append(refs, dirName+"/")
	}

	return refs, nil
}

// copyFile copia un archivo
func copyFile(src io.Reader, dstPath string) error {
	dst, err := os.Create(dstPath)
	if err != nil {
		return err
	}
	defer func() { _ = dst.Close() }()

	_, err = io.Copy(dst, src)
	return err
}

// ApplyConfig aplica una configuración guardada a un directorio de proyecto
func (r *Repository) ApplyConfig(agentType models.AgentType, configName string, targetPath string) error {
	// Verificar que existe la configuración guardada
	configDir := filepath.Join(r.baseDir, "agents", string(agentType), "configs", configName)
	if !dirExists(configDir) {
		return fmt.Errorf("configuration '%s' not found for agent '%s'", configName, agentType)
	}

	// Leer metadata
	metadataPath := filepath.Join(configDir, "agent.json")
	metadataData, err := os.ReadFile(metadataPath)
	if err != nil {
		return fmt.Errorf("error reading metadata: %w", err)
	}

	var config models.AgentConfig
	if err := json.Unmarshal(metadataData, &config); err != nil {
		return fmt.Errorf("error parsing metadata: %w", err)
	}

	// Crear directorio .opencode/ en destino si no existe
	opencodeDir := filepath.Join(targetPath, ".opencode")
	if err := os.MkdirAll(opencodeDir, 0755); err != nil {
		return fmt.Errorf("error creating .opencode directory: %w", err)
	}

	// Aplicar archivos no-portables
	// 1. Config principal (config.json → opencode.json)
	savedConfigPath := filepath.Join(configDir, "config.json")
	if fileExists(savedConfigPath) {
		targetConfigPath := filepath.Join(targetPath, "opencode.json")
		if err := copyFilePath(savedConfigPath, targetConfigPath); err != nil {
			return fmt.Errorf("error applying config: %w", err)
		}
		fmt.Printf("✓ Applied opencode.json\n")
	}

	// 2. Reglas (rules.md → AGENTS.md)
	savedRulesPath := filepath.Join(configDir, "rules.md")
	if fileExists(savedRulesPath) {
		targetRulesPath := filepath.Join(targetPath, "AGENTS.md")
		if err := copyFilePath(savedRulesPath, targetRulesPath); err != nil {
			return fmt.Errorf("error applying rules: %w", err)
		}
		fmt.Printf("✓ Applied AGENTS.md\n")
	}

	// 3. TUI config (tui.json)
	savedTUIPath := filepath.Join(configDir, "tui.json")
	if fileExists(savedTUIPath) {
		targetTUIPath := filepath.Join(targetPath, "tui.json")
		if err := copyFilePath(savedTUIPath, targetTUIPath); err != nil {
			return fmt.Errorf("error applying tui config: %w", err)
		}
		fmt.Printf("✓ Applied tui.json\n")
	}

	// 4. Aplicar elementos portables desde shared/
	sharedDir := filepath.Join(r.baseDir, "agents", string(agentType), "shared")
	if dirExists(sharedDir) {
		if err := r.applyPortables(sharedDir, opencodeDir, config.SharedRefs); err != nil {
			return fmt.Errorf("error applying portables: %w", err)
		}
	}

	fmt.Printf("✓ Configuration '%s' applied successfully to %s\n", configName, targetPath)
	return nil
}

// applyPortables aplica los elementos portables a un proyecto
func (r *Repository) applyPortables(sharedDir, targetOpenCodeDir string, sharedRefs []string) error {
	for _, ref := range sharedRefs {
		// ref es algo como "skills/" o "agents/"
		ref = filepath.Clean(ref)
		srcDir := filepath.Join(sharedDir, ref)
		
		if !dirExists(srcDir) {
			continue
		}

		dstDir := filepath.Join(targetOpenCodeDir, ref)
		
		// Copiar todo el contenido
		if err := copyDir(srcDir, dstDir); err != nil {
			return fmt.Errorf("error applying %s: %w", ref, err)
		}

		fmt.Printf("✓ Applied %s\n", ref)
	}

	return nil
}

// copyFilePath copia un archivo de origen a destino
func copyFilePath(srcPath, dstPath string) error {
	src, err := os.Open(srcPath)
	if err != nil {
		return err
	}
	defer src.Close()

	return copyFile(src, dstPath)
}

// copyDir copia un directorio recursivamente
func copyDir(src, dst string) error {
	return filepath.Walk(src, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return err
		}

		// Calcular ruta relativa
		relPath, err := filepath.Rel(src, path)
		if err != nil {
			return err
		}

		dstPath := filepath.Join(dst, relPath)

		if info.IsDir() {
			return os.MkdirAll(dstPath, info.Mode())
		}

		// Verificar si el archivo ya existe en destino (para no sobrescribir)
		if fileExists(dstPath) {
			return nil // Skip existing files
		}

		srcFile, err := os.Open(path)
		if err != nil {
			return err
		}
		defer func() { _ = srcFile.Close() }()

		return copyFile(srcFile, dstPath)
	})
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
