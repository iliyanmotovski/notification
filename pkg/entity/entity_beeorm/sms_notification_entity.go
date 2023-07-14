package entitybeeorm

import (
	"time"

	"github.com/latolukasz/beeorm"
)

type SMSNotificationEntity struct {
	beeorm.ORM        `orm:"table=sms_notifications;redisCache"`
	ID                uint64
	Text              string `orm:"required"`
	MobileNumber      string `orm:"required"`
	Status            string `orm:"required;enum=entitybeeorm.NotificationStatusAll"`
	StatusFromGateway string
	CreatedAt         time.Time `orm:"time=true"`
}
