# xiaoxi - CLI Task Manager

A simple, fast CLI task manager for everyday tasks. Built with Go.

## Features

- Add, list, complete, and delete tasks
- Sequential integer IDs (0, 1, 2, ...)
- Interactive REPL mode or single commands
- JSON-based storage with safe file writes
- Configurable task file path
- Path resolution priority: CLI flag → env var → saved config → default

## Installation

```bash
go build -o xiaoxi ./cmd/main.go
```

Optionally, move it to your PATH:

```bash
mv xiaoxi /usr/local/bin/
```

## Usage

### Interactive REPL Mode

```bash
./xiaoxi
```

```
Task CLI - Interactive mode (type 'exit' to quit)
Commands: add, list, done, delete, config, exit

> add "Buy groceries"
Added task: Buy groceries (ID: 0)
> add "Walk the dog"
Added task: Walk the dog (ID: 1)
> list
[ ] 0 - Buy groceries
[ ] 1 - Walk the dog
> done 1
Completed task: Walk the dog
> list --all
[ ] 0 - Buy groceries
[x] 1 - Walk the dog
> exit
Goodbye!
```

### Single Commands

```bash
# Add a task
./xiaoxi add "Buy groceries"

# List pending tasks
./xiaoxi list

# List all tasks
./xiaoxi list --all

# List completed tasks
./xiaoxi list --completed

# Mark task as completed (keeps the record)
./xiaoxi done 0

# Delete a task permanently
./xiaoxi delete 1
```

### Options

```bash
-f, --file string   Path to tasks file
```

## Task File Location

Tasks are stored in a JSON file. The path is resolved in this order:

1. CLI flag: `-f /path/to/tasks.json`
2. Environment variable: `TASKS_FILE=/path/to/tasks.json`
3. Saved config: `~/.config/.task-config`
4. Default: `~/tasks.json`

### Config Commands

```bash
# Set default tasks file
./xiaoxi config set-file /path/to/tasks.json

# Get current tasks file path
./xiaoxi config get-file
```

## Data Format

Tasks are stored in JSON format:

```json
{
  "tasks": [
    {
      "id": 0,
      "title": "Buy groceries",
      "completed": false,
      "created_at": "2024-01-01T00:00:00Z"
    }
  ],
  "next_id": 1
}
```

## Development

```bash
# Run tests
go test ./...

# Build
go build -o xiaoxi ./cmd/main.go
```

## Architecture

```
cmd/main.go           # CLI entry point and command handlers
internal/
  config/config.go    # Path resolution
  storage/storage.go   # Safe file I/O
  task/task.go         # Task model and operations
```

## License

MIT