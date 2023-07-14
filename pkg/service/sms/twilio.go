package sms

import (
	"fmt"

	"github.com/iliyanm/notification/pkg/service/config"
	"github.com/kevinburke/twilio-go"
)

type gatewayTwilio struct {
	sid        string
	token      string
	fromNumber string
}

func NewSMSGatewayTwilio(configService config.Config) (Gateway, error) {
	sid, ok := configService.Get("sms.twilio.sid")
	if !ok {
		return nil, fmt.Errorf("sms.twilio.sid missing from config")
	}

	token, ok := configService.Get("sms.twilio.token")
	if !ok {
		return nil, fmt.Errorf("sms.twilio.token missing from config")
	}

	fromNumber, ok := configService.Get("sms.twilio.from_number")
	if !ok {
		return nil, fmt.Errorf("sms.twilio.from_number missing from config")
	}

	return &gatewayTwilio{
		sid:        sid.(string),
		token:      token.(string),
		fromNumber: fromNumber.(string),
	}, nil
}

func (g *gatewayTwilio) SendSMS(message *Message) (string, error) {
	api := twilio.NewClient(g.sid, g.token, nil)

	result, err := api.Messages.SendMessage(g.fromNumber, message.MobileNumber, message.Text, nil)
	if err != nil {
		return "", err
	}

	return result.Status.Friendly(), nil
}
