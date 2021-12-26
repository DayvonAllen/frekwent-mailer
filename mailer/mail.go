package mailer

import "github.com/robfig/cron/v3"

type Mail struct {
	// where mail is coming from
	Domain string
	// path to html templates
	Templates   string
	Host        string
	Port        int
	Username    string
	Password    string
	Encryption  string
	FromName    string
	FromAddress string
	Jobs        chan Email
	// what happened when we tried to send mail
	Results   chan Result
	Scheduler *cron.Cron
}

var Instance *Mail
