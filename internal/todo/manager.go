package todo

import (
	"fmt"
	"slices"

	"github.com/vladiakimenko/go_project_planner/internal/logging"
)

type TaskStateFilter string

const (
	FilterAll     TaskStateFilter = "all"
	FilterDone    TaskStateFilter = "done"
	FilterPending TaskStateFilter = "pending"
)

var FilterConditionsMap = map[TaskStateFilter]func(Task) bool{
	FilterAll:     func(t Task) bool { return true },
	FilterDone:    func(t Task) bool { return t.Done },
	FilterPending: func(t Task) bool { return !t.Done },
}

func Add(tasks []Task, desc string) []Task {
	return append(tasks, Task{
		ID:          len(tasks),
		Description: desc,
		Done:        false,
	})
}

func List(tasks []Task, filter string) []Task {
	result := []Task{}
	filterType := TaskStateFilter(filter)
	for _, item := range tasks {
		if FilterConditionsMap[filterType](item) {
			result = append(result, item)
		}
	}
	return result
}

func Complete(tasks []Task, id int) ([]Task, error) {
	if err := validateId(tasks, id); err != nil {
		return []Task{}, fmt.Errorf("could not complete task with requested id: %w", err)
	}
	tasks[id].Done = true
	return tasks, nil
}

func Delete(tasks []Task, id int) ([]Task, error) {
	if err := validateId(tasks, id); err != nil {
		return []Task{}, fmt.Errorf("could not complete task with requested id: %w", err)
	}
	tasks = slices.Delete(tasks, 1, 1)
	return tasks, nil
}

func validateId(tasks []Task, id int) error {
	if id < 0 {
		logging.Logger.Error("The provided id is below 0", "id", id)
		return fmt.Errorf("invalid id provided: %d, must be a positive integer", id)
	}
	tasksCount := len(tasks)
	if id > tasksCount {
		logging.Logger.Error("The provided id is missing in the slice", "id", id, "count", tasksCount)
		return fmt.Errorf("the requested id %d is out of range", id)
	}
	return nil
}
