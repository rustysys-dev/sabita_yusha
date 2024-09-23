package usecase

import (
	"fmt"
	"log"
	"strings"

	"github.com/grafov/evdev"
	"github.com/pilebones/go-udev/netlink"
	"github.com/rustysys-dev/sabita_yusha/pkg/config"
	"github.com/rustysys-dev/sabita_yusha/pkg/entity"
)

const (
	EVT_CODE_MIN = 656
	EVT_CODE_MAX = 685
)

func HandleUdevEvents(q chan netlink.UEvent, e chan error, cfg *config.Config) {
	// Handling message from queue
	// Setup Tracking Var for DEVNAME
	var evtDir string
	for {
		select {
		case uevent := <-q:
			if evtDir != "" && strings.Contains(uevent.Env["DEVPATH"], evtDir) {
				if uevent.Action.String() == "add" {
					if strings.Contains(uevent.Env["DEVNAME"], "event") {
						err := HandleDevInputEventFile(uevent.Env["DEVNAME"], cfg)
						if err != nil {
							fmt.Println(err)
						}
					}
				}
				continue
			}

			if uevent.Env["NAME"] == cfg.UdevDeviceName() {
				evtDir = uevent.Env["DEVPATH"]
			}

		case err := <-e:
			log.Println("ERROR:", err)
		}
	}
}

func HandleDevInputEventFile(name string, cfg *config.Config) error {
	fmt.Println("/// Begin Reading Events:", name, "///")
	f, err := evdev.Open(name)
	if err != nil {
		return err
	}
	if err = HandleDevInputEvents(f, cfg); err != nil {
		return err
	}
	return nil
}

func HandleDevInputEvents(d *evdev.InputDevice, cfg *config.Config) error {
	for {
		evts, err := d.Read()
		if err != nil {
			return err
		}

		for _, evt := range evts {
			if evt.Code >= EVT_CODE_MIN && evt.Code <= EVT_CODE_MAX {
				macro, ok := cfg.CodeMacroMap[int(evt.Code)]
				if !ok {
					continue
				}
				if command, ok := macro.HandlerMap[int(evt.Value)]; ok {
					fmt.Println("Handling:", macro.Name)
					if err := HandleCustomCommand(command); err != nil {
						fmt.Println(err)
					}
				}
			}
		}
	}
}

func HandleCustomCommand(cmd *entity.Command) error {
	return cmd.Runner.Execute()
}
