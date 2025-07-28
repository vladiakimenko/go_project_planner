package storage

import (
	"encoding/json"
	"errors"
	"fmt"
	"os"

	"github.com/vladiakimenko/go_project_planner/internal/logging"
	"github.com/vladiakimenko/go_project_planner/internal/todo"
)

func LoadJSON(path string) ([]todo.Task, error) {
	tasks := []todo.Task{}
	if _, err := os.Stat(path); err != nil {
		if errors.Is(err, os.ErrNotExist) {
			logging.Logger.Debug(fmt.Sprintf("The json storage file is missing. Creating %s", path))
			SaveJSON(path, tasks)
		} else {
			logging.Logger.Error("File stats request failed unexpectively", "error", err.Error(), "file", path)
			return tasks, fmt.Errorf("failed to access json storage: %w", err)
		}
	}
	fileBytes, err := os.ReadFile(path)
	if err != nil {
		logging.Logger.Error("Error reading the json storage file", "error", err.Error(), "file", path)
		return tasks, fmt.Errorf("failed to read json storage: %w", err)
	}
	if err := json.Unmarshal(fileBytes, &tasks); err != nil {
		logging.Logger.Error("Error unmarshalling the json storage file", "error", err.Error(), "file", path)
		return tasks, fmt.Errorf("failed to parse json: %w", err)
	}
	logging.Logger.Debug("Loaded tasks from json", "amount", len(tasks))
	return tasks, nil
}

func SaveJSON(path string, tasks []todo.Task) error {
	resultBytes, err := json.MarshalIndent(tasks, "", "  ")
	if err != nil {
		logging.Logger.Error("Error marshalling tasks to json", "error", err.Error(), "tasks", tasks)
		return fmt.Errorf("failed to dump json: %w", err)
	}
	if err := os.WriteFile(path, resultBytes, 0644); err != nil {
		logging.Logger.Error("Failed writing to file", "error", err.Error(), "tasks", tasks, "path", path)
		return fmt.Errorf("failed to write to json storage: %w", err)
	}
	logging.Logger.Debug("Save tasks to json", "amount", len(tasks))
	return nil
}
