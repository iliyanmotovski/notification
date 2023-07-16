package fake

import (
	"github.com/iliyanm/notification/pkg/service/sms"
	"github.com/stretchr/testify/mock"
)

type MockSMSGateway struct {
	mock.Mock
}

func (m *MockSMSGateway) SendSMS(message *sms.Message) (string, error) {
	args := m.Called(message)
	return args.String(0), args.Error(1)
}
