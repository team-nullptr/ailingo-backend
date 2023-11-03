package mysql

import "errors"

var (
	ErrDuplicateRow = errors.New("duplicate row")
)
