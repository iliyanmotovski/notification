package slack

import (
	"fmt"

	"github.com/iliyanmotovski/notification/pkg/service/config"
	"github.com/slack-go/slack"
)

type Service interface {
	SendToChannel(*Message) error
}

type Message struct {
	BotName     string
	ChannelName string
	Message     string
}

func NewSlackService(configService config.Config) (Service, error) {
	botTokensConfig, ok := configService.Get("slack.bot_tokens")
	if !ok {
		return nil, fmt.Errorf("slack.bot_tokens missing from config")
	}

	clients := make(map[string]*slack.Client)

	for name, token := range botTokensConfig.(map[interface{}]interface{}) {
		clients[name.(string)] = slack.New(token.(string))
	}

	return &slackService{clients: clients}, nil
}

type slackService struct {
	clients map[string]*slack.Client
}

func (s *slackService) SendToChannel(message *Message) error {
	client, ok := s.clients[message.BotName]
	if !ok {
		return fmt.Errorf("slack bot '%s' not defined", message.BotName)
	}

	_, _, err := client.PostMessage(message.ChannelName, slack.MsgOptionText(message.Message, true))
	return err
}
