package config

import (
	"os"

	"github.com/BurntSushi/toml"
)

const Test = "TEST"
const TestRunner = "TEST_RUNNER"

func read(fpath string) (c Config, err error) {
	_, err = toml.DecodeFile(fpath, &c)

	return
}

func Environment() string {
	return os.Getenv("ENV")
}

func GetConfig() (*Config, error) {
	env := Environment()

	path := "./configs/local.toml"
	switch env {
	case Test:
		path = "./configs/test.toml"
	// TODO: hack, try absolute path instead
	case TestRunner:
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
