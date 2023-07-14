package email

type Gateway interface {
	SendTemplate(*Message) error
}

type Message struct {
	From         string
	FromName     string
	ReplyTo      string
	To           string
	Subject      string
	TemplateName string
	TemplateData map[string]interface{}
}
