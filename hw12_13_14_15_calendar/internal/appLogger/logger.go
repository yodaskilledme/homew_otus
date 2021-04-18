package appLogger

import (
	"log"

	"go.uber.org/zap"
)

type Logger struct {
	Zap *zap.Logger
}

func New(logOutput string) *Logger {
	logConf := zap.NewProductionConfig()
	logConf.OutputPaths = []string{logOutput}

	zapLogger, err := logConf.Build()
	if err != nil {
		log.Fatalf("can't initialize zap logger: %v", err)
	}

	return &Logger{
		Zap: zapLogger,
	}
}

func (l *Logger) Error(msg string) {
	l.Zap.Error(msg)
}

func (l *Logger) Info(msg string) {
	l.Zap.Info(msg)
}

func (l *Logger) Debug(msg string) {
	l.Zap.Debug(msg)
}
