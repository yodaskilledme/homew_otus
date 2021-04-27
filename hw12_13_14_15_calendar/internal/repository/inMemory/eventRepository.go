package inMemory

import (
	"context"
	"sync"
	"time"

	"github.com/yodaskilledme/homew_otus/hw12_13_14_15_calendar/internal/appError"
	"github.com/yodaskilledme/homew_otus/hw12_13_14_15_calendar/internal/domain"
)

type Repo struct {
	mutex      sync.RWMutex
	Events     map[uint64]domain.Event
	idSequence uint64
}

func New() *Repo {
	return &Repo{
		mutex:      sync.RWMutex{},
		Events:     make(map[uint64]domain.Event),
		idSequence: 0,
	}
}

func (r *Repo) Create(ctx context.Context, event domain.Event) (domain.Event, error) {
	const op = "EventRepository.Create"
	r.mutex.Lock()
	defer r.mutex.Unlock()

	event.DateStart = event.DateStart.Round(time.Minute)
	event.DateEnd = event.DateEnd.Round(time.Minute)
	if err := r.isTimeBusy(event); err != nil {
		return domain.Event{}, appError.OpError(op, err)
	}
	r.idSequence++
	event.ID = r.idSequence
	r.Events[event.ID] = event

	return event, nil
}

func (r *Repo) Update(ctx context.Context, event domain.Event) (domain.Event, error) {
	const op = "EventRepository.Update"
	r.mutex.Lock()
	defer r.mutex.Unlock()

	if _, ok := r.Events[event.ID]; !ok {
		return event, appError.OpError(op, domain.ErrNotFound)
	}
	if err := r.isTimeBusy(event); err != nil {
		return event, appError.OpError(op, err)
	}

	r.Events[event.ID] = event

	return event, nil
}

func (r *Repo) Delete(ctx context.Context, id uint64) error {
	const op = "EventRepository.Delete"
	r.mutex.Lock()
	defer r.mutex.Unlock()

	if _, ok := r.Events[id]; !ok {
		return appError.OpError(op, domain.ErrNotFound)
	}
	delete(r.Events, id)

	return nil
}

func (r *Repo) List(ctx context.Context, userID uint64, dateFrom, dateTo time.Time) ([]domain.Event, error) {
	const op = "EventRepository.List"

	var events []domain.Event

	for _, storedEvent := range r.Events {
		if storedEvent.UserID != userID {
			continue
		}
		if storedEvent.DateStart.Unix() >= dateFrom.Unix() && storedEvent.DateEnd.Unix() <= dateTo.Unix() {
			events = append(events, storedEvent)
		}
	}

	if len(events) == 0 {
		return events, appError.OpError(op, domain.ErrNotFound)
	}

	return events, nil
}

func (r *Repo) isTimeBusy(event domain.Event) error {
	for _, storedEvent := range r.Events {
		if storedEvent.UserID != event.UserID {
			continue
		}
		if storedEvent.DateStart.Before(event.DateEnd) && storedEvent.DateEnd.After(event.DateStart) {
			return domain.ErrTimeBusy
		}
	}
	return nil
}
