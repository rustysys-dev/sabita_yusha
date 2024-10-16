package usecase

import (
	"context"
	"fmt"
	"log"
	"os/signal"
	"sync"
	"syscall"

	"github.com/pilebones/go-udev/netlink"
	"github.com/pkg/errors"
	"github.com/rustysys-dev/sabita_yusha/pkg/entity"
)

type App struct {
	ctx    context.Context
	isTest bool
	cancel func()
	conn   *netlink.UEventConn
	rule   netlink.RuleDefinition
	cfg    *entity.Config
	wg     *sync.WaitGroup
}

func New(cfg *entity.Config, isTest bool) (*App, error) {
	ctx, cancel := signal.NotifyContext(context.Background(), syscall.SIGINT, syscall.SIGTERM, syscall.SIGQUIT)
	conn := new(netlink.UEventConn)

	if err := conn.Connect(netlink.UdevEvent); err != nil {
		return nil, errors.New("Unable to connect to Netlink Kobject UEvent socket")
	}

	rule := netlink.RuleDefinition{
		Env: map[string]string{
			"SUBSYSTEM": "input",
		},
	}
	if err := rule.Compile(); err != nil {
		return nil, fmt.Errorf("Wrong matcher, err: %w", err)
	}

	wg := &sync.WaitGroup{}

	return &App{
		isTest: isTest,
		ctx:    ctx,
		conn:   conn,
		rule:   rule,
		cfg:    cfg,
		wg:     wg,
		cancel: func() {
			cancel()
			conn.Close()
		},
	}, nil
}

func (a *App) Run() {
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

func (a *App) startEventQueue() chan netlink.UEvent {
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
