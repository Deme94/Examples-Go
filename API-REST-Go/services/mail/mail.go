package mail

import (
	"API-REST/services/conf"
	"os"

	"github.com/wneessen/go-mail"
)

var client *mail.Client

type Mail struct {
	From        string
	To          []string
	Subject     string
	Body        string
	Attachments []*os.File
}

func Setup() error {
	var err error
	client, err = mail.NewClient(
		conf.Conf.GetString("mailHost"),
		mail.WithPort(conf.Conf.GetInt("mailPort")),
		//mail.WithSMTPAuth(mail.SMTPAuthPlain), // uncomment this when using real host with required smtp auth
		mail.WithUsername(conf.Conf.GetString("mailUsername")),
		mail.WithPassword(conf.Conf.GetString("mailPassword")),
		mail.WithTLSPolicy(mail.TLSOpportunistic),
	)
	return err
}

func Send(m *Mail) error {
	msg := mail.NewMsg()
	// msg.SetMessageID()
	// msg.SetDate()
	// msg.SetBulk()
	err := msg.From(m.From)
	if err != nil {
		return err
	}
	err = msg.To(m.To...)
	if err != nil {
		return err
	}
	msg.Subject(m.Subject)
	msg.SetBodyString(mail.TypeTextHTML, m.Body)
	for _, f := range m.Attachments {
		msg.AttachReader(f.Name(), f)
	}

	return client.DialAndSend(msg)
}
