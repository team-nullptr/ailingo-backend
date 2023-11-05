package mysql

import (
	"context"
	"database/sql"
	"errors"
	"fmt"

	"ailingo/internal/domain"
)

// getAllStudySets queries for all study sets that have ever been created.
const getStudySets = `
SELECT study_set.id,
       study_set.name,
       study_set.description,
       study_set.phrase_language,
       study_set.definition_language,
       user.id,
       user.username,
       user.image_url
FROM study_set
         INNER JOIN user ON user.id = study_set.author_id
`

// getStudySetsCreatedBy queries for all study sets created by the specified user.
const getStudySetsCreatedBy = `
SELECT id, name, description, phrase_language, definition_language
FROM study_set
WHERE author_id = ?
`

// getStudySetsStarredBy queries for all study sets starred by the specified user.
const getStudySetsStarredBy = `
SELECT study_set.id,
       study_set.name,
       study_set.description,
       study_set.phrase_language,
       study_set.definition_language,
       user.id,
       user.username,
       user.image_url
FROM star
         INNER JOIN study_set ON star.study_set_id = study_set.id
         INNER JOIN user ON user.id = study_set.author_id
WHERE star.user_id = ?
`

// getStudySetById queries for a study set with the given id
const getStudySetById = `
SELECT study_set.id,
       study_set.name,
       study_set.description,
       study_set.phrase_language,
       study_set.definition_language,
       user.id,
       user.username,
       user.image_url
FROM study_set
         INNER JOIN user ON user.id = study_set.author_id
WHERE study_set.id = ?
LIMIT 1
`

// insertStudySets inserts a new study sets into the db.
const insertStudySet = `
INSERT INTO study_set (author_id, name, description, phrase_language, definition_language)
VALUES (?, ?, ?, ?, ?)
`

// updateStudySet updates the given study set.
const updateStudySet = `
UPDATE study_set
SET name                = ?,
    description         = ?,
    phrase_language     = ?,
    definition_language = ?
WHERE id = ?
`

// deleteStudySet deletes the specified study set.
const deleteStudySet = `
DELETE
FROM study_set
WHERE id = ?
`

// studySetExists checks if a study set with the specified id exists
const studySetExists = `
SELECT EXISTS(SELECT 1 FROM study_set WHERE study_set.id = ?)
`

const deleteStudySetDefinitions = `
DELETE
FROM definition
WHERE study_set_id = ?

`

const deleteStudySetStars = `
DELETE 
FROM star 
WHERE study_set_id = ?
`

const deleteStudySetStudySessions = `
DELETE 
FROM study_session 
WHERE study_set_id = ?
`

type studySetRepo struct {
	db DBTX
}

func NewStudySetRepo(db DBTX) domain.StudySetRepo {
	return &studySetRepo{
		db: db,
	}
}

func (r *studySetRepo) GetAll(ctx context.Context) ([]*domain.StudySetWithAuthor, error) {
	studySets := make([]*domain.StudySetWithAuthor, 0)

	rows, err := r.db.QueryContext(ctx, getStudySets)
	if err != nil {
		return nil, fmt.Errorf("failed to query: %w", err)
	}

	for rows.Next() {
		var studySet domain.StudySetWithAuthor
		if err := rows.Scan(
			// study set
			&studySet.Id, &studySet.Name, &studySet.Description, &studySet.PhraseLanguage, &studySet.DefinitionLanguage,
			// author
			&studySet.Author.Id, &studySet.Author.Username, &studySet.Author.ImageURL,
		); err != nil {
			return nil, fmt.Errorf("failed to scan: %w", err)
		}
		studySets = append(studySets, &studySet)
	}

	return studySets, nil
}

func (r *studySetRepo) GetById(ctx context.Context, studySetID int64) (*domain.StudySetWithAuthor, error) {
	var studySet domain.StudySetWithAuthor

	if err := r.db.QueryRowContext(ctx, getStudySetById, studySetID).Scan(
		// study set
		&studySet.Id, &studySet.Name, &studySet.Description, &studySet.PhraseLanguage, &studySet.DefinitionLanguage,
		// author
		&studySet.Author.Id, &studySet.Author.Username, &studySet.Author.ImageURL,
	); err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, nil
		}
		return nil, fmt.Errorf("failed to query: %w", err)
	}

	return &studySet, nil
}

