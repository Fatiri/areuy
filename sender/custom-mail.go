package sender

import (
	"strconv"

	"gopkg.in/gomail.v2"
)

type CustomMailConfig struct {
	SenderAddress string
	Host          string
	Port          string
	Username      string
	Password      string
}

type CustomMailPayload struct {
	ReceiverEmail string
	Subject       string
	Message       string
}

type CustomMail interface {
	V1(payload CustomMailPayload) error
}

type CustomMailCtx struct {
	config CustomMailConfig
}

func NewCustomMail(config CustomMailConfig) CustomMail {
	return &CustomMailCtx{
		config: config,
	}
}

// CustomMailEmailSender CustomMail smtp sender
func (cm *CustomMailCtx) V1(payload CustomMailPayload) error {
	m := gomail.NewMessage()
	m.SetAddressHeader("From", cm.config.Username, cm.config.SenderAddress)
	m.SetHeader("To",
		m.FormatAddress(payload.ReceiverEmail, ""),
	)
	m.SetHeader("Subject", payload.Subject)
	m.SetBody("text/html", payload.Message)

	port, _ := strconv.Atoi(cm.config.Port)

	d := gomail.NewDialer(cm.config.Host, port, cm.config.Username, cm.config.Password)
	if err := d.DialAndSend(m); err != nil {
		return err
	}

	return nil
}
