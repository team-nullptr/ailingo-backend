package mysql

import (
	"context"
	"database/sql"
	"errors"
	"fmt"

	"ailingo/internal/domain"
)

const getTask = `
SELECT task.id, task.state
FROM task
WHERE id = ?
`

const insertTask = `
INSERT INTO task () VALUES ()
`

const failTask = `
UPDATE task SET state = 'FAILED' WHERE id = ?
`

const completeTask = `
UPDATE task SET state = 'DONE' WHERE id = ?
`

type taskRepo struct {
	db DBTX
}

func NewTaskRepo(db DBTX) domain.TaskRepo {
	return &taskRepo{
		db: db,
	}
}

func (r *taskRepo) Get(ctx context.Context, taskID int64) (*domain.Task, error) {
	var task domain.Task
	if err := r.db.QueryRowContext(ctx, getTask, taskID).Scan(&task.Id, &task.State); err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, nil
		}
		return nil, fmt.Errorf("failed to scan: %w", err)
	}
	return &task, nil
}

func (r *taskRepo) Insert(ctx context.Context) (int64, error) {
	res, err := r.db.ExecContext(ctx, insertTask)
	if err != nil {
		return 0, fmt.Errorf("failed to exec: %w", err)
	}

	lastInsertedId, err := res.LastInsertId()
	if err != nil {
		return 0, fmt.Errorf("failed to get last inserted id: %w", err)
	}

	return lastInsertedId, nil
}

func (r *taskRepo) Complete(ctx context.Context, taskID int64) error {
	if _, err := r.db.ExecContext(ctx, completeTask, taskID); err != nil {
		return fmt.Errorf("failed to exec query: %w", err)
	}
	return nil
}

func (r *taskRepo) Fail(ctx context.Context, taskID int64) error {
	if _, err := r.db.ExecContext(ctx, failTask, taskID); err != nil {
		return fmt.Errorf("failed to exec query: %w", err)
	}
	return nil
}
