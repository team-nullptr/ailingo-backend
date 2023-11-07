package domain

import "context"

const TaskStateDone = "DONE"

type Task struct {
	Id    int64  `json:"id"`
	State string `json:"finished"`
}

type TaskUseCase interface {
	Get(ctx context.Context, taskID int64) (*Task, error)
}

type TaskRepo interface {
	Get(ctx context.Context, taskID int64) (*Task, error)
	Insert(ctx context.Context) (int64, error)
	Complete(ctx context.Context, taskID int64) error
	Fail(ctx context.Context, taskID int64) error
}
