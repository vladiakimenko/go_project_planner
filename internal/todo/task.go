package todo

import "fmt"

type Task struct {
	ID          int
	Description string
	Done        bool
}

func (t Task) String() string {
	return fmt.Sprintf("%d. %s: %t", t.ID, t.Description, t.Done)
}
