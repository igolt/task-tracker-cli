package main

import (
	"encoding/json"
	"fmt"
	"os"
	"strconv"
	"time"
)

const TasksFile = "tasks.json"

// Task status
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
		updateTaskCommand(cmdArgs)
	case "delete":
		deleteTaskCommand(cmdArgs)
	case "mark-in-progress":
		markCommand(InProgress, cmdArgs)
	case "mark-done":
		markCommand(Done, cmdArgs)
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
  add                 Create new a task
  list                List existing tasks
  update              Update an existing task
  delete              Delete a task
  mark-done           Mark a task as done
  mark-in-progress    Mark a task as in progress`)
}

func getTasks() (map[uint32]*Task, error) {
	fileContent, err := os.ReadFile(TasksFile)
	if err != nil {
		if os.IsNotExist(err) {
			return make(map[uint32]*Task), nil
		}
		return nil, err
	}

	var tasks map[uint32]*Task
	if err = json.Unmarshal(fileContent, &tasks); err != nil {
		return nil, err
	}
	return tasks, nil
}

func saveTasks(tasks map[uint32]*Task) error {
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
	tasks[taskId] = &newTask
	if err := saveTasks(tasks); err != nil {
		fmt.Fprintf(os.Stderr, "ERROR: \"task-cli add\": %v\n", err)
		os.Exit(1)
	}
}

func updateTaskCommand(cmdArgs []string) {
	if len(cmdArgs) != 2 {
		fmt.Fprintln(os.Stderr, "ERROR: \"task-cli update\" requires exactly two arguments")
		os.Exit(1)
	}

	taskIdStr, newDesc := cmdArgs[0], cmdArgs[1]

	taskId, err := parseTaskId(taskIdStr)
	if err != nil {
		fmt.Fprintf(os.Stderr, "ERROR: \"task-cli update\" %v", err)
		os.Exit(1)
	}

	tasks, err := getTasks()
	if err != nil {
		fmt.Fprintf(os.Stderr, "ERROR: \"task-cli update\" %v", err)
		os.Exit(1)
	}

	task, ok := tasks[taskId]
	if !ok {
		fmt.Fprintf(os.Stderr, "ERROR: \"task-cli update\": task with ID %d does not exists\n", taskId)
		os.Exit(1)
	}

	oldDesc := task.Description
	task.Description = newDesc
	task.UpdatedAt = time.Now()

	if err := saveTasks(tasks); err != nil {
		fmt.Fprintf(os.Stderr, "ERROR: \"task-cli update\": %v\n", err)
		os.Exit(1)
	}

	fmt.Printf("Updated: %#v -> %#v\n", oldDesc, newDesc)
}

func parseTaskId(s string) (uint32, error) {
	parsedStr, err := strconv.ParseUint(s, 10, 32)
	if err != nil {
		return 0, fmt.Errorf("invalid task ID %#v", s)
	}
	return uint32(parsedStr), nil
}

func deleteTaskCommand(cmdArgs []string) {
	if len(cmdArgs) != 1 {
		fmt.Fprintf(os.Stderr, "ERROR: \"task-cli delete\" requires exactly one argument")
		os.Exit(1)
	}

	taskId, err := parseTaskId(cmdArgs[0])
	if err != nil {
		fmt.Fprintf(os.Stderr, "ERROR: \"task-cli delete\": %v\n", err)
		os.Exit(1)
	}

	tasks, err := getTasks()
	if err != nil {
		fmt.Fprintf(os.Stderr, "ERROR: \"task-cli delete\": %v\n", err)
		os.Exit(1)
	}

	task, exists := tasks[taskId]
	if !exists {
		fmt.Fprintf(os.Stderr, "ERROR: \"task-cli delete\": task with ID %d does not exists\n", taskId)
		os.Exit(1)
	}
	fmt.Printf("%#v deleted\n", task.Description)
	delete(tasks, taskId)
	if err := saveTasks(tasks); err != nil {
		fmt.Fprintf(os.Stderr, "ERROR: \"task-cli delete\": %v\n", err)
		os.Exit(1)
	}
}

func markCommand(status string, cmdArgs []string) {
	if len(cmdArgs) != 1 {
		fmt.Fprintf(os.Stderr, "ERROR: \"task-cli mark-%s\": requires exactly one argument\n", status)
		os.Exit(1)
	}

	taskId, err := parseTaskId(cmdArgs[0])
	if err != nil {
		fmt.Fprintf(os.Stderr, "ERROR: \"task-cli mark-%s\": %v\n", status, err)
		os.Exit(1)
	}

	tasks, err := getTasks()
	if err != nil {
		fmt.Fprintf(os.Stderr, "ERROR: \"task-cli mark-%s\": %v\n", status, err)
		os.Exit(1)
	}

	task, ok := tasks[taskId]
	if !ok {
		fmt.Fprintf(os.Stderr, "ERROR: \"task-cli mark-%s\": task with ID %d does not exists\n", status, taskId)
		os.Exit(1)
	}

	if task.Status == status {
		fmt.Printf("%#v already %s\n", task.Description, status)
		return
	}

	task.Status = status
	if err := saveTasks(tasks); err != nil {
		fmt.Fprintf(os.Stderr, "ERROR: \"task-cli mark-%s\": %v\n", status, err)
		os.Exit(1)
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
			printTask(task)
		}
	} else {
		for _, task := range tasks {
			if task.Status == statusFilter {
				printTask(task)
			}
		}
	}
}

func printTask(task *Task) {
	fmt.Printf("%-10d%-15s%-10s\n", task.Id, task.Status, task.Description)
}
