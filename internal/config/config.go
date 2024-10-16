package config

import (
	"encoding/json"
	"log"
	"os"
	"path/filepath"
)

type Config struct {
	DBurl           string `json:"db_url"`
	CurrentUserName string `json:"current_user_name"`
}

func Read() Config {
	data, err := os.ReadFile(getConfigPath())
	if err != nil {
		log.Fatal(err)
	}
	config := Config{}
	err = json.Unmarshal(data, &config)
	if err != nil {
		log.Fatal(err)
	}
	return config
}

func (c Config) SetUser(name string) {
	c.CurrentUserName = name
	configJson, err := json.Marshal(c)
	if err != nil {
		log.Fatal(err)
	}
	writeConfig(configJson)
}

func getConfigPath() string {
	root, err := os.UserHomeDir()
	if err != nil {
		log.Fatal(err)
	}
	return filepath.Join(root, configFileName)
}

func writeConfig(configJson []byte) {
	path := getConfigPath()
	permissions := os.FileMode(0600)
	if err := os.WriteFile(path, configJson, permissions); err != nil {
		log.Fatal(err)
	}
}
