package fake

import (
	"encoding/json"

	entitybeeorm "github.com/iliyanmotovski/notification/pkg/entity/entity_beeorm"
	"github.com/stretchr/testify/mock"
)

type MockFlusher struct {
	mock.Mock
}

func (m *MockFlusher) Track(entity interface{}) {
	switch v := entity.(type) {
	case *entitybeeorm.NotificationEntity:
		m.Called(v.CreatedAt)
	case *entitybeeorm.SMSNotificationEntity:
		m.Called(
			v.Text,
			v.MobileNumber,
			v.Status,
			v.StatusFromGateway,
			v.CreatedAt,
		)
	case *entitybeeorm.EmailNotificationEntity:
		var templateData map[string]interface{}
		if v.TemplateData != nil {
			if err := json.Unmarshal(v.TemplateData, &templateData); err != nil {
				panic(err)
			}
		}

		m.Called(
			v.From,
			v.FromName,
			v.ReplyTo,
			v.To,
			v.Subject,
			v.TemplateName,
			templateData,
			v.Status,
			v.StatusFromGateway,
			v.CreatedAt,
		)
	case *entitybeeorm.SlackNotificationEntity:
		m.Called(
			v.BotName,
			v.ChannelName,
			v.Message,
			v.Status,
			v.StatusFromGateway,
			v.CreatedAt,
		)
	}
}

func (m *MockFlusher) Flush() {
	m.Called()
}
