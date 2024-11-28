package main

import (
	"encoding/json"
	"fmt"
	"os"
	"time"
)

const TasksFile = "tasks.json"

const (
	Todo       = "todo"
	InProgress = "in-progress"
	Done       = "done"
)

type Task struct {
	Id          uint32    `json:"id"`
	Description string    `json:"description"`
	Status      string    `json:"status"`
	CreatedAt   time.Time `json:"createdAt"`
	UpdatedAt   time.Time `json:"updatedAt"`
}

func main() {
	args := os.Args

	if len(args) == 1 {
		printUsage()
		return
	}

	cmdArgs := args[2:]
	switch command := args[1]; command {
	case "-h", "--help":
		printUsage()
	case "add":
		addTaskCommand(cmdArgs)
	case "update":
		fmt.Println("update")
	case "delete":
		fmt.Println("delete")
	case "mark-in-progress":
		fmt.Println("mark-in-progress")
	case "mark-done":
		fmt.Println("mark-in-progress")
	case "list":
		listTasksCommand(cmdArgs)
	default:
		fmt.Fprintf(os.Stderr, "task-cli: invalid command %#v\n", command)
		os.Exit(1)
	}
}

func printUsage() {
	fmt.Println(`
Usage: task-cli [OPTIONS] COMMAND [ARG...]

A command line task manager

Commands:

add         Create new a task
update      Update an existing task
delete      Delete a task
list        List existing tasks`)
}

func getTasks() (map[uint32]Task, error) {
	fileContent, err := os.ReadFile(TasksFile)
	if err != nil {
		if os.IsNotExist(err) {
			return make(map[uint32]Task), nil
		}
		return nil, err
	}

	var tasks map[uint32]Task
	err = json.Unmarshal(fileContent, &tasks)
	if err != nil {
		return nil, err
	}
	return tasks, nil
}

func saveTasks(tasks map[uint32]Task) error {
	tasksJson, err := json.Marshal(tasks)
	if err != nil {
		return err
	}
	return os.WriteFile(TasksFile, tasksJson, 0644)
}

func addTaskCommand(cmdArgs []string) {
	if len(cmdArgs) != 1 {
		fmt.Fprintln(os.Stderr, "ERROR: \"task-cli add\" requires exactly one argument")
		os.Exit(1)
	}

	taskDesc := cmdArgs[0]
	tasks, err := getTasks()
	if err != nil {
		fmt.Fprintf(os.Stderr, "ERROR: \"task-cli add\": %v\n", err)
		os.Exit(1)
	}

	var maxId uint32 = 0
	for id := range tasks {
		if id > maxId {
			maxId = id
		}
	}
	taskId := maxId + 1
	now := time.Now()
	newTask := Task{Id: taskId, Description: taskDesc, Status: Todo, CreatedAt: now, UpdatedAt: now}
	tasks[taskId] = newTask
	if err := saveTasks(tasks); err != nil {
		fmt.Fprintf(os.Stderr, "ERROR: \"task-cli add\": %v\n", err)
	}
}

func listTasksCommand(cmdArgs []string) {
	statusFilter := ""
	if len(cmdArgs) > 1 {
		fmt.Fprintln(os.Stderr, "ERROR: \"task-cli list\" accepts zero or one argument")
		os.Exit(1)
	}

	if len(cmdArgs) == 1 {
		switch statusFilter = cmdArgs[0]; statusFilter {
		case Todo, InProgress, Done:
		default:
			fmt.Fprintf(os.Stderr, "ERROR: \"task-cli list\": invalid filter %#v\n", statusFilter)
			os.Exit(1)
		}
	}

	tasks, err := getTasks()
	if err != nil {
		fmt.Fprintf(os.Stderr, "ERROR: \"task-cli list\": %v\n", err)
		os.Exit(1)
	}

	fmt.Printf("%-10s%-15s%-10s\n", "ID", "STATUS", "DESCRIPTION")
	if statusFilter == "" {
		for _, task := range tasks {
			printTask(&task)
		}
	} else {
		for _, task := range tasks {
			if task.Status == statusFilter {
				printTask(&task)
			}
		}
	}
}

func printTask(task *Task) {
	fmt.Printf("%-10d%-15s%-10s\n", task.Id, task.Status, task.Description)
}
