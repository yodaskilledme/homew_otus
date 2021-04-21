package inMemory

import (
	"context"
	"sync"
	"testing"
	"time"

	"github.com/stretchr/testify/require"
	"github.com/yodaskilledme/homew_otus/hw12_13_14_15_calendar/internal/appError"
	"github.com/yodaskilledme/homew_otus/hw12_13_14_15_calendar/internal/domain"
)

func prepareMocks() (*Repo, context.Context) {
	r := New()
	e1 := domain.Event{
		ID:          1,
		Title:       "123",
		DateStart:   time.Now(),
		DateEnd:     time.Now().Add(time.Minute * time.Duration(30)),
		Description: "test",
		UserID:      1,
	}
	e2 := domain.Event{
		ID:          2,
		Title:       "321",
		DateStart:   time.Now().Add(time.Minute * time.Duration(31)),
		DateEnd:     time.Now().Add(time.Minute * time.Duration(61)),
		Description: "test1",
		UserID:      1,
	}
	r.Events[e1.ID] = e1
	r.Events[e2.ID] = e2
	r.idSequence = 2

	return r, context.Background()
}

func Test_Repo(t *testing.T) {
	t.Run("test create event, success", func(t *testing.T) {
		r, ctx := prepareMocks()
		event := domain.Event{
			Title:       "555",
			DateStart:   time.Now().Add(time.Minute * time.Duration(62)),
			DateEnd:     time.Now().Add(time.Minute * time.Duration(121)),
			Description: "test1",
			UserID:      1,
		}

		got, err := r.Create(ctx, event)
		require.NoError(t, err)
		require.Equal(t, 3, int(got.ID))
	})

	t.Run("create event, errTimeBusy", func(t *testing.T) {
		r, ctx := prepareMocks()
		event := domain.Event{
			Title:       "555",
			DateStart:   time.Now().Add(time.Minute * time.Duration(15)),
			DateEnd:     time.Now().Add(time.Minute * time.Duration(30)),
			Description: "test1",
			UserID:      1,
		}

		_, err := r.Create(ctx, event)
		require.Error(t, err)
		require.Equal(t, appError.OpError("EventRepository.Create", domain.ErrTimeBusy), err)
	})

	t.Run("test update event, success", func(t *testing.T) {
		r, ctx := prepareMocks()
		event := domain.Event{
			ID:          1,
			Title:       "555",
			DateStart:   time.Now().Add(time.Minute * time.Duration(62)),
			DateEnd:     time.Now().Add(time.Minute * time.Duration(121)),
			Description: "test1",
			UserID:      1,
		}

		got, err := r.Update(ctx, event)
		require.NoError(t, err)
		require.Equal(t, 1, int(got.ID))
	})

	t.Run("test update event, errTimeBusy", func(t *testing.T) {
		r, ctx := prepareMocks()
		event := domain.Event{
			ID:          1,
			Title:       "555",
			DateStart:   time.Now().Add(time.Minute * time.Duration(60)),
			DateEnd:     time.Now().Add(time.Minute * time.Duration(80)),
			Description: "test1",
			UserID:      1,
		}

		_, err := r.Update(ctx, event)
		require.Error(t, err)
		require.Equal(t, appError.OpError("EventRepository.Update", domain.ErrTimeBusy), err)
	})

	t.Run("test delete event, success", func(t *testing.T) {
		r, ctx := prepareMocks()
		var eventID uint64

		eventID = 1
		err := r.Delete(ctx, eventID)
		require.NoError(t, err)
	})

	t.Run("test delete event, errNotFound", func(t *testing.T) {
		r, ctx := prepareMocks()
		var eventID uint64

		eventID = 99
		err := r.Delete(ctx, eventID)
		require.Error(t, err)
		require.Equal(t, appError.OpError("EventRepository.Delete", domain.ErrNotFound), err)
	})

	t.Run("create event concurrent", func(t *testing.T) {
		r, ctx := prepareMocks()
		event := domain.Event{Title: "test", UserID: 2}
		var wg sync.WaitGroup
		for i := 0; i < 10; i++ {
			wg.Add(1)
			go func(i int, e domain.Event) {
				defer wg.Done()
				e.DateStart = time.Now().Add(time.Hour * time.Duration(i+i))
				e.DateEnd = time.Now().Add(time.Minute * time.Duration(i+i))
				_, err := r.Create(ctx, e)
				require.NoError(t, err)
			}(i, event)
		}
		wg.Wait()
		// assert that all 10 events has been added
		require.Equal(t, len(r.Events), 12)
	})
}
