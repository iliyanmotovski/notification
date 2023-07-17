package notification

import (
	"encoding/json"
	"errors"

	entitybeeorm "github.com/iliyanmotovski/notification/pkg/entity/entity_beeorm"
	"github.com/iliyanmotovski/notification/pkg/service/clock"
	"github.com/iliyanmotovski/notification/pkg/service/orm"
)

type Engine interface {
	SendAsync(*Notification) error
}

type Notification struct {
	SMS          *SMS
	Email        *Email
	SlackMessage *SlackMessage
}

func NewNotificationService(ormService orm.Engine, clock clock.IClock) Engine {
	return &engine{ormService: ormService, clock: clock}
}

type engine struct {
	ormService orm.Engine
	clock      clock.IClock
}

func (e *engine) SendAsync(notification *Notification) error {
	if notification.SMS == nil && notification.Email == nil && notification.SlackMessage == nil {
		return errors.New("at least one notification type needs to be sent")
	}

	now := e.clock.Now()
	flusher := e.ormService.NewFlusher()

	notificationEntity := &entitybeeorm.NotificationEntity{
		CreatedAt: now,
	}

	flusher.Track(notificationEntity)

	if notification.SMS != nil {
		if err := notification.SMS.validate(); err != nil {
			return err
		}

		smsNotificationEntity := &entitybeeorm.SMSNotificationEntity{
			Text:         notification.SMS.Text,
			MobileNumber: notification.SMS.MobileNumber,
			Status:       entitybeeorm.NotificationStatusEnqueued.String(),
			CreatedAt:    now,
		}

		notificationEntity.SMSNotificationID = smsNotificationEntity

		flusher.Track(smsNotificationEntity)
	}

	if notification.Email != nil {
		if err := notification.Email.validate(); err != nil {
			return err
		}

		emailNotificationEntity := &entitybeeorm.EmailNotificationEntity{
			From:         notification.Email.From,
			FromName:     notification.Email.FromName,
			ReplyTo:      notification.Email.ReplyTo,
			To:           notification.Email.To,
			Subject:      notification.Email.Subject,
			TemplateName: notification.Email.TemplateName,
			Status:       entitybeeorm.NotificationStatusEnqueued.String(),
			CreatedAt:    now,
		}

		if notification.Email.TemplateData != nil {
			marshaled, err := json.Marshal(notification.Email.TemplateData)
			if err != nil {
				return err
			}

			emailNotificationEntity.TemplateData = marshaled
		}

		notificationEntity.EmailNotificationID = emailNotificationEntity

		flusher.Track(emailNotificationEntity)
	}

	if notification.SlackMessage != nil {
		if err := notification.SlackMessage.validate(); err != nil {
			return err
		}

		slackNotificationEntity := &entitybeeorm.SlackNotificationEntity{
			BotName:     notification.SlackMessage.BotName,
			ChannelName: notification.SlackMessage.ChannelName,
			Message:     notification.SlackMessage.Message,
			Status:      entitybeeorm.NotificationStatusEnqueued.String(),
			CreatedAt:   now,
		}

		notificationEntity.SlackNotificationID = slackNotificationEntity

		flusher.Track(slackNotificationEntity)
	}

	flusher.Flush()
	return nil
}

type SMS struct {
	Text         string
	MobileNumber string
}

func (s *SMS) validate() error {
	if s.Text == "" {
		return errors.New("empty required field Text")
	}
	if s.MobileNumber == "" {
		return errors.New("empty required field MobileNumber")
	}

	return nil
}

type Email struct {
	From         string
	FromName     string
	ReplyTo      string
	To           string
	Subject      string
	TemplateName string
	TemplateData map[string]interface{}
}

func (e *Email) validate() error {
	if e.From == "" {
		return errors.New("empty required field From")
	}
	if e.FromName == "" {
		return errors.New("empty required field FromName")
	}
	if e.To == "" {
		return errors.New("empty required field To")
	}
	if e.Subject == "" {
		return errors.New("empty required field Subject")
	}
	if e.TemplateName == "" {
		return errors.New("empty required field TemplateName")
	}

	return nil
}

type SlackMessage struct {
	BotName     string
	ChannelName string
	Message     string
}

func (e *SlackMessage) validate() error {
	if e.BotName == "" {
		return errors.New("empty required field BotName")
	}
	if e.ChannelName == "" {
		return errors.New("empty required field ChannelName")
	}
	if e.Message == "" {
		return errors.New("empty required field Message")
	}

	return nil
}
