package config

import (
	"os"

	"github.com/yodaskilledme/homew_otus/hw12_13_14_15_calendar/internal/appError"
	"gopkg.in/yaml.v3"
)

type Http struct {
	Host string `yaml:"host"`
	Port string `yaml:"port"`
}

type Database struct {
	Host     string `yaml:"host"`
	Port     string `yaml:"port"`
	Name     string `yaml:"name"`
	User     string `yaml:"user"`
	Password string `yaml:"password"`
}

type Logger struct {
	Level  string `yaml:"level"`
	Output string `yaml:"output"`
}

type Grpc struct {
	Host string `yaml:"host"`
	Port string `yaml:"port"`
}

type Amqp struct {
	Uri   string `yaml:"uri"`
	Queue string `yaml:"queue"`
}

type Storage struct {
	Type string `yaml:"type"`
}

type Config struct {
	Http     `yaml:"http"`
	Database `yaml:"database"`
	Grpc     `yaml:"grpc"`
	Amqp     `yaml:"amqp"`
	Logger   `yaml:"logger"`
	Storage  `yaml:"storage"`
}

func ParseConfig(path string) (*Config, error) {
	const op = "config.parseConfig"

	f, err := os.Open(path)
	if err != nil {
		return nil, &appError.AppError{
			Op:      op,
			Err:     err,
			Message: "failed to parse config",
			Code:    "",
		}
	}
	defer f.Close()

	var cfg Config
	decoder := yaml.NewDecoder(f)
	err = decoder.Decode(&cfg)
	if err != nil {
		return nil, &appError.AppError{
			Op:      op,
			Err:     err,
			Message: "failed to parse config",
			Code:    "",
		}
	}

	return &cfg, nil
}
