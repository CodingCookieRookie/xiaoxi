package storage

import (
	"encoding/json"
	"errors"
	"os"
	"path/filepath"

	"task/internal/task"
)

var (
	ErrFileNotFound = errors.New("task file not found")
	ErrInvalidJSON  = errors.New("invalid JSON in task file")
)

func Load(path string) (*task.TaskList, error) {
	data, err := os.ReadFile(path)
	if err != nil {
		if os.IsNotExist(err) {
			return &task.TaskList{Tasks: []task.Task{}}, nil
		}
		return nil, err
	}

	if len(data) == 0 {
		return &task.TaskList{Tasks: []task.Task{}}, nil
	}

	var tl task.TaskList
	if err := json.Unmarshal(data, &tl); err != nil {
		return nil, ErrInvalidJSON
	}

	tl.SyncNextID()

	return &tl, nil
}

func Save(path string, tl *task.TaskList) error {
	data, err := json.MarshalIndent(tl, "", "  ")
	if err != nil {
		return err
	}

	dir := filepath.Dir(path)
	if err := os.MkdirAll(dir, 0755); err != nil {
		return err
	}

	tmpPath := path + ".tmp"
	if err := os.WriteFile(tmpPath, data, 0644); err != nil {
		return err
	}

	if err := os.Rename(tmpPath, path); err != nil {
		os.Remove(tmpPath)
		return err
	}

	return nil
}
