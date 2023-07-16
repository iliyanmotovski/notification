package consumerbeeorm

import (
	"errors"
	"testing"
	"time"

	entitybeeorm "github.com/iliyanm/notification/pkg/entity/entity_beeorm"
	"github.com/iliyanm/notification/pkg/service/email"
	fakeEmail "github.com/iliyanm/notification/pkg/service/email/fake"
	fakeORM "github.com/iliyanm/notification/pkg/service/orm/fake"
	"github.com/iliyanm/notification/pkg/service/slack"
	fakeSlack "github.com/iliyanm/notification/pkg/service/slack/fake"
	"github.com/iliyanm/notification/pkg/service/sms"
	fakeSMS "github.com/iliyanm/notification/pkg/service/sms/fake"
	"github.com/stretchr/testify/assert"
)

func TestNotificationsDirtyConsumer(t *testing.T) {
	now := time.Unix(1, 1).UTC()

	var nilMap map[string]interface{}

	mockFlusher := &fakeORM.MockFlusher{}
	mockFlusher.On("Track",
		"text",
		"mobile",
		entitybeeorm.NotificationStatusInGateway.String(),
		"",
		now,
	).Once()
	mockFlusher.On("Track",
		"from",
		"from name",
		"reply to",
		"to",
		"subject",
		"template name",
		nilMap,
		entitybeeorm.NotificationStatusInGateway.String(),
		"",
		now,
	).Once()
	mockFlusher.On("Track",
		"bot name",
		"channel name",
		"message",
		entitybeeorm.NotificationStatusInGateway.String(),
		"",
		now,
	).Once()
	mockFlusher.On("Flush").Once()

	mockFlusher.On("Track",
		"text",
		"mobile",
		entitybeeorm.NotificationStatusSuccess.String(),
		"success",
		now,
	).Once()
	mockFlusher.On("Track",
		"from",
		"from name",
		"reply to",
		"to",
		"subject",
		"template name",
		nilMap,
		entitybeeorm.NotificationStatusSuccess.String(),
		"",
		now,
	).Once()
	mockFlusher.On("Track",
		"bot name",
		"channel name",
		"message",
		entitybeeorm.NotificationStatusSuccess.String(),
		"",
		now,
	).Once()
	mockFlusher.On("Flush").Once()

	mockORMEngine := &fakeORM.MockORMEngine{}
	mockORMEngine.On("NewFlusher").Return(mockFlusher)

	notificationEntity := &entitybeeorm.NotificationEntity{
		ID: 1,
		SMSNotificationID: &entitybeeorm.SMSNotificationEntity{
			ID:           1,
			Text:         "text",
			MobileNumber: "mobile",
			Status:       entitybeeorm.NotificationStatusEnqueued.String(),
			CreatedAt:    now,
		},
		EmailNotificationID: &entitybeeorm.EmailNotificationEntity{
			ID:           1,
			From:         "from",
			FromName:     "from name",
			ReplyTo:      "reply to",
			To:           "to",
			Subject:      "subject",
			TemplateName: "template name",
			Status:       entitybeeorm.NotificationStatusEnqueued.String(),
			CreatedAt:    now,
		},
		SlackNotificationID: &entitybeeorm.SlackNotificationEntity{
			ID:          1,
			BotName:     "bot name",
			ChannelName: "channel name",
			Message:     "message",
			Status:      entitybeeorm.NotificationStatusEnqueued.String(),
			CreatedAt:   now,
		},
		CreatedAt: now,
	}

	mockORMEngine.TestData().Set("entitybeeorm.NotificationEntity", notificationEntity)

	mockORMEngine.On("LoadByID", uint64(1), []string{"SMSNotificationID", "EmailNotificationID", "SlackNotificationID"}).Return(true)

	messageSMS := &sms.Message{
		Text:         "text",
		MobileNumber: "mobile",
	}

	messageEmail := &email.Message{
		From:         "from",
		FromName:     "from name",
		ReplyTo:      "reply to",
		To:           "to",
		Subject:      "subject",
		TemplateName: "template name",
	}

	messageSlack := &slack.Message{
		BotName:     "bot name",
		ChannelName: "channel name",
		Message:     "message",
	}

	mockSMS := &fakeSMS.MockSMSGateway{}
	mockSMS.On("SendSMS", messageSMS).Return("success", nil).Once()

	mockEmail := &fakeEmail.MockEmailGateway{}
	mockEmail.On("SendTemplate", messageEmail).Return(nil).Once()

	mockSlack := &fakeSlack.MockSlackService{}
	mockSlack.On("SendToChannel", messageSlack).Return(nil).Once()

	consumerFunc := NotificationsDirtyConsumer(mockSMS, mockEmail, mockSlack)
	err := consumerFunc(mockORMEngine, []uint64{1})

	assert.Nil(t, err)

	mockFlusher.AssertExpectations(t)
	mockORMEngine.AssertExpectations(t)
	mockSMS.AssertExpectations(t)
	mockEmail.AssertExpectations(t)
	mockSlack.AssertExpectations(t)
}

