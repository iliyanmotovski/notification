package queue

const (
	// orm dirty queues
	OrmDirtyNotificationEntity = "orm.dirty-notification-entity"
)

func GetConsumerGroupName(queue string) string {
	return queue + "_group"
}
