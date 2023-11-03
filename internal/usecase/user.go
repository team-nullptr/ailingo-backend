package usecase

import (
	"context"
	"fmt"

	"ailingo/internal/domain"
)

type UserUseCase struct {
	userRepo domain.UserRepo
}

func NewUserUseCase(userRepo domain.UserRepo) *UserUseCase {
	return &UserUseCase{
		userRepo: userRepo,
	}
}

func (uc *UserUseCase) Insert(ctx context.Context, insertData *domain.InsertUserData) error {
	if err := uc.userRepo.Insert(ctx, insertData); err != nil {
		return fmt.Errorf("%w: failed to insert the user: %w", err)
	}

	return nil
}

func (uc *UserUseCase) Update(ctx context.Context, updateData *domain.UpdateUserData) error {
	if err := uc.userRepo.Update(ctx, updateData); err != nil {
		return fmt.Errorf("%w: failed to update the user: %w", err)
	}

	return nil
}

func (uc *UserUseCase) Delete(ctx context.Context, userID string) error {
	if err := uc.userRepo.Delete(ctx, userID); err != nil {
		return fmt.Errorf("%w: failed to delete the user: %w", err)
	}

	return nil
}
