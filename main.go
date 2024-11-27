package main

import (
	"encoding/json"
	"fmt"
	"os"
	"time"
)

const TasksFile = "tasks.json"

type TaskStatus string;

const (
	NotStarted TaskStatus = "not-started"
	InProgress = "in-progress"
	Done = "done"
)

type Task struct {
	Id uint32 `json:"id"`;
	Description string `json:"description"`;
	Status TaskStatus `json:"status"`;
	CreatedAt time.Time `json:"createdAt"`;
	UpdatedAt time.Time `json:"updatedAt"`;
}

func main() {
	args := os.Args

	if len(args) == 1 {
		printUsage()
		return
	}

	switch command := args[1]; command {
	case "-h", "--help":
		printUsage()
	case "add":
		addTaskCommand(args[:1])
	case "update":
		fmt.Println("update")
	case "delete":
		fmt.Println("delete")
	case "mark-in-progress":
		fmt.Println("mark-in-progress")
	case "list":
		fmt.Println("list")
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
list        List existing tasks
`)
}

func addTaskCommand(cmdArgs []string) {
	if len(cmdArgs) != 1 {
		fmt.Fprintln(os.Stderr,cmdArgs)
		fmt.Fprintln(os.Stderr, "ERROR: \"task-cli add\" requires exactly one argument")
		os.Exit(1)
	}
	taskDesc := cmdArgs[0]
	tasks, err := getTasks()
	if err != nil {
		fmt.Fprintf(os.Stderr, "ERROR: \"task-cli add\": %v", err)
		os.Exit(1)
	}

	var maxId uint32 = 0
	for id, _ := range tasks {
		if id > maxId {
			maxId = id
		}
	}
	taskId := maxId + 1
	createdAt := time.Now()
	newTask := Task{Id: taskId, Description: taskDesc, Status:
		NotStarted, CreatedAt: createdAt, UpdatedAt: createdAt}
	tasks[taskId] = newTask
	if err := saveTasks(tasks); err != nil {
		fmt.Fprintf(os.Stderr, "ERROR: \"task-cli add\": %v", err)
	}
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
