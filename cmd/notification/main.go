package main

import (
	"log"

	"github.com/iliyanm/notification/pkg/service/clock"
	"github.com/iliyanm/notification/pkg/service/config"
	"github.com/iliyanm/notification/pkg/service/notification"
	"github.com/iliyanm/notification/pkg/service/orm"
)

func main() {
	configService, err := config.NewConfigService("notification", "../../config")
	if err != nil {
		log.Fatal(err)
	}

	ormRegistry, deferFunc, err := orm.NewORMRegistryService(configService)
	if err != nil {
		log.Fatal(err)
	}

	defer deferFunc()

	clockService := clock.NewClockService()

	ormService := ormRegistry.GetORMService()

	notificationService := notification.NewNotificationService(ormService, clockService)

	if err := notificationService.SendAsync(&notification.Notification{
		SMS: &notification.SMS{
			Text:         "test",
			MobileNumber: "+359878697929",
		},
		Email: &notification.Email{
			From:         "test@gmail.com",
			FromName:     "test",
			ReplyTo:      "test@gmail.com",
			To:           "test_to@gmail.com",
			Subject:      "subject",
			TemplateName: "template",
			TemplateData: nil,
		},
		SlackMessage: &notification.SlackMessage{
			BotName:     "bot-name",
			ChannelName: "channel-name",
			Message:     "message",
		},
	}); err != nil {
		log.Fatal(err)
	}
}
