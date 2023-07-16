package dto

type RequestDTONotification struct {
	SMS          *RequestDTOSMS          `json:"sms,omitempty"`
	Email        *RequestDTOEmail        `json:"email,omitempty"`
	SlackMessage *RequestDTOSlackMessage `json:"slack_message,omitempty"`
}

type RequestDTOSMS struct {
	Text         string `json:"text" validate:"required"`
	MobileNumber string `json:"mobile_number" validate:"required"`
}

type RequestDTOEmail struct {
	From         string                    `json:"from" validate:"required,email"`
	FromName     string                    `json:"from_name" validate:"required"`
	ReplyTo      string                    `json:"reply_to"`
	To           string                    `json:"to" validate:"required,email"`
	Subject      string                    `json:"subject" validate:"required"`
	TemplateName string                    `json:"template_name" validate:"required"`
	TemplateData []*RequestDTOTemplateDate `json:"template_data,omitempty" validate:"omitempty,dive"`
}

type RequestDTOTemplateDate struct {
	Key   string `json:"key" validate:"required"`
	Value string `json:"value" validate:"required"`
}

type RequestDTOSlackMessage struct {
	BotName     string `json:"bot_name" validate:"required"`
	ChannelName string `json:"channel_name" validate:"required"`
	Message     string `json:"message" validate:"required"`
}
