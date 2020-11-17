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
	// TODO: encapsulate env var management probably
	env := os.Getenv("ENV")

	path := "./configs/local.toml"
	switch env {
	case "TEST":
		path = "./configs/test.toml"
	case "TEST_RUNNER":
		path = "../configs/test.toml"
	}

	c, err := read(path)
	if err != nil {
		return nil, err
	}

	return &c, nil
}

type Config struct {
	PSQL       PSQLConfig
	Logger     LoggerConfig
	Queue      QueueConfig
	Scheduler  SchedulerConfig
	GRPCServer GRPCConfig
	HTTPServer HTTPConfig
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

type GRPCConfig struct {
	Address string
}

type HTTPConfig struct {
	Address string
}
