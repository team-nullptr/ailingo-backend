package studyset

import (
	"errors"
	"fmt"

	"github.com/go-playground/validator/v10"
)

type Service struct {
	studySetRepo Repo
	validate     *validator.Validate
}

// NewService creates a new instance of ServiceImpl.
func NewService(studySetRepo Repo, validate *validator.Validate) *Service {
	return &Service{
		studySetRepo: studySetRepo,
		validate:     validate,
	}
}

var (
	ErrRepoFailed = errors.New("repo failed")
	ErrValidation = errors.New("validation failed")
)

// Create is a service for creating a new study set.
// If the validation fails returns ErrValidationFailed.
func (s *Service) Create(data *StudySetCreate) (*StudySet, error) {
	if err := s.validate.Struct(data); err != nil {
		return nil, fmt.Errorf("%w: %w", ErrValidation, err)
	}

	studySet, err := s.studySetRepo.Insert(data)
	if err != nil {
		return nil, fmt.Errorf("%w: failed to create the user: %w", ErrRepoFailed, err)
	}

	return studySet, nil
}

// GetById is a service for retrieving a study set by id.
func (s *Service) GetById(id int64) (*StudySet, error) {
	studySet, err := s.studySetRepo.GetById(id)
	if err != nil {
		return nil, fmt.Errorf("%w: failed to get the user: %w", ErrRepoFailed, err)
	}

	return studySet, nil
}
