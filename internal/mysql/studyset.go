package mysql

import (
	"context"
	"database/sql"
	"errors"
	"fmt"

	"ailingo/internal/domain"
)

// studySetQueries provides all the necessary queries used by study set repo.
type studySetQueries struct {
	getById      *sql.Stmt
	getAll       *sql.Stmt
	getCreatedBy *sql.Stmt
	getStarred   *sql.Stmt
	insert       *sql.Stmt
	update       *sql.Stmt
	delete       *sql.Stmt
	exists       *sql.Stmt
}

// newStudySetQueries creates a new prepared query provider for study set db operations.
func newStudySetQueries(ctx context.Context, db *sql.DB) (*studySetQueries, error) {
	getByIdStmt, err := db.PrepareContext(ctx, "SELECT study_set.id, study_set.name, study_set.description, study_set.phrase_language, study_set.definition_language, user.id, user.username, user.image_url FROM study_set INNER JOIN user ON user.id = study_set.author_id WHERE study_set.id = ? LIMIT 1")
	if err != nil {
		return nil, fmt.Errorf("getById query: %w", err)
	}

	getAllStmt, err := db.PrepareContext(ctx, "SELECT study_set.id, study_set.name, study_set.description, study_set.phrase_language, study_set.definition_language, user.id, user.username, user.image_url FROM study_set INNER JOIN user ON user.id = study_set.author_id")
	if err != nil {
		return nil, fmt.Errorf("getAll query: %w", err)
	}

	getCreatedByStmt, err := db.PrepareContext(ctx, "SELECT id, name, description, phrase_language, definition_language FROM study_set WHERE author_id = ?")
	if err != nil {
		return nil, fmt.Errorf("getAllBy query: %w", err)
	}

	getStarredStmt, err := db.PrepareContext(ctx, "SELECT study_set.id, study_set.name, study_set.description, study_set.phrase_language, study_set.definition_language, user.id, user.username, user.image_url FROM star INNER JOIN study_set ON star.study_set_id = study_set.id INNER JOIN user ON user.id = study_set.author_id WHERE star.user_id = ?")
	if err != nil {
		return nil, fmt.Errorf("getAllStarred query: %w", err)
	}

	insertStmt, err := db.PrepareContext(ctx, "INSERT INTO study_set (author_id, name, description, phrase_language, definition_language) VALUES (?, ?, ?, ?, ?)")
	if err != nil {
		return nil, fmt.Errorf("insert query: %w", err)
	}

	updateStmt, err := db.PrepareContext(ctx, "UPDATE study_set SET name = ?, description = ?, phrase_language = ?, definition_language = ? WHERE id = ?")
	if err != nil {
		return nil, fmt.Errorf("upate query: %w", err)
	}

	deleteStmt, err := db.Prepare("DELETE FROM study_set WHERE id = ?")
	if err != nil {
		return nil, fmt.Errorf("delete query: %w", err)
	}

	existsStmt, err := db.Prepare("SELECT EXISTS(SELECT 1 FROM study_set WHERE study_set.id = ?)")
	if err != nil {
		return nil, fmt.Errorf("exists query: %w", err)
	}

	return &studySetQueries{
		getById:      getByIdStmt,
		getAll:       getAllStmt,
		getCreatedBy: getCreatedByStmt,
		getStarred:   getStarredStmt,
		insert:       insertStmt,
		update:       updateStmt,
		delete:       deleteStmt,
		exists:       existsStmt,
	}, nil
}

type StudySetRepo struct {
	query *studySetQueries
}

func NewStudySetRepo(ctx context.Context, db *sql.DB) (domain.StudySetRepo, error) {
	query, err := newStudySetQueries(ctx, db)
	if err != nil {
		return nil, err
	}

	return &StudySetRepo{
		query: query,
	}, nil
}

func (r *StudySetRepo) GetAll(ctx context.Context) ([]*domain.StudySetWithAuthor, error) {
	studySets := make([]*domain.StudySetWithAuthor, 0)

	rows, err := r.query.getAll.QueryContext(ctx)
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

func (r *StudySetRepo) GetById(ctx context.Context, studySetID int64) (*domain.StudySetWithAuthor, error) {
	var studySet domain.StudySetWithAuthor

	if err := r.query.getById.QueryRowContext(ctx, studySetID).Scan(
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

func (r *StudySetRepo) GetCreatedBy(ctx context.Context, userID string) ([]*domain.StudySet, error) {
	studySets := make([]*domain.StudySet, 0)

	rows, err := r.query.getCreatedBy.QueryContext(ctx, userID)
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

func (r *StudySetRepo) GetStarredBy(ctx context.Context, userID string) ([]*domain.StudySetWithAuthor, error) {
	studySets := make([]*domain.StudySetWithAuthor, 0)

	rows, err := r.query.getStarred.QueryContext(ctx, userID)
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

func (r *StudySetRepo) Insert(ctx context.Context, insertData *domain.InsertStudySetData) (int64, error) {
	res, err := r.query.insert.ExecContext(
		ctx,
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

func (r *StudySetRepo) Update(ctx context.Context, studySetID int64, updateData *domain.UpdateStudySetData) error {
	if _, err := r.query.update.ExecContext(
		ctx,
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

func (r *StudySetRepo) Delete(ctx context.Context, studySetID int64) error {
	if _, err := r.query.delete.ExecContext(ctx, studySetID); err != nil {
		return fmt.Errorf("failed to exec: %w", err)
	}
	return nil
}

func (r *StudySetRepo) Exists(ctx context.Context, studySetID int64) (bool, error) {
	var exists int

	res := r.query.exists.QueryRowContext(ctx, studySetID)
	if err := res.Scan(&exists); err != nil {
		return false, fmt.Errorf("failed to exec exists statement: %w", err)
	}

	return exists == 1, nil
}
