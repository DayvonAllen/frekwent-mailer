package main

import (
	"fmt"
	"github.com/gofiber/fiber/v2"
	"github.com/robfig/cron/v3"
	"log"
	"myapp/database"
	"myapp/mailer"
	"os"
	"os/signal"
)

func init() {
	_ = database.ConnectToDB()
	mailer.Instance = mailer.CreateMailer()
	go mailer.Instance.ListenForMail()

	// cron job for resending failed emails
	scheduler := cron.New(cron.WithChain(
		cron.Recover(cron.DefaultLogger),
	))

	_, err := scheduler.AddFunc("@every 4h30m", func() {
		// TODO check every 4 hours and 30 minutes for failed emails and try to resend them
		fmt.Println("test")
	})

	if err != nil {
		panic(err)
	}
}
func main() {
	app := fiber.New()

	//go mailer.SendTestMessage()

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