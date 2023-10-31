package usecase

import (
	"context"
	"fmt"

	"ailingo/internal/domain"
)

// DefinitionUseCase implements methods required by domain.DefinitionUseCase interface.
type DefinitionUseCase struct {
	definitionRepo domain.DefinitionRepo
	studySetRepo   domain.StudySetRepo
}

// NewDefinitionUseCase creates a new DefinitionUseCase.
func NewDefinitionUseCase(definitionRepo domain.DefinitionRepo, studySetRepo domain.StudySetRepo) domain.DefinitionUseCase {
	return &DefinitionUseCase{
		definitionRepo: definitionRepo,
		studySetRepo:   studySetRepo,
	}
}

func (uc *DefinitionUseCase) Create(ctx context.Context, userID string, studySetID int64, insertData *domain.InsertDefinitionData) error {
	parentStudySet, err := uc.studySetRepo.GetById(ctx, studySetID)
	if err != nil {
		return fmt.Errorf("%w: failed to get the parent study set: %w", ErrRepoFailed, err)
	}
	if parentStudySet == nil {
		return ErrNotFound
	}
	if parentStudySet.AuthorId != userID {
		return ErrForbidden
	}

	if err := uc.definitionRepo.Insert(ctx, studySetID, insertData); err != nil {
		return fmt.Errorf("%w: failed to insert a new definition: %w", ErrRepoFailed, err)
	}

	return nil
}
