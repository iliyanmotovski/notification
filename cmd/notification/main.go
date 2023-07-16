package main

import (
	"log"

	"github.com/iliyanm/notification/pkg/service/clock"
	"github.com/iliyanm/notification/pkg/service/config"
	"github.com/iliyanm/notification/pkg/service/notification"
	"github.com/iliyanm/notification/pkg/service/orm"
)

func main() {
	configService, err := config.NewConfigService("notification", "../../config/config.yaml")
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
	ormService.ExecuteAlters()

	notificationService := notification.NewNotificationService(ormService, clockService)

	if err := notificationService.SendAsync(&notification.Notification{
		SMS: &notification.SMS{
			Text:         "test",
			MobileNumber: "+359878697929",
		},
		Email: &notification.Email{
			From:         "ilqnskiq@abv.bg",
			FromName:     "iliyan",
			To:           "iliyan.motovski@gmail.com",
			Subject:      "test",
			TemplateName: "4954974",
			TemplateData: nil,
		},
		SlackMessage: &notification.SlackMessage{
			BotName:     "bot",
			ChannelName: "bot",
			Message:     "test message",
		},
	}); err != nil {
		log.Fatal(err)
	}
}
