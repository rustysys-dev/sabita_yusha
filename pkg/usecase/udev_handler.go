package usecase

import (
	"fmt"
	"log"
	"strings"

	"github.com/grafov/evdev"
	"github.com/pilebones/go-udev/netlink"
	"github.com/rustysys-dev/sabita_yusha/pkg/entity"
)

func (a *App) HandleUdevEvents(q chan netlink.UEvent) {
	// Handling message from queue Setup Tracking Var for DEVNAME
	var evtDir string
	for {
		select {
		case <-a.ctx.Done():
			return
		case uevent := <-q:
			if evtDir != "" && strings.Contains(uevent.Env["DEVPATH"], evtDir) {
				if uevent.Action.String() == "add" {
					if strings.Contains(uevent.Env["DEVNAME"], "event") {
						err := a.HandleDevInputEventFile(uevent.Env["DEVNAME"])
						if err != nil {
							fmt.Println(err)
						}
					}
				}
				continue
			}

			if uevent.Env["NAME"] == a.cfg.UdevDeviceName() {
				evtDir = uevent.Env["DEVPATH"]
			}
		}
	}
}

func (a *App) HandleDevInputEventFile(name string) error {
	log.Println("/// Begin Reading Events:", name, "///")
	f, err := evdev.Open(name)
	if err != nil {
		return err
	}
	h := entity.EVFileProcessorFactory(&entity.FactoryProps{
		RunType:  entity.EV_FILE_TYPE_EXEC,
		Name:     f.Name,
		Device:   f,
		Cfg:      a.cfg,
		Ctx:      a.ctx,
		Cancelfn: a.cancel,
	})
	return h.Run()
}
