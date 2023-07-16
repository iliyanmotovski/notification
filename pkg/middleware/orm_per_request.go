package middleware

import (
	"github.com/gin-gonic/gin"
	"github.com/iliyanm/notification/pkg/constant"
	"github.com/iliyanm/notification/pkg/service/orm"
)

func ORMPerRequestMiddleware(ormRegistry orm.RegistryService) gin.HandlerFunc {
	return func(ctx *gin.Context) {
		ctx.Set(constant.ORMService, ormRegistry.GetORMService())
	}
}
