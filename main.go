//go:generate go run cmd/package_app.go

package main

import (
	"flag"
	"log"
	"os"
	"path/filepath"

	"github.com/rustysys-dev/sabita_yusha/pkg/core"
	"github.com/rustysys-dev/sabita_yusha/pkg/entity"
	"github.com/rustysys-dev/sabita_yusha/pkg/usecase"
)

const (
	CONFIG_PATH = ".config/sabita_yusha/config.yaml"
)

func main() {
	isTest := flag.Bool("test", false, "defines whether to start the daemon, or print out udev input events")
	isSystemd := flag.Bool("systemd", false, "defines if this service is run from systemd")
	isUpdate := flag.Bool("update", false, "update the installed systemd service binary and restart")
	flag.Parse()

	if !*isSystemd {
		switch result, err := core.EnsureSystemdService(*isTest, *isUpdate); result {
		case core.SystemdStatusRunning:
			log.Println("running in systemd closing this instance")
			return
		case core.SystemdStatusUpdated:
			log.Println("service updated and restarted closing this instance")
			return
		case core.SystemdStatusFailed:
			log.Println("error failed to install to systemd:", err)
			return
		case core.SystemdStatusNotFound, core.SystemdStatusSkipped:
		}
	}

	cfg, err := entity.GetConfig(filepath.Join(os.Getenv("HOME"), CONFIG_PATH))
	if err != nil {
		log.Fatalln(err)
	}

	app, err := usecase.New(cfg, *isTest)
	if err != nil {
		log.Fatalln(err)
	}
	app.Run()
}
