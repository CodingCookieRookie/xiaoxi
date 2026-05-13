package storage

import (
	"os"
	"path/filepath"
	"testing"

	"task/internal/task"
)

func TestLoad_MissingFile(t *testing.T) {
	tl, err := Load("/nonexistent/path/tasks.json")
	if err != nil {
		t.Errorf("expected no error for missing file, got %v", err)
	}

	if len(tl.Tasks) != 0 {
		t.Errorf("expected empty task list, got %d tasks", len(tl.Tasks))
	}
}

func TestLoad_EmptyFile(t *testing.T) {
	tmpDir := t.TempDir()
	tmpFile := filepath.Join(tmpDir, "tasks.json")

	if err := os.WriteFile(tmpFile, []byte(""), 0644); err != nil {
		t.Fatalf("failed to create temp file: %v", err)
	}

	tl, err := Load(tmpFile)
	if err != nil {
		t.Errorf("expected no error for empty file, got %v", err)
	}

	if len(tl.Tasks) != 0 {
		t.Errorf("expected empty task list, got %d tasks", len(tl.Tasks))
	}
}

func TestLoad_ValidJSON(t *testing.T) {
	tmpDir := t.TempDir()
	tmpFile := filepath.Join(tmpDir, "tasks.json")

	data := `{"tasks":[{"id":1,"title":"Test task","completed":false,"created_at":"2024-01-01T00:00:00Z"}],"next_id":2}`
	if err := os.WriteFile(tmpFile, []byte(data), 0644); err != nil {
		t.Fatalf("failed to create temp file: %v", err)
	}

	tl, err := Load(tmpFile)
	if err != nil {
		t.Errorf("expected no error, got %v", err)
	}

	if len(tl.Tasks) != 1 {
		t.Errorf("expected 1 task, got %d", len(tl.Tasks))
	}

	if tl.Tasks[0].Title != "Test task" {
		t.Errorf("expected title %q, got %q", "Test task", tl.Tasks[0].Title)
	}
}

func TestLoad_InvalidJSON(t *testing.T) {
	tmpDir := t.TempDir()
	tmpFile := filepath.Join(tmpDir, "tasks.json")

	data := `{invalid json}`
	if err := os.WriteFile(tmpFile, []byte(data), 0644); err != nil {
		t.Fatalf("failed to create temp file: %v", err)
	}

	_, err := Load(tmpFile)
	if err != ErrInvalidJSON {
		t.Errorf("expected ErrInvalidJSON, got %v", err)
	}
}

func TestSave(t *testing.T) {
	tmpDir := t.TempDir()
	tmpFile := filepath.Join(tmpDir, "tasks.json")

	tl := task.NewTaskList()
	tl.Add("Task 1")
	tl.Add("Task 2")

	if err := Save(tmpFile, tl); err != nil {
		t.Errorf("expected no error, got %v", err)
	}

	loaded, err := Load(tmpFile)
	if err != nil {
		t.Errorf("expected no error loading saved file, got %v", err)
	}

	if len(loaded.Tasks) != 2 {
		t.Errorf("expected 2 tasks, got %d", len(loaded.Tasks))
	}
}

func TestSave_CreatesDirectory(t *testing.T) {
	tmpDir := t.TempDir()
	tmpFile := filepath.Join(tmpDir, "subdir", "tasks.json")

	tl := task.NewTaskList()
	tl.Add("Test task")

	if err := Save(tmpFile, tl); err != nil {
		t.Errorf("expected no error, got %v", err)
	}

	loaded, err := Load(tmpFile)
	if err != nil {
		t.Errorf("expected no error, got %v", err)
	}

	if len(loaded.Tasks) != 1 {
		t.Errorf("expected 1 task, got %d", len(loaded.Tasks))
	}
}
