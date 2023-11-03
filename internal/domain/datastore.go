package domain

import (
	"context"
)

type DataStore interface {
	Atomic(ctx context.Context, fn func(ds DataStore) error) error
	GetStudySetRepo() StudySetRepo
	GetDefinitionRepo() DefinitionRepo
	GetProfileRepo() ProfileRepo
	GetUserRepo() UserRepo
}
