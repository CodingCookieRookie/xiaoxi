package main

import (
	"bufio"
	"fmt"
	"os"
	"os/signal"
	"strconv"
	"strings"
	"syscall"

	"github.com/spf13/cobra"

	"task/internal/config"
	"task/internal/storage"
	"task/internal/task"
)

var tasksFile string

func main() {
	rootCmd := &cobra.Command{
		Use:   "task",
		Short: "A simple CLI task manager",
	}

	rootCmd.PersistentFlags().StringVarP(&tasksFile, "file", "f", "", "Path to tasks file")

	addCmd := newAddCmd()
	listCmd := newListCmd()
	doneCmd := newDoneCmd()
	deleteCmd := newDeleteCmd()
	configCmd := newConfigCmd()

	rootCmd.AddCommand(addCmd, listCmd, doneCmd, deleteCmd, configCmd)

	if len(os.Args) == 1 || (len(os.Args) == 3 && os.Args[1] == "-f") {
		runREPL(rootCmd)
	} else if err := rootCmd.Execute(); err != nil {
		fmt.Fprintln(os.Stderr, "Error:", err)
		os.Exit(1)
	}
}

func runREPL(rootCmd *cobra.Command) {
	fmt.Println("Task CLI - Interactive mode (type 'exit' to quit)")
	fmt.Println("Commands: add, list, done, delete, config, exit")
	fmt.Println()

	scanner := bufio.NewScanner(os.Stdin)

	sigChan := make(chan os.Signal, 1)
	signal.Notify(sigChan, syscall.SIGINT, syscall.SIGTERM)

	go func() {
		<-sigChan
		fmt.Println("\nGoodbye!")
		os.Exit(0)
	}()

	for {
		fmt.Print("> ")

		if !scanner.Scan() {
			break
		}

		input := strings.TrimSpace(scanner.Text())
		if input == "" {
			continue
		}

		if input == "exit" || input == "quit" || input == "q" {
			fmt.Println("Goodbye!")
			return
		}

		args := parseREPLInput(input)
		os.Args = append([]string{"task"}, args...)

		// Reset flag changed state by creating fresh command instances
		for _, cmd := range rootCmd.Commands() {
			if cmd.Name() == "list" {
				cmd.Flags().Set("all", "false")
				cmd.Flags().Set("completed", "false")
				cmd.Flags().Lookup("all").Changed = false
				cmd.Flags().Lookup("completed").Changed = false
			}
		}

		if err := rootCmd.Execute(); err != nil {
			fmt.Fprintln(os.Stderr, "Error:", err)
		}
	}
}

func parseREPLInput(input string) []string {
	var args []string
	var current strings.Builder
	inQuote := false

	for i := 0; i < len(input); i++ {
		ch := input[i]

		if ch == '"' {
			if inQuote {
				args = append(args, current.String())
				current.Reset()
				inQuote = false
			} else {
				inQuote = true
			}
		} else if ch == ' ' && !inQuote {
			if current.Len() > 0 {
				args = append(args, current.String())
				current.Reset()
			}
		} else {
			current.WriteByte(ch)
		}
	}

	if current.Len() > 0 {
		args = append(args, current.String())
	}

	return args
}

func resolveTasksFilePath() string {
	path, err := config.ResolveTasksFilePath(tasksFile)
	if err != nil {
		fmt.Fprintln(os.Stderr, "Error resolving tasks file path:", err)
		os.Exit(1)
	}
	return path
}

func loadTasks() *task.TaskList {
	path := resolveTasksFilePath()
	tl, err := storage.Load(path)
	if err != nil {
		fmt.Fprintln(os.Stderr, "Error loading tasks:", err)
		os.Exit(1)
	}
	return tl
}

func saveTasks(tl *task.TaskList) {
	path := resolveTasksFilePath()
	if err := storage.Save(path, tl); err != nil {
		fmt.Fprintln(os.Stderr, "Error saving tasks:", err)
		os.Exit(1)
	}
}

func newAddCmd() *cobra.Command {
	return &cobra.Command{
		Use:   "add [title]",
		Short: "Add a new task",
		Args:  cobra.MinimumNArgs(1),
		Run: func(cmd *cobra.Command, args []string) {
			title := strings.Join(args, " ")
			tl := loadTasks()
			t := tl.Add(title)
			saveTasks(tl)
			fmt.Printf("Added task: %s (ID: %d)\n", t.Title, t.ID)
		},
	}
}

func newListCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "list",
		Short: "List tasks",
		Run: func(cmd *cobra.Command, args []string) {
			tl := loadTasks()

			all := cmd.Flags().Changed("all")
			completed := cmd.Flags().Changed("completed")

			var tasks []task.Task
			switch {
			case completed:
				tasks = tl.ListCompleted()
			case all:
				tasks = tl.Tasks
			default:
				tasks = tl.ListPending()
			}

			if len(tasks) == 0 {
				fmt.Println("No tasks found.")
				return
			}

			for _, t := range tasks {
				status := "[ ]"
				if t.Completed {
					status = "[x]"
				}
				fmt.Printf("%s %d - %s\n", status, t.ID, t.Title)
			}
		},
	}

	cmd.Flags().Bool("all", false, "List all tasks")
	cmd.Flags().Bool("completed", false, "List completed tasks")

	return cmd
}

func newDoneCmd() *cobra.Command {
	return &cobra.Command{
		Use:   "done <task-id>",
		Short: "Mark a task as completed",
		Args:  cobra.ExactArgs(1),
		Run: func(cmd *cobra.Command, args []string) {
			id, err := strconv.Atoi(args[0])
			if err != nil {
				fmt.Fprintln(os.Stderr, "Invalid task ID:", args[0])
				os.Exit(1)
			}
			tl := loadTasks()

			t, found := tl.Complete(id)
			if !found {
				fmt.Fprintln(os.Stderr, "Task not found:", id)
				os.Exit(1)
			}

			saveTasks(tl)
			fmt.Printf("Completed task: %s\n", t.Title)
		},
	}
}

func newDeleteCmd() *cobra.Command {
	return &cobra.Command{
		Use:   "delete <task-id>",
		Short: "Delete a task",
		Args:  cobra.ExactArgs(1),
		Run: func(cmd *cobra.Command, args []string) {
			id, err := strconv.Atoi(args[0])
			if err != nil {
				fmt.Fprintln(os.Stderr, "Invalid task ID:", args[0])
				os.Exit(1)
			}
			tl := loadTasks()

			if !tl.Delete(id) {
				fmt.Fprintln(os.Stderr, "Task not found:", id)
				os.Exit(1)
			}

			saveTasks(tl)
			fmt.Printf("Deleted task: %d\n", id)
		},
	}
}

func newConfigCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "config",
		Short: "Manage configuration",
	}

	setFileCmd := &cobra.Command{
		Use:   "set-file <path>",
		Short: "Set the tasks file path",
		Args:  cobra.ExactArgs(1),
		Run: func(cmd *cobra.Command, args []string) {
			path := args[0]
			if err := config.SaveTasksFilePath(path); err != nil {
				fmt.Fprintln(os.Stderr, "Error saving config:", err)
				os.Exit(1)
			}
			fmt.Println("Tasks file path saved:", path)
		},
	}

	getFileCmd := &cobra.Command{
		Use:   "get-file",
		Short: "Get the tasks file path",
		Run: func(cmd *cobra.Command, args []string) {
			path := resolveTasksFilePath()
			fmt.Println(path)
		},
	}

	cmd.AddCommand(setFileCmd, getFileCmd)
	return cmd
}
