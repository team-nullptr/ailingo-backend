package mysql

import (
	"context"
	"errors"
	"fmt"

	"github.com/go-sql-driver/mysql"

	"ailingo/internal/domain"
)

const insertStar = `
INSERT INTO star (user_id, study_set_id) VALUES (?, ?)
`

const deleteStar = `
DELETE
FROM star
WHERE user_id = ?
  AND study_set_id = ?
`

type profileRepo struct {
	db DBTX
}

func NewProfileRepo(db DBTX) domain.ProfileRepo {
	return &profileRepo{
		db: db,
	}
}

func (r *profileRepo) InsertStar(ctx context.Context, userID string, studySetID int64) error {
	if _, err := r.db.ExecContext(ctx, insertStar, userID, studySetID); err != nil {
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

func (r *profileRepo) DeleteStar(ctx context.Context, userID string, studySetID int64) error {
	if _, err := r.db.ExecContext(ctx, deleteStar, userID, studySetID); err != nil {
		return fmt.Errorf("failed to exec: %w", err)
	}
	return nil
}
