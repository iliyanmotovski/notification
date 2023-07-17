package middleware

import (
	"github.com/gin-gonic/gin"
	"github.com/iliyanmotovski/notification/pkg/constant"
	"github.com/iliyanmotovski/notification/pkg/service/orm"
)

func ORMPerRequestMiddleware(ormRegistry orm.RegistryService) gin.HandlerFunc {
	return func(ctx *gin.Context) {
		ctx.Set(constant.ORMService, ormRegistry.GetORMService())
	}
}
