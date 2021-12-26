package mailer

import (
	"errors"
	"testing"
)

func TestMail_SendSMTPMessage(t *testing.T) {
	msg := Message{
		From: "test@example.com",
		To: "you@there.com",
		Subject: "test",
		Template: "test",
		Attachments: []string{"./testdata/mail/test.html.tmpl"},
		Data: nil,
	}

	err := mailer.SendSMTPMessage(msg)

	if err != nil {
		t.Error(err)
	}
}

func TestMail_Send(t *testing.T) {
	msg := Message{
		From: "test@example.com",
		To: "you@there.com",
		Subject: "test",
		Template: "test",
		Attachments: []string{"./testdata/mail/test.html.tmpl"},
		Data: nil,
	}

	mailer.Jobs <- msg
	res := <-mailer.Results

	if res.Error != nil {
		t.Error(res.Error)
	}

	msg.To = "wrong"
	mailer.Jobs <-msg
	res = <-mailer.Results

	if res.Error == nil {
		t.Error(errors.New("no error received with invalid to address"))
	}
}
