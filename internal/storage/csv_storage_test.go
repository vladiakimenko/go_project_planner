package storage

import (
	"os"
	"path/filepath"
	"testing"

	"github.com/vladiakimenko/go_project_planner/internal/todo"
)

func TestLoadCSV(t *testing.T) {
	tests := []struct {
		name          string
		setup         func() string
		expected      []todo.Task
		errorExpected bool
	}{
		{
			name: "successful load",
			setup: func() string {
				path := filepath.Join(t.TempDir(), "test.csv")
				file, err := os.Create(path)
				if err != nil {
					t.Fatal(err)
				}
				defer file.Close()

				content := "ID,Description,Done\n0,Task A,false\n1,Task B,true\n2,Task C,false\n"
				if _, err := file.WriteString(content); err != nil {
					t.Fatal(err)
				}
				return path
			},
			expected: []todo.Task{
				{ID: 0, Description: "Task A", Done: false},
				{ID: 1, Description: "Task B", Done: true},
				{ID: 2, Description: "Task C", Done: false},
			},
			errorExpected: false,
		},
		{
			name: "file does not exist - creates new",
			setup: func() string {
				return filepath.Join(t.TempDir(), "missing.csv")
			},
			expected:      []todo.Task{},
			errorExpected: false,
		},
		{
			name: "invalid csv format",
			setup: func() string {
				path := filepath.Join(t.TempDir(), "invalid.csv")
				file, err := os.Create(path)
				if err != nil {
					t.Fatal(err)
				}
				defer file.Close()

				content := "ID,Description\n0,Task A\n1,Task B\n" // missing Done
				if _, err := file.WriteString(content); err != nil {
					t.Fatal(err)
				}
				return path
			},
			expected:      []todo.Task{},
			errorExpected: true,
		},
		{
			name: "invalid ID format",
			setup: func() string {
				path := filepath.Join(t.TempDir(), "invalid_id.csv")
				file, err := os.Create(path)
				if err != nil {
					t.Fatal(err)
				}
				defer file.Close()

				content := "ID,Description,Done\nA,Task A,false\n" // non-numeric ID
				if _, err := file.WriteString(content); err != nil {
					t.Fatal(err)
				}
				return path
			},
			expected:      []todo.Task{},
			errorExpected: true,
		},
		{
			name: "invalid Done format",
			setup: func() string {
				path := filepath.Join(t.TempDir(), "invalid_done.csv")
				file, err := os.Create(path)
				if err != nil {
					t.Fatal(err)
				}
				defer file.Close()

				content := "ID,Description,Done\n0,Task A,notaboolean\n" // invalid boolean
				if _, err := file.WriteString(content); err != nil {
					t.Fatal(err)
				}
				return path
			},
			expected:      []todo.Task{},
			errorExpected: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			path := tt.setup()
			result, err := LoadCSV(path)
			if (err != nil) != tt.errorExpected {
				t.Errorf("Test failed: error occured %v, error expected %v", err, tt.errorExpected)
				return
			}

			if len(result) != len(tt.expected) {
				t.Errorf("Test failed: got %d tasks, expected %d", len(result), len(tt.expected))
				return
			}

			for i := range result {
				if result[i].ID != tt.expected[i].ID ||
					result[i].Description != tt.expected[i].Description ||
					result[i].Done != tt.expected[i].Done {
					t.Errorf("Test failed: task %d = %v, expected %v", i, result[i], tt.expected[i])
				}
			}
		})
	}
}

func TestSaveCSV(t *testing.T) {
	tests := []struct {
		name  string
		tasks []todo.Task
	}{
		{
			name: "successful save",
			tasks: []todo.Task{
				{ID: 0, Description: "Task A", Done: false},
				{ID: 1, Description: "Task B", Done: true},
			},
		},
		{
			name:  "empty tasks",
			tasks: []todo.Task{},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			path := filepath.Join(t.TempDir(), "test_save.csv")

			if err := SaveCSV(path, tt.tasks); err != nil {
				t.Errorf("Test failed: Unexpected error: %v", err)
			}
			loadedTasks, err := LoadCSV(path)
			if err != nil {
				t.Errorf("Test failed: couldn't load the file back: %v", err)
			}

			if len(loadedTasks) != len(tt.tasks) {
				t.Errorf("Test failed: saved %d tasks, but loaded %d", len(tt.tasks), len(loadedTasks))
			}

			for i := range loadedTasks {
				if loadedTasks[i].ID != tt.tasks[i].ID ||
					loadedTasks[i].Description != tt.tasks[i].Description ||
					loadedTasks[i].Done != tt.tasks[i].Done {
					t.Errorf("Test failed: saved task %d = %v, but loaded %v", i, tt.tasks[i], loadedTasks[i])
				}
			}
		})
	}
}
