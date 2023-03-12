package sms

import (
	"API-REST/services/conf"

	twilio "github.com/twilio/twilio-go"
	openapi "github.com/twilio/twilio-go/rest/api/v2010"
)

var client *twilio.RestClient

func Setup() error {
	client = twilio.NewRestClientWithParams(twilio.ClientParams{
		Username: conf.Env.GetString("TWILIO_ACCOUNT_SID"),
		Password: conf.Env.GetString("TWILIO_AUTH_TOKEN"),
	})

	return nil
}

func Send(to string, body string) error {
	params := &openapi.CreateMessageParams{}
	params.SetFrom(conf.Env.GetString("TWILIO_PHONE"))
	params.SetTo(to) // example: +34633444555
	params.SetBody(body)

	// Send sms
	_, err := client.Api.CreateMessage(params)
	// response, _ := json.Marshal(*resp)
	// fmt.Println("Response: " + string(response))

	return err
}
