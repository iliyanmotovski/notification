package entitybeeorm

import (
	"encoding/json"
	"time"

	"github.com/latolukasz/beeorm"
)

type EmailNotificationEntity struct {
	beeorm.ORM        `orm:"table=email_notifications;redisCache"`
	ID                uint64
	From              string `orm:"required"`
	FromName          string `orm:"required"`
	ReplyTo           string
	To                string `orm:"required"`
	Subject           string `orm:"required"`
	TemplateName      string `orm:"required"`
	TemplateData      json.RawMessage
	Status            string `orm:"required;enum=entitybeeorm.NotificationStatusAll"`
	StatusFromGateway string
	CreatedAt         time.Time `orm:"time=true"`
}
