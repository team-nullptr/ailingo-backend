package domain

import "context"

type UserRow struct {
	Id       string
	Username string
	ImageURL string
}

type InsertUserData struct {
	Id       string
	Username string
	ImageURL string
}

type UpdateUserData InsertUserData

type UserUseCase interface {
	Insert(ctx context.Context, insertData *InsertUserData) error
	Update(ctx context.Context, updateData *UpdateUserData) error
	Delete(ctx context.Context, userID string) error
}

type UserRepo interface {
	GetById(ctx context.Context, userID string) (*UserRow, error)
	Insert(ctx context.Context, insertData *InsertUserData) error
	Update(ctx context.Context, updateData *UpdateUserData) error
	Delete(ctx context.Context, userID string) error
}
