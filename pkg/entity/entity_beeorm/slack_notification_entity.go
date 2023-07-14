package entitybeeorm

import (
	"time"

	"github.com/latolukasz/beeorm"
)

type SlackNotificationEntity struct {
	beeorm.ORM        `orm:"table=slack_notifications;redisCache"`
	ID                uint64
	BotName           string `orm:"required"`
	ChannelName       string `orm:"required"`
	Message           string `orm:"required"`
	Status            string `orm:"required;enum=entitybeeorm.NotificationStatusAll"`
	StatusFromGateway string
	CreatedAt         time.Time `orm:"time=true"`
}
