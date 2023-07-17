package fake

import (
	entitybeeorm "github.com/iliyanmotovski/notification/pkg/entity/entity_beeorm"
	"github.com/iliyanmotovski/notification/pkg/service/orm"
	"github.com/stretchr/testify/mock"
)

type MockORMEngine struct {
	mock.Mock
}

func (m *MockORMEngine) Flush(entity interface{}) {
	m.Called(entity)
}

func (m *MockORMEngine) NewFlusher() orm.Flusher {
	return m.Called().Get(0).(orm.Flusher)
}

func (m *MockORMEngine) LoadByID(id uint64, entity interface{}, references ...string) bool {
	switch v := entity.(type) {
	case *entitybeeorm.NotificationEntity:
		testData := m.Mock.TestData().Get("entitybeeorm.NotificationEntity")
		if data, ok := testData.Data().(*entitybeeorm.NotificationEntity); ok {
			v.ID = data.ID
			v.SMSNotificationID = data.SMSNotificationID
			v.EmailNotificationID = data.EmailNotificationID
			v.SlackNotificationID = data.SlackNotificationID
			v.CreatedAt = data.CreatedAt
		}
	}

	return m.Called(id, references).Bool(0)
}

func (m *MockORMEngine) ExecuteAlters() {
	m.Called()
}

func (m *MockORMEngine) TruncateTables() {
	m.Called()
}

func (m *MockORMEngine) GetEventBroker() orm.EventBroker {
	return m.Called().Get(0).(orm.EventBroker)
}

func (m *MockORMEngine) GetCacheService(namespace ...string) orm.CacheService {
	return m.Called(namespace).Get(0).(orm.CacheService)
}

func (m *MockORMEngine) GetRegistry() orm.RegistryService {
	return m.Called().Get(0).(orm.RegistryService)
}
