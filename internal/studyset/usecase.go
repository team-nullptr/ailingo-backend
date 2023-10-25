package studyset

import (
	"errors"
	"fmt"

	"github.com/go-playground/validator/v10"

	"ailingo/internal/models"
)

type UseCase interface {
	// Create creates a new study set. It is responsible for data validation.
	Create(data *InsertStudySetData) (*models.StudySet, error)

	// GetById gets a study set by its id.
	GetById(id int64) (*models.StudySet, error)

	// GetAll gets all the created study sets.
	GetAll() ([]*models.StudySet, error)
}

type DefaultUseCase struct {
	studySetRepo Repo
	validate     *validator.Validate
}

// NewUseCase creates a new instance of StudySetUseCaseImpl.
func NewUseCase(studySetRepo Repo, validate *validator.Validate) UseCase {
	return &DefaultUseCase{
		studySetRepo: studySetRepo,
		validate:     validate,
	}
}

var (
	ErrRepoFailed = errors.New("repo failed")
	ErrValidation = errors.New("validation failed")
)

func (uc *DefaultUseCase) Create(data *InsertStudySetData) (*models.StudySet, error) {
	if err := uc.validate.Struct(data); err != nil {
		return nil, fmt.Errorf("%w: %w", ErrValidation, err)
	}
	studySet, err := uc.studySetRepo.Insert(data)
	if err != nil {
		return nil, fmt.Errorf("%w: failed to create the study set: %w", ErrRepoFailed, err)
	}
	return studySet, nil
}

func (uc *DefaultUseCase) GetById(id int64) (*models.StudySet, error) {
	studySet, err := uc.studySetRepo.GetById(id)
	if err != nil {
		return nil, fmt.Errorf("%w: failed to get the user: %w", ErrRepoFailed, err)
	}
	return studySet, nil
}

func (uc *DefaultUseCase) GetAll() ([]*models.StudySet, error) {
	studySets, err := uc.studySetRepo.GetAll()
	if err != nil {
		return nil, fmt.Errorf("%w: failed to get study sets: %w", ErrRepoFailed, err)
	}
	return studySets, nil
}
