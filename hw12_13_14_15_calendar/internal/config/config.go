package config

import (
	"flag"
	"log"
	"os"

	"github.com/BurntSushi/toml"
)

func read(fpath string) (c Config, err error) {
	_, err = toml.DecodeFile(fpath, &c)

	return
}

type Args struct {
	configPath string
}

func getArgs() *Args {
	configPath := flag.String("config", "", "path to config file")
	flag.Parse()

	args := Args{
		configPath: *configPath,
	}

	return &args
}

func GetConfig() (*Config, error) {
	// TODO: encapsulate env var management probably
	env := os.Getenv("ENV")
	path := ""
	//args := getArgs()
	//
	//if args.configPath != "" {
	//	path = args.configPath
	//} else if env == "TEST" {
	//	path = "../configs/local.toml"
	//} else {
	//	path = "./configs/local.toml"
	//}
	if env == "TEST" {
		path = "./configs/test.toml"
	} else if env == "TEST_RUNNER" {
		path = "../configs/test.toml"
	} else {
		path = "./configs/local.toml"
	}

	log.Println(env)

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
