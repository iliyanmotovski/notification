package consumerbeeorm

import (
	"encoding/json"
	"fmt"

	entitybeeorm "github.com/iliyanmotovski/notification/pkg/entity/entity_beeorm"
	"github.com/iliyanmotovski/notification/pkg/service/email"
	"github.com/iliyanmotovski/notification/pkg/service/orm"
	"github.com/iliyanmotovski/notification/pkg/service/slack"
	"github.com/iliyanmotovski/notification/pkg/service/sms"
)

func NotificationsDirtyConsumer(smsService sms.Gateway, emailService email.Gateway, slackService slack.Service) orm.ConsumerHandler {
	return func(ormService orm.Engine, entityIDs []uint64) error {
		for _, entityID := range entityIDs {
			notificationEntity := &entitybeeorm.NotificationEntity{}

			if !ormService.LoadByID(entityID, notificationEntity, "SMSNotificationID", "EmailNotificationID", "SlackNotificationID") {
				return fmt.Errorf("notification entity with ID: %d not found", entityID)
			}

			flusher := ormService.NewFlusher()

			sendSMS := notificationEntity.SMSNotificationID != nil
			sendEmail := notificationEntity.EmailNotificationID != nil
			sendToSlack := notificationEntity.SlackNotificationID != nil

			if sendSMS {
				if checkStatus(notificationEntity.SMSNotificationID.Status) {
					sendSMS = false
				} else {
					notificationEntity.SMSNotificationID.Status = entitybeeorm.NotificationStatusInGateway.String()
					flusher.Track(notificationEntity.SMSNotificationID)
				}
			}
			if sendEmail {
				if checkStatus(notificationEntity.EmailNotificationID.Status) {
					sendEmail = false
				} else {
					notificationEntity.EmailNotificationID.Status = entitybeeorm.NotificationStatusInGateway.String()
					flusher.Track(notificationEntity.EmailNotificationID)
				}
			}
			if sendToSlack {
				if checkStatus(notificationEntity.SlackNotificationID.Status) {
					sendToSlack = false
				} else {
					notificationEntity.SlackNotificationID.Status = entitybeeorm.NotificationStatusInGateway.String()
					flusher.Track(notificationEntity.SlackNotificationID)
				}
			}

			flush := sendSMS || sendEmail || sendToSlack

			if flush {
				flusher.Flush()
			}

			if sendSMS {
				doSendSMS(smsService, notificationEntity.SMSNotificationID)

				flusher.Track(notificationEntity.SMSNotificationID)
			}
			if sendEmail {
				if err := doSendEmail(emailService, notificationEntity.EmailNotificationID); err != nil {
					return err
				}

				flusher.Track(notificationEntity.EmailNotificationID)
			}
			if sendToSlack {
				doSendToSlack(slackService, notificationEntity.SlackNotificationID)

				flusher.Track(notificationEntity.SlackNotificationID)
			}

			if flush {
				flusher.Flush()
			}
		}

		return nil
	}
}

func doSendSMS(smsService sms.Gateway, smsNotificationEntity *entitybeeorm.SMSNotificationEntity) {
	statusFromGateway, err := smsService.SendSMS(&sms.Message{
		Text:         smsNotificationEntity.Text,
		MobileNumber: smsNotificationEntity.MobileNumber,
	})
	if err != nil {
		smsNotificationEntity.Status = entitybeeorm.NotificationStatusFailed.String()
		smsNotificationEntity.StatusFromGateway = err.Error()
	} else {
		smsNotificationEntity.Status = entitybeeorm.NotificationStatusSuccess.String()
		smsNotificationEntity.StatusFromGateway = statusFromGateway
	}
}

func doSendEmail(emailService email.Gateway, emailNotificationEntity *entitybeeorm.EmailNotificationEntity) error {
	emailMessage := &email.Message{
		From:         emailNotificationEntity.From,
		FromName:     emailNotificationEntity.FromName,
		ReplyTo:      emailNotificationEntity.ReplyTo,
		To:           emailNotificationEntity.To,
		Subject:      emailNotificationEntity.Subject,
		TemplateName: emailNotificationEntity.TemplateName,
	}

	var templateData map[string]interface{}

	if len(emailNotificationEntity.TemplateData) > 0 {
		if err := json.Unmarshal(emailNotificationEntity.TemplateData, &templateData); err != nil {
			return err
		}

		emailMessage.TemplateData = templateData
	}

	err := emailService.SendTemplate(emailMessage)
	if err != nil {
		emailNotificationEntity.Status = entitybeeorm.NotificationStatusFailed.String()
		emailNotificationEntity.StatusFromGateway = err.Error()
	} else {
		emailNotificationEntity.Status = entitybeeorm.NotificationStatusSuccess.String()
	}

	return nil
}

func doSendToSlack(slackService slack.Service, slackNotificationEntity *entitybeeorm.SlackNotificationEntity) {
	err := slackService.SendToChannel(&slack.Message{
		BotName:     slackNotificationEntity.BotName,
		ChannelName: slackNotificationEntity.ChannelName,
		Message:     slackNotificationEntity.Message,
	})
	if err != nil {
		slackNotificationEntity.Status = entitybeeorm.NotificationStatusFailed.String()
		slackNotificationEntity.StatusFromGateway = err.Error()
	} else {
		slackNotificationEntity.Status = entitybeeorm.NotificationStatusSuccess.String()
	}
}

func checkStatus(status string) bool {
	return status == entitybeeorm.NotificationStatusSuccess.String()
}
