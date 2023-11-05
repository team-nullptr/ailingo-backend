package domain

import (
	"context"
	"time"
)

type StudySession struct {
	LastSessionAt *time.Time `json:"lastSessionAt"`
}

type StudySessionWithStudySet struct {
	LastSessionAt *time.Time         `json:"lastSessionAt"`
	StudySet      StudySetWithAuthor `json:"studySet"`
}

type StudySessionUseCase interface {
	GetRecent(ctx context.Context, userID string) ([]*StudySessionWithStudySet, error)
	GetForStudySet(ctx context.Context, userID string, studySetID int64) (*StudySession, error)
	// Refresh is responsible for refreshing the given study session.
	// Refreshing means updating last session timestamp. If the session does not exist a new session is created.
	Refresh(ctx context.Context, userID string, studySetID int64) error
}

type StudySessionRepo interface {
	GetRecent(ctx context.Context, userID string) ([]*StudySessionWithStudySet, error)
	GetForStudySet(ctx context.Context, userID string, studySetID int64) (*StudySession, error)
	Create(ctx context.Context, userID string, studySetID int64) error
	Refresh(ctx context.Context, userID string, studySetID int64) error
	Exists(ctx context.Context, userID string, studySetID int64) (bool, error)
}
