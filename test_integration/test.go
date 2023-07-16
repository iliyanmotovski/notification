package test_integration

import (
	"log"

	"github.com/iliyanm/notification/pkg/service/config"
	"github.com/iliyanm/notification/pkg/service/orm"
)

var ormService orm.Engine

func createTestEngine() {
	// make orm instance only once
	if ormService == nil {
		configService, err := config.NewConfigService("test", "../config/config_test.yaml")
		if err != nil {
			log.Fatal(err)
		}

		registry, _, err := orm.NewORMRegistryService(configService)
		if err != nil {
			log.Fatal(err)
		}

		ormService = registry.GetORMService()
		ormService.ExecuteAlters()
	}

	ormService.TruncateTables()
	ormService.GetCacheService().Clear()
	ormService.GetCacheService(orm.StreamsPool).Clear()
}
