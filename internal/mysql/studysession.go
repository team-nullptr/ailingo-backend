package mysql

import (
	"context"
	"database/sql"
	"errors"
	"fmt"

	"ailingo/internal/domain"
)

const getRecentStudySessions = `
SELECT study_session.last_session_at,
       study_set.id,
       study_set.name,
       study_set.description,
       study_set.phrase_language,
       study_set.definition_language,
       study_set.icon,
       study_set.color,
       user.id,
       user.username,
       user.image_url
FROM study_session
	     INNER JOIN study_set ON study_session.study_set_id = study_set.id
	     INNER JOIN user ON study_set.author_id = user.id
WHERE study_session.user_id = ?
ORDER BY study_session.last_session_at DESC 
`

const getStudySessionForStudySet = `
SELECT last_session_at
FROM study_session
WHERE user_id = ?
  AND study_set_id = ?
`

const refreshStudySession = `
UPDATE study_session
SET last_session_at = NOW()
WHERE user_id = ?
  AND study_set_id = ?
`

const insertStudySession = `
INSERT INTO study_session (user_id, study_set_id)
VALUES (?, ?)
`

const studySessionExists = `
SELECT EXISTS(SELECT 1 FROM study_session WHERE user_id = ? AND study_set_id = ?) 
`

type studySessionRepo struct {
	db DBTX
}

func NewStudySessionRepo(db DBTX) domain.StudySessionRepo {
	return &studySessionRepo{
		db: db,
	}
}

func (r *studySessionRepo) GetRecent(ctx context.Context, userID string) ([]*domain.StudySessionWithStudySet, error) {
	studySessions := make([]*domain.StudySessionWithStudySet, 0)

	rows, err := r.db.QueryContext(ctx, getRecentStudySessions, userID)
	if err != nil {
		return nil, fmt.Errorf("failed to query: %w", err)
	}

	for rows.Next() {
		var studySession domain.StudySessionWithStudySet

		if err := rows.Scan(
			&studySession.LastSessionAt, &studySession.StudySet.Id, &studySession.StudySet.Name, &studySession.StudySet.Description, &studySession.StudySet.PhraseLanguage, &studySession.StudySet.DefinitionLanguage, &studySession.StudySet.Icon, &studySession.StudySet.Color,
			&studySession.StudySet.Author.Id, &studySession.StudySet.Author.Username, &studySession.StudySet.Author.ImageURL,
		); err != nil {
			return nil, fmt.Errorf("failed to scan: %w", err)
		}

		studySessions = append(studySessions, &studySession)
	}

	return studySessions, nil
}

func (r *studySessionRepo) GetForStudySet(ctx context.Context, userID string, studySetID int64) (*domain.StudySession, error) {
	row := r.db.QueryRowContext(ctx, getStudySessionForStudySet, userID, studySetID)

	var studySession domain.StudySession
	if err := row.Scan(&studySession.LastSessionAt); err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, nil
		}
		return nil, fmt.Errorf("")
	}

	return &studySession, nil
}

func (r *studySessionRepo) Create(ctx context.Context, userID string, studySetID int64) error {
	if _, err := r.db.ExecContext(ctx, insertStudySession, userID, studySetID); err != nil {
		return fmt.Errorf("failed to exec insert query: %w", err)
	}
	return nil
}

func (r *studySessionRepo) Refresh(ctx context.Context, userID string, studySetID int64) error {
	if _, err := r.db.ExecContext(ctx, refreshStudySession, userID, studySetID); err != nil {
		return fmt.Errorf("failed to exec a refresh query: %w", err)
	}
	return nil
}

func (r *studySessionRepo) Exists(ctx context.Context, userID string, studySetID int64) (bool, error) {
	var exists int

	res := r.db.QueryRowContext(ctx, studySessionExists, userID, studySetID)
	if err := res.Scan(&exists); err != nil {
		return false, fmt.Errorf("failed to exec exists statement: %w", err)
	}

	return exists == 1, nil
}
