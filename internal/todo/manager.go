package todo

import (
	"fmt"

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
	for i, task := range tasks {
		if task.ID == id {
			tasks[i].Done = true
			return tasks, nil
		}
	}
	logging.Logger.Error("Could not find a task with specified id", "id", id)
	return []Task{}, fmt.Errorf("task with requested id=%d is missing", id)
}

func Delete(tasks []Task, id int) ([]Task, error) {
	for i, task := range tasks {
		if task.ID == id {
			tasks = append(tasks[:i], tasks[i+1:]...)
			return tasks, nil
		}
	}
	logging.Logger.Error("Could not find a task with specified id", "id", id)
	return []Task{}, fmt.Errorf("task with requested id=%d is missing", id)
}
