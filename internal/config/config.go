package config

import (
	"encoding/json"
	"os"
)

type Config struct {
	EncryptionEnabled bool   `json:"encryption_enabled"`
	ServerHost        string `json:"server_host"`
	ServerPort        int    `json:"server_port"`
	SecretKey         []byte `json:"secret_key"`
}

const DefaultConfigFilename = "config.json"

func ReadConfig(path string) (*Config, error) {
	data, err := os.ReadFile(path)
	if err != nil {
		return nil, err
	}

	var config Config
	err = json.Unmarshal(data, &config)
	if err != nil {
		return nil, err
	}

	return &config, nil
}

func DefaultConfig() *Config {
	return &Config{
		EncryptionEnabled: true,
		ServerHost:        "localhost",
		ServerPort:        8080,
		SecretKey:         []byte("A0KWJW3qRCiYcEj3"),
	}
}

func EnsureConfig(path string) (bool, error) {
	if _, err := os.Stat(path); !os.IsNotExist(err) {
		return false, nil
	}

	defaultConfig := DefaultConfig()
	data, err := json.MarshalIndent(defaultConfig, "", "    ")
	if err != nil {
		return false, err
	}

	err = os.WriteFile(path, data, 0644)
	if err != nil {
		return false, err
	}
	return true, nil
}
