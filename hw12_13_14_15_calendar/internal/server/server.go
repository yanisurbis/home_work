package server

import (
	domain2 "calendar/internal/domain/services"
	"context"
)

type Server interface {
	Start(eventService domain2.EventService, address string) error
	Stop(ctx context.Context) error
}
