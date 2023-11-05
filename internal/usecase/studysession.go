package usecase

import (
	"context"
	"fmt"

	"ailingo/internal/domain"
)

type studySessionUseCase struct {
	datastore domain.DataStore
}

func NewStudySessionUseCase(datastore domain.DataStore) domain.StudySessionUseCase {
	return &studySessionUseCase{
		datastore: datastore,
	}
}

func (uc *studySessionUseCase) GetRecent(ctx context.Context, userID string) ([]*domain.StudySessionWithStudySet, error) {
	studySessionRepo := uc.datastore.GetStudySessionRepo()

	studySessions, err := studySessionRepo.GetRecent(ctx, userID)
	if err != nil {
		return nil, fmt.Errorf("%w: failed to get study sessions: %w", ErrRepoFailed, err)
	}

	return studySessions, nil
}

func (uc *studySessionUseCase) GetForStudySet(ctx context.Context, userID string, studySetID int64) (*domain.StudySession, error) {
	studySessionRepo := uc.datastore.GetStudySessionRepo()

	studySession, err := studySessionRepo.GetForStudySet(ctx, userID, studySetID)
	if err != nil {
		return nil, fmt.Errorf("%w: failed to get the study session: %w", ErrRepoFailed, err)
	}

	return studySession, nil
}

func (uc *studySessionUseCase) Refresh(ctx context.Context, userID string, studySetID int64) error {
	if err := uc.datastore.Atomic(ctx, func(ds domain.DataStore) error {
		studySessionRepo := ds.GetStudySessionRepo()

		studySessionExists, err := studySessionRepo.Exists(ctx, userID, studySetID)
		if err != nil {
			return fmt.Errorf("%w: failed to check if study session exists: %w", ErrRepoFailed, err)
		}

		if studySessionExists {
			// If study session already exists we want to update last session timestamp
			if err := studySessionRepo.Refresh(ctx, userID, studySetID); err != nil {
				return fmt.Errorf("%w: failed to refresh existing study session: %w", ErrRepoFailed, err)
			}
		} else {
			// Otherwise we want to create a new study session if the study set exists.
			studySetExists, err := ds.GetStudySetRepo().Exists(ctx, studySetID)
			if err != nil {
				return fmt.Errorf("%w: failed to check if study set exists: %w", ErrRepoFailed, err)
			}
			if !studySetExists {
				return &ErrNotFound{
					Resource: StudySetResource,
				}
			}
			if err := studySessionRepo.Create(ctx, userID, studySetID); err != nil {
				return fmt.Errorf("%w: failed to create a new study session: %w", ErrRepoFailed, err)
			}
		}

		return nil
	}); err != nil {
		return fmt.Errorf("atomic operation failed: %w", err)
	}

	return nil
}
