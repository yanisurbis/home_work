package domain

import (
	"errors"
)

var (
	ErrForbidden = errors.New("access denied")
	ErrNotFound = errors.New("not found")
)