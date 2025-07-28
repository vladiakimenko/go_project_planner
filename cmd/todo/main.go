package main

import (
	"os"
)

const JsonStoragePath string = "../tasks.json"

// parse args and call funcs
// TODO: add tests

func main() {
	if len(os.Args) < 1 {
		printHelp()
		os.Exit(1)
	}
}
