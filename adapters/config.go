package adapters

import (
	"fmt"
	"gopkg.in/yaml.v2"
	"os"
	"path/filepath"
)

type AuthType string

const (
	AuthTypeOAuth2 AuthType = "oauth2"
	AuthTypeToken  AuthType = "token"
)

type TodoistConfig struct {
	AuthType     AuthType  `yaml:"auth_type"`
	AuthToken    *string   `yaml:"token"`
	ClientID     *string   `yaml:"client_id"`
	ClientSecret *string   `yaml:"client_secret"`
	Scopes       *[]string `yaml:"scopes"`
}

type Settings struct {
	Todoist TodoistConfig `yaml:"todoist"`
}

func GetSettings() Settings {
	home, err := os.UserHomeDir()
	if err != nil {
		fmt.Println("Error getting user's home directory:", err)
		os.Exit(1)
	}

	configPath := filepath.Join(home, ".mgtd", "config.yaml")

	data, err := os.ReadFile(configPath)
	if err != nil {
		fmt.Println("Error reading config file:", err)
		os.Exit(1)
	}

	var settings Settings
	if err := yaml.Unmarshal(data, &settings); err != nil {
		fmt.Println("Error unmarshalling YAML:", err)
		os.Exit(1)
	}

	return settings
}