func (r *studySetRepo) GetCreatedBy(ctx context.Context, userID string) ([]*domain.StudySet, error) {
	studySets := make([]*domain.StudySet, 0)

	rows, err := r.db.QueryContext(ctx, getStudySetsCreatedBy, userID)
	if err != nil {
		return nil, fmt.Errorf("failed to query: %w", err)
	}

	for rows.Next() {
		var studySet domain.StudySet
		if err := rows.Scan(&studySet.Id, &studySet.Name, &studySet.Description, &studySet.PhraseLanguage, &studySet.DefinitionLanguage); err != nil {
			return nil, fmt.Errorf("failed to scan: %w", err)
		}
		studySets = append(studySets, &studySet)
	}

	return studySets, nil
}

func (r *studySetRepo) GetStarredBy(ctx context.Context, userID string) ([]*domain.StudySetWithAuthor, error) {
	studySets := make([]*domain.StudySetWithAuthor, 0)

	rows, err := r.db.QueryContext(ctx, getStudySetsStarredBy, userID)
	if err != nil {
		return nil, fmt.Errorf("failed to query: %w", err)
	}

	for rows.Next() {
		var studySet domain.StudySetWithAuthor
		if err := rows.Scan(
			// study set
			&studySet.Id, &studySet.Name, &studySet.Description, &studySet.PhraseLanguage, &studySet.DefinitionLanguage,
			// author
			&studySet.Author.Id, &studySet.Author.Username, &studySet.Author.ImageURL,
		); err != nil {
			return nil, fmt.Errorf("failed to scan: %w", err)
		}
		studySets = append(studySets, &studySet)
	}

	return studySets, nil
}

func (r *studySetRepo) Insert(ctx context.Context, insertData *domain.InsertStudySetData) (int64, error) {
	res, err := r.db.ExecContext(
		ctx,
		insertStudySet,
		insertData.AuthorId,
		insertData.Name,
		insertData.Description,
		insertData.PhraseLanguage,
		insertData.DefinitionLanguage,
	)
	if err != nil {
		return 0, fmt.Errorf("failed to exec: %w", err)
	}

	lastInsertId, err := res.LastInsertId()
	if err != nil {
		return 0, fmt.Errorf("failed to get last insert id: %w", err)
	}

	return lastInsertId, nil
}

func (r *studySetRepo) Update(ctx context.Context, studySetID int64, updateData *domain.UpdateStudySetData) error {
	if _, err := r.db.ExecContext(
		ctx,
		updateStudySet,
		updateData.Name,
		updateData.Description,
		updateData.PhraseLanguage,
		updateData.DefinitionLanguage,
		studySetID,
	); err != nil {
		return fmt.Errorf("failed to exec: %w", err)
	}
	return nil
}

func (r *studySetRepo) Delete(ctx context.Context, studySetID int64) error {
	// TODO: This could be split into separate repo functions and run with DataStore.Atomic

	if _, err := r.db.ExecContext(ctx, deleteStudySetStars, studySetID); err != nil {
		return fmt.Errorf("failed to exec: %w", err)
	}

	if _, err := r.db.ExecContext(ctx, deleteStudySetStudySessions, studySetID); err != nil {
		return fmt.Errorf("failed to exec: %w", err)
	}

	if _, err := r.db.ExecContext(ctx, deleteStudySetDefinitions, studySetID); err != nil {
		return fmt.Errorf("failed to exec: %w", err)
	}

	if _, err := r.db.ExecContext(ctx, deleteStudySet, studySetID); err != nil {
		return fmt.Errorf("failed to exec: %w", err)
	}

	return nil
}

func (r *studySetRepo) Exists(ctx context.Context, studySetID int64) (bool, error) {
	var exists int

	res := r.db.QueryRowContext(ctx, studySetExists, studySetID)
	if err := res.Scan(&exists); err != nil {
		return false, fmt.Errorf("failed to exec exists statement: %w", err)
	}

	return exists == 1, nil
}
