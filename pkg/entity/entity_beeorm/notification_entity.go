package entitybeeorm

import (
	"time"

	"github.com/latolukasz/beeorm"
)

type NotificationEntity struct {
	beeorm.ORM `orm:"table=notifications;redisCache"`
	ID         uint64
	Asd        string
	CreatedAt  time.Time `orm:"time=true;dirty=orm.dirty-notification-entity"` // orm will push inside this queue only when this field is changed
}
