package entity

import (
	"context"
	"fmt"
	"log"
	"sync"

	"github.com/grafov/evdev"
)

type ExecHandler struct {
	name     string
	device   *evdev.InputDevice
	cfg      *Config
	ctx      context.Context
	cancelfn func()
	evtChan  chan evdev.InputEvent
	wg       *sync.WaitGroup
}

func NewExecHandler(f *FactoryProps) EventHandler {
	return &ExecHandler{
		name:     f.Name,
		device:   f.Device,
		cfg:      f.Cfg,
		ctx:      f.Ctx,
		cancelfn: f.Cancelfn,
		evtChan:  make(chan evdev.InputEvent),
	}
}

func (h *ExecHandler) DeviceName() string {
	return h.device.Name
}

func (h *ExecHandler) FileName() string {
	return h.device.Fn
}

func (h *ExecHandler) Run() error {
	errs := make(chan error, 1)
	go func() {
		for {
			evts, err := h.device.Read()
			if err != nil {
				log.Println("Error while reading device:", err)
				errs <- err
				return
			}
			for _, evt := range evts {
				h.evtChan <- evt
			}
		}
	}()

	for {
		select {
		// app ctx
		case <-h.ctx.Done():
			return nil
		// err chan
		case <-errs:
			return nil
		case evt := <-h.evtChan:
			if evt.Code >= EVT_CODE_MIN && evt.Code <= EVT_CODE_MAX {
				h.handleEvent(evt)
			}
		}
	}
}

func (h *ExecHandler) handleEvent(ev evdev.InputEvent) {
	macro, ok := h.cfg.CodeMacroMap[int(ev.Code)]
	if !ok {
		return
	}
	if command, ok := macro.HandlerMap[int(ev.Value)]; ok {
		fmt.Println("Handling:", macro.Name)
		if err := command.Runner.Execute(); err != nil {
			fmt.Println(err)
		}
	}
}