func TestNotificationsDirtyConsumerWithErrorFromGateways(t *testing.T) {
	now := time.Unix(1, 1).UTC()

	var nilMap map[string]interface{}

	mockFlusher := &fakeORM.MockFlusher{}
	mockFlusher.On("Track",
		"text",
		"mobile",
		entitybeeorm.NotificationStatusInGateway.String(),
		"",
		now,
	).Once()
	mockFlusher.On("Track",
		"from",
		"from name",
		"reply to",
		"to",
		"subject",
		"template name",
		nilMap,
		entitybeeorm.NotificationStatusInGateway.String(),
		"",
		now,
	).Once()
	mockFlusher.On("Track",
		"bot name",
		"channel name",
		"message",
		entitybeeorm.NotificationStatusInGateway.String(),
		"",
		now,
	).Once()
	mockFlusher.On("Flush").Once()

	mockFlusher.On("Track",
		"text",
		"mobile",
		entitybeeorm.NotificationStatusFailed.String(),
		"fail",
		now,
	).Once()
	mockFlusher.On("Track",
		"from",
		"from name",
		"reply to",
		"to",
		"subject",
		"template name",
		nilMap,
		entitybeeorm.NotificationStatusFailed.String(),
		"fail",
		now,
	).Once()
	mockFlusher.On("Track",
		"bot name",
		"channel name",
		"message",
		entitybeeorm.NotificationStatusFailed.String(),
		"fail",
		now,
	).Once()
	mockFlusher.On("Flush").Once()

	mockORMEngine := &fakeORM.MockORMEngine{}
	mockORMEngine.On("NewFlusher").Return(mockFlusher)

	notificationEntity := &entitybeeorm.NotificationEntity{
		ID: 1,
		SMSNotificationID: &entitybeeorm.SMSNotificationEntity{
			ID:           1,
			Text:         "text",
			MobileNumber: "mobile",
			Status:       entitybeeorm.NotificationStatusEnqueued.String(),
			CreatedAt:    now,
		},
		EmailNotificationID: &entitybeeorm.EmailNotificationEntity{
			ID:           1,
			From:         "from",
			FromName:     "from name",
			ReplyTo:      "reply to",
			To:           "to",
			Subject:      "subject",
			TemplateName: "template name",
			Status:       entitybeeorm.NotificationStatusEnqueued.String(),
			CreatedAt:    now,
		},
		SlackNotificationID: &entitybeeorm.SlackNotificationEntity{
			ID:          1,
			BotName:     "bot name",
			ChannelName: "channel name",
			Message:     "message",
			Status:      entitybeeorm.NotificationStatusEnqueued.String(),
			CreatedAt:   now,
		},
		CreatedAt: now,
	}

	mockORMEngine.TestData().Set("entitybeeorm.NotificationEntity", notificationEntity)

	mockORMEngine.On("LoadByID", uint64(1), []string{"SMSNotificationID", "EmailNotificationID", "SlackNotificationID"}).Return(true)

	messageSMS := &sms.Message{
		Text:         "text",
		MobileNumber: "mobile",
	}

	messageEmail := &email.Message{
		From:         "from",
		FromName:     "from name",
		ReplyTo:      "reply to",
		To:           "to",
		Subject:      "subject",
		TemplateName: "template name",
	}

	messageSlack := &slack.Message{
		BotName:     "bot name",
		ChannelName: "channel name",
		Message:     "message",
	}

	err := errors.New("fail")

	mockSMS := &fakeSMS.MockSMSGateway{}
	mockSMS.On("SendSMS", messageSMS).Return("fail", err).Once()

	mockEmail := &fakeEmail.MockEmailGateway{}
	mockEmail.On("SendTemplate", messageEmail).Return(err).Once()

	mockSlack := &fakeSlack.MockSlackService{}
	mockSlack.On("SendToChannel", messageSlack).Return(err).Once()

	consumerFunc := NotificationsDirtyConsumer(mockSMS, mockEmail, mockSlack)
	err = consumerFunc(mockORMEngine, []uint64{1})

	assert.Nil(t, err)

	mockFlusher.AssertExpectations(t)
	mockORMEngine.AssertExpectations(t)
	mockSMS.AssertExpectations(t)
	mockEmail.AssertExpectations(t)
	mockSlack.AssertExpectations(t)
}

