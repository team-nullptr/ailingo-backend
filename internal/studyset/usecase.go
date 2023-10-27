package studyset

import (
	"context"
	"errors"
	"fmt"

	"github.com/go-playground/validator/v10"

	"ailingo/internal/models"
)

var (
	ErrForbidden = errors.New("not enough permission to perform the operation")
	ErrNotFound  = errors.New("study set not found")
)

type UseCase interface {
	// GetAllSummary gets summaries of all the created study sets.
	GetAllSummary(ctx context.Context) ([]*models.StudySetSummary, error)
	// GetById gets a study set by its id.
	GetById(ctx context.Context, studySetID int64) (*models.StudySet, error)
	// Create creates a new study set. It is responsible for data validation.
	Create(ctx context.Context, createData *insertStudySetData) (*models.StudySet, error)
	// Update updates the study set with the given id.
	Update(ctx context.Context, studySetID int64, userID string, updateData *updateStudySetData) error
	// Delete deletes the given study set.
	Delete(ctx context.Context, studySetID int64, userID string) error
}

type DefaultUseCase struct {
	studySetRepo Repo
	validate     *validator.Validate
}

// NewUseCase creates a new instance of StudySetUseCaseImpl.
func NewUseCase(studySetRepo Repo, validate *validator.Validate) UseCase {
	return &DefaultUseCase{
		studySetRepo: studySetRepo,
		validate:     validate,
	}
}

var (
	ErrRepoFailed = errors.New("repository failed")
	ErrValidation = errors.New("validation error")
)

func (uc *DefaultUseCase) Create(ctx context.Context, insertData *insertStudySetData) (*models.StudySet, error) {
	if err := uc.validate.Struct(insertData); err != nil {
		return nil, fmt.Errorf("%w: %w", ErrValidation, err)
	}

	studySet, err := uc.studySetRepo.Insert(ctx, insertData)
	if err != nil {
		return nil, fmt.Errorf("%w: failed to create the study set: %w", ErrRepoFailed, err)
	}

	return studySet, nil
}

func (uc *DefaultUseCase) GetById(ctx context.Context, studySetID int64) (*models.StudySet, error) {
	studySet, err := uc.studySetRepo.GetById(ctx, studySetID)
	if err != nil {
		return nil, fmt.Errorf("%w: failed to get the study set: %w", ErrRepoFailed, err)
	}
	if studySet == nil {
		return nil, ErrNotFound
	}

	return studySet, nil
}

func (uc *DefaultUseCase) GetAllSummary(ctx context.Context) ([]*models.StudySetSummary, error) {
	studySets, err := uc.studySetRepo.GetAllSummary(ctx)
	if err != nil {
		return nil, fmt.Errorf("%w: failed to get study sets: %w", ErrRepoFailed, err)
	}

	return studySets, nil
}

func (uc *DefaultUseCase) Update(ctx context.Context, studySetID int64, userID string, updateData *updateStudySetData) error {
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

func (uc *DefaultUseCase) Delete(ctx context.Context, studySetID int64, userID string) error {
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
