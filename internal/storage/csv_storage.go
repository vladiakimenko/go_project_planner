package storage

import (
	"encoding/csv"
	"errors"
	"fmt"
	"os"
	"strconv"

	"github.com/vladiakimenko/go_project_planner/internal/logging"
	"github.com/vladiakimenko/go_project_planner/internal/todo"
)

func LoadCSV(path string) ([]todo.Task, error) {
	tasks := []todo.Task{}
	if _, err := os.Stat(path); err != nil {
		if errors.Is(err, os.ErrNotExist) {
			logging.Logger.Debug(fmt.Sprintf("The csv storage file is missing. Creating %s", path))
			SaveCSV(path, tasks)
		} else {
			logging.Logger.Error("File stats request failed unexpectively", "error", err.Error(), "file", path)
			return tasks, fmt.Errorf("failed to access csv storage: %w", err)
		}
	}
	file, err := os.Open(path)
	if err != nil {
		logging.Logger.Error("Error reading the csv storage file", "error", err.Error(), "file", path)
		return tasks, fmt.Errorf("failed to read csv storage: %w", err)
	}
	defer file.Close()

	data, err := csv.NewReader(file).ReadAll()
	if err != nil {
		logging.Logger.Error("Error parsing the csv storage file", "error", err.Error(), "file", path)
		return tasks, fmt.Errorf("failed to parse csv: %w", err)
	}
	for _, row := range data[1:] {
		if len(row) < 3 {
			logging.Logger.Error("Error desierializing a row. Wrong number of values", "row", row)
			return []todo.Task{}, errors.New("could not create a Task from a row, wrong number of values found")
		}
		convertedId, err := strconv.Atoi(row[0])
		if err != nil {
			return []todo.Task{}, fmt.Errorf("invalid ID format: %v", row[0])
		}
		convertedDone, err := strconv.ParseBool(row[2])
		if err != nil {
			return []todo.Task{}, fmt.Errorf("invalid Done format: %v", row[2])
		}
		tasks = append(tasks, todo.Task{
			ID:          convertedId,
			Description: row[1],
			Done:        convertedDone,
		})
	}
	logging.Logger.Debug("Loaded tasks from csv", "amount", len(tasks))
	return tasks, nil
}

func SaveCSV(path string, tasks []todo.Task) error {

	file, err := os.Create(path)
	if err != nil {
		logging.Logger.Error("Error creating file", "error", err.Error(), "path", path)
		return fmt.Errorf("failed to create a csv storage: %w", err)
	}
	defer file.Close()

	writer := csv.NewWriter(file)
	defer writer.Flush()

	headers := []string{"ID", "Description", "Done"}
	if err := writer.Write(headers); err != nil {
		logging.Logger.Error("Error writing headers", "error", err.Error(), "headers", headers)
		return fmt.Errorf("failed to write headers to csv: %w", err)
	}
	for _, task := range tasks {
		serializedRow := []string{strconv.Itoa(task.ID), task.Description, strconv.FormatBool(task.Done)}
		if err := writer.Write(serializedRow); err != nil {
			logging.Logger.Error("Error writing a row", "error", err.Error(), "task", task)
			return fmt.Errorf("failed to write a row to csv: %w", err)
		}
	}
	logging.Logger.Debug("Save tasks to csv", "amount", len(tasks))
	return nil
}
