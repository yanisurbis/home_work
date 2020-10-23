package config

import (
	"fmt"

	"github.com/BurntSushi/toml"
)

func Read(fpath string) (c Config, err error) {
	_, err = toml.DecodeFile(fpath, &c)
	if err != nil {
		fmt.Println(err)
	}

	return
}

type Config struct {
	PSQL   PSQLConfig
	Logger LoggerConfig
	Queue  QueueConfig
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
