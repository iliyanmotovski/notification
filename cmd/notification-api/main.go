package main

import (
	"context"
	"log"

	"github.com/iliyanmotovski/notification/api/notification-api/rest/router"
	"github.com/iliyanmotovski/notification/pkg/constant"
	"github.com/iliyanmotovski/notification/pkg/service/clock"
	"github.com/iliyanmotovski/notification/pkg/service/config"
	gracefulshutdown "github.com/iliyanmotovski/notification/pkg/service/graceful_shutdown"
	"github.com/iliyanmotovski/notification/pkg/service/orm"
	"github.com/iliyanmotovski/notification/pkg/service/server"
)

// @title NOTIFICATION API
// @version 1.0
// @contact.name Iliyan Motovski
// @contact.email iliyan.motovski@gmail.com

// @BasePath /v1
func main() {
	appContext := context.Background()

	waitChan := make(chan struct{})
	gracefulshutdown.GracefulShutdown(func() {
		waitChan <- struct{}{}
	})

	configService, err := config.NewConfigService(constant.AppNameNotificationAPI, "../../config/config.yaml")
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

	clockService := clock.NewClockService()

	server.StartGinHTTPServer(
		appContext,
		configService,
		router.Init(ormRegistry, clockService),
		waitChan,
	)
}
