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

func (uc *DefinitionUseCase) GetAllFor(ctx context.Context, studySetID int64) ([]*domain.Definition, error) {
	// TODO: In theory study set can be deleted between checking if it exists and getting it's definitions. Should we use transaction for that?
	parentExists, err := uc.studySetRepo.Exists(ctx, studySetID)
	if err != nil {
		return nil, fmt.Errorf("%w: failed to check if parent study set exists: %w", ErrRepoFailed, err)
	}
	if !parentExists {
		return nil, ErrNotFound
	}

	definitionRows, err := uc.definitionRepo.GetAllFor(ctx, studySetID)
	if err != nil {
		return nil, fmt.Errorf("%w: failed to get all definitions for the study set: %w", ErrRepoFailed, err)
	}

	definitions := make([]*domain.Definition, 0)
	for _, definitionRow := range definitionRows {
		definitions = append(definitions, definitionRow.Populate())
	}

	return definitions, nil
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
