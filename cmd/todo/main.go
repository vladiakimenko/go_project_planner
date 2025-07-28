package main

import (
	"flag"
	"fmt"
	"log"
	"os"
	"slices"

	"github.com/vladiakimenko/go_project_planner/internal/storage"
	"github.com/vladiakimenko/go_project_planner/internal/todo"
)

const JsonStoragePath string = "tasks.json"

const (
	AddCmd      string = "add"
	ListCmd     string = "list"
	CompleteCmd string = "complete"
	DeleteCmd   string = "delete"
	ExportCmd   string = "export"
	LoadCmd     string = "import"
)

// parse args and call funcs
// TODO: add tests

func main() {
	tasks, err := storage.LoadJSON(JsonStoragePath)
	if err != nil {
		log.Fatal(err.Error())
	}

	rawArgs := os.Args
	if len(rawArgs) < 2 {
		log.Fatal("No command provided")
	}
	command, args := rawArgs[1], rawArgs[2:]

	switch command {
	case AddCmd:
		flagSet := flag.NewFlagSet(AddCmd, flag.ExitOnError)
		desc := flagSet.String("desc", "", "Task description")
		if err := flagSet.Parse(args); err != nil {
			log.Fatal(err.Error())
		}
		if *desc == "" {
			log.Fatal("description is required")
		}
		tasks := todo.Add(tasks, *desc)
		fmt.Printf("Successfully added:\n%v\n", tasks[len(tasks)-1])
		storage.SaveJSON(JsonStoragePath, tasks)
	case ListCmd:
		flagSet := flag.NewFlagSet(ListCmd, flag.ExitOnError)
		filter := flagSet.String(
			"filter", string(todo.FilterAll),
			fmt.Sprintf("One of: %s, %s, %s", todo.FilterAll, todo.FilterDone, todo.FilterPending),
		)
		if err := flagSet.Parse(args); err != nil {
			log.Fatal(err.Error())
		}
		if !slices.Contains([]string{string(todo.FilterAll), string(todo.FilterDone), string(todo.FilterPending)}, *filter) {
			log.Fatalf("Invalid filter value: %s", *filter)
		}
		tasks := todo.List(tasks, *filter)
		for _, task := range tasks {
			fmt.Println(task)
		}
	}
	os.Exit(0)
}
