package agents

import (
	"github.com/komyzi/komyzi/pkg/models"
)

// Detector define la interfaz para detectar configuraciones de agentes
type Detector interface {
	// Name retorna el nombre del agente
	Name() string
	
	// DetectProject busca configuración de agente en un directorio de proyecto
	DetectProject(path string) (*models.AgentPaths, error)
	
	// DetectGlobal busca configuración global del agente en el sistema
	DetectGlobal() (*models.AgentPaths, error)
	
	// IsInstalled verifica si el agente está instalado en el sistema
	IsInstalled() bool
}
