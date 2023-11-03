package usecase

import (
	"context"
	"fmt"

	"github.com/go-playground/validator/v10"

	"ailingo/internal/domain"
	"ailingo/pkg/auth"
)

type StudySetUseCase struct {
	studySetRepo domain.StudySetRepo
	userService  *auth.UserService
	validate     *validator.Validate
}

// NewStudySetUseCase creates a new instance of StudySetUseCaseImpl.
func NewStudySetUseCase(studySetRepo domain.StudySetRepo, userService *auth.UserService, validate *validator.Validate) domain.StudySetUseCase {
	return &StudySetUseCase{
		studySetRepo: studySetRepo,
		userService:  userService,
		validate:     validate,
	}
}

func (uc *StudySetUseCase) GetAll(ctx context.Context) ([]*domain.StudySetWithAuthor, error) {
	studySets, err := uc.studySetRepo.GetAll(ctx)
	if err != nil {
		return nil, fmt.Errorf("%w: failed to get study sets: %w", ErrRepoFailed, err)
	}
	return studySets, nil
}

func (uc *StudySetUseCase) GetById(ctx context.Context, studySetID int64) (*domain.StudySetWithAuthor, error) {
	studySet, err := uc.studySetRepo.GetById(ctx, studySetID)
	if err != nil {
		return nil, fmt.Errorf("%w: failed to get the study set: %w", ErrRepoFailed, err)
	}
	if studySet == nil {
		return nil, ErrNotFound
	}
	return studySet, nil
}

func (uc *StudySetUseCase) Create(ctx context.Context, insertData *domain.InsertStudySetData) (int64, error) {
	if err := uc.validate.Struct(insertData); err != nil {
		return 0, fmt.Errorf("%w: invalid insert data: %w", ErrValidation, err)
	}
	insertedId, err := uc.studySetRepo.Insert(ctx, insertData)
	if err != nil {
		return 0, fmt.Errorf("%w: failed to create the study set: %w", ErrRepoFailed, err)
	}
	return insertedId, nil
}

func (uc *StudySetUseCase) Update(ctx context.Context, userID string, studySetID int64, updateData *domain.UpdateStudySetData) error {
	if err := uc.validate.Struct(updateData); err != nil {
		return fmt.Errorf("%w: invalid update data: %w", ErrValidation, err)
	}
	if err := uc.checkStudySetOwnership(ctx, userID, studySetID); err != nil {
		return err
	}
	if err := uc.studySetRepo.Update(ctx, studySetID, updateData); err != nil {
		return fmt.Errorf("%w: Update failed: %w", ErrRepoFailed, err)
	}
	return nil
}

func (uc *StudySetUseCase) Delete(ctx context.Context, userID string, studySetID int64) error {
	if err := uc.checkStudySetOwnership(ctx, userID, studySetID); err != nil {
		return err
	}
	if err := uc.studySetRepo.Delete(ctx, studySetID); err != nil {
		return fmt.Errorf("%w: Delete failed: %w", ErrRepoFailed, err)
	}
	return nil
}

func (uc *StudySetUseCase) checkStudySetOwnership(ctx context.Context, userID string, studySetID int64) error {
	studySet, err := uc.studySetRepo.GetById(ctx, studySetID)
	if err != nil {
		return fmt.Errorf("%w: failed to get the study set: %w", ErrRepoFailed, err)
	}
	if studySet == nil {
		return ErrNotFound
	}
	if studySet.Author.Id != userID {
		return ErrForbidden
	}
	return nil
}
