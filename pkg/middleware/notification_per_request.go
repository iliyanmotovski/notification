package middleware

import (
	"github.com/gin-gonic/gin"
	"github.com/iliyanm/notification/pkg/constant"
	"github.com/iliyanm/notification/pkg/service/clock"
	"github.com/iliyanm/notification/pkg/service/notification"
	"github.com/iliyanm/notification/pkg/service/orm"
)

func NotificationServicePerRequestMiddleware(clockService clock.IClock) gin.HandlerFunc {
	return func(ctx *gin.Context) {
		ormService := ctx.MustGet(constant.ORMService).(orm.Engine)
		ctx.Set(constant.NotificationService, notification.NewNotificationService(ormService, clockService))
	}
}
