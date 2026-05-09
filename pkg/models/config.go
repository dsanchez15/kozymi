package models

import "time"

// AgentType representa el tipo de agente de IA
type AgentType string

const (
	AgentOpenCode AgentType = "opencode"
	AgentKiro     AgentType = "kiro"
	AgentClaude   AgentType = "claude"
)

// ConfigSource indica si la config viene de un proyecto o es global
type ConfigSource string

const (
	SourceProject ConfigSource = "project"
	SourceGlobal  ConfigSource = "global"
)

// AgentConfig representa una configuración guardada de un agente
type AgentConfig struct {
	Agent      AgentType    `json:"agent"`
	Name       string       `json:"name"`
	Source     ConfigSource `json:"source"`
	SourcePath string       `json:"source_path,omitempty"`
	DateSaved  time.Time    `json:"date_saved"`
	SharedRefs []string     `json:"shared_refs,omitempty"`
}

// AgentPaths contiene las rutas relevantes de un agente detectado
type AgentPaths struct {
	AgentType      AgentType
	ConfigPath     string // ruta a opencode.json
	RulesPath      string // ruta a AGENTS.md
	TUIPath        string // ruta a tui.json
	OpenCodeDir    string // ruta a .opencode/
	IsGlobal       bool
}
