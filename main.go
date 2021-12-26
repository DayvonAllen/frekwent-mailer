package main

import (
	"fmt"
	"github.com/gofiber/fiber/v2"
	"github.com/robfig/cron/v3"
	"log"
	"myapp/mailer"
	"os"
	"os/signal"
	"runtime"
	"sync"
	"time"
)

var scheduler *cron.Cron

func init() {
	//conn := database.ConnectToDB()
	mailer.Instance = mailer.CreateMailer()
	go mailer.Instance.ListenForMail()
	go mailer.CreateConsumer()

	// cron job for resending failed emails
	scheduler = cron.New(cron.WithChain(
		cron.Recover(cron.DefaultLogger),
	))

	_, err := scheduler.AddFunc("@every 4h30m", func() {
		// TODO check every 4 hours and 30 minutes for failed emails and try to resend them
		stat := mailer.Failed
		emails, err := mailer.EmailRepoImpl{}.FindAllByStatus(&stat)
		if err != nil {
			fmt.Println("none found")
		}

		var wg sync.WaitGroup
		wg.Add(len(*emails))

		for _, email := range *emails {
			email := email
			go func(e mailer.Email) {
				defer wg.Done()
				mailer.SendMessage(&email)
			}(email)
		}
	})

	if err != nil {
		panic(err)
	}

	scheduler.Start()

	//_, err = conn.EmailCollection.DeleteMany(context.TODO(), bson.M{})
	//if err != nil {
	//	return
	//}
}

func main() {
	app := fiber.New()

	//go mailer.SendTestMessage()

	c := make(chan os.Signal, 1)
	signal.Notify(c, os.Interrupt)

	go func() {
		_ = <-c
		fmt.Println("Shutting down...")
		scheduler.Stop()
		time.Sleep(20 * time.Second)
		fmt.Println(runtime.NumGoroutine())
		_ = app.Shutdown()
	}()

	if err := app.Listen(":8081"); err != nil {
		log.Panic(err)
	}
}
