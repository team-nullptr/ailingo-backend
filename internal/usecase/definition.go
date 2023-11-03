package usecase

import (
	"context"
	"fmt"

	"github.com/go-playground/validator/v10"

	"ailingo/internal/domain"
)

// DefinitionUseCase implements methods required by domain.DefinitionUseCase interface.
type DefinitionUseCase struct {
	definitionRepo domain.DefinitionRepo
	studySetRepo   domain.StudySetRepo
	validate       *validator.Validate
}

// NewDefinitionUseCase creates a new DefinitionUseCase.
func NewDefinitionUseCase(definitionRepo domain.DefinitionRepo, studySetRepo domain.StudySetRepo, validate *validator.Validate) domain.DefinitionUseCase {
	return &DefinitionUseCase{
		definitionRepo: definitionRepo,
		studySetRepo:   studySetRepo,
		validate:       validate,
	}
}

func (uc *DefinitionUseCase) GetAllFor(ctx context.Context, parentStudySetID int64) ([]*domain.Definition, error) {
	// TODO: In theory study set can be deleted between checking if it exists and getting it's definitions. Should we use transaction for that?
	parentExists, err := uc.studySetRepo.Exists(ctx, parentStudySetID)
	if err != nil {
		return nil, fmt.Errorf("%w: failed to check if parent study set exists: %w", ErrRepoFailed, err)
	}
	if !parentExists {
		return nil, ErrNotFound
	}

	definitionRows, err := uc.definitionRepo.GetAllFor(ctx, parentStudySetID)
	if err != nil {
		return nil, fmt.Errorf("%w: failed to get all definitions for the study set: %w", ErrRepoFailed, err)
	}

	definitions := make([]*domain.Definition, 0)
	for _, definitionRow := range definitionRows {
		definitions = append(definitions, definitionRow.Populate())
	}

	return definitions, nil
}

func (uc *DefinitionUseCase) Create(ctx context.Context, userID string, parentStudySetID int64, insertData *domain.InsertDefinitionData) error {
	if err := uc.validate.Struct(insertData); err != nil {
		return fmt.Errorf("%w: invalid insert data: %w", ErrValidation, err)
	}

	if err := uc.checkStudySetOwnership(ctx, userID, parentStudySetID); err != nil {
		return err
	}

	if err := uc.definitionRepo.Insert(ctx, parentStudySetID, insertData); err != nil {
		return fmt.Errorf("%w: failed to insert a new definition: %w", ErrRepoFailed, err)
	}

	return nil
}

func (uc *DefinitionUseCase) Update(ctx context.Context, userID string, parentStudySetID int64, definitionID int64, updateData *domain.UpdateDefinitionData) error {
	if err := uc.validate.Struct(updateData); err != nil {
		return fmt.Errorf("%w: invalid insert data: %w", ErrValidation, err)
	}

	// TODO: In theory study set can be deleted between checking if it exists and deleting it's definition. Should we use transaction for that?
	if err := uc.checkStudySetOwnership(ctx, userID, parentStudySetID); err != nil {
		return err
	}

	if err := uc.definitionRepo.Update(ctx, parentStudySetID, definitionID, updateData); err != nil {
		return fmt.Errorf("failed to update the definition: %w", err)
	}

	return nil
}

func (uc *DefinitionUseCase) Delete(ctx context.Context, userID string, parentStudySetID int64, definitionID int64) error {
	// TODO: In theory study set can be deleted between checking if it exists and deleting it's definition. Should we use transaction for that?
	if err := uc.checkStudySetOwnership(ctx, userID, parentStudySetID); err != nil {
		return fmt.Errorf("failed to check study set ownership: %w", err)
	}

	if err := uc.definitionRepo.Delete(ctx, parentStudySetID, definitionID); err != nil {
		return fmt.Errorf("%w: failed to delete the definition: %w", ErrRepoFailed, err)
	}

	return nil
}

func (uc *DefinitionUseCase) checkStudySetOwnership(ctx context.Context, userID string, studySetID int64) error {
	parentStudySet, err := uc.studySetRepo.GetById(ctx, studySetID)
	if err != nil {
		return fmt.Errorf("%w: failed to get the parent study set: %w", ErrRepoFailed, err)
	}

	if parentStudySet == nil {
		return ErrNotFound
	}

	if parentStudySet.Author.Id != userID {
		return ErrForbidden
	}

	return nil
}
