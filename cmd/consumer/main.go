package main

import (
	"context"
	"log"

	"github.com/iliyanm/notification/pkg/queue"
	consumerbeeorm "github.com/iliyanm/notification/pkg/queue/consumer/consumer_beeorm"
	"github.com/iliyanm/notification/pkg/service/config"
	"github.com/iliyanm/notification/pkg/service/email"
	gracefulshutdown "github.com/iliyanm/notification/pkg/service/graceful_shutdown"
	"github.com/iliyanm/notification/pkg/service/orm"
	"github.com/iliyanm/notification/pkg/service/slack"
	"github.com/iliyanm/notification/pkg/service/sms"
)

func main() {
	appContext, cancel := context.WithCancel(context.Background())

	gracefulshutdown.GracefulShutdown(func() {
		cancel()
	})

	configService, err := config.NewConfigService("consumer", "../../config/config.yaml")
	if err != nil {
		log.Fatal(err)
	}

	smsService, err := sms.NewSMSGatewayTwilio(configService)
	if err != nil {
		log.Fatal(err)
	}

	emailService, err := email.NewEmailGatewayMailjet(configService)
	if err != nil {
		log.Fatal(err)
	}

	slackService, err := slack.NewSlackService(configService)
	if err != nil {
		log.Fatal(err)
	}

	ormRegistry, deferFunc, err := orm.NewORMRegistryService(configService)
	if err != nil {
		log.Fatal(err)
	}

	defer deferFunc()

	ormService := ormRegistry.GetORMService()
	ormService.ExecuteAlters()

	consumerRunner := orm.NewConsumerRunner(appContext, ormRegistry)

	consumerRunner.RunConsumerMany(
		consumerbeeorm.NotificationsDirtyConsumer(smsService, emailService, slackService),
		queue.OrmDirtyNotificationEntity,
		100,
	)

	consumerRunner.Wait()
}
