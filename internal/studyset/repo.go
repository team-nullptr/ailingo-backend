package studyset

import (
	"database/sql"
	"encoding/json"
	"fmt"

	"ailingo/internal/models"
)

// InsertStudySetData represents data required to create a new study set.
// TODO: Possible values for languages are `pl-PL` and `en-US`. Can we validate that in a neat way?
type InsertStudySetData struct {
	AuthorId           string        `json:"-" validate:"required"`
	Name               string        `json:"name" validate:"required,ascii,max=128"`
	Description        string        `json:"description" validate:"required,ascii,max=512"`
	PhraseLanguage     string        `json:"phraseLanguage" validate:"required"`
	DefinitionLanguage string        `json:"definitionLanguage" validate:"required"`
	Definitions        []models.Word `json:"definitions"`
}

// Repo is an interface describing methods available on study set repository.
type Repo interface {
	// Insert inserts a new study set to the db and returns the created row.
	Insert(data *InsertStudySetData) (*models.StudySet, error)

	// GetAll gets all the study sets.
	// TODO: Do we need pagination?
	GetAll() ([]*models.StudySet, error)

	// GetById gets the study set by id.
	GetById(id int64) (*models.StudySet, error)
}

// DefaultRepo is the default implementation for Repo.
type DefaultRepo struct {
	query *queryProvider
}

func NewRepo(db *sql.DB) (Repo, error) {
	query, err := newQueryProvider(db)
	if err != nil {
		return nil, err
	}

	return &DefaultRepo{
		query: query,
	}, nil
}

func (r *DefaultRepo) Insert(data *InsertStudySetData) (*models.StudySet, error) {
	definitionsJson, err := json.Marshal(data.Definitions)
	if err != nil {
		return nil, fmt.Errorf("failed to marshal study set's definitions: %w", err)
	}

	res, err := r.query.insert.Exec(
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

func (r *DefaultRepo) GetById(id int64) (*models.StudySet, error) {
	var studySet models.StudySet
	var definitionsRaw json.RawMessage

	if err := r.query.getById.QueryRow(id).Scan(&studySet.Id, &studySet.AuthorId, &studySet.Name, &studySet.Description, &studySet.PhraseLanguage, &studySet.DefinitionLanguage, &definitionsRaw); err != nil {
		return nil, fmt.Errorf("failed to query: %w", err)
	}

	if err := json.Unmarshal(definitionsRaw, &studySet.Definitions); err != nil {
		return nil, fmt.Errorf("failed to unmarshal scanned definitions: %w", err)
	}

	return &studySet, nil
}

func (r *DefaultRepo) GetAll() ([]*models.StudySet, error) {
	rows, err := r.query.getAll.Query()
	if err != nil {
		return nil, fmt.Errorf("failed to query: %w", err)
	}

	var studySets []*models.StudySet
	for rows.Next() {
		var (
			studySet       models.StudySet
			definitionsRaw json.RawMessage
		)

		if err := rows.Scan(&studySet.Id, &studySet.AuthorId, &studySet.Name, &studySet.Description, &studySet.PhraseLanguage, &studySet.DefinitionLanguage, &definitionsRaw); err != nil {
			return nil, fmt.Errorf("failed to scan the row: %w", err)
		}

		if err := json.Unmarshal(definitionsRaw, &studySet.Definitions); err != nil {
			return nil, fmt.Errorf("failed to unmarshal scanned definitions: %w", err)
		}

		studySets = append(studySets, &studySet)
	}

	return studySets, nil
}
