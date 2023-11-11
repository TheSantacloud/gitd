package adapters

// TODO: rename this module into a more fitting name

import (
	"time"
)

type TaskManagerAdapter interface {
	Initialize(Settings) error
	FetchTasks() ([]Task, error)
	UpdateTasks(*[]TaskAction) error
}

type ArchiverAdapter interface {
	Initialize(Settings) error
}

type Priority int8
type Status int8
type Action int8

const (
	// Priority
	PriorityCritical Priority = 1
	PriorityHigh     Priority = 2
	PriorityMedium   Priority = 3
	PriorityLow      Priority = 4

	// Status
	StatusActive    Status = 1
	StatusCompleted Status = 2
	StatusDeleted   Status = 3
	StatusNext      Status = 4
	StatusSomeday   Status = 5

	// Actions
	ActionIgnore     Action = 0
	ActionComplete   Action = 1
	ActionDelete     Action = 2
	ActionRevalidate Action = 3
	ActionDefer      Action = 4
)

type Task struct {
	ID          string    `json:"id"`
	Project     string    `json:"project"`
	Content     string    `json:"content"`
	CreatedDate time.Time `json:"created_at"`
	UpdatedDate time.Time `json:"modified_at"`
	Tags        []string  `json:"tags"`
	Status      Status    `json:"status"`
	Priority    Priority  `json:"priority"`
	TaskManger  string    `json:"taskmanager"`
}

type FilterRequest struct {
	Project       *string
	Status        *Status
	AfterTimeSpan *TimeSpan
	Tags          *[]string
}

type TaskAction struct {
	Task   *Task
	Action Action
}
