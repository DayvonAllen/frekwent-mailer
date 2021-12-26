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

func init() {
	_ = database.ConnectToDB()
	mailer.Instance = mailer.CreateMailer()
	go mailer.Instance.ListenForMail()
}
func main() {
	app := fiber.New()

	go mailer.SendTestMessage()

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