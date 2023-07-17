package fake

import (
	"github.com/iliyanmotovski/notification/pkg/service/slack"
	"github.com/stretchr/testify/mock"
)

type MockSlackService struct {
	mock.Mock
}

func (m *MockSlackService) SendToChannel(message *slack.Message) error {
	return m.Called(message).Error(0)
}
