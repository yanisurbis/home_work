package logger

import (
	"log"
	"os"
)

type Logger interface {
	Init(path string) error
}

type Instance struct {
}

func (i *Instance) Init(_ string) error {
	log.SetOutput(os.Stdout)

	return nil
}
