package main

import (
	"fmt"
	"github.com/gofiber/fiber/v2"
	"log"
	"myapp/database"
	"myapp/mailer"
	"os"
	"os/signal"
)

var MailerInstance *mailer.Mail

func init() {
	_ = database.ConnectToDB()
	MailerInstance = mailer.CreateMailer()
}
func main() {
	app := fiber.New()

	fmt.Println(MailerInstance)

	c := make(chan os.Signal, 1)
	signal.Notify(c, os.Interrupt)

	go func() {
		_ = <-c
		fmt.Println("Shutting down...")
		_ = app.Shutdown()
	}()

	if err := app.Listen(":8081"); err != nil {
		log.Panic(err)
	}
}