package sql

import (
	"context"
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
		query = `INSERT INTO events (title, date_start, date_end, description, user_id)
SELECT ?, ?, ?, ?, ?
    WHERE NOT exists(SELECT 1
FROM events 
WHERE user_id = ?
AND date_start < ?
AND date_end > ?)
ON CONFLICT DO NOTHING RETURNING id;`
	)
	var id int

	tx, err := r.db.Beginx()
	if err != nil {
		return -1, appError.OpError(op, err)
	}
	defer func() {
		_ = tx.Rollback()
	}()

	err = tx.QueryRowContext(ctx, query, event.Title, event.DateStart, event.DateEnd, event.Description, event.UserID).Scan(&id)
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
		query = `UPDATE events (title, date_start, date_end, description, user_id) 
VALUES (?, ?, ?, ?, ?) 
WHERE id = ? 
AND NOT EXISTS(
SELECT 
FROM events 
WHERE user_id = ? 
AND date_start < ?
AND date_end > ?)`
	)
	tx, err := r.db.Beginx()
	if err != nil {
		return domain.Event{}, appError.OpError(op, err)
	}
	defer func() {
		_ = tx.Rollback()
	}()

	res, err := tx.ExecContext(
		ctx,
		query,
		event.Title,
		event.DateStart,
		event.DateEnd,
		event.Description,
		event.UserID,
		event.ID,
		event.UserID,
		event.DateEnd,
		event.DateStart)
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

	tx, err := r.db.Beginx()
	if err != nil {
		return appError.OpError(op, err)
	}
	defer func() {
		_ = tx.Rollback()
	}()

	res, err := tx.ExecContext(ctx, query, id)
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

	tx, err := r.db.Beginx()
	if err != nil {
		return nil, appError.OpError(op, err)
	}
	defer func() {
		_ = tx.Rollback()
	}()

	var events []domain.Event
	rows, err := tx.QueryContext(ctx, query, userID, dateFrom, dateTo)
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
