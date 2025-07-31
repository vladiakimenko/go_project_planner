# Task Manager CLI
A simple CLI task manager with JSON storage supporting CRUD operations, filtering, and import/export.  Stores data in a `tasks.json`.

## Usage
Run from project root:
```bash
go run cmd/todo/main.go <command> [flags]
```

## Commands
**add** - Add a new task  
Flags:  
*-desc* - Task description (required)

**list**- List all tasks  
Flags:  
*-filter* - Filter tasks (values: all, done, pending)

**complete** - Mark a task as completed  
Flags:  
*-id* - Task ID to complete (required)

**delete** - Delete a task  
Flags:  
*-id* - Task ID to delete (required)

**export** - Export tasks to file  
Flags:  
*-format* - Output format (json or csv)  
*-out* - Output file path (required)

**load** - Import tasks from file
Flags:  
*-file* - Input file path (json or csv) (required)

## Logging
Set stodout logging verbosity with LOG_LEVEL (default: INFO/0):
| DEBUG | INFO | WARN | ERROR |
|-------|------|------|-------|
| -4 | 0 | 4 | 8 |


```bash
export LOG_LEVEL=DEBUG
```

## Requirements
Go 1.21 or later

## Testing
```bash
go test ./...
```

## Examples

```
$ go run cmd/todo/main.go add --desc "Buy milk"
Successfully added:
0. Buy milk: false
done
$ go run cmd/todo/main.go add --desc "Do homework"
Successfully added:
1. Do homework: false
done
$ go run cmd/todo/main.go add --desc "Clean room"
Successfully added:
2. Clean room: false
done
$ go run cmd/todo/main.go list
0. Buy milk: false
1. Do homework: false
2. Clean room: false
done
$ go run cmd/todo/main.go complete --id 1
done
$ go run cmd/todo/main.go list --filter pending
0. Buy milk: false
2. Clean room: false
done
$ go run cmd/todo/main.go list --filter done
1. Do homework: true
done
$ go run cmd/todo/main.go list --filter all
0. Buy milk: false
1. Do homework: true
2. Clean room: false
done
$ go run cmd/todo/main.go delete --id 0
done
$ go run cmd/todo/main.go list
1. Do homework: true
2. Clean room: false
done
$ go run cmd/todo/main.go export --format csv --out "output.csv"
done
$ cat output.csv 
ID,Description,Done
1,Do homework,true
2,Clean room,false
$ go run cmd/todo/main.go export --format json --out "output.json"
done
$ cat output.json
[
  {
    "ID": 1,
    "Description": "Do homework",
    "Done": true
  },
  {
    "ID": 2,
    "Description": "Clean room",
    "Done": false
  }
]$rm tasks.json 
$ go run cmd/todo/main.go list
done
$ go run cmd/todo/main.go load --file output.csv
done
$ go run cmd/todo/main.go list
1. Do homework: true
2. Clean room: false
done
$ rm tasks.json 
$ go run cmd/todo/main.go list
done
$ go run cmd/todo/main.go load --file output.json
done
$ go run cmd/todo/main.go list
1. Do homework: true
2. Clean room: false
done
```