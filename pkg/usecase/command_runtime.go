package usecase

import (
	"fmt"

	"github.com/rustysys-dev/sabita_yusha/pkg/entity"
)

type CommandApp struct {
	command string
	cfg     *entity.Config
}

func (c *CommandApp) Run() {
	cmd, ok := c.cfg.MacroMap[c.command]
	if !ok {
		fmt.Println("Please check your configuration: .config/sabita_yusha/config.yaml")
		fmt.Println("Invalid Command: ", c.command)
		return
	}

	err := cmd.OnPress.Runner.Execute()
	if err != nil {
		fmt.Println("Error Executing Command: ", err)
	}
}
