package usecase

import (
	"context"
	"fmt"

	"github.com/go-playground/validator/v10"

	"ailingo/internal/domain"
)

// DefinitionUseCase implements methods required by domain.DefinitionUseCase interface.
type DefinitionUseCase struct {
	dataStore domain.DataStore
	validate  *validator.Validate
}

// NewDefinitionUseCase creates a new DefinitionUseCase.
func NewDefinitionUseCase(dataStore domain.DataStore, validate *validator.Validate) domain.DefinitionUseCase {
	return &DefinitionUseCase{
		dataStore: dataStore,
		validate:  validate,
	}
}

func (uc *DefinitionUseCase) GetAllFor(ctx context.Context, parentStudySetID int64) ([]*domain.Definition, error) {
	var definitions []*domain.Definition

	err := uc.dataStore.Atomic(ctx, func(ds domain.DataStore) error {
		studySetRepo := ds.GetStudySetRepo()
		definitionRepo := ds.GetDefinitionRepo()

		parentExists, err := studySetRepo.Exists(ctx, parentStudySetID)
		if err != nil {
			return fmt.Errorf("%w: failed to check if parent study set exists: %w", ErrRepoFailed, err)
		}
		if !parentExists {
			return &ErrNotFound{
				Resource: DefinitionResource,
			}
		}

		definitionRows, err := definitionRepo.GetAllFor(ctx, parentStudySetID)
		if err != nil {
			return fmt.Errorf("%w: failed to get all definitions for the study set: %w", ErrRepoFailed, err)
		}

		definitions = make([]*domain.Definition, 0)
		for _, definitionRow := range definitionRows {
			definitions = append(definitions, definitionRow.Populate())
		}

		return nil
	})

	if err != nil {
		return nil, fmt.Errorf("atomic operation failed: %w", err)
	}

	return definitions, nil
}

func (uc *DefinitionUseCase) Create(ctx context.Context, userID string, parentStudySetID int64, insertData *domain.InsertDefinitionData) error {
	if err := uc.validate.Struct(insertData); err != nil {
		return fmt.Errorf("%w: invalid insert data: %w", ErrValidation, err)
	}

	err := uc.dataStore.Atomic(ctx, func(ds domain.DataStore) error {
		studySetRepo := ds.GetStudySetRepo()
		definitionRepo := ds.GetDefinitionRepo()

		if err := uc.checkStudySetOwnership(ctx, studySetRepo, userID, parentStudySetID); err != nil {
			return err
		}

		if err := definitionRepo.Insert(ctx, parentStudySetID, insertData); err != nil {
			return fmt.Errorf("%w: failed to insert a new definition: %w", ErrRepoFailed, err)
		}

		return nil
	})

	if err != nil {
		return fmt.Errorf("atomic operation failed: %w", err)
	}

	return nil
}

func (uc *DefinitionUseCase) Update(ctx context.Context, userID string, parentStudySetID int64, definitionID int64, updateData *domain.UpdateDefinitionData) error {
	if err := uc.validate.Struct(updateData); err != nil {
		return fmt.Errorf("%w: invalid insert data: %w", ErrValidation, err)
	}

	err := uc.dataStore.Atomic(ctx, func(ds domain.DataStore) error {
		studySetRepo := ds.GetStudySetRepo()
		definitionRepo := ds.GetDefinitionRepo()

		if err := uc.checkStudySetOwnership(ctx, studySetRepo, userID, parentStudySetID); err != nil {
			return err
		}

		if err := definitionRepo.Update(ctx, definitionID, updateData); err != nil {
			return fmt.Errorf("failed to update the definition: %w", err)
		}

		return nil
	})

	if err != nil {
		return fmt.Errorf("atomic operation failed: %w", err)
	}

	return nil
}

func (uc *DefinitionUseCase) Delete(ctx context.Context, userID string, parentStudySetID int64, definitionID int64) error {
	err := uc.dataStore.Atomic(ctx, func(ds domain.DataStore) error {
		studySetRepo := ds.GetStudySetRepo()
		definitionRepo := ds.GetDefinitionRepo()

		// TODO: In theory study set can be deleted between checking if it exists and deleting it's definition. Should we use transaction for that?
		if err := uc.checkStudySetOwnership(ctx, studySetRepo, userID, parentStudySetID); err != nil {
			return fmt.Errorf("failed to check study set ownership: %w", err)
		}

		if err := definitionRepo.Delete(ctx, definitionID); err != nil {
			return fmt.Errorf("%w: failed to delete the definition: %w", ErrRepoFailed, err)
		}

		return nil
	})

	if err != nil {
		return fmt.Errorf("atomic operation failed: %w", err)
	}

	return nil
}

func (uc *DefinitionUseCase) checkStudySetOwnership(ctx context.Context, studySetRepo domain.StudySetRepo, userID string, studySetID int64) error {
	parentStudySet, err := studySetRepo.GetById(ctx, studySetID)
	if err != nil {
		return fmt.Errorf("%w: failed to get the parent study set: %w", ErrRepoFailed, err)
	}
	if parentStudySet == nil {
		return &ErrNotFound{
			Resource: StudySetResource,
		}
	}
	if parentStudySet.Author.Id != userID {
		return ErrForbidden
	}
	return nil
}
