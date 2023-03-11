package sms

import (
	"API-REST/services/conf"
	"os"

	twilio "github.com/twilio/twilio-go"
	openapi "github.com/twilio/twilio-go/rest/api/v2010"
)

var client *twilio.RestClient

func Setup() error {
	sid := conf.Conf.GetString("twilioAccountSID")
	token := conf.Conf.GetString("twilioAuthToken")

	err := os.Setenv("TWILIO_ACCOUNT_SID", sid)
	if err != nil {
		return err
	}
	err = os.Setenv("TWILIO_AUTH_TOKEN", token)
	if err != nil {
		return err
	}

	client = twilio.NewRestClient()

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
