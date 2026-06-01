package storage

import (
	"context"
	"time"
)

type Task struct {
	ID          string                 `json:"id"`
	Type        string                 `json:"type"`
	Payload     map[string]interface{} `json:"payload"`
	Priority    int                    `json:"priority"`
	Status      string                 `json:"status"`
	RetryCount  int                    `json:"retry_count"`
	MaxRetries  int                    `json:"max_retries"`
	CreatedAt   time.Time              `json:"created_at"`
	StartedAt   *time.Time             `json:"started_at"`
	CompletedAt *time.Time             `json:"completed_at"`
}

type Storage interface {
	CreateTask(ctx context.Context, task Task) error
	GetTask(ctx context.Context, userId string) (error, Task)
	GetTaskStatus(ctx context.Context, taskId string) (string, error)
	GetAllTasks(ctx context.Context) ([]Task, error)
	DeleteTask(ctx context.Context, taskId string) error
}
