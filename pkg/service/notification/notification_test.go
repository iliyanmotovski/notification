package notification

import (
	"testing"
	"time"

	entitybeeorm "github.com/iliyanmotovski/notification/pkg/entity/entity_beeorm"
	fakeClock "github.com/iliyanmotovski/notification/pkg/service/clock/fake"
	fakeORM "github.com/iliyanmotovski/notification/pkg/service/orm/fake"
	"github.com/stretchr/testify/assert"
)

func TestNotification(t *testing.T) {
	now := time.Unix(1, 1).UTC()

	mockClock := &fakeClock.MockClock{}
	mockClock.On("Now").Return(now)

	mockFlusher := &fakeORM.MockFlusher{}
	mockFlusher.On("Track", now).Once()
	mockFlusher.On("Track",
		"text",
		"mobile",
		entitybeeorm.NotificationStatusEnqueued.String(),
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
		map[string]interface{}{"k": "v"},
		entitybeeorm.NotificationStatusEnqueued.String(),
		"",
		now,
	).Once()
	mockFlusher.On("Track",
		"bot name",
		"channel name",
		"message",
		entitybeeorm.NotificationStatusEnqueued.String(),
		"",
		now,
	).Once()
	mockFlusher.On("Flush").Once()

	mockORMEngine := &fakeORM.MockORMEngine{}
	mockORMEngine.On("NewFlusher").Return(mockFlusher)

	notificationService := NewNotificationService(mockORMEngine, mockClock)
	err := notificationService.SendAsync(&Notification{
		SMS: &SMS{
			Text:         "text",
			MobileNumber: "mobile",
		},
		Email: &Email{
			From:         "from",
			FromName:     "from name",
			ReplyTo:      "reply to",
			To:           "to",
			Subject:      "subject",
			TemplateName: "template name",
			TemplateData: map[string]interface{}{"k": "v"},
		},
		SlackMessage: &SlackMessage{
			BotName:     "bot name",
			ChannelName: "channel name",
			Message:     "message",
		},
	})

	assert.Nil(t, err)

	mockClock.AssertExpectations(t)
	mockFlusher.AssertExpectations(t)
	mockORMEngine.AssertExpectations(t)
}

func TestNotificationWithValidationError(t *testing.T) {
	now := time.Unix(1, 1).UTC()

	tables := []*tableTest{
		{
			input: &Notification{
				SMS: &SMS{
					Text:         "",
					MobileNumber: "",
				},
			},
			want: "empty required field Text",
		},
		{
			input: &Notification{
				SMS: &SMS{
					Text:         "text",
					MobileNumber: "",
				},
			},
			want: "empty required field MobileNumber",
		},
		{
			input: &Notification{
				Email: &Email{
					From:         "",
					FromName:     "",
					To:           "",
					Subject:      "",
					TemplateName: "",
				},
			},
			want: "empty required field From",
		},
		{
			input: &Notification{
				Email: &Email{
					From:         "from",
					FromName:     "",
					To:           "",
					Subject:      "",
					TemplateName: "",
				},
			},
			want: "empty required field FromName",
		},
		{
			input: &Notification{
				Email: &Email{
					From:         "from",
					FromName:     "from name",
					To:           "",
					Subject:      "",
					TemplateName: "",
				},
			},
			want: "empty required field To",
		},
		{
			input: &Notification{
				Email: &Email{
					From:         "from",
					FromName:     "from name",
					To:           "to",
					Subject:      "",
					TemplateName: "",
				},
			},
			want: "empty required field Subject",
		},
		{
			input: &Notification{
				Email: &Email{
					From:         "from",
					FromName:     "from name",
					To:           "to",
					Subject:      "subject",
					TemplateName: "",
				},
			},
			want: "empty required field TemplateName",
		},
		{
			input: &Notification{
				SlackMessage: &SlackMessage{
					BotName:     "",
					ChannelName: "",
					Message:     "",
				},
			},
			want: "empty required field BotName",
		},
		{
			input: &Notification{
				SlackMessage: &SlackMessage{
					BotName:     "bot name",
					ChannelName: "",
					Message:     "",
				},
			},
			want: "empty required field ChannelName",
		},
		{
			input: &Notification{
				SlackMessage: &SlackMessage{
					BotName:     "bot name",
					ChannelName: "channel name",
					Message:     "",
				},
			},
			want: "empty required field Message",
		},
	}

	for _, table := range tables {
		mockClock := &fakeClock.MockClock{}
		mockClock.On("Now").Return(now)

		mockFlusher := &fakeORM.MockFlusher{}
		mockFlusher.On("Track", now).Once()

		mockORMEngine := &fakeORM.MockORMEngine{}
		mockORMEngine.On("NewFlusher").Return(mockFlusher)

		notificationService := NewNotificationService(mockORMEngine, mockClock)
		err := notificationService.SendAsync(table.input)

		assert.Equal(t, table.want, err.Error())

		mockClock.AssertExpectations(t)
		mockFlusher.AssertExpectations(t)
		mockORMEngine.AssertExpectations(t)
	}
}

func TestNotificationWithNoMethodSpecified(t *testing.T) {
	notificationService := NewNotificationService(nil, nil)
	err := notificationService.SendAsync(&Notification{})

	assert.Equal(t, "at least one notification type needs to be sent", err.Error())
}

type tableTest struct {
	input *Notification
	want  string
}
