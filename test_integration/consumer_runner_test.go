package test_integration

import (
	"context"
	"testing"
	"time"

	entitybeeorm "github.com/iliyanm/notification/pkg/entity/entity_beeorm"
	"github.com/iliyanm/notification/pkg/queue"
	"github.com/iliyanm/notification/pkg/service/orm"
	"github.com/stretchr/testify/assert"
)

func TestConsumerRunner(t *testing.T) {
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

	flusher := ormService.NewFlusher()

	notificationEntity := &entitybeeorm.NotificationEntity{
		CreatedAt: now,
	}

	flusher.Track(notificationEntity)

	smsNotificationEntity := &entitybeeorm.SMSNotificationEntity{
		Text: "text",
	}
	notificationEntity.SMSNotificationID = smsNotificationEntity

	flusher.Track(smsNotificationEntity)

	emailNotificationEntity := &entitybeeorm.EmailNotificationEntity{
		From: "from",
	}
	notificationEntity.EmailNotificationID = emailNotificationEntity

	flusher.Track(emailNotificationEntity)

	slackNotificationEntity := &entitybeeorm.SlackNotificationEntity{
		Message: "message",
	}
	notificationEntity.SlackNotificationID = slackNotificationEntity

	flusher.Track(slackNotificationEntity)

	flusher.Flush()

	runner.Wait()
}
