package main

import (
	"fmt"
)

type CommandType int

const (
	AddTask CommandType = iota
	ListTasks
	CompleteTask
	DeleteTask
	ExportTasks
	LoadTasks
)

func (ct CommandType) String() string {
	return [...]string{
		"add",
		"list",
		"complete",
		"delete",
		"export",
		"import",
	}[ct]
}

var ArgsApplicable = map[CommandType][]string{
	AddTask:      {"--desc"},
	ListTasks:    {"--filter"},
	CompleteTask: {"--id"},
	DeleteTask:   {"--id"},
	ExportTasks:  {},
	LoadTasks:    {},
}

func printHelp() {
	fmt.Println(`
Usage:
  task [command] [flags]

Available Commands:
  add       Add new task
    --desc string    Task description (required)

  list      List tasks
    --filter string  Filter tasks (all/done/pending, default "all")

  complete  Mark task as completed
    --id int         Task ID to complete (required)

  delete    Delete task
    --id int         Task ID to delete (required)

  export    Export tasks
    --format string  Export format (required)
    --out string     Output file path (required)

  load      Import tasks
    --file string    File to import (required)`,
	)
}
