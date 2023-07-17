package controller

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/iliyanmotovski/notification/api/notification-api/rest/dto"
	"github.com/iliyanmotovski/notification/pkg/constant"
	notificationModel "github.com/iliyanmotovski/notification/pkg/domain/model/notification"
	"github.com/iliyanmotovski/notification/pkg/service/notification"
	"github.com/iliyanmotovski/notification/pkg/service/server"
)

type NotificationController struct {
}

// @Description Send notification async
// @Tags Notification
// @Param body body dto.RequestDTONotification true "Request in body"
// @Router /notification/send-async [post]
// @Success 200
// @Failure 400 {object} server.Error
// @Failure 500
func (n *NotificationController) PostSendNotificationAsyncAction(ctx *gin.Context) {
	requestDTO := &dto.RequestDTONotification{}

	err := server.ShouldBindJSON(ctx, requestDTO)
	if server.HandleError(ctx, err) {
		return
	}

	notificationService := ctx.MustGet(constant.NotificationService).(notification.Engine)

	err = notificationModel.SendNotificationAsync(notificationService, requestDTO)
	if server.HandleError(ctx, err) {
		return
	}

	ctx.JSON(http.StatusOK, gin.H{})
}
