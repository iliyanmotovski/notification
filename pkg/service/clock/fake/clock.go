package fake

import (
	"time"

	"github.com/stretchr/testify/mock"
)

type MockClock struct {
	mock.Mock
}

func (m *MockClock) Now() time.Time {
	return m.Called().Get(0).(time.Time)
}
