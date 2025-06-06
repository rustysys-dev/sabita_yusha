package usecase

import (
	"context"
	"fmt"
	"os/signal"
	"sync"
	"syscall"

	"github.com/pilebones/go-udev/netlink"
	"github.com/pkg/errors"
	"github.com/rustysys-dev/sabita_yusha/pkg/entity"
)

type App interface {
	Run()
}

func New(cfg *entity.Config, command string, isTest bool) (App, error) {
	if command != "" {
		return &CommandApp{
			command, cfg,
		}, nil
	}

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

	return &UDevApp{
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
