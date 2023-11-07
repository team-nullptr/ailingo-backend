package usecase

import (
	"context"
	"fmt"
	"log/slog"
	"time"

	"github.com/go-playground/validator/v10"

	"ailingo/internal/domain"
)

// definitionUseCase implements methods required by domain.DefinitionUseCase interface.
type definitionUseCase struct {
	l         *slog.Logger
	dataStore domain.DataStore
	aiService domain.AiService
	validate  *validator.Validate
}

// NewDefinitionUseCase creates a new definitionUseCase.
func NewDefinitionUseCase(l *slog.Logger, dataStore domain.DataStore, aiService domain.AiService, validate *validator.Validate) domain.DefinitionUseCase {
	return &definitionUseCase{
		l:         l,
		dataStore: dataStore,
		aiService: aiService,
		validate:  validate,
	}
}

func (uc *definitionUseCase) GetAllFor(ctx context.Context, parentStudySetID int64) ([]*domain.Definition, error) {
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

func (uc *definitionUseCase) Create(ctx context.Context, userID string, parentStudySetID int64, insertData *domain.InsertDefinitionData) error {
	if err := uc.validate.Struct(insertData); err != nil {
		return fmt.Errorf("%w: invalid insert data: %w", ErrValidation, err)
	}

	err := uc.dataStore.Atomic(ctx, func(ds domain.DataStore) error {
		studySetRepo := ds.GetStudySetRepo()
		definitionRepo := ds.GetDefinitionRepo()

		if _, err := uc.checkStudySetOwnership(ctx, studySetRepo, userID, parentStudySetID); err != nil {
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

func (uc *definitionUseCase) Update(ctx context.Context, userID string, parentStudySetID int64, definitionID int64, updateData *domain.UpdateDefinitionData) error {
	if err := uc.validate.Struct(updateData); err != nil {
		return fmt.Errorf("%w: invalid insert data: %w", ErrValidation, err)
	}

	err := uc.dataStore.Atomic(ctx, func(ds domain.DataStore) error {
		studySetRepo := ds.GetStudySetRepo()
		definitionRepo := ds.GetDefinitionRepo()

		if _, err := uc.checkStudySetOwnership(ctx, studySetRepo, userID, parentStudySetID); err != nil {
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

func (uc *definitionUseCase) Delete(ctx context.Context, userID string, parentStudySetID int64, definitionID int64) error {
	err := uc.dataStore.Atomic(ctx, func(ds domain.DataStore) error {
		studySetRepo := ds.GetStudySetRepo()
		definitionRepo := ds.GetDefinitionRepo()

		// TODO: In theory study set can be deleted between checking if it exists and deleting it's definition. Should we use transaction for that?
		if _, err := uc.checkStudySetOwnership(ctx, studySetRepo, userID, parentStudySetID); err != nil {
			return err
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

func (uc *definitionUseCase) AiFill(ctx context.Context, userID string, parentStudySetID int64) (int64, error) {
	studySetRepo := uc.dataStore.GetStudySetRepo()
	taskRepo := uc.dataStore.GetTaskRepo()

	parentStudySet, err := uc.checkStudySetOwnership(ctx, studySetRepo, userID, parentStudySetID)
	if err != nil {
		return 0, err
	}

	// Create task.
	taskId, err := taskRepo.Insert(ctx)
	if err != nil {
		return 0, fmt.Errorf("%w: failed to create a new task: %w", ErrRepoFailed, err)
	}

	// Start a new task in a go routine.
	go func() {
		ctx, cancel := context.WithTimeout(context.Background(), time.Minute)
		defer cancel()

		// Process the task
		err := func() error {
			definitions, err := uc.aiService.GenerateDefinitions(ctx, &domain.SetGenerationRequest{Name: parentStudySet.Name})
			if err != nil {
				return fmt.Errorf("could not generate definitions: %w", err)
			}

			err = uc.dataStore.Atomic(ctx, func(ds domain.DataStore) error {
				definitionRepo := ds.GetDefinitionRepo()
				for _, definition := range definitions {
					definition.Sentences = []string{}
					if err := definitionRepo.Insert(ctx, parentStudySetID, definition); err != nil {
						return fmt.Errorf("failed to insert a definition: %w", err)
					}
				}
				return nil
			})
			if err != nil {
				return fmt.Errorf("atomic insert failed: %w", err)
			}

			return nil
		}()

		// If task has failed then fail it
		if err != nil {
			uc.l.Error(fmt.Sprintf("failed to generate definitions: %s", err))
			if err := taskRepo.Fail(ctx, taskId); err != nil {
				uc.l.Error(fmt.Sprintf("failed to fail the task: %s", err))
			}
			return
		}

		// Otherwise complete
		if err := taskRepo.Complete(ctx, taskId); err != nil {
			uc.l.Error(fmt.Sprintf("failed to complete the task: %s", err))
			return
		}

		uc.l.Info(fmt.Sprintf("task %d has been completed successfully", taskId))
	}()

	return taskId, nil
}

func (uc *definitionUseCase) checkStudySetOwnership(ctx context.Context, studySetRepo domain.StudySetRepo, userID string, studySetID int64) (*domain.StudySetWithAuthor, error) {
	parentStudySet, err := studySetRepo.GetById(ctx, studySetID)
	if err != nil {
		return nil, fmt.Errorf("%w: failed to get the parent study set: %w", ErrRepoFailed, err)
	}
	if parentStudySet == nil {
		return nil, &ErrNotFound{
			Resource: StudySetResource,
		}
	}
	if parentStudySet.Author.Id != userID {
		return nil, ErrForbidden
	}
	return parentStudySet, nil
}
