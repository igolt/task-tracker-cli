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
	Id int32 `json:"id"`;
	Description string `json:"description"`;
	Status TaskStatus `json:"status"`;
	CreatedAt time.Time `json:"createdAt"`;
	UpdatedAt time.Time `json:"updatedAt"`;
}

func main() {
	args := os.Args

	if len(args) == 1 {
		fmt.Fprintln(os.Stderr, "command expected")
		os.Exit(1)
	}

	command := args[1]

	switch command {
	case "add":
	fmt.Println("add")
	case "update":
	fmt.Println("update")
	case "delete":
	fmt.Println("delete")
	case "mark-in-progress":
	fmt.Println("mark-in-progress")
	case "list":
	fmt.Println("list")
	}
}

func getTasks() (map[int]Task, error) {
	fileContent, err := os.ReadFile(TasksFile)

	if err != nil {
		return nil, err
	}

	var tasks map[int]Task

	err = json.Unmarshal(fileContent, &tasks)
	if err != nil {
		return nil, err
	}
	return tasks, nil
}
