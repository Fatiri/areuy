package sender

import (
	"github.com/twilio/twilio-go"
	api "github.com/twilio/twilio-go/rest/api/v2010"
)

type Twilio interface {
	SendChatWhatsApp(to, body string) error
}

type twilioCtx struct {
	from   string
	client *twilio.RestClient
}

func NewTwilio(username, password, from string) Twilio {
	client := twilio.NewRestClientWithParams(twilio.ClientParams{
		Username: username,
		Password: password,
	})

	return &twilioCtx{
		client: client,
		from:   from,
	}
}

func (tw *twilioCtx) SendChatWhatsApp(to, body string) error {
	params := &api.CreateMessageParams{}

	params.SetFrom("whatsapp:+" + tw.from)
	params.SetTo("whatsapp:+" + to)
	params.SetBody(body)

	_, err := tw.client.Api.CreateMessage(params)
	if err != nil {
		return err
	}

	return nil
}
