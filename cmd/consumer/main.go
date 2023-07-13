package main

import (
	"context"
	"log"

	"github.com/iliyanm/notification/pkg/queue"
	"github.com/iliyanm/notification/pkg/service/config"
	gracefulshutdown "github.com/iliyanm/notification/pkg/service/graceful_shutdown"
	"github.com/iliyanm/notification/pkg/service/orm"
	"github.com/latolukasz/beeorm"
)

func main() {
	appContext, cancel := context.WithCancel(context.Background())

	gracefulshutdown.GracefulShutdown(func() {
		cancel()
	})

	configService, err := config.NewConfigService("consumer", "../../config")
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

	consumerRunner.RunConsumerMany(func(ormService orm.Engine, events []beeorm.Event) error {
		for _, event := range events {
			log.Println(beeorm.EventDirtyEntity(event).ID())
		}

		return nil
	}, queue.OrmDirtyNotificationEntity, 100)

	consumerRunner.Wait()
}
