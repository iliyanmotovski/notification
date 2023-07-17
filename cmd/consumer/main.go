package main

import (
	"context"
	"log"

	"github.com/iliyanmotovski/notification/pkg/constant"
	"github.com/iliyanmotovski/notification/pkg/queue"
	consumerbeeorm "github.com/iliyanmotovski/notification/pkg/queue/consumer/consumer_beeorm"
	"github.com/iliyanmotovski/notification/pkg/service/config"
	"github.com/iliyanmotovski/notification/pkg/service/email"
	gracefulshutdown "github.com/iliyanmotovski/notification/pkg/service/graceful_shutdown"
	"github.com/iliyanmotovski/notification/pkg/service/orm"
	"github.com/iliyanmotovski/notification/pkg/service/slack"
	"github.com/iliyanmotovski/notification/pkg/service/sms"
)

func main() {
	appContext, cancel := context.WithCancel(context.Background())

	gracefulshutdown.GracefulShutdown(func() {
		cancel()
	})

	configService, err := config.NewConfigService(constant.AppNameConsumer, "../../config/config.yaml")
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
