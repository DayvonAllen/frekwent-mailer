package mailer

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
	Jobs        chan Message
	// what happened when we tried to send mail
	Results chan Result
	// which email api we are using(optional, get rid of for an app that is using only one API)
	API    string
	APIKey string
	APIUrl string
}