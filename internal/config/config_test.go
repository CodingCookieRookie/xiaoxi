package config

import (
	"os"
	"path/filepath"
	"testing"
)

func TestResolveTasksFilePath_CLIFlag(t *testing.T) {
	path, err := ResolveTasksFilePath("/custom/path/tasks.json")
	if err != nil {
		t.Errorf("expected no error, got %v", err)
	}

	if path != "/custom/path/tasks.json" {
		t.Errorf("expected /custom/path/tasks.json, got %s", path)
	}
}

func TestResolveTasksFilePath_CLIFlag_Priority(t *testing.T) {
	os.Setenv(envTasksFile, "/env/path/tasks.json")
	defer os.Unsetenv(envTasksFile)

	path, err := ResolveTasksFilePath("/cli/path/tasks.json")
	if err != nil {
		t.Errorf("expected no error, got %v", err)
	}

	if path != "/cli/path/tasks.json" {
		t.Errorf("expected CLI flag path, got %s", path)
	}
}

func TestResolveTasksFilePath_EnvVar(t *testing.T) {
	os.Setenv(envTasksFile, "/env/path/tasks.json")
	defer os.Unsetenv(envTasksFile)

	path, err := ResolveTasksFilePath("")
	if err != nil {
		t.Errorf("expected no error, got %v", err)
	}

	if path != "/env/path/tasks.json" {
		t.Errorf("expected env var path, got %s", path)
	}
}

func TestResolveTasksFilePath_SavedConfig(t *testing.T) {
	tmpDir := t.TempDir()
	configDir := filepath.Join(tmpDir, ".config")
	if err := os.MkdirAll(configDir, 0755); err != nil {
		t.Fatalf("failed to create config dir: %v", err)
	}

	savedPath := filepath.Join(configDir, configFileName)
	expectedPath := "/saved/path/tasks.json"
	if err := os.WriteFile(savedPath, []byte(expectedPath), 0644); err != nil {
		t.Fatalf("failed to write config file: %v", err)
	}

	originalHome := os.Getenv("HOME")
	os.Setenv("HOME", tmpDir)
	defer os.Setenv("HOME", originalHome)

	os.Unsetenv(envTasksFile)

	path, err := ResolveTasksFilePath("")
	if err != nil {
		t.Errorf("expected no error, got %v", err)
	}

	if path != expectedPath {
		t.Errorf("expected saved config path, got %s", path)
	}
}

func TestResolveTasksFilePath_DefaultPath(t *testing.T) {
	os.Unsetenv(envTasksFile)

	path, err := ResolveTasksFilePath("")
	if err != nil {
		t.Errorf("expected no error, got %v", err)
	}

	homeDir, _ := os.UserHomeDir()
	expected := filepath.Join(homeDir, defaultFileName)

	if path != expected {
		t.Errorf("expected default path %s, got %s", expected, path)
	}
}

func TestSaveTasksFilePath(t *testing.T) {
	tmpDir := t.TempDir()
	homeDir := filepath.Join(tmpDir, "home")
	if err := os.MkdirAll(homeDir, 0755); err != nil {
		t.Fatalf("failed to create home dir: %v", err)
	}

	originalHome := os.Getenv("HOME")
	os.Setenv("HOME", homeDir)
	defer os.Setenv("HOME", originalHome)

	testPath := "/test/path/tasks.json"
	if err := SaveTasksFilePath(testPath); err != nil {
		t.Errorf("expected no error, got %v", err)
	}

	savedPath := filepath.Join(homeDir, ".config", configFileName)
	data, err := os.ReadFile(savedPath)
	if err != nil {
		t.Errorf("expected to read saved config, got %v", err)
	}

	if string(data) != testPath {
		t.Errorf("expected %s, got %s", testPath, string(data))
	}
}
