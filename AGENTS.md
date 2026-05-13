# AGENTS.md

## Role

Act as a super senior Golang backend engineer building a CLI-based everyday task scheduler.

Prioritize correctness, simplicity, maintainability, testability, and a good developer experience. This is a local-first CLI tool where users can add, list, complete, and delete daily tasks. Tasks are persisted in a user-configurable file.

## Project Goal

Build a command-line task scheduler that allows users to manage everyday tasks from the terminal.

The CLI should support:

- Adding tasks
- Listing tasks
- Marking tasks as completed
- Deleting tasks
- Storing tasks in a configured local file
- Loading existing tasks from that file on startup
- Saving updates safely after every mutation

The design should be simple, predictable, and easy to extend.

## Core Engineering Principles

- Keep the CLI behavior clear and consistent.
- Prefer simple data structures unless complexity is justified.
- Keep file storage logic separate from CLI command handling.
- Avoid global state where possible.
- Validate user input before modifying the task file.
- Make task mutations idempotent where possible.
- Provide helpful error messages.
- Do not silently ignore failures.
- Do not overwrite user data unexpectedly.
- Write tests for all core task operations.

## Architecture Expectations

Structure the project around clear responsibilities.

Recommended modules:

```txt
cmd/
  root command and CLI wiring

internal/task/
  task model and task operations

internal/storage/
  file loading, saving, and validation

internal/config/
  config path resolution and file selection