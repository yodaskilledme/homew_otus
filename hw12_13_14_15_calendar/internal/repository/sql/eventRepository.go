package sql

import (
	"context"
	"database/sql"
	"time"

	"github.com/jmoiron/sqlx"
	"github.com/yodaskilledme/homew_otus/hw12_13_14_15_calendar/internal/appError"
	"github.com/yodaskilledme/homew_otus/hw12_13_14_15_calendar/internal/domain"
)

type Repo struct {
	db *sqlx.DB
}

func New(db *sqlx.DB) *Repo {
	return &Repo{
		db: db,
	}
}

func (r *Repo) Create(ctx context.Context, event domain.Event) (int, error) {
	const (
		op    = "EventRepository.Create"
		query = `INSERT INTO events (title, date_start, date_end, description, user_id) values (?, ?, ?, ?, ?) ON CONFLICT DO NOTHING RETURNING id`
	)
	var id int

	tx, err := r.db.Begin()
	if err != nil {
		return -1, appError.OpError(op, err)
	}

	ok, err := r.isTimeBusy(ctx, event)
	if err != nil {
		return -1, appError.OpError(op, err)
	}
	if !ok {
		return -1, appError.OpError(op, domain.ErrTimeBusy)
	}

	err = r.db.QueryRowContext(ctx, query, event.Title, event.DateStart, event.DateEnd, event.Description, event.UserID).Scan(&id)
	if err != nil {
		return -1, appError.OpError(op, err)
	}

	if err := tx.Commit(); err != nil {
		return -1, appError.OpError(op, err)
	}

	return id, nil
}

func (r *Repo) Update(ctx context.Context, event domain.Event) (domain.Event, error) {
	const (
		op    = "EventRepository.Update"
		query = "UPDATE events (title, date_start, date_end, description, user_id) values (?, ?, ?, ?, ?) WHERE id = ?"
	)

	tx, err := r.db.Begin()
	if err != nil {
		return domain.Event{}, appError.OpError(op, err)
	}

	ok, err := r.isTimeBusy(ctx, event)
	if err != nil {
		return domain.Event{}, appError.OpError(op, err)
	}
	if !ok {
		return domain.Event{}, appError.OpError(op, domain.ErrTimeBusy)
	}

	res, err := r.db.ExecContext(ctx, query, event.Title, event.DateStart, event.DateEnd, event.Description, event.UserID, event.ID)
	if err != nil {
		return domain.Event{}, appError.OpError(op, err)
	}
	affectedRows, err := res.RowsAffected()
	if err != nil {
		return domain.Event{}, appError.OpError(op, err)
	}
	if 0 == affectedRows {
		return domain.Event{}, appError.OpError(op, err)
	}

	if err := tx.Commit(); err != nil {
		return domain.Event{}, appError.OpError(op, err)
	}

	return event, nil
}

func (r *Repo) Delete(ctx context.Context, id uint64) error {
	const (
		op    = "EventRepository.Delete"
		query = "DELETE FROM events WHERE id = ?"
	)

	tx, err := r.db.Begin()
	if err != nil {
		return appError.OpError(op, err)
	}

	res, err := r.db.ExecContext(ctx, query, id)
	if err != nil {
		return appError.OpError(op, err)
	}
	affectedRows, err := res.RowsAffected()
	if err != nil {
		return appError.OpError(op, err)
	}
	if 0 == affectedRows {
		return appError.OpError(op, domain.ErrNotFound)
	}

	if err := tx.Commit(); err != nil {
		return appError.OpError(op, err)
	}

	return nil
}

func (r *Repo) List(ctx context.Context, userID uint64, dateFrom, dateTo time.Time) ([]domain.Event, error) {
	const (
		op    = "EventRepository.List"
		query = "SELECT id, title, date_start, date_end, description, user_id F" +
			"ROM events WHERE date_start >= ? " +
			"AND date_end <= ?" +
			"AND user_id = ?"
	)

	tx, err := r.db.Begin()
	if err != nil {
		return nil, appError.OpError(op, err)
	}

	var events []domain.Event
	rows, err := r.db.QueryContext(ctx, query, userID, dateFrom, dateTo)
	if err != nil {
		return nil, appError.OpError(op, err)
	}
	defer rows.Close()

	var event domain.Event
	for rows.Next() {
		if err := rows.Scan(&event); err != nil {
			return nil, appError.OpError(op, err)
		}

		events = append(events, event)
	}

	if err := tx.Commit(); err != nil {
		return nil, appError.OpError(op, err)
	}

	return events, nil
}

func (r *Repo) isTimeBusy(ctx context.Context, event domain.Event) (bool, error) {
	const (
		op    = "EventRepository.isTimeBusy"
		query = `SELECT exists
FROM events 
WHERE user_id = ?
AND date_start < ?
AND date_end > ?`
	)

	var exists bool
	err := r.db.QueryRowContext(ctx, query, event.UserID, event.DateEnd, event.DateStart).Scan(exists)
	if err != nil && err != sql.ErrNoRows {
		return false, appError.OpError(op, err)
	}

	return exists, nil
}
