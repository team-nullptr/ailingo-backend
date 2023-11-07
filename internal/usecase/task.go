package usecase

import (
	"context"
	"fmt"

	"ailingo/internal/domain"
)

type taskUseCase struct {
	dataStore domain.DataStore
}

func NewTaskUseCase(dataStore domain.DataStore) domain.TaskUseCase {
	return &taskUseCase{
		dataStore: dataStore,
	}
}

func (uc *taskUseCase) Get(ctx context.Context, taskID int64) (*domain.Task, error) {
	task, err := uc.dataStore.GetTaskRepo().Get(ctx, taskID)
	if err != nil {
		return nil, fmt.Errorf("%w: failed to get the task: %w", ErrRepoFailed, err)
	}
	if task == nil {
		return nil, &ErrNotFound{
			Resource: TaskResource,
		}
	}
	return task, nil
}
