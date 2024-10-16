package entity

import (
	"context"
	"log"

	"github.com/grafov/evdev"
)

type SearchHandler struct {
	name     string
	device   *evdev.InputDevice
	cfg      *Config
	ctx      context.Context
	cancelfn func()
	evtChan  chan evdev.InputEvent
}

func NewSearchHandler(f *FactoryProps) EventHandler {
	return &SearchHandler{
		name:     f.Name,
		device:   f.Device,
		cfg:      f.Cfg,
		ctx:      f.Ctx,
		cancelfn: f.Cancelfn,
		evtChan:  make(chan evdev.InputEvent),
	}
}

func (h *SearchHandler) DeviceName() string {
	return h.device.Name
}

func (h *SearchHandler) FileName() string {
	return h.device.Fn
}

func (h *SearchHandler) Run() error {
	go func() {
		evts, err := h.device.Read()
		if err != nil {
			log.Println("Error while reading device:", err)
			return
		}
		for _, evt := range evts {
			h.evtChan <- evt
		}
	}()

	for {
		select {
		case <-h.ctx.Done():
			return nil
		case evt := <-h.evtChan:
			if evt.Code >= EVT_CODE_MIN && evt.Code <= EVT_CODE_MAX {
				h.handleEvent(evt)
				// cancel the context
				h.cancelfn()
				return nil
			}
		}
	}
}

func (h *SearchHandler) handleEvent(ev evdev.InputEvent) {
	log.Println("DeviceName:", h.DeviceName())
	log.Println("FileName:", h.FileName())
}
