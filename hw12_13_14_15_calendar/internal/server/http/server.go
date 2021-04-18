package internalhttp

import (
	"fmt"
	"net/http"
	"time"

	"github.com/yodaskilledme/homew_otus/hw12_13_14_15_calendar/internal/appLogger"
	"github.com/yodaskilledme/homew_otus/hw12_13_14_15_calendar/internal/config"
)

func New(config config.Config, logger *appLogger.Logger) *http.Server {
	handler := http.HandlerFunc(handle)
	http.Handle("/", logMiddleware(handler, logger))
	return &http.Server{
		Addr: config.Http.Host + ":" + config.Http.Port,
	}
}

func handle(w http.ResponseWriter, r *http.Request) {
	time.Sleep(1 * time.Second)
	_, _ = fmt.Fprintf(w, "Hello!")
}

func logMiddleware(h http.Handler, logger *appLogger.Logger) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		startTime := time.Now()
		h.ServeHTTP(w, r)
		latency := time.Since(startTime)

		info := fmt.Sprintf("%s [%s] %s %s %s %d %s \"%s\"",
			r.RemoteAddr,
			time.Now().Format("2006-01-02 15:04:05 -0700"),
			r.Method,
			r.URL.Path,
			r.Proto,
			http.StatusOK,
			latency,
			r.UserAgent(),
		)
		logger.Info(info)
	})
}
