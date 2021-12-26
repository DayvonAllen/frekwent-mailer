package mailer


type Message struct {
	From string
	FromName string
	To string
	Subject string
	// what template we want to use(go or jet)
	Template string
	Attachments []string
	// any data we want to pass to the email, it's an interface because we don't know what it's going to be
	Data interface{}
}
