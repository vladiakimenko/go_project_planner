package main

import (
	"flag"
	"fmt"
	"log"
	"os"
	"path/filepath"
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
	LoadCmd     string = "load"
)

func main() {
	var tasks []todo.Task

	rawArgs := os.Args
	if len(rawArgs) < 2 {
		log.Fatal("No command provided")
	}
	command, args := rawArgs[1], rawArgs[2:]

	if command != LoadCmd {
		result, err := storage.LoadJSON(JsonStoragePath)
		if err != nil {
			log.Fatal(err)
		}
		tasks = result
	}

	switch command {
	case AddCmd:
		flagSet := flag.NewFlagSet(AddCmd, flag.ExitOnError)
		desc := flagSet.String("desc", "", "Task description")
		if err := flagSet.Parse(args); err != nil {
			log.Fatal(err)
		}
		if *desc == "" {
			log.Fatal("description is required")
		}
		updatedTasks := todo.Add(tasks, *desc)
		fmt.Printf("Successfully added:\n%v\n", updatedTasks[len(updatedTasks)-1])
		if err := storage.SaveJSON(JsonStoragePath, updatedTasks); err != nil {
			log.Fatal(err)
		}
	case ListCmd:
		flagSet := flag.NewFlagSet(ListCmd, flag.ExitOnError)
		filter := flagSet.String(
			"filter", string(todo.FilterAll),
			fmt.Sprintf("One of: %s, %s, %s", todo.FilterAll, todo.FilterDone, todo.FilterPending),
		)
		if err := flagSet.Parse(args); err != nil {
			log.Fatal(err)
		}
		if !slices.Contains([]string{string(todo.FilterAll), string(todo.FilterDone), string(todo.FilterPending)}, *filter) {
			log.Fatalf("Invalid filter value: %s", *filter)
		}
		filteredTasks := todo.List(tasks, *filter)
		for _, task := range filteredTasks {
			fmt.Println(task)
		}
	case CompleteCmd:
		flagSet := flag.NewFlagSet(CompleteCmd, flag.ExitOnError)
		id := flagSet.Int("id", -1, "Task id")
		if err := flagSet.Parse(args); err != nil {
			log.Fatal(err)
		}
		if *id == -1 {
			log.Fatal("id is required")
		}
		updatedTasks, err := todo.Complete(tasks, *id)
		if err != nil {
			log.Fatal(err)
		}
		if err := storage.SaveJSON(JsonStoragePath, updatedTasks); err != nil {
			log.Fatal(err)
		}
	case DeleteCmd:
		flagSet := flag.NewFlagSet(DeleteCmd, flag.ExitOnError)
		id := flagSet.Int("id", -1, "Task id")
		if err := flagSet.Parse(args); err != nil {
			log.Fatal(err)
		}
		if *id == -1 {
			log.Fatal("id is required")
		}
		updatedTasks, err := todo.Delete(tasks, *id)
		if err != nil {
			log.Fatal(err)
		}
		if err := storage.SaveJSON(JsonStoragePath, updatedTasks); err != nil {
			log.Fatal(err)
		}
	case ExportCmd:
		flagSet := flag.NewFlagSet(ExportCmd, flag.ExitOnError)
		format := flagSet.String("format", "", "Output format: json or csv")
		out := flagSet.String("out", "", "Output filepath")
		if err := flagSet.Parse(args); err != nil {
			log.Fatal(err)
		}
		if *format == "" || *out == "" {
			log.Fatal("both format and out flags are required")
		}
		var action func(string, []todo.Task) error
		switch *format {
		case "csv":
			action = storage.SaveCSV
		case "json":
			action = storage.SaveJSON
		default:
			log.Fatalf("incorrect format value provided: %s", *format)
		}
		if err := action(*out, tasks); err != nil {
			log.Fatal(err)
		}
	case LoadCmd:
		flagSet := flag.NewFlagSet(LoadCmd, flag.ExitOnError)
		file := flagSet.String("file", "", "Filepath to import")
		if err := flagSet.Parse(args); err != nil {
			log.Fatal(err)
		}
		var action func(string) ([]todo.Task, error)
		switch filepath.Ext(*file) {
		case ".json":
			action = storage.LoadJSON
		case ".csv":
			action = storage.LoadCSV
		default:
			log.Fatal("the output file must be either .json or .csv")
		}
		updatedTasks, err := action(*file)
		if err != nil {
			log.Fatal(err)
		}
		if err := storage.SaveJSON(JsonStoragePath, updatedTasks); err != nil {
			log.Fatal(err)
		}
	}
	fmt.Println("done")
	os.Exit(0)

}
