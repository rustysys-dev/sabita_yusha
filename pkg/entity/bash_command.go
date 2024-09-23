package entity

import (
	"errors"
	"fmt"
	"os/exec"

	"gopkg.in/yaml.v3"
)

type BashCommand struct {
	Entrypoint string   `yaml:"entrypoint"`
	Args       []string `yaml:"args,omitempty"`
}

func (c *BashCommand) Execute() error {
	cmd := exec.Command(c.Entrypoint, c.Args...)
	out, err := cmd.Output()
	if err != nil {
		return err
	}
	fmt.Println(string(out))
	return nil
}

func (c *BashCommand) FromYAML(y *yaml.Node) error {
	if y == nil {
		return errors.New("y cannot be nil")
	}
	for di := 1; di <= len(y.Content); di += 2 {
		dkNode := y.Content[di-1]
		dvNode := y.Content[di]

		switch dkNode.Value {
		case "entrypoint":
			c.Entrypoint = dvNode.Value
		case "args":
			for _, v := range dvNode.Content {
				c.Args = append(c.Args, v.Value)
			}
		default:
			return errors.New("unexpected node value, please fix bash command config")
		}
	}

	return nil
}
