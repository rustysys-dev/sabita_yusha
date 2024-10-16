package entity

import (
	"context"
	"fmt"

	"github.com/grafov/evdev"
)

type EventHandler interface {
	DeviceName() string
	FileName() string
	Run() error
}

type FactoryProps struct {
	RunType  EVFileType
	Name     string
	Device   *evdev.InputDevice
	Cfg      *Config
	Ctx      context.Context
	Cancelfn func()
}

func EVFileProcessorFactory(f *FactoryProps) EventHandler {
	switch f.RunType {
	case EV_FILE_TYPE_SEARCH:
		return NewSearchHandler(f)
	case EV_FILE_TYPE_EXEC:
		return NewExecHandler(f)
	default:
		fmt.Println("ERROR: invalid runtype:", f.RunType)
		return nil
	}
}
