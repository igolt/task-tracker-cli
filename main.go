package main

import (
	"encoding/json"
	"errors"
	"fmt"
	"os"
	"strconv"
	"time"
)

type (
	TaskId   uint32
	TasksMap map[TaskId]*Task
)

type Task struct {
	Id          TaskId    `json:"id"`
	Description string    `json:"description"`
	Status      string    `json:"status"`
	CreatedAt   time.Time `json:"createdAt"`
	UpdatedAt   time.Time `json:"updatedAt"`
}

// Task status
const (
	Todo       = "todo"
	InProgress = "in-progress"
	Done       = "done"
)

const TasksFile = "tasks.json"

func main() {
	args := os.Args

	if len(args) == 1 {
		printUsage()
		return
	}

	var cmd func([]string) error
	command := args[1]
	switch command {
	case "-h", "--help":
		printUsage()
		return
	case "add":
		cmd = addTaskCommand
	case "update":
		cmd = updateTaskCommand
	case "delete":
		cmd = deleteTaskCommand
	case "mark-in-progress":
		cmd = markInProgressCommand
	case "mark-done":
		cmd = markDoneCommand
	case "list":
		cmd = listTasksCommand
	default:
		fmt.Fprintf(os.Stderr, "task-cli: invalid command %#v\n", command)
		os.Exit(1)
	}

	cmdArgs := args[2:]
	if err := cmd(cmdArgs); err != nil {
		fmt.Fprintf(os.Stderr, "ERROR: \"task-cli %s\": %v\n", command, err)
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

func loadTasks() (TasksMap, error) {
	fileContent, err := os.ReadFile(TasksFile)
	if err != nil {
		if os.IsNotExist(err) {
			return make(TasksMap), nil
		}
		return nil, err
	}

	var tasks TasksMap
	if err = json.Unmarshal(fileContent, &tasks); err != nil {
		return nil, err
	}
	return tasks, nil
}

func saveTasks(tasks TasksMap) error {
	tasksJson, err := json.Marshal(tasks)
	if err != nil {
		return err
	}
	return os.WriteFile(TasksFile, tasksJson, 0644)
}

func addTaskCommand(cmdArgs []string) error {
	if len(cmdArgs) != 1 {
		return requiresExactlyOneArgumentError()
	}

	taskDesc := cmdArgs[0]
	tasks, err := loadTasks()
	if err != nil {
		return err
	}

	var maxId TaskId = 0
	for id := range tasks {
		if id > maxId {
			maxId = id
		}
	}
	taskId := maxId + 1
	now := time.Now()
	newTask := Task{Id: taskId, Description: taskDesc, Status: Todo, CreatedAt: now, UpdatedAt: now}
	tasks[taskId] = &newTask
	return saveTasks(tasks)
}

func requiresExactlyOneArgumentError() error {
	return errors.New("exactly one argument is required")
}

func updateTaskCommand(cmdArgs []string) error {
	if len(cmdArgs) != 2 {
		return errors.New("exactly two arguments are required")
	}

	taskIdStr, newDesc := cmdArgs[0], cmdArgs[1]

	taskId, err := parseTaskId(taskIdStr)
	if err != nil {
		return err
	}

	tasks, err := loadTasks()
	if err != nil {
		return err
	}

	task, err := getTask(tasks, taskId)
	if err != nil {
		return err
	}

	oldDesc := task.Description
	task.Description = newDesc
	task.UpdatedAt = time.Now()

	if err := saveTasks(tasks); err != nil {
		return err
	}

	fmt.Printf("Updated: %#v -> %#v\n", oldDesc, newDesc)
	return nil
}

func getTask(tasks TasksMap, taskId TaskId) (*Task, error) {
	task, ok := tasks[taskId]
	if ok {
		return task, nil
	}
	return nil, fmt.Errorf("task with ID %d does not exists", taskId)
}

func parseTaskId(s string) (TaskId, error) {
	parsedStr, err := strconv.ParseUint(s, 10, 32)
	if err != nil {
		return 0, fmt.Errorf("invalid task ID %#v", s)
	}
	return TaskId(parsedStr), nil
}

func deleteTaskCommand(cmdArgs []string) error {
	if len(cmdArgs) != 1 {
		return requiresExactlyOneArgumentError()
	}

	taskId, err := parseTaskId(cmdArgs[0])
	if err != nil {
		return err
	}

	tasks, err := loadTasks()
	if err != nil {
		return err
	}

	task, err := getTask(tasks, taskId)
	if err != nil {
		return err
	}

	fmt.Printf("%#v deleted\n", task.Description)
	delete(tasks, taskId)
	return saveTasks(tasks)
}

func markCommand(status string, cmdArgs []string) error {
	if len(cmdArgs) != 1 {
		return requiresExactlyOneArgumentError()
	}

	taskId, err := parseTaskId(cmdArgs[0])
	if err != nil {
		return err
	}

	tasks, err := loadTasks()
	if err != nil {
		return err
	}

	task, err := getTask(tasks, taskId)
	if err != nil {
		return err
	}

	if task.Status == status {
		if status == InProgress {
			status = "in progress"
		}
		fmt.Printf("%#v already %s\n", task.Description, status)
		return nil
	}
	task.Status = status
	return saveTasks(tasks)
}

func markInProgressCommand(cmdArgs []string) error {
	return markCommand(InProgress, cmdArgs)
}

func markDoneCommand(cmdArgs []string) error {
	return markCommand(Done, cmdArgs)
}

func listTasksCommand(cmdArgs []string) error {
	statusFilter := ""
	if len(cmdArgs) > 1 {
		return errors.New("accepts zero or one argument")
	}

	if len(cmdArgs) == 1 {
		switch statusFilter = cmdArgs[0]; statusFilter {
		case Todo, InProgress, Done:
		default:
			return fmt.Errorf("invalid filter %#v", statusFilter)
		}
	}

	tasks, err := loadTasks()
	if err != nil {
		return err
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
	return nil
}

func printTask(task *Task) {
	fmt.Printf("%-10d%-15s%-10s\n", task.Id, task.Status, task.Description)
}
