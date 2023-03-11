package sms

import (
	"API-REST/services/conf"

	twilio "github.com/twilio/twilio-go"
	openapi "github.com/twilio/twilio-go/rest/api/v2010"
)

var client *twilio.RestClient

func Setup() error {
	client = twilio.NewRestClientWithParams(twilio.ClientParams{
		Username: conf.Conf.GetString("twilioAccountSID"),
		Password: conf.Conf.GetString("twilioAuthToken"),
	})

	return nil
}

func Send(to string, body string) error {
	params := &openapi.CreateMessageParams{}
	params.SetFrom(conf.Conf.GetString("twilioPhone"))
	params.SetTo(to)
	params.SetBody(body)

	// Send sms
	_, err := client.Api.CreateMessage(params)
	// response, _ := json.Marshal(*resp)
	// fmt.Println("Response: " + string(response))

	return err
}
