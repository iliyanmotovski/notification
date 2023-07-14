package email

import (
	"fmt"
	"strconv"

	"github.com/iliyanm/notification/pkg/service/config"
	"github.com/mailjet/mailjet-apiv3-go/v4"
)

type mailjetGateway struct {
	*mailjet.Client
}

func NewEmailGatewayMailjet(configService config.Config) (Gateway, error) {
	apiKeyPublic, ok := configService.Get("email.mailjet.api_key_public")
	if !ok {
		return nil, fmt.Errorf("email.mailjet.api_key_public missing from config")
	}

	apiKeyPrivate, ok := configService.Get("email.mailjet.api_key_private")
	if !ok {
		return nil, fmt.Errorf("email.mailjet.api_key_private missing from config")
	}

	return &mailjetGateway{mailjet.NewMailjetClient(apiKeyPublic.(string), apiKeyPrivate.(string))}, nil
}

func (m *mailjetGateway) SendTemplate(message *Message) error {
	return m.sendTemplate(
		message.From,
		message.FromName,
		message.To,
		message.ReplyTo,
		message.Subject,
		message.TemplateName,
		message.TemplateData,
	)
}

func (m *mailjetGateway) sendTemplate(
	fromEmail string,
	fromName string,
	to string,
	replyTo string,
	subject string,
	templateName string,
	templateData map[string]interface{},
) error {
	templateID, err := strconv.Atoi(templateName)
	if err != nil {
		return err
	}

	messageInfo := mailjet.InfoMessagesV31{
		From: &mailjet.RecipientV31{
			Email: fromEmail,
			Name:  fromName,
		},
		To: &mailjet.RecipientsV31{
			mailjet.RecipientV31{
				Email: to,
			},
		},
		Subject:          subject,
		Variables:        templateData,
		TemplateID:       templateID,
		TemplateLanguage: true,
	}

	if len(replyTo) > 0 {
		messageInfo.ReplyTo = &mailjet.RecipientV31{
			Email: replyTo,
		}
	}

	message := &mailjet.MessagesV31{
		Info: []mailjet.InfoMessagesV31{messageInfo},
	}

	results, err := m.SendMailV31(message)
	if err != nil {
		return err
	}

	if results != nil {
		for _, response := range results.ResultsV31 {
			if response.Status != "success" {
				return fmt.Errorf("mailjet returned status: %s", response.Status)
			}
		}
	}

	return nil
}
