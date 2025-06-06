package usecase

import (
	"context"
	"log"
	"sync"

	"github.com/pilebones/go-udev/netlink"
	"github.com/rustysys-dev/sabita_yusha/pkg/entity"
)

type UDevApp struct {
	ctx    context.Context
	isTest bool
	cancel func()
	conn   *netlink.UEventConn
	rule   netlink.RuleDefinition
	cfg    *entity.Config
	wg     *sync.WaitGroup
}

func (a *UDevApp) Run() {
	evtQueue := a.startEventQueue()

	if a.isTest {
		log.Println("Starting Search Please Press Your Macro Key...")
		a.HandleSearchForMacroDevice()
	} else {
		log.Println("Starting Macro Daemon")
		log.Println("Attempting to connect to target device")
		a.HandleExistingTarget()
		log.Println("Attempting to scan udev events for target device")
		a.HandleUdevEvents(evtQueue)
		log.Println("Application execution canceled")
	}
}

func (a *UDevApp) startEventQueue() chan netlink.UEvent {
	queue := make(chan netlink.UEvent)

	go func() {
		defer a.cancel()
	loop:
		for {
			select {
			case <-a.ctx.Done():
				close(queue)
				log.Println("udev event queue closed")
				break loop // stop iteration in case of stop signal received
			default:
				buf, err := a.conn.ReadMsg()
				if err != nil {
					log.Println("Unable to read uevent, err:", err)
					break loop // stop iteration in case of error
				}

				uevent, err := netlink.ParseUEvent(buf)
				if err != nil {
					log.Println("Unable to parse uevent, err:", err)
					continue loop // Drop uevent if not known
				}

				if !a.rule.Evaluate(*uevent) {
					continue loop // Drop uevent if not match
				}

				queue <- *uevent
			}
		}
	}()
	return queue
}
