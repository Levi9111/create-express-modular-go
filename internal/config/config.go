package config

import (
	"encoding/json"
	"os"
	"path/filepath"
)

type CemProject struct {
	Name              string `json:"name"`
	DB                string `json:"db"`
	Validator         string `json:"validator"`
	Auth              bool   `json:"auth"`
	AuthTokenDelivery string `json:"authTokenDelivery,omitempty"`
	Docker            bool   `json:"docker"`
	PackageManager    string `json:"packageManager"`
}

type CemStructure struct {
	SrcDir         string `json:"srcDir"`
	AppDir         string `json:"appDir"`
	ModulesDir     string `json:"modulesDir"`
	MiddlewaresDir string `json:"middlewaresDir"`
	UtilsDir       string `json:"utilsDir"`
	ConfigDir      string `json:"configDir"`
}

type CemFeatures struct {
	Testing   bool `json:"testing"`
	Websocket bool `json:"websocket"`
	Swagger   bool `json:"swagger"`
}

type CemState struct {
	Modules     []string `json:"modules"`
	Middlewares []string `json:"middlewares"`
	Env         []string `json:"env"`
}

type CemConfig struct {
	Version   string       `json:"version"`
	Project   CemProject   `json:"project"`
	Structure CemStructure `json:"structure"`
	Features  CemFeatures  `json:"features"`
	State     *CemState    `json:"state,omitempty"`
}

const ConfigFileName = "cem-cli.json"

// LoadConfig reads and parses cem-cli.json from the project root.
func LoadConfig(dir string) (*CemConfig, error) {
	path := filepath.Join(dir, ConfigFileName)
	data, err := os.ReadFile(path)
	if err != nil {
		return nil, err
	}

	var cfg CemConfig
	if err := json.Unmarshal(data, &cfg); err != nil {
		return nil, err
	}
	return &cfg, nil
}

// SaveConfig writes the CemConfig to cem-cli.json.
func SaveConfig(dir string, cfg *CemConfig) error {
	path := filepath.Join(dir, ConfigFileName)
	data, err := json.MarshalIndent(cfg, "", "  ")
	if err != nil {
		return err
	}
	return os.WriteFile(path, data, 0644)
}

// UpdateState updates cem-cli.json's live state arrays (modules, middlewares, env).
// No-op if cem-cli.json does not exist.
func UpdateState(dir, kind, action string, items []string) {
	if len(items) == 0 {
		return
	}
	cfg, err := LoadConfig(dir)
	if err != nil || cfg == nil {
		return
	}
	if cfg.State == nil {
		cfg.State = &CemState{
			Modules:     []string{},
			Middlewares: []string{},
			Env:         []string{},
		}
	}

	var target *[]string
	switch kind {
	case "modules":
		target = &cfg.State.Modules
	case "middlewares":
		target = &cfg.State.Middlewares
	case "env":
		target = &cfg.State.Env
	default:
		return
	}

	seen := make(map[string]bool)
	var next []string
	for _, existing := range *target {
		if !seen[existing] {
			seen[existing] = true
			next = append(next, existing)
		}
	}

	if action == "add" {
		for _, item := range items {
			if !seen[item] {
				seen[item] = true
				next = append(next, item)
			}
		}
	} else if action == "remove" {
		rm := make(map[string]bool, len(items))
		for _, item := range items {
			rm[item] = true
		}
		var filtered []string
		for _, existing := range next {
			if !rm[existing] {
				filtered = append(filtered, existing)
			}
		}
		next = filtered
	}

	if next == nil {
		next = []string{}
	}
	*target = next
	_ = SaveConfig(dir, cfg)
}
