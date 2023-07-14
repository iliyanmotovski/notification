package entitybeeorm

type NotificationStatus string

func (u NotificationStatus) String() string {
	return string(u)
}

const (
	NotificationStatusEnqueued  NotificationStatus = "enqueued"
	NotificationStatusInGateway NotificationStatus = "in_gateway"
	NotificationStatusFailed    NotificationStatus = "failed"
	NotificationStatusSuccess   NotificationStatus = "success"
)

type notificationStatus struct {
	Enqueued  string
	InGateway string
	Failed    string
	Success   string
}

var NotificationStatusAll = notificationStatus{
	Enqueued:  NotificationStatusEnqueued.String(),
	InGateway: NotificationStatusInGateway.String(),
	Failed:    NotificationStatusFailed.String(),
	Success:   NotificationStatusSuccess.String(),
}
