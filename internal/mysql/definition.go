package mysql

import (
	"context"
	"database/sql"
	"encoding/json"
	"errors"
	"fmt"

	"ailingo/internal/domain"
)

// getDefinitionsForStudySet queries for all definitions connected with the given study set.
const getDefinitionsForStudySet = `
SELECT id, phrase, meaning, sentences
FROM definition
WHERE study_set_id = ?
`

// insertDefinition inserts a new definitions.
const insertDefinition = `
INSERT INTO definition (study_set_id, phrase, meaning, sentences)
VALUES (?, ?, ?, ?)
`

// updateDefinitionById updates the specified definition.
const updateDefinitionById = `
UPDATE definition
SET phrase    = ?,
    meaning   = ?,
    sentences = ?
WHERE id = ?
`

// deleteDefinitionById deletes the specified definition.
const deleteDefinitionById = `
DELETE
FROM definition
WHERE id = ?
`

type DefinitionRepo struct {
	db DBTX
}

func NewDefinitionRepo(db DBTX) domain.DefinitionRepo {
	return &DefinitionRepo{
		db: db,
	}
}

func (r *DefinitionRepo) GetAllFor(ctx context.Context, parentStudySetID int64) ([]*domain.DefinitionRow, error) {
	definitions := make([]*domain.DefinitionRow, 0)

	rows, err := r.db.QueryContext(ctx, getDefinitionsForStudySet, parentStudySetID)
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

	if _, err = r.db.ExecContext(ctx, insertDefinition, parentStudySetID, insertData.Phrase, insertData.Meaning, sentencesJson); err != nil {
		return fmt.Errorf("failed to exec: %w", err)
	}

	return nil
}

func (r *DefinitionRepo) Update(ctx context.Context, definitionID int64, updateData *domain.UpdateDefinitionData) error {
	sentencesJson, err := json.Marshal(updateData.Sentences)
	if err != nil {
		return fmt.Errorf("failed to marshal sentences array")
	}

	if _, err := r.db.ExecContext(ctx, updateDefinitionById, updateData.Phrase, updateData.Meaning, sentencesJson, definitionID); err != nil {
		return fmt.Errorf("failed to update the definition: %w", err)
	}

	return nil
}

func (r *DefinitionRepo) Delete(ctx context.Context, definitionID int64) error {
	// TODO: We could inform if any rows were removed or not.
	if _, err := r.db.ExecContext(ctx, deleteDefinitionById, definitionID); err != nil {
		return fmt.Errorf("failed to execute delete query: %w", err)
	}
	return nil
}
