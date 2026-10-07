package config

import (
	"encoding/json"
	"os"
	"path/filepath"
)

type CemConfig struct {
	Name           string `json:"name"`
	DB             string `json:"db"`
	Validator      string `json:"validator"`
	Auth           bool   `json:"auth"`
	TokenDelivery  string `json:"tokenDelivery,omitempty"`
	Docker         bool   `json:"docker"`
	Swagger        bool   `json:"swagger"`
	PackageManager string `json:"packageManager"`
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
