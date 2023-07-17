package notification

import (
	"github.com/iliyanmotovski/notification/api/notification-api/rest/dto"
	"github.com/iliyanmotovski/notification/pkg/service/notification"
)

func SendNotificationAsync(notificationService notification.Engine, requestDTO *dto.RequestDTONotification) error {
	notificationRequest := &notification.Notification{}

	if requestDTO.SMS != nil {
		notificationRequest.SMS = &notification.SMS{
			Text:         requestDTO.SMS.Text,
			MobileNumber: requestDTO.SMS.MobileNumber,
		}
	}

	if requestDTO.Email != nil {
		email := &notification.Email{
			From:         requestDTO.Email.From,
			FromName:     requestDTO.Email.FromName,
			ReplyTo:      requestDTO.Email.ReplyTo,
			To:           requestDTO.Email.To,
			Subject:      requestDTO.Email.Subject,
			TemplateName: requestDTO.Email.TemplateName,
		}

		if len(requestDTO.Email.TemplateData) > 0 {
			templateData := map[string]interface{}{}

			for _, data := range requestDTO.Email.TemplateData {
				templateData[data.Key] = data.Value
			}

			email.TemplateData = templateData
		}

		notificationRequest.Email = email
	}

	if requestDTO.SlackMessage != nil {
		notificationRequest.SlackMessage = &notification.SlackMessage{
			BotName:     requestDTO.SlackMessage.BotName,
			ChannelName: requestDTO.SlackMessage.ChannelName,
			Message:     requestDTO.SlackMessage.Message,
		}
	}

	if err := notificationService.SendAsync(notificationRequest); err != nil {
		return err
	}

	return nil
}
