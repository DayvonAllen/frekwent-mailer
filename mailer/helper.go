package mailer

import (
	"bytes"
	"fmt"
	"github.com/sendgrid/sendgrid-go"
	"github.com/sendgrid/sendgrid-go/helpers/mail"
	"github.com/vanng822/go-premailer/premailer"
	legacyMail "github.com/xhit/go-simple-mail/v2"
	"html/template"
	"log"
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
	formattedMessage, err := m.buildHTMLMessage(msg)
	if err != nil {
		return err
	}

	plainMessage, err := m.buildPlainTextMessage(msg)
	if err != nil {
		return err
	}

	msg.From = config.Config("BUSINESS_EMAIL")
	msg.To = config.Config("MY_EMAIL")

	from := mail.NewEmail("Example User", msg.From)
	subject := msg.Subject
	to := mail.NewEmail("Example User", msg.To)
	plainTextContent := plainMessage
	htmlContent := formattedMessage
	message := mail.NewSingleEmail(from, subject, to, plainTextContent, htmlContent)
	client := sendgrid.NewSendClient(config.Config("SENDGRID_API_KEY"))
	response, err := client.Send(message)

	if err != nil {
		log.Println(err)
		return err
	} else {
		fmt.Println(response.StatusCode)
		fmt.Println(response.Body)
		fmt.Println(response.Headers)
	}

	return nil
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

	server := legacyMail.NewSMTPClient()

	server.Host = m.Host
	server.Port = m.Port
	server.Username = m.Username
	server.Password = m.Password
	server.Encryption = m.getEncryption(m.Encryption)
	// keepAlive will keep a connection to the legacyMail server alive at all times
	server.KeepAlive = false
	server.ConnectTimeout = 10 * time.Second
	server.SendTimeout = 10 * time.Second

	smtpClient, err := server.Connect()

	if err != nil {
		return err
	}

	email := legacyMail.NewMSG()
	email.SetFrom(msg.From).
		AddTo(msg.To).
		SetSubject(msg.Subject)

	email.SetBody(legacyMail.TextHTML, formattedMessage)
	// alternative body, if html message fails to work properly
	email.AddAlternative(legacyMail.TextPlain, plainMessage)

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

func (m *Mail) getEncryption(encryption string) legacyMail.Encryption {
	// constants for encryption types in legacyMail.Encryption from the simple legacyMail library
	switch encryption {
	// most common
	case "tls":
		return legacyMail.EncryptionTLS
	case "ssl":
		return legacyMail.EncryptionSSL
	// for development only
	case "none":
		return legacyMail.EncryptionNone
	default:
		return legacyMail.EncryptionTLS
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
		From: config.Config("BUSINESS_EMAIL"),
		To: config.Config("MY_EMAIL"),
		Subject: "test subject - sent by send grid",
		Template: "test",
		Attachments: nil,
		Data: nil,
	}

	Instance.Jobs <- msg
	res := <-Instance.Results

	if res.Error != nil {
		fmt.Println(res.Error)
		fmt.Println("couldn't send email")
	}
}