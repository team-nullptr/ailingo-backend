package domain

import "context"

type ProfileUseCase interface {
	GetCreatedStudySets(ctx context.Context, userID string) ([]*StudySet, error)
	GetStarredStudySets(ctx context.Context, userID string) ([]*StudySetWithAuthor, error)
	StarStudySet(ctx context.Context, userID string, studySetID int64) error
	InstarStudySet(ctx context.Context, userID string, studySetID int64) error
}

type ProfileRepo interface {
	InsertStar(ctx context.Context, userID string, studySetID int64) error
	DeleteStar(ctx context.Context, userID string, studySetID int64) error
}
