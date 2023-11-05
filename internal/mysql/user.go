package mysql

import (
	"context"
	"database/sql"
	"fmt"

	"github.com/clerkinc/clerk-sdk-go/clerk"

	"ailingo/internal/domain"
)

type userQueries struct {
	getById *sql.Stmt
	insert  *sql.Stmt
	update  *sql.Stmt
	delete  *sql.Stmt
}

// getUserById queries for a user with the given id.
const getUserById = `
SELECT id, username, image_url
FROM user
WHERE id = ?
`

// insertUser inserts a new user.
const insertUser = `
INSERT IGNORE INTO user (id, username, image_url)
VALUES (?, ?, ?)
`

// updateUserById updates user by id.
const updateUserById = `
UPDATE user
SET username  = ?,
    image_url = ?
WHERE id = ?
`

// deleteUserById deletes user by id.
const deleteUserById = `
DELETE
FROM user
WHERE id = ?
`

type userRepo struct {
	db DBTX
}

func NewUserRepo(db DBTX) domain.UserRepo {
	return &userRepo{
		db: db,
	}
}

func (r *userRepo) GetById(ctx context.Context, userID string) (*domain.UserRow, error) {
	row := r.db.QueryRowContext(ctx, getUserById, userID)

	var user domain.UserRow
	if err := row.Scan(&user.Id, &user.Username, &user.ImageURL); err != nil {
		return nil, fmt.Errorf("failed to query: %w", err)
	}

	return &user, nil
}

func (r *userRepo) Insert(ctx context.Context, insertData *domain.InsertUserData) error {
	if _, err := r.db.ExecContext(ctx, insertUser, insertData.Id, insertData.Username, insertData.ImageURL); err != nil {
		return fmt.Errorf("failed to exec: %w", err)
	}
	return nil
}

func (r *userRepo) Update(ctx context.Context, updateData *domain.UpdateUserData) error {
	if _, err := r.db.ExecContext(ctx, updateUserById, updateData.Username, updateData.ImageURL, updateData.Id); err != nil {
		return fmt.Errorf("failed to exec: %w", err)
	}
	return nil
}

func (r *userRepo) Delete(ctx context.Context, userID string) error {
	if _, err := r.db.ExecContext(ctx, deleteUserById, userID); err != nil {
		return fmt.Errorf("failed to exec: %w", err)
	}
	return nil
}

func (r *userRepo) SyncUsers(ctx context.Context, users []clerk.User) error {
	stmt, err := r.db.PrepareContext(ctx, insertUser)
	if err != nil {
		return fmt.Errorf("failed to prepare insert stmt: %w", err)
	}

	for _, user := range users {
		if _, err := stmt.ExecContext(ctx, user.ID, user.Username, user.ImageURL); err != nil {
			return fmt.Errorf("failed to insert a user: %w", err)
		}
	}

	return nil
}
