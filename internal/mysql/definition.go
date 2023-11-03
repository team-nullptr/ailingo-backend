package mysql

import (
	"context"
	"database/sql"
	"encoding/json"
	"errors"
	"fmt"

	"ailingo/internal/domain"
)

type definitionQueries struct {
	getAll *sql.Stmt
	insert *sql.Stmt
	update *sql.Stmt
	delete *sql.Stmt
}

func newDefinitionQueries(ctx context.Context, db *sql.DB) (*definitionQueries, error) {
	getAllStmt, err := db.PrepareContext(ctx, "SELECT id, phrase, meaning, sentences FROM definition WHERE study_set_id = ?")
	if err != nil {
		return nil, fmt.Errorf("getAll query: %w", err)
	}

	insertStmt, err := db.PrepareContext(ctx, `INSERT INTO definition (study_set_id, phrase, meaning, sentences) VALUES (?, ?, ?, ?)`)
	if err != nil {
		return nil, fmt.Errorf("insert query: %w", err)
	}

	updateStmt, err := db.PrepareContext(ctx, "UPDATE definition SET phrase = ?, meaning = ?, sentences = ? WHERE study_set_id = ? AND id = ?")
	if err != nil {
		return nil, fmt.Errorf("update query: %w", err)
	}

	deleteStmt, err := db.PrepareContext(ctx, "DELETE FROM definition WHERE study_set_id = ? AND id = ?")
	if err != nil {
		return nil, fmt.Errorf("delete query: %w", err)
	}

	return &definitionQueries{
		getAll: getAllStmt,
		insert: insertStmt,
		update: updateStmt,
		delete: deleteStmt,
	}, nil
}

type DefinitionRepo struct {
	db    *sql.DB
	query *definitionQueries
}

func NewDefinitionRepo(ctx context.Context, db *sql.DB) (domain.DefinitionRepo, error) {
	query, err := newDefinitionQueries(ctx, db)
	if err != nil {
		return nil, err
	}
	return &DefinitionRepo{
		db:    db,
		query: query,
	}, nil
}

func (r *DefinitionRepo) GetAllFor(ctx context.Context, parentStudySetID int64) ([]*domain.DefinitionRow, error) {
	definitions := make([]*domain.DefinitionRow, 0)

	rows, err := r.query.getAll.QueryContext(ctx, parentStudySetID)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return definitions, nil
		}
		return nil, fmt.Errorf("failed to query: %w", err)
	}

	for rows.Next() {
		var definition domain.DefinitionRow
		var sentencesRaw json.RawMessage

		if err := rows.Scan(&definition.Id, &definition.Phrase, &definition.Meaning, &sentencesRaw); err != nil {
			return nil, fmt.Errorf("failed to scan: %w", err)
		}

		if err := json.Unmarshal(sentencesRaw, &definition.Sentences); err != nil {
			return nil, fmt.Errorf("failed to unmarshal sentences: %w", err)
		}

		definitions = append(definitions, &definition)
	}

	return definitions, nil
}

func (r *DefinitionRepo) Insert(ctx context.Context, parentStudySetID int64, insertData *domain.InsertDefinitionData) error {
	sentencesJson, err := json.Marshal(insertData.Sentences)
	if err != nil {
		return fmt.Errorf("failed to marshal sentences array")
	}

	if _, err = r.query.insert.ExecContext(ctx, parentStudySetID, insertData.Phrase, insertData.Meaning, sentencesJson); err != nil {
		return fmt.Errorf("failed to exec: %w", err)
	}

	return nil
}

func (r *DefinitionRepo) Update(ctx context.Context, parentStudySetID int64, definitionID int64, updateData *domain.UpdateDefinitionData) error {
	sentencesJson, err := json.Marshal(updateData.Sentences)
	if err != nil {
		return fmt.Errorf("failed to marshal sentences array")
	}

	if _, err := r.query.update.ExecContext(ctx, updateData.Phrase, updateData.Meaning, sentencesJson, parentStudySetID, definitionID); err != nil {
		return fmt.Errorf("failed to update the definition: %w", err)
	}

	return nil
}

func (r *DefinitionRepo) Delete(ctx context.Context, parentStudySetID int64, definitionID int64) error {
	// TODO: We could inform if any rows were removed or not.
	if _, err := r.query.delete.ExecContext(ctx, parentStudySetID, definitionID); err != nil {
		return fmt.Errorf("failed to execute delete query: %w", err)
	}
	return nil
}
