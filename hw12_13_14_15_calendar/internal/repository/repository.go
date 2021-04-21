package repository

import (
	"context"
	"time"

	"github.com/yodaskilledme/homew_otus/hw12_13_14_15_calendar/internal/domain"
)

type EventsRepo interface {
	Create(ctx context.Context, event domain.Event) (domain.Event, error)
	Update(ctx context.Context, event domain.Event) (domain.Event, error)
	Delete(ctx context.Context, id uint64) error
	List(ctx context.Context, userID uint64, dateFrom, dateTo time.Time) ([]domain.Event, error)
}
