package fake

import (
	"github.com/iliyanmotovski/notification/pkg/service/email"
	"github.com/stretchr/testify/mock"
)

type MockEmailGateway struct {
	mock.Mock
}

func (m *MockEmailGateway) SendTemplate(message *email.Message) error {
	return m.Called(message).Error(0)
}
