package config

import (
	"os"
	"path/filepath"
)

const (
	configFileName  = ".task-config"
	envTasksFile    = "TASKS_FILE"
	defaultFileName = "tasks.json"
)

type Config struct {
	TasksFilePath string
}

func ResolveTasksFilePath(cliFlagPath string) (string, error) {
	if cliFlagPath != "" {
		return cliFlagPath, nil
	}

	if envPath := os.Getenv(envTasksFile); envPath != "" {
		return envPath, nil
	}

	savedPath, err := loadSavedConfigPath()
	if err == nil && savedPath != "" {
		return savedPath, nil
	}

	return getDefaultPath(), nil
}

func getDefaultPath() string {
	homeDir, err := os.UserHomeDir()
	if err != nil {
		return defaultFileName
	}
	return filepath.Join(homeDir, defaultFileName)
}

func SaveTasksFilePath(path string) error {
	homeDir, err := os.UserHomeDir()
	if err != nil {
		return err
	}

	configDir := filepath.Join(homeDir, ".config")
	if err := os.MkdirAll(configDir, 0755); err != nil {
		return err
	}

	configPath := filepath.Join(configDir, configFileName)
	return os.WriteFile(configPath, []byte(path), 0644)
}

func loadSavedConfigPath() (string, error) {
	homeDir, err := os.UserHomeDir()
	if err != nil {
		return "", err
	}

	configPath := filepath.Join(homeDir, ".config", configFileName)
	data, err := os.ReadFile(configPath)
	if err != nil {
		return "", err
	}

	return string(data), nil
}
