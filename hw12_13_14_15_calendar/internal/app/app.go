package app

import (
	"context"
	"fmt"
	"log"
	"net/http"
	"os"
	"os/signal"
	"time"

	"github.com/jmoiron/sqlx"
	"github.com/yodaskilledme/homew_otus/hw12_13_14_15_calendar/internal/appLogger"
	"github.com/yodaskilledme/homew_otus/hw12_13_14_15_calendar/internal/config"
	"github.com/yodaskilledme/homew_otus/hw12_13_14_15_calendar/internal/repository/inMemory"
	"github.com/yodaskilledme/homew_otus/hw12_13_14_15_calendar/internal/repository/sql"
	internalhttp "github.com/yodaskilledme/homew_otus/hw12_13_14_15_calendar/internal/server/http"
)

type App struct {
	Logger *appLogger.Logger
	Repo   interface{}
	Server *http.Server
}

const ctxCancelTime = 5 * time.Second

func New(config config.Config) *App {
	var repo interface{}

	switch config.Storage.Type {
	case "sql":
		conn, err := sqlx.Connect("postgres", fmt.Sprintf("host=%s port=%s user=%s password=%s dbname=%s sslmode=disable",
			config.Database.Host, config.Database.Port, config.Database.User, config.Database.Password, config.Database.Name))
		if err != nil {
			log.Fatalln(err)
		}

		repo = sql.New(conn)
	default:
		repo = inMemory.New()
	}

	logger := appLogger.New(config.Logger.Output)
	server := internalhttp.New(config, logger)

	return &App{
		Logger: logger,
		Repo:   repo,
		Server: server,
	}
}

func (a *App) Run() error {
	eChan := make(chan error)
	sigChan := make(chan os.Signal, 1)

	go func() {
		err := a.Server.ListenAndServe()
		if err != nil {
			eChan <- err
		}
	}()

	signal.Notify(sigChan, os.Interrupt)

	select {
	case err := <-eChan:
		return err
	case <-sigChan:
	}

	ctx, cancelFn := context.WithTimeout(context.Background(), ctxCancelTime)
	defer cancelFn()
	if err := a.Server.Shutdown(ctx); err != nil {
		return err
	}

	return nil
}
