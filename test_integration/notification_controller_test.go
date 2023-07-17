package test_integration

import (
	"context"
	"net/http"
	"testing"
	"time"

	"github.com/iliyanm/notification/api/notification-api/rest/dto"
	entitybeeorm "github.com/iliyanm/notification/pkg/entity/entity_beeorm"
	"github.com/iliyanm/notification/pkg/queue"
	"github.com/iliyanm/notification/pkg/service/clock/fake"
	"github.com/iliyanm/notification/pkg/service/orm"
	"github.com/stretchr/testify/assert"
)

func TestPostSendNotificationAsyncAction(t *testing.T) {
	createTestEngine()

	now := time.Unix(11111111, 11111111).UTC()

	cancelChan := make(chan struct{})

	testContext, cancel := context.WithCancel(context.Background())

	go func() {
		<-cancelChan
		cancel()
	}()

	runner := orm.NewConsumerRunner(testContext, ormService.GetRegistry())
	runner.RunConsumerMany(func(ormService orm.Engine, entityIDs []uint64) error {
		for _, entityID := range entityIDs {
			notificationEntity := &entitybeeorm.NotificationEntity{}
			if !ormService.LoadByID(entityID, notificationEntity, "SMSNotificationID", "EmailNotificationID", "SlackNotificationID") {
				t.Fatalf("notification entity with ID: %d not found", entityID)
			}

			assert.Equal(t, now.Year(), notificationEntity.CreatedAt.Year())
			assert.Equal(t, now.Month(), notificationEntity.CreatedAt.Month())
			assert.Equal(t, now.Day(), notificationEntity.CreatedAt.Day())

			assert.Equal(t, "text", notificationEntity.SMSNotificationID.Text)
			assert.Equal(t, "from", notificationEntity.EmailNotificationID.From)
			assert.Equal(t, "message", notificationEntity.SlackNotificationID.Message)

			cancelChan <- struct{}{}
		}
		return nil
	},
		queue.OrmDirtyNotificationEntity,
		1)

	clockService := &fake.MockClock{}
	clockService.On("Now").Return(now)

	body := &dto.RequestDTONotification{
		SMS: &dto.RequestDTOSMS{
			Text:         "text",
			MobileNumber: "mobile",
		},
		Email: &dto.RequestDTOEmail{
			From:         "from",
			FromName:     "from name",
			ReplyTo:      "reply to",
			To:           "to",
			Subject:      "subject",
			TemplateName: "template name",
		},
		SlackMessage: &dto.RequestDTOSlackMessage{
			BotName:     "bot name",
			ChannelName: "channel name",
			Message:     "message",
		},
	}

	sendHTTPRequest(clockService, http.MethodPost, "/v1/notification/send-async", body, nil)

	runner.Wait()

	clockService.AssertExpectations(t)
}
