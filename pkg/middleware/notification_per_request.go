package middleware

import (
	"github.com/gin-gonic/gin"
	"github.com/iliyanmotovski/notification/pkg/constant"
	"github.com/iliyanmotovski/notification/pkg/service/clock"
	"github.com/iliyanmotovski/notification/pkg/service/notification"
	"github.com/iliyanmotovski/notification/pkg/service/orm"
)

func NotificationServicePerRequestMiddleware(clockService clock.IClock) gin.HandlerFunc {
	return func(ctx *gin.Context) {
		ormService := ctx.MustGet(constant.ORMService).(orm.Engine)
		ctx.Set(constant.NotificationService, notification.NewNotificationService(ormService, clockService))
	}
}
