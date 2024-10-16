package usecase

import (
	"log"
	"path/filepath"
	"sync"

	"github.com/grafov/evdev"
	"github.com/rustysys-dev/sabita_yusha/pkg/entity"
)

const MatchAny = ""

func (a *App) HandleExistingTarget() {
	a.handleEVFiles(entity.EV_FILE_TYPE_EXEC, a.cfg.TargetDeviceName)
}

func (a *App) HandleSearchForMacroDevice() {
	a.handleEVFiles(entity.EV_FILE_TYPE_SEARCH, MatchAny)
}

func (a *App) handleEVFiles(runType entity.EVFileType, match string) error {
	files, err := filepath.Glob("/dev/input/event*")
	if err != nil {
		return err
	}

	wg := &sync.WaitGroup{}
	wg.Add(len(files))
	for _, f := range files {
		ev, err := evdev.Open(f)
		if err != nil {
			log.Println("failed to open file:", f, "error:", err)
		}
		h := entity.EVFileProcessorFactory(&entity.FactoryProps{
			RunType:  runType,
			Name:     ev.Name,
			Device:   ev,
			Cfg:      a.cfg,
			Ctx:      a.ctx,
			Cancelfn: a.cancel,
		})
		go func() {
			defer wg.Done()
			// run indiscriminately if there is no match string
			if match == MatchAny {
				h.Run()
			} else if match == ev.Name {
				log.Println("Handling target:", match)
				h.Run()
			}
		}()
	}

	wg.Wait()

	return nil
}
