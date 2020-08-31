package server

import (
	"calendar/internal/repository"
	"context"
)

type Server interface {
	Start(repo repository.BaseRepo) error
	Stop(ctx context.Context) error
}
