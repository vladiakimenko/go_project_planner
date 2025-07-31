package todo

import (
	"testing"
)

var testTasks = []Task{
	{ID: 0, Description: "Test task A", Done: false},
	{ID: 1, Description: "Test task B", Done: true},
	{ID: 2, Description: "Test task C", Done: false},
}

func TestAdd(t *testing.T) {
	t.Run("add", func(t *testing.T) {
		tasks := append([]Task{}, testTasks...) // copy slice
		originalLength := len(tasks)
		tasks = Add(tasks, "Test Task X")

		if len(tasks) != originalLength+1 {
			t.Fatalf("Test failed: incorrect number of tasks: %d", len(tasks))
		}

		if tasks[len(tasks)-1].Done != false {
			t.Errorf("Test failed: Incorrect initial 'Done' value: %+v", tasks[0])
		}
	})
}

func TestList(t *testing.T) {
	tests := []struct {
		name     string
		filter   TaskStateFilter
		expected int
	}{
		{"all tasks", FilterAll, 3},
		{"done", FilterDone, 1},
		{"pending tasks", FilterPending, 2},
	}
	for _, tt := range tests {
		tasks := append([]Task{}, testTasks...)
		t.Run(tt.name, func(t *testing.T) {
			got := List(tasks, string(tt.filter))
			if len(got) != tt.expected {
				t.Errorf("Test failed: Expected %d tasks, got %d", tt.expected, len(got))
			}
		})
	}

	t.Run("unknown filter defaults to all", func(t *testing.T) {
		tasks := append([]Task{}, testTasks...)
		got := List(tasks, "unknown")
		if len(got) != 3 {
			t.Errorf("Test failed: Expected 3 tasks with unknown filter, got %d", len(got))
		}
	})
}

func TestComplete(t *testing.T) {
	tests := []struct {
		name          string
		id            int
		errorExpected bool
	}{
		{"complete existing task", 1, false},
		{"complete non-existent task", 999, true},
	}
	for _, tt := range tests {
		tasks := append([]Task{}, testTasks...)
		t.Run(tt.name, func(t *testing.T) {
			updatedTasks, err := Complete(tasks, tt.id)

			if tt.errorExpected {
				if err == nil {
					t.Error("Test failed: Expected an error but didn't get one")
				}
				if len(updatedTasks) != 0 {
					t.Error("Test failed: Expected an empty slice along with an error")
				}
			} else {
				if err != nil {
					t.Errorf("Test failed: Unexpected error: %v", err)
				}
				// testTasks indices match the ID's
				if !updatedTasks[tt.id].Done {
					t.Errorf("Task %d was not marked as done", tt.id)
				}
			}
		})
	}
}

func TestDelete(t *testing.T) {
	tests := []struct {
		name          string
		id            int
		errorExpected bool
	}{
		{"delete existing task", 1, false},
		{"delete non-existent task", 999, true},
	}
	for _, tt := range tests {
		tasks := append([]Task{}, testTasks...)
		t.Run(tt.name, func(t *testing.T) {
			updatedTasks, err := Delete(tasks, tt.id)
			originalLength := len(tasks)
			if tt.errorExpected {
				if err == nil {
					t.Error("Test failed: Expected an error but didn't get one")
				}
				if len(updatedTasks) != 0 {
					t.Error("Test failed: Expected an empty slice along with an error")
				}
			} else {
				if err != nil {
					t.Errorf("Test failed: Unexpected error: %v", err)
				}
				if len(updatedTasks) != originalLength-1 {
					t.Errorf("Test failed: Unexpected number of tasks after deletion: %d", len(updatedTasks))
				}
				for _, task := range updatedTasks {
					if task.ID == tt.id {
						t.Errorf("Test failed: Task with ID %d was not deleted", tt.id)
					}
				}
			}
		})
	}
}
