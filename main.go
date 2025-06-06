//go:generate go run cmd/package_app.go

package main

import (
	"flag"
	"log"
	"os"
	"path/filepath"

	"github.com/rustysys-dev/sabita_yusha/pkg/entity"
	"github.com/rustysys-dev/sabita_yusha/pkg/usecase"
)

const (
	CONFIG_PATH = ".config/sabita_yusha/config.yaml"
)

func main() {
	isTest := flag.Bool("test", false, "defines whether to start the daemon, or print out udev input events")
	command := flag.String("command", "", "defines the command to execute for CLI mode execution")
	flag.Parse()

	cfg, err := entity.GetConfig(filepath.Join(os.Getenv("HOME"), CONFIG_PATH))
	if err != nil {
		log.Fatalln(err)
	}

	app, err := usecase.New(cfg, *command, *isTest)
	if err != nil {
		log.Fatalln(err)
	}
	app.Run()
}
