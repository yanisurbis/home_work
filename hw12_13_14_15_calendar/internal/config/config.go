package config

import (
	"os"

	"github.com/BurntSushi/toml"
)

func read(fpath string) (c Config, err error) {
	_, err = toml.DecodeFile(fpath, &c)

	return
}

func GetConfig() (*Config, error) {
	// TODO: encapsulate probably
	env := os.Getenv("ENV")
	path := ""
	if env == "TEST" {
		path = "../configs/local.toml"
	} else {
		path = "./configs/local.toml"
	}

	c, err := read(path)
	if err != nil {
		return nil, err
	}

	return &c, nil
}

type Config struct {
	PSQL      PSQLConfig
	Logger    LoggerConfig
	Queue     QueueConfig
	Scheduler SchedulerConfig
}

type PSQLConfig struct {
	DSN string
}

type LoggerConfig struct {
	Path string
}

type QueueConfig struct {
	ConsumerTag  string
	URI          string
	ExchangeName string
	ExchangeType string
	Queue        string
	BindingKey   string
}

type SchedulerConfig struct {
	FetchIntervalSeconds int
}
