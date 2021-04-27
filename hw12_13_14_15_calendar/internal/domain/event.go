package domain

import (
	"errors"
	"time"
)

var (
	ErrTimeBusy = errors.New("event's time busy")
	ErrNotFound = errors.New("event not found")
)

type Event struct {
	ID          uint64
	Title       string
	DateStart   time.Time
	DateEnd     time.Time
	Description string
	UserID      uint64
}
