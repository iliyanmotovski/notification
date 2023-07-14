package sms

type Gateway interface {
	SendSMS(*Message) (string, error)
}

type Message struct {
	Text         string
	MobileNumber string
}
