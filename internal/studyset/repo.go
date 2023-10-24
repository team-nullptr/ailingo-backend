package studyset

import (
	"database/sql"
	"encoding/json"
	"fmt"
)

// Repo is an interface describing methods available on study set repository.
type Repo interface {
	// Insert inserts a new study set to the db and returns the created row.
	Insert(data *StudySetCreate) (*StudySet, error)
	// GetById gets the study set by id
	GetById(id int64) (*StudySet, error)
}

// RepoImpl is the default implementation for Repo.
type RepoImpl struct {
	db *sql.DB
}

func NewRepo(db *sql.DB) Repo {
	return &RepoImpl{
		db: db,
	}
}

func (r *RepoImpl) Insert(data *StudySetCreate) (*StudySet, error) {
	query := "INSERT INTO study_sets (author_id, name, description, phrase_language, definition_language, definitions) VALUES (?, ?, ?, ?, ?, ?);"

	stmt, err := r.db.Prepare(query)
	if err != nil {
		return nil, fmt.Errorf("failed to prepare the query: %w", err)
	}

	definitionsJson, err := json.Marshal(data.Definitions)
	if err != nil {
		return nil, fmt.Errorf("failed to marshal study set's definitions: %w", err)
	}

	res, err := stmt.Exec(
		data.AuthorId,
		data.Name,
		data.Description,
		data.PhraseLanguage,
		data.DefinitionLanguage,
		definitionsJson,
	)
	if err != nil {
		return nil, fmt.Errorf("failed to exec the query: %w", err)
	}

	lastInsertId, err := res.LastInsertId()
	if err != nil {
		return nil, fmt.Errorf("failed to get last insert id: %w", err)
	}

	return r.GetById(lastInsertId)
}

func (r *RepoImpl) GetById(id int64) (*StudySet, error) {
	query := "SELECT * FROM study_sets WHERE id = ?"

	stmt, err := r.db.Prepare(query)
	if err != nil {
		return nil, fmt.Errorf("failed to prepare the query: %w", err)
	}

	var definitionsRaw json.RawMessage
	var studySet StudySet

	if err := stmt.QueryRow(id).Scan(&studySet.Id, &studySet.AuthorId, &studySet.Name, &studySet.Description, &studySet.PhraseLanguage, &studySet.DefinitionLanguage, &definitionsRaw); err != nil {
		return nil, fmt.Errorf("failed to query: %w", err)
	}

	if err := json.Unmarshal(definitionsRaw, &studySet.Definitions); err != nil {
		return nil, fmt.Errorf("failed to unmarshal scanned definitions: %w", err)
	}

	return &studySet, nil
}
