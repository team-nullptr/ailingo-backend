package mysql

import (
	"context"
	"database/sql"
	"fmt"

	"ailingo/internal/domain"
)

type userQueries struct {
	getById *sql.Stmt
	insert  *sql.Stmt
	update  *sql.Stmt
	delete  *sql.Stmt
}

func newUserQueries(ctx context.Context, db *sql.DB) (*userQueries, error) {
	getByIdStmt, err := db.PrepareContext(ctx, "SELECT id, username, image_url FROM user WHERE id = ?")
	if err != nil {
		return nil, fmt.Errorf("getById query: %w", err)
	}

	insertStmt, err := db.PrepareContext(ctx, "INSERT INTO user (id, username, image_url) VALUES (?, ?, ?)")
	if err != nil {
		return nil, fmt.Errorf("insert query: %w", err)
	}

	updateStmt, err := db.PrepareContext(ctx, "UPDATE user SET username = ?, image_url = ? WHERE id = ?")
	if err != nil {
		return nil, fmt.Errorf("update query: %w", err)
	}

	deleteStmt, err := db.PrepareContext(ctx, "DELETE FROM user WHERE id = ?")
	if err != nil {
		return nil, fmt.Errorf("delete query: %w", err)
	}

	return &userQueries{
		getById: getByIdStmt,
		insert:  insertStmt,
		update:  updateStmt,
		delete:  deleteStmt,
	}, nil
}

type UserRepo struct {
	query *userQueries
}

func NewUserRepo(ctx context.Context, db *sql.DB) (*UserRepo, error) {
	query, err := newUserQueries(ctx, db)
	if err != nil {
		return nil, err
	}

	return &UserRepo{
		query: query,
	}, nil
}

func (r *UserRepo) GetById(ctx context.Context, userID string) (*domain.UserRow, error) {
	row := r.query.getById.QueryRowContext(ctx, userID)

	var user domain.UserRow
	if err := row.Scan(&user.Id, &user.Username, &user.ImageURL); err != nil {
		return nil, fmt.Errorf("failed to query: %w", err)
	}

	return &user, nil
}

func (r *UserRepo) Insert(ctx context.Context, insertData *domain.InsertUserData) error {
	if _, err := r.query.insert.ExecContext(ctx, insertData.Id, insertData.Username, insertData.ImageURL); err != nil {
		return fmt.Errorf("failed to exec: %w", err)
	}
	return nil
}

func (r *UserRepo) Update(ctx context.Context, updateData *domain.UpdateUserData) error {
	if _, err := r.query.update.ExecContext(ctx, updateData.Username, updateData.ImageURL, updateData.Id); err != nil {
		return fmt.Errorf("failed to exec: %w", err)
	}
	return nil
}

func (r *UserRepo) Delete(ctx context.Context, userID string) error {
	if _, err := r.query.delete.ExecContext(ctx, userID); err != nil {
		return fmt.Errorf("failed to exec: %w", err)
	}
	return nil
}