func TestNotificationsDirtyConsumerWithAlreadySentNotifications(t *testing.T) {
	now := time.Unix(1, 1).UTC()

	mockFlusher := &fakeORM.MockFlusher{}
	mockORMEngine := &fakeORM.MockORMEngine{}
	mockORMEngine.On("NewFlusher").Return(mockFlusher)

	notificationEntity := &entitybeeorm.NotificationEntity{
		ID: 1,
		SMSNotificationID: &entitybeeorm.SMSNotificationEntity{
			ID:           1,
			Text:         "text",
			MobileNumber: "mobile",
			Status:       entitybeeorm.NotificationStatusSuccess.String(),
			CreatedAt:    now,
		},
		EmailNotificationID: &entitybeeorm.EmailNotificationEntity{
			ID:           1,
			From:         "from",
			FromName:     "from name",
			ReplyTo:      "reply to",
			To:           "to",
			Subject:      "subject",
			TemplateName: "template name",
			Status:       entitybeeorm.NotificationStatusSuccess.String(),
			CreatedAt:    now,
		},
		SlackNotificationID: &entitybeeorm.SlackNotificationEntity{
			ID:          1,
			BotName:     "bot name",
			ChannelName: "channel name",
			Message:     "message",
			Status:      entitybeeorm.NotificationStatusSuccess.String(),
			CreatedAt:   now,
		},
		CreatedAt: now,
	}

	mockORMEngine.TestData().Set("entitybeeorm.NotificationEntity", notificationEntity)

	mockORMEngine.On("LoadByID", uint64(1), []string{"SMSNotificationID", "EmailNotificationID", "SlackNotificationID"}).Return(true)

	mockSMS := &fakeSMS.MockSMSGateway{}
	mockEmail := &fakeEmail.MockEmailGateway{}
	mockSlack := &fakeSlack.MockSlackService{}

	consumerFunc := NotificationsDirtyConsumer(mockSMS, mockEmail, mockSlack)
	err := consumerFunc(mockORMEngine, []uint64{1})

	assert.Nil(t, err)

	mockFlusher.AssertExpectations(t)
	mockORMEngine.AssertExpectations(t)
	mockSMS.AssertExpectations(t)
	mockEmail.AssertExpectations(t)
	mockSlack.AssertExpectations(t)
}

func TestNotificationsDirtyConsumerWithEntityIDNotFound(t *testing.T) {

	mockORMEngine := &fakeORM.MockORMEngine{}
	mockORMEngine.On("LoadByID", uint64(1), []string{"SMSNotificationID", "EmailNotificationID", "SlackNotificationID"}).Return(false)

	mockSMS := &fakeSMS.MockSMSGateway{}
	mockEmail := &fakeEmail.MockEmailGateway{}
	mockSlack := &fakeSlack.MockSlackService{}

	consumerFunc := NotificationsDirtyConsumer(mockSMS, mockEmail, mockSlack)
	err := consumerFunc(mockORMEngine, []uint64{1})

	assert.Equal(t, "notification entity with ID: 1 not found", err.Error())

	mockORMEngine.AssertExpectations(t)
	mockSMS.AssertExpectations(t)
	mockEmail.AssertExpectations(t)
	mockSlack.AssertExpectations(t)
}
