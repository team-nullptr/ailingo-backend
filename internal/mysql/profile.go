package mysql

import (
	"context"
	"database/sql"
	"errors"
	"fmt"

	"github.com/go-sql-driver/mysql"
)

type profileQueries struct {
	insertStar *sql.Stmt
	deleteStar *sql.Stmt
}

func newProfileQueries(ctx context.Context, db *sql.DB) (*profileQueries, error) {
	insertStarStmt, err := db.PrepareContext(ctx, "INSERT INTO star (user_id, study_set_id) VALUE (?, ?)")
	if err != nil {
		return nil, fmt.Errorf("insertStar query: %w", err)
	}

	deleteStarStmt, err := db.PrepareContext(ctx, "DELETE FROM star WHERE user_id = ? AND study_set_id = ?")
	if err != nil {
		return nil, fmt.Errorf("delete query: %w", err)
	}

	return &profileQueries{
		insertStar: insertStarStmt,
		deleteStar: deleteStarStmt,
	}, nil
}

type ProfileRepo struct {
	query *profileQueries
}

func NewProfileRepo(ctx context.Context, db *sql.DB) (*ProfileRepo, error) {
	query, err := newProfileQueries(ctx, db)
	if err != nil {
		return nil, err
	}

	return &ProfileRepo{
		query: query,
	}, nil
}

func (r *ProfileRepo) InsertStar(ctx context.Context, userID string, studySetID int64) error {
	if _, err := r.query.insertStar.ExecContext(ctx, userID, studySetID); err != nil {
		var mysqlerr *mysql.MySQLError
		if errors.As(err, &mysqlerr) {
			// 1062 error number stands for duplicate entry code
			if mysqlerr.Number == 1062 {
				return ErrDuplicateRow
			}
		}
		return fmt.Errorf("failed to exec: %w", err)
	}
	return nil
}

func (r *ProfileRepo) DeleteStar(ctx context.Context, userID string, studySetID int64) error {
	if _, err := r.query.deleteStar.ExecContext(ctx, userID, studySetID); err != nil {
		return fmt.Errorf("failed to exec: %w", err)
	}
	return nil
}
