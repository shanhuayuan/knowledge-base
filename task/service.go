package task

import "context"

type CreateTaskOption struct {
	Type    string
	Payload any
}

type TaskService interface {
	Create(ctx context.Context, option CreateTaskOption) (*Task, error)
	Get(ctx context.Context, taskId string) (*Task, error)
	Cancel(ctx context.Context, taskId string) error
	Retry(ctx context.Context, taskId string) (*Task, error)
}
