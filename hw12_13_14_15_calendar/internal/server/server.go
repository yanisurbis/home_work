package server

import (
	domain "calendar/internal/domain/interfaces"
	"context"
)

type Server interface {
	Start(storage domain.EventStorage) error
	Stop(ctx context.Context) error
}
