package usecase

import (
	"context"
	"fmt"

	"github.com/go-playground/validator/v10"

	"ailingo/internal/domain"
)

type StudySetUseCase struct {
	studySetRepo domain.StudySetRepo
	validate     *validator.Validate
}

// NewUseCase creates a new instance of StudySetUseCaseImpl.
func NewStudySetUseCase(studySetRepo domain.StudySetRepo, validate *validator.Validate) domain.StudySetUseCase {
	return &StudySetUseCase{
		studySetRepo: studySetRepo,
		validate:     validate,
	}
}

func (uc *StudySetUseCase) Create(ctx context.Context, insertData *domain.InsertStudySetData) (*domain.StudySet, error) {
	if err := uc.validate.Struct(insertData); err != nil {
		return nil, fmt.Errorf("%w: %w", ErrValidation, err)
	}

	studySet, err := uc.studySetRepo.Insert(ctx, insertData)
	if err != nil {
		return nil, fmt.Errorf("%w: failed to create the study set: %w", ErrRepoFailed, err)
	}

	return studySet, nil
}

func (uc *StudySetUseCase) GetById(ctx context.Context, studySetID int64) (*domain.StudySet, error) {
	studySet, err := uc.studySetRepo.GetById(ctx, studySetID)
	if err != nil {
		return nil, fmt.Errorf("%w: failed to get the study set: %w", ErrRepoFailed, err)
	}
	if studySet == nil {
		return nil, ErrNotFound
	}

	return studySet, nil
}

func (uc *StudySetUseCase) GetAllSummary(ctx context.Context) ([]*domain.StudySetSummary, error) {
	studySets, err := uc.studySetRepo.GetAllSummary(ctx)
	if err != nil {
		return nil, fmt.Errorf("%w: failed to get study sets: %w", ErrRepoFailed, err)
	}

	return studySets, nil
}

func (uc *StudySetUseCase) Update(ctx context.Context, studySetID int64, userID string, updateData *domain.UpdateStudySetData) error {
	targetStudySet, err := uc.studySetRepo.GetById(ctx, studySetID)
	if err != nil {
		return fmt.Errorf("%w: GetById failed: %w", ErrRepoFailed, err)
	}
	if targetStudySet == nil {
		return ErrNotFound
	}
	if targetStudySet.AuthorId != userID {
		return ErrForbidden
	}

	if err := uc.studySetRepo.Update(ctx, studySetID, updateData); err != nil {
		return fmt.Errorf("%w: Update failed: %w", ErrRepoFailed, err)
	}

	return nil
}

func (uc *StudySetUseCase) Delete(ctx context.Context, studySetID int64, userID string) error {
	targetStudySet, err := uc.studySetRepo.GetById(ctx, studySetID)
	if err != nil {
		return fmt.Errorf("%w: GetById failed: %w", ErrRepoFailed, err)
	}
	if targetStudySet == nil {
		return ErrNotFound
	}
	if targetStudySet.AuthorId != userID {
		return ErrForbidden
	}

	if err := uc.studySetRepo.Delete(ctx, studySetID); err != nil {
		return fmt.Errorf("%w: Delete failed: %w", ErrRepoFailed, err)
	}

	return nil
}
