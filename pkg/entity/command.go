package entity

import (
	"errors"
	"fmt"

	"gopkg.in/yaml.v3"
)

type CommandRunner interface {
	Execute() error
}

type Command struct {
	Type   string        `yaml:"type"`
	Runner CommandRunner `yaml:"details"`
}

func (c *Command) UnmarshalYAML(node *yaml.Node) error {
	if node.Kind != yaml.MappingNode {
		fmt.Println(node.Kind, yaml.MappingNode)
		return errors.New("expected mapping node")
	}

	for i := 1; i <= len(node.Content); i += 2 {
		kNode := node.Content[i-1]
		vNode := node.Content[i]

		switch kNode.Value {
		case "type":
			c.Type = vNode.Value
		case "details":
			switch c.Type {
			case "bash":
				tmpCmd := &BashCommand{}
				if err := tmpCmd.FromYAML(vNode); err != nil {
					return err
				}
				c.Runner = tmpCmd
			case "slack":
				tmpCmd := &SlackCommand{}
				if err := tmpCmd.FromYAML(vNode); err != nil {
					return err
				}
				c.Runner = tmpCmd
			default:
				return fmt.Errorf("expected valid command type, got: %s", c.Type)
			}
		}
	}

	return nil
}
