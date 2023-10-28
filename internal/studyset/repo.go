package studyset

import (
	"context"
	"database/sql"
	"encoding/json"
	"errors"
	"fmt"

	"ailingo/internal/models"
)

type insertStudySetData struct {
	AuthorId           string        `json:"-" validate:"required"`
	Name               string        `json:"name" validate:"required,ascii,max=128"`
	Description        string        `json:"description" validate:"required,ascii,max=512"`
	PhraseLanguage     string        `json:"phraseLanguage" validate:"required"`
	DefinitionLanguage string        `json:"definitionLanguage" validate:"required"`
	Definitions        []models.Word `json:"definitions"`
}

type updateStudySetData struct {
	Name               string        `json:"name" validate:"required,ascii,max=128"`
	Description        string        `json:"description" validate:"required,ascii,max=512"`
	PhraseLanguage     string        `json:"phraseLanguage" validate:"required"`
	DefinitionLanguage string        `json:"definitionLanguage" validate:"required"`
	Definitions        []models.Word `json:"definitions"`
}

// Repo is an interface describing methods available on study set repository.
type Repo interface {
	// GetAllSummary gets all the study sets without heavy fields like `definitions`.
	GetAllSummary(ctx context.Context) ([]*models.StudySetSummary, error)
	// GetById gets the study set by id.
	GetById(ctx context.Context, studySetID int64) (*models.StudySet, error)
	// Insert inserts a new study set to the db.
	Insert(ctx context.Context, insertData *insertStudySetData) (*models.StudySet, error)
	// Update updates the study set.
	Update(ctx context.Context, studySetID int64, updateData *updateStudySetData) error
	// Delete deletes the study set with the given id.
	Delete(ctx context.Context, studySetID int64) error
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

func (r *DefaultRepo) GetAllSummary(ctx context.Context) ([]*models.StudySetSummary, error) {
	rows, err := r.query.getAllSummary.QueryContext(ctx)
	if err != nil {
		return nil, fmt.Errorf("failed to query: %w", err)
	}

	studySets := make([]*models.StudySetSummary, 0)

	for rows.Next() {
		var studySet models.StudySetSummary
		if err := rows.Scan(&studySet.Id, &studySet.AuthorId, &studySet.Name, &studySet.Description, &studySet.PhraseLanguage, &studySet.DefinitionLanguage); err != nil {
			return nil, fmt.Errorf("failed to scan the row: %w", err)
		}

		studySets = append(studySets, &studySet)
	}

	return studySets, nil
}

func (r *DefaultRepo) GetById(ctx context.Context, studySetID int64) (*models.StudySet, error) {
	var (
		studySet       models.StudySet
		definitionsRaw json.RawMessage
	)

	if err := r.query.getById.QueryRowContext(ctx, studySetID).Scan(&studySet.Id, &studySet.AuthorId, &studySet.Name, &studySet.Description, &studySet.PhraseLanguage, &studySet.DefinitionLanguage, &definitionsRaw); err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, nil
		}
		return nil, fmt.Errorf("failed to query: %w", err)
	}

	if err := json.Unmarshal(definitionsRaw, &studySet.Definitions); err != nil {
		return nil, fmt.Errorf("failed to unmarshal scanned definitions: %w", err)
	}

	return &studySet, nil
}

func (r *DefaultRepo) Insert(ctx context.Context, insertData *insertStudySetData) (*models.StudySet, error) {
	definitionsJson, err := json.Marshal(insertData.Definitions)
	if err != nil {
		return nil, fmt.Errorf("failed to marshal study set's definitions: %w", err)
	}

	res, err := r.query.insert.ExecContext(
		ctx,
		insertData.AuthorId,
		insertData.Name,
		insertData.Description,
		insertData.PhraseLanguage,
		insertData.DefinitionLanguage,
		definitionsJson,
	)
	if err != nil {
		return nil, fmt.Errorf("failed to exec the query: %w", err)
	}

	lastInsertId, err := res.LastInsertId()
	if err != nil {
		return nil, fmt.Errorf("failed to get last insert id: %w", err)
	}

	return r.GetById(ctx, lastInsertId)
}

func (r *DefaultRepo) Update(ctx context.Context, studySetID int64, updateData *updateStudySetData) error {
	definitionsJson, err := json.Marshal(updateData.Definitions)
	if err != nil {
		return fmt.Errorf("failed to marshal study set's definitions: %w", err)
	}

	if _, err := r.query.update.ExecContext(
		ctx,
		updateData.Name,
		updateData.Description,
		updateData.PhraseLanguage,
		updateData.DefinitionLanguage,
		definitionsJson,
		studySetID,
	); err != nil {
		return fmt.Errorf("failed to exec the query: %w", err)
	}

	return nil
}

func (r *DefaultRepo) Delete(ctx context.Context, studySetID int64) error {
	if _, err := r.query.delete.ExecContext(ctx, studySetID); err != nil {
		return fmt.Errorf("failed to exec the query: %w", err)
	}

	return nil
}
