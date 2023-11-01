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

func (uc *StudySetUseCase) Create(ctx context.Context, insertData *domain.InsertStudySetData) (*domain.StudySet, error) {
	if err := uc.validate.Struct(insertData); err != nil {
		return nil, fmt.Errorf("%w: %w", ErrValidation, err)
	}

	author, err := uc.userService.GetUserById(insertData.AuthorId)
	if err != nil {
		return nil, fmt.Errorf("failed to get the author: %w", err)
	}

	studySetRow, err := uc.studySetRepo.Insert(ctx, insertData)
	if err != nil {
		return nil, fmt.Errorf("%w: failed to create the study set: %w", ErrRepoFailed, err)
	}

	return studySetRow.Populate(domain.AuthorFromClerkUser(author)), nil
}

func (uc *StudySetUseCase) GetById(ctx context.Context, studySetID int64) (*domain.StudySet, error) {
	studySetRow, err := uc.studySetRepo.GetById(ctx, studySetID)
	if err != nil {
		return nil, fmt.Errorf("%w: failed to get the study set: %w", ErrRepoFailed, err)
	}
	if studySetRow == nil {
		return nil, ErrNotFound
	}

	author, err := uc.userService.GetUserById(studySetRow.AuthorId)
	if err != nil {
		return nil, fmt.Errorf("failed to get the study set author: %w", err)
	}

	return studySetRow.Populate(domain.AuthorFromClerkUser(author)), nil
}

func (uc *StudySetUseCase) GetAll(ctx context.Context) ([]*domain.StudySet, error) {
	studySetRows, err := uc.studySetRepo.GetAll(ctx)
	if err != nil {
		return nil, fmt.Errorf("%w: failed to get study sets: %w", ErrRepoFailed, err)
	}

	studySets := make([]*domain.StudySet, 0)

	for _, studySetRow := range studySetRows {
		author, err := uc.userService.GetUserById(studySetRow.AuthorId)
		if err != nil {
			// TODO: We could add some logs here
			continue
		}
		studySets = append(studySets, studySetRow.Populate(domain.AuthorFromClerkUser(author)))
	}

	return studySets, nil
}

func (uc *StudySetUseCase) Update(ctx context.Context, studySetID int64, userID string, updateData *domain.UpdateStudySetData) error {
	target, err := uc.studySetRepo.GetById(ctx, studySetID)
	if err != nil {
		return fmt.Errorf("%w: GetById failed: %w", ErrRepoFailed, err)
	}
	if target == nil {
		return ErrNotFound
	}
	if target.AuthorId != userID {
		return ErrForbidden
	}

	if err := uc.studySetRepo.Update(ctx, studySetID, updateData); err != nil {
		return fmt.Errorf("%w: Update failed: %w", ErrRepoFailed, err)
	}

	return nil
}

func (uc *StudySetUseCase) Delete(ctx context.Context, studySetID int64, userID string) error {
	target, err := uc.studySetRepo.GetById(ctx, studySetID)
	if err != nil {
		return fmt.Errorf("%w: GetById failed: %w", ErrRepoFailed, err)
	}
	if target == nil {
		return ErrNotFound
	}
	if target.AuthorId != userID {
		return ErrForbidden
	}

	if err := uc.studySetRepo.Delete(ctx, studySetID); err != nil {
		return fmt.Errorf("%w: Delete failed: %w", ErrRepoFailed, err)
	}

	return nil
}
