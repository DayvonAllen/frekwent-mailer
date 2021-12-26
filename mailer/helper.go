package mailer

import (
	"bytes"
	"fmt"
	"github.com/vanng822/go-premailer/premailer"
	mail "github.com/xhit/go-simple-mail/v2"
	"html/template"
	"myapp/config"
	"os"
	"strconv"
	"time"
)

func (m *Mail) ListenForMail() {
	// endless for loop that runs in the background
	for {
		// take anything we get from the jobs type and do something with it
		// msg listens for any incoming jobs on the jobs channel
		msg := <-m.Jobs
		// send message
		err := m.Send(msg)
		if err != nil {
			// send an error to the result channel and also set success to false
			m.Results <- Result{false, err}
		} else {
			m.Results <- Result{true, nil}
		}
	}
}

func (m *Mail) Send(msg Message) error {
	// make a decision on whether we are using an SMTP server or API

	return m.SendSMTPMessage(msg)
}

func (m *Mail) SendSMTPMessage(msg Message) error {
	formattedMessage, err := m.buildHTMLMessage(msg)
	if err != nil {
		return err
	}

	plainMessage, err := m.buildPlainTextMessage(msg)
	if err != nil {
		return err
	}

	server := mail.NewSMTPClient()

	server.Host = m.Host
	server.Port = m.Port
	server.Username = m.Username
	server.Password = m.Password
	server.Encryption = m.getEncryption(m.Encryption)
	// keepAlive will keep a connection to the mail server alive at all times
	server.KeepAlive = false
	server.ConnectTimeout = 10 * time.Second
	server.SendTimeout = 10 * time.Second

	smtpClient, err := server.Connect()

	if err != nil {
		return err
	}

	email := mail.NewMSG()
	email.SetFrom(msg.From).
		AddTo(msg.To).
		SetSubject(msg.Subject)

	email.SetBody(mail.TextHTML, formattedMessage)
	// alternative body, if html message fails to work properly
	email.AddAlternative(mail.TextPlain, plainMessage)

	if len(msg.Attachments) > 0 {
		for _, x := range msg.Attachments {
			email.AddAttachment(x)
		}
	}

	// try sending email
	err = email.Send(smtpClient)

	if err != nil {
		return err
	}

	return nil
}

func (m *Mail) buildHTMLMessage(msg Message) (string, error) {
	// using go templates
	templateToRender := fmt.Sprintf("%s/%s.html.tmpl", m.Templates, msg.Template)

	t, err := template.New("email-html").ParseFiles(templateToRender)

	if err != nil {
		return "", err
	}

	// we need this to execute the template
	var tpl bytes.Buffer

	if err = t.ExecuteTemplate(&tpl, "body", msg.Data); err != nil {
		return "", err
	}

	formattedMessage := tpl.String()

	// inline CSS to make sure the email renders the way it's supposed to on all email clients
	formattedMessage, err = m.inlineCSS(formattedMessage)

	if err != nil {
		return "", err
	}

	return formattedMessage, nil
}

func (m *Mail) buildPlainTextMessage(msg Message) (string, error) {
	// using go templates
	templateToRender := fmt.Sprintf("%s/%s.plain.tmpl", m.Templates, msg.Template)

	t, err := template.New("email-html").ParseFiles(templateToRender)

	if err != nil {
		return "", err
	}

	// we need this to execute the template
	var tpl bytes.Buffer

	if err = t.ExecuteTemplate(&tpl, "body", msg.Data); err != nil {
		return "", err
	}

	plainMessage := tpl.String()

	return plainMessage, nil
}

func (m *Mail) getEncryption(encryption string) mail.Encryption {
	// constants for encryption types in mail.Encryption from the simple mail library
	switch encryption {
	// most common
	case "tls":
		return mail.EncryptionTLS
	case "ssl":
		return mail.EncryptionSSL
	// for development only
	case "none":
		return mail.EncryptionNone
	default:
		return mail.EncryptionTLS
	}
}

func (m *Mail) inlineCSS(s string) (string, error) {
	// after building html, we want to use the CSS inliner
	options := premailer.Options{
		RemoveClasses:     false,
		CssToAttributes:   false,
		KeepBangImportant: true,
	}

	prem, err := premailer.NewPremailerFromString(s, &options)

	if err != nil {
		return "", err
	}

	html, err := prem.Transform()

	if err != nil {
		return "", err
	}

	return html, nil

}

func CreateMailer() *Mail {
	port, err := strconv.Atoi(config.Config("PORT"))

	if err != nil {
		panic(err)
	}

	// get working directory
	rootPath, err := os.Getwd()

	if err != nil {
		panic(err)
	}

	m := Mail{
		Domain:      config.Config("MAIL_DOMAIN"),
		Templates:   rootPath + "/mail",
		Host:        config.Config("HOST"),
		Port:        port,
		Username:    config.Config("USERNAME"),
		Password:    config.Config("PASSWORD"),
		Encryption:  config.Config("ENCRYPTION"),
		FromName:    config.Config("FROM_NAME"),
		FromAddress: config.Config("FROM_ADDRESS"),
		Jobs:        make(chan Message, 20),
		Results:     make(chan Result, 20),
		API: config.Config("API"),
		APIKey: config.Config("API_KEY"),
		APIUrl: config.Config("API_URL"),
	}

	return &m
}

func SendTestMessage() {
	msg := Message{
		From: "test@example.com",
		To: "you@there.com",
		Subject: "test subject - sent in dev mode",
		Template: "test",
		Attachments: nil,
		Data: nil,
	}

	// message for 3rd party API
	//msg := Message{
	//	From: "admin@secretstash.tech",
	//	To: config.Config("MY_EMAIL"),
	//	Subject: "test subject - sent using an api",
	//	Template: "test",
	//	Attachments: nil,
	//	Data: nil,
	//}

	Instance.Jobs <- msg
	res := <-Instance.Results

	if res.Error != nil {
		fmt.Println("couldn't send email")
	}
}