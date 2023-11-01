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
	insert        *sql.Stmt
	getById       *sql.Stmt
	getAllSummary *sql.Stmt
	update        *sql.Stmt
	delete        *sql.Stmt
	exists        *sql.Stmt
}

// newStudySetQueries creates a new prepared query provider for study set db operations.
func newStudySetQueries(db *sql.DB) (*studySetQueries, error) {
	insertStmt, err := db.Prepare("INSERT INTO study_set (author_id, name, description, phrase_language, definition_language) VALUES (?, ?, ?, ?, ?);")
	if err != nil {
		return nil, fmt.Errorf("insert query: %w", err)
	}

	getByIdStmt, err := db.Prepare("SELECT id, author_id, name, description, phrase_language, definition_language FROM study_set WHERE id = ?")
	if err != nil {
		return nil, fmt.Errorf("getById query: %w", err)
	}

	getAllSummaryStmt, err := db.Prepare("SELECT id, author_id, name, description, phrase_language, definition_language FROM study_set")
	if err != nil {
		return nil, fmt.Errorf("getAll query: %w", err)
	}

	updateStmt, err := db.Prepare("UPDATE study_set SET name = ?, description = ?, phrase_language = ?, definition_language = ? WHERE id = ?")
	if err != nil {
		return nil, fmt.Errorf("upate query: %w", err)
	}

	deleteStmt, err := db.Prepare("DELETE FROM study_set WHERE id = ?")
	if err != nil {
		return nil, fmt.Errorf("delete query: %w", err)
	}

	existsStmt, err := db.Prepare("SELECT EXISTS(SELECT 1 FROM study_set WHERE id = ?)")
	if err != nil {
		return nil, fmt.Errorf("exists query: %w", err)
	}

	return &studySetQueries{
		insert:        insertStmt,
		getById:       getByIdStmt,
		getAllSummary: getAllSummaryStmt,
		update:        updateStmt,
		delete:        deleteStmt,
		exists:        existsStmt,
	}, nil
}

type StudySetRepo struct {
	query *studySetQueries
}

func NewStudySetRepo(db *sql.DB) (domain.StudySetRepo, error) {
	query, err := newStudySetQueries(db)
	if err != nil {
		return nil, err
	}
	return &StudySetRepo{
		query: query,
	}, nil
}

func (r *StudySetRepo) GetAll(ctx context.Context) ([]*domain.StudySetRow, error) {
	studySets := make([]*domain.StudySetRow, 0)

	rows, err := r.query.getAllSummary.QueryContext(ctx)
	if err != nil {
		return nil, fmt.Errorf("failed to query: %w", err)
	}

	for rows.Next() {
		var studySet domain.StudySetRow
		if err := rows.Scan(&studySet.Id, &studySet.AuthorId, &studySet.Name, &studySet.Description, &studySet.PhraseLanguage, &studySet.DefinitionLanguage); err != nil {
			return nil, fmt.Errorf("failed to scan: %w", err)
		}
		studySets = append(studySets, &studySet)
	}

	return studySets, nil
}

func (r *StudySetRepo) GetById(ctx context.Context, studySetID int64) (*domain.StudySetRow, error) {
	var studySet domain.StudySetRow
	if err := r.query.getById.QueryRowContext(ctx, studySetID).Scan(&studySet.Id, &studySet.AuthorId, &studySet.Name, &studySet.Description, &studySet.PhraseLanguage, &studySet.DefinitionLanguage); err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, nil
		}
		return nil, fmt.Errorf("failed to query: %w", err)
	}
	return &studySet, nil
}

func (r *StudySetRepo) Insert(ctx context.Context, insertData *domain.InsertStudySetData) (*domain.StudySetRow, error) {
	res, err := r.query.insert.ExecContext(
		ctx,
		insertData.AuthorId,
		insertData.Name,
		insertData.Description,
		insertData.PhraseLanguage,
		insertData.DefinitionLanguage,
	)
	if err != nil {
		return nil, fmt.Errorf("failed to exec: %w", err)
	}

	lastInsertId, err := res.LastInsertId()
	if err != nil {
		return nil, fmt.Errorf("failed to get last insert id: %w", err)
	}

	return r.GetById(ctx, lastInsertId)
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
