package usecase

import (
	"context"
	"fmt"

	"ailingo/internal/domain"
)

type UserUseCase struct {
	dataStore domain.DataStore
}

func NewUserUseCase(dataStore domain.DataStore) *UserUseCase {
	return &UserUseCase{
		dataStore: dataStore,
	}
}

func (uc *UserUseCase) Insert(ctx context.Context, insertData *domain.InsertUserData) error {
	userRepo := uc.dataStore.GetUserRepo()

	if err := userRepo.Insert(ctx, insertData); err != nil {
		return fmt.Errorf("%w: failed to insert the user: %w", err)
	}

	return nil
}

func (uc *UserUseCase) Update(ctx context.Context, updateData *domain.UpdateUserData) error {
	userRepo := uc.dataStore.GetUserRepo()

	if err := userRepo.Update(ctx, updateData); err != nil {
		return fmt.Errorf("%w: failed to update the user: %w", err)
	}

	return nil
}

func (uc *UserUseCase) Delete(ctx context.Context, userID string) error {
	userRepo := uc.dataStore.GetUserRepo()

	if err := userRepo.Delete(ctx, userID); err != nil {
		return fmt.Errorf("%w: failed to delete the user: %w", err)
	}

	return nil
}
