package usecase

import (
	"context"
	"errors"
	"fmt"

	"ailingo/internal/domain"
	"ailingo/internal/mysql"
	"ailingo/pkg/auth"
)

var (
	ErrAlreadyStarred = errors.New("study set has already been starred")
)

type ProfileUseCase struct {
	studySetRepo domain.StudySetRepo
	profileRepo  domain.ProfileRepo
	userService  *auth.UserService
}

func NewProfileUseCase(studySetRepo domain.StudySetRepo, profileRepo domain.ProfileRepo, userService *auth.UserService) *ProfileUseCase {
	return &ProfileUseCase{
		studySetRepo: studySetRepo,
		profileRepo:  profileRepo,
		userService:  userService,
	}
}

func (uc *ProfileUseCase) GetStarredStudySets(ctx context.Context, userID string) ([]*domain.StudySetWithAuthor, error) {
	studySets, err := uc.studySetRepo.GetStarredBy(ctx, userID)
	if err != nil {
		return nil, fmt.Errorf("%w: failed to get all starred study sets: %w", ErrRepoFailed, err)
	}

	return studySets, nil
}

func (uc *ProfileUseCase) GetCreatedStudySets(ctx context.Context, userID string) ([]*domain.StudySet, error) {
	studySets, err := uc.studySetRepo.GetCreatedBy(ctx, userID)
	if err != nil {
		return nil, fmt.Errorf("failed to get created study sets: %w", err)
	}

	return studySets, nil
}

func (uc *ProfileUseCase) StarStudySet(ctx context.Context, userID string, studySetID int64) error {
	parentExists, err := uc.studySetRepo.Exists(ctx, studySetID)
	if err != nil {
		return fmt.Errorf("%w: failed to check if parent study set exists: %w", ErrRepoFailed, err)
	}
	if !parentExists {
		return ErrNotFound
	}

	if err := uc.profileRepo.InsertStar(ctx, userID, studySetID); err != nil {
		if errors.Is(err, mysql.ErrDuplicateRow) {
			return ErrAlreadyStarred
		}
		return fmt.Errorf("%w: failed to insert a star: %w", ErrRepoFailed, err)
	}

	return nil
}

func (uc *ProfileUseCase) InstarStudySet(ctx context.Context, userID string, studySetID int64) error {
	if err := uc.profileRepo.DeleteStar(ctx, userID, studySetID); err != nil {
		return fmt.Errorf("failed to delete the start: %w", err)
	}
	return nil
}
