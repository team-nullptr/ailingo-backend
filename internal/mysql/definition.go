package mysql

import (
	"context"
	"database/sql"
	"encoding/json"
	"fmt"

	"ailingo/internal/domain"
)

type definitionQueries struct {
	getAll *sql.Stmt
	insert *sql.Stmt
}

func newDefinitionQueries(db *sql.DB) (*definitionQueries, error) {
	getAllStmt, err := db.Prepare("SELECT id, phrase, meaning, sentences FROM definition WHERE study_set_id = ?")
	if err != nil {
		return nil, fmt.Errorf("getAll query: %w", err)
	}

	insertStmt, err := db.Prepare(`INSERT INTO definition (study_set_id, phrase, meaning, sentences) VALUES (?, ?, ?, ?)`)
	if err != nil {
		return nil, fmt.Errorf("insert query: %w", err)
	}

	return &definitionQueries{
		getAll: getAllStmt,
		insert: insertStmt,
	}, nil
}

type DefinitionRepo struct {
	db    *sql.DB
	query *definitionQueries
}

func NewDefinitionRepo(db *sql.DB) (domain.DefinitionRepo, error) {
	query, err := newDefinitionQueries(db)
	if err != nil {
		return nil, err
	}

	return &DefinitionRepo{
		db:    db,
		query: query,
	}, nil
}

func (r *DefinitionRepo) GetAllFor(ctx context.Context, studySetID int64) ([]*domain.Definition, error) {
	rows, err := r.query.getAll.QueryContext(ctx, studySetID)
	if err != nil {
		return nil, fmt.Errorf("failed to query: %w", err)
	}

	definitions := make([]*domain.Definition, 0)

	for rows.Next() {
		var (
			definition   domain.Definition
			sentencesRaw json.RawMessage
		)

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

func (r *DefinitionRepo) Insert(ctx context.Context, studySetID int64, insertData *domain.InsertDefinitionData) error {
	sentencesJson, err := json.Marshal(insertData.Sentences)
	if err != nil {
		return fmt.Errorf("failed to marshal sentences array")
	}

	if _, err = r.query.insert.ExecContext(ctx, studySetID, insertData.Phrase, insertData.Meaning, sentencesJson); err != nil {
		return fmt.Errorf("failed to exec: %w", err)
	}

	return nil
}
