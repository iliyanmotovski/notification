package router

import (
	"github.com/gin-gonic/gin"
	"github.com/iliyanm/notification/api/notification-api/rest/controller"
	"github.com/iliyanm/notification/pkg/constant"
	"github.com/iliyanm/notification/pkg/domain/view/doc"
	"github.com/iliyanm/notification/pkg/middleware"
	"github.com/iliyanm/notification/pkg/service/clock"
	"github.com/iliyanm/notification/pkg/service/orm"
	"github.com/iliyanm/notification/pkg/service/server"
)

func Init(ormRegistry orm.RegistryService, clockService clock.IClock) server.GinInitHandler {
	return func(ginEngine *gin.Engine) {
		ginEngine.GET("/doc", doc.GetAPIDoc(constant.AppNameNotificationAPI))
		ginEngine.Static("static", "../../static")

		v1Router := ginEngine.Group("/v1")

		v1Router.Use(middleware.ORMPerRequestMiddleware(ormRegistry))

		notificationController := &controller.NotificationController{}
		notificationGroup := v1Router.Group("/notification")
		notificationGroup.Use(middleware.NotificationServicePerRequestMiddleware(clockService))
		{
			notificationGroup.POST("/send-async", notificationController.PostSendNotificationAsyncAction)
		}
	}
}
