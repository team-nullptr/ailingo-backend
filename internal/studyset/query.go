package studyset

import (
	"database/sql"
	"fmt"
)

// queryProvider provides all the necessary queries used by study set repo.
// To keep things simple the queryProvider does not expose queryProvider.Close function
// as the queryProvider is meant to be long-lived.
type queryProvider struct {
	insert        *sql.Stmt
	getById       *sql.Stmt
	getAllSummary *sql.Stmt
	update        *sql.Stmt
	delete        *sql.Stmt
}

// newQueryProvider creates a new prepared query provider for study set db operations.
func newQueryProvider(db *sql.DB) (*queryProvider, error) {
	insertStmt, err := db.Prepare("INSERT INTO study_sets (author_id, name, description, phrase_language, definition_language, definitions) VALUES (?, ?, ?, ?, ?, ?);")
	if err != nil {
		return nil, fmt.Errorf("failed to prepare insert query: %w", err)
	}

	getByIdStmt, err := db.Prepare("SELECT id, author_id, name, description, phrase_language, definition_language, definitions FROM study_sets WHERE id = ?")
	if err != nil {
		return nil, fmt.Errorf("faile to prepare getById query: %w", err)
	}

	getAllSummaryStmt, err := db.Prepare("SELECT id, author_id, name, description, phrase_language, definition_language FROM study_sets")
	if err != nil {
		return nil, fmt.Errorf("failed to prepare getAll query: %w", err)
	}

	updateStmt, err := db.Prepare("UPDATE study_sets SET name = ?, description = ?, phrase_language = ?, definition_language = ?, definitions = ? WHERE id = ?")
	if err != nil {
		return nil, fmt.Errorf("failed to prepare upate query: %w", err)
	}

	deleteStmt, err := db.Prepare("DELETE FROM study_sets WHERE id = ?")
	if err != nil {
		return nil, fmt.Errorf("failed to prepare delete query: %w", err)
	}

	return &queryProvider{
		insert:        insertStmt,
		getById:       getByIdStmt,
		getAllSummary: getAllSummaryStmt,
		update:        updateStmt,
		delete:        deleteStmt,
	}, nil
}
