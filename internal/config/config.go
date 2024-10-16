package config

import (
	"encoding/json"
	"os"
	"path/filepath"
)

type Config struct {
	DBurl           string `json:"db_url"`
	CurrentUserName string `json:"current_user_name"`
}

func Read() (Config, error) {
	path, err := getConfigPath()
	if err != nil {
		return Config{}, err
	}

	data, err := os.ReadFile(path)
	if err != nil {
		return Config{}, err
	}

	config := Config{}
	err = json.Unmarshal(data, &config)
	if err != nil {
		return Config{}, err
	}

	return config, nil
}

func (c Config) SetUser(name string) error {
	c.CurrentUserName = name

	configJson, err := json.Marshal(c)
	if err != nil {
		return err
	}
	writeConfig(configJson)
	return nil
}

func getConfigPath() (string, error) {
	root, err := os.UserHomeDir()
	if err != nil {
		return "", err
	}
	return filepath.Join(root, configFileName), nil
}

func writeConfig(configJson []byte) error {
	path, err := getConfigPath()
	if err != nil {
		return err
	}

	permissions := os.FileMode(0600)

	if err := os.WriteFile(path, configJson, permissions); err != nil {
		return err
	}

	return nil
}
