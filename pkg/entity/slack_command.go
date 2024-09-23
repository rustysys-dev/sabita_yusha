package entity

import (
	"errors"

	"github.com/slack-go/slack"
	"gopkg.in/yaml.v3"
)

type SlackCommand struct {
	SlackToken string `yaml:"slack_token"`
	ChannelID  string `yaml:"channel_id"`
	AsUser     bool   `yaml:"as_user,omitempty"`
	Message    string `yaml:"message"`
}

func (c *SlackCommand) Execute() error {
	api := slack.New(c.SlackToken)
	if _, _, err := api.PostMessage(
		c.ChannelID,
		slack.MsgOptionText(c.Message, false),
		slack.MsgOptionAsUser(c.AsUser),
	); err != nil {
		return err
	}
	return nil
}

func (c *SlackCommand) FromYAML(y *yaml.Node) error {
	if y == nil {
		return errors.New("y cannot be nil")
	}
	for di := 1; di <= len(y.Content); di += 2 {
		dkNode := y.Content[di-1]
		dvNode := y.Content[di]

		switch dkNode.Value {
		case "slack_token":
			c.SlackToken = dvNode.Value
		case "channel_id":
			c.ChannelID = dvNode.Value
		case "as_user":
			if dvNode.Value == "true" {
				c.AsUser = true
			} else {
				c.AsUser = false
			}
		case "message":
			c.Message = dvNode.Value
		default:
			return errors.New("unexpected node value, please fix bash command config")
		}
	}

	return nil
}
