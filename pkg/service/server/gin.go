package server

import (
	"context"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/iliyanm/notification/pkg/service/config"
)

type GinInitHandler func(ginEngine *gin.Engine)

func StartGinHTTPServer(ctx context.Context, configService config.Config, init GinInitHandler, waitChan chan struct{}) {
	port, ok := configService.GetString("server.port")
	if !ok {
		panic("config missing server.port")
	}

	ginEngine := gin.New()
	init(ginEngine)

	srv := &http.Server{
		Addr:    ":" + port,
		Handler: ginEngine,
	}

	go func() {
		if err := srv.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			panic(err)
		}
	}()

	<-waitChan

	if err := srv.Shutdown(ctx); err != nil {
		panic(err)
	}
}
