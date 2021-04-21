package domain

import (
	"errors"
	"time"
)

var (
	ErrTimeBusy = errors.New("event's time busy")
	ErrNotFound = errors.New("event not found")
)

type EventID uint64

type UserID uint64

type UserTakenTime map[uint64]map[time.Time]struct{}

type Event struct {
	ID          uint64
	Title       string
	DateStart   time.Time
	DateEnd     time.Time
	Description string
	UserID      uint64
}
