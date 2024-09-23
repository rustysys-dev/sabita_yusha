package usecase

import (
	"fmt"
	"path/filepath"

	"github.com/grafov/evdev"
	"github.com/rustysys-dev/sabita_yusha/pkg/config"
)

func HandleExistingTarget(cfg *config.Config) {
	files, err := filepath.Glob("/dev/input/event*")
	if err != nil {
		fmt.Println(err)
		return
	}

	for _, f := range files {
		ev, err := evdev.Open(f)
		if err != nil {
			fmt.Println("failed to open file:", f, "error:", err)
		}
		go func() {
			if ev.Name == cfg.TargetDeviceName {
				println(ev.Name)
				if err := HandleDevInputEvents(ev, cfg); err != nil {
					fmt.Println(err)
				}
			}
		}()
	}
}
