package server

import (
	domain2 "calendar/internal/domain/services"
	"context"
)

type Server interface {
	Start(eventService domain2.EventService) error
	Stop(ctx context.Context) error
}
