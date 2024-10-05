//go:generate go run cmd/package_app.go

package main

import (
	"log"
	"os"
	"os/signal"
	"syscall"

	"github.com/davecgh/go-spew/spew"
	"github.com/pilebones/go-udev/netlink"
	"github.com/rustysys-dev/sabita_yusha/pkg/config"
	"github.com/rustysys-dev/sabita_yusha/pkg/usecase"
)

const (
	CONFIG_PATH = ".config/sabita_yusha/config.yaml"
)

func main() {
	log.Println("Monitoring UEvent kernel message to user-space...")
	spew.Config.DisableMethods = true
	cfg, err := config.GetConfig(os.Getenv("HOME") + "/" + CONFIG_PATH)
	if err != nil {
		panic(err)
	}

	conn := new(netlink.UEventConn)
	if err = conn.Connect(netlink.UdevEvent); err != nil {
		log.Fatalln("Unable to connect to Netlink Kobject UEvent socket")
	}
	defer conn.Close()

	queue := make(chan netlink.UEvent)
	errors := make(chan error)
	quit := conn.Monitor(queue, errors, &netlink.RuleDefinition{
		Env: map[string]string{
			"SUBSYSTEM": "input",
		},
	})

	// Signal handler to quit properly monitor mode
	signals := make(chan os.Signal, 1)
	signal.Notify(signals, syscall.SIGINT, syscall.SIGTERM, syscall.SIGQUIT)
	go func() {
		<-signals
		log.Println("Exiting monitor mode...")
		close(quit)
		os.Exit(0)
	}()

	usecase.HandleExistingTarget(cfg)

	usecase.HandleUdevEvents(queue, errors, cfg)
}
