package main

import (
	entitybeeorm "github.com/iliyanm/notification/pkg/entity/entity_beeorm"
	"github.com/iliyanm/notification/pkg/service/clock"
	"github.com/iliyanm/notification/pkg/service/config"
	"github.com/iliyanm/notification/pkg/service/orm"
	"log"
)

func main() {
	configService, err := config.NewConfigService("api", "../../config")
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

	ormService.Flush(&entitybeeorm.NotificationEntity{
		Asd:       "aaaaaaa",
		CreatedAt: clockService.Now(),
	})
}
