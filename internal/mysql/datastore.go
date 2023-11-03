package mysql

import (
	"context"
	"database/sql"
	"fmt"

	"ailingo/internal/domain"
)

type dataStore struct {
	// conn is the original *sql.DB connection
	conn *sql.DB
	// db is either the original *sql.DB or created *sql.Tx
	db DBTX
}

func NewDataStore(db *sql.DB) domain.DataStore {
	return &dataStore{
		conn: db,
		db:   db,
	}
}

func (ds *dataStore) GetStudySetRepo() domain.StudySetRepo {
	return NewStudySetRepo(ds.db)
}

func (ds *dataStore) GetDefinitionRepo() domain.DefinitionRepo {
	return NewDefinitionRepo(ds.db)
}

func (ds *dataStore) GetProfileRepo() domain.ProfileRepo {
	return NewProfileRepo(ds.db)
}

func (ds *dataStore) GetUserRepo() domain.UserRepo {
	return NewUserRepo(ds.db)
}

func (ds *dataStore) Atomic(ctx context.Context, cb func(ds domain.DataStore) error) error {
	var err error

	tx, err := ds.conn.BeginTx(ctx, nil)
	if err != nil {
		return fmt.Errorf("failed to begin a tx: %w", err)
	}

	defer func() {
		// TODO: In the given article there is something more in this function
		// https://blog.devgenius.io/go-golang-clean-architecture-repositories-vs-transactions-9b3b7c953463
		err = tx.Commit()
	}()

	newStore := &dataStore{
		conn: ds.conn,
		db:   tx,
	}

	err = cb(newStore)

	return err
}
