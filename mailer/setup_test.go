package mailer

import (
	"fmt"
	"github.com/ory/dockertest"
	"github.com/ory/dockertest/docker"
	"log"
	"os"
	"testing"
	"time"
)

var pool *dockertest.Pool
var resource *dockertest.Resource

var mailer = Mail{
	Domain:      "localhost",
	Templates:   "./testdata/mail",
	Host:        "localhost",
	Port:        1026,
	Encryption:  "none",
	FromName:    "Joe",
	FromAddress: "test@test.com",
	Jobs:        make(chan Message, 1),
	Results:     make(chan Result, 1),
}

func TestMain(m *testing.M) {
	p, err := dockertest.NewPool("")

	if err != nil {
		log.Fatal("Could not connect to docker")
	}

	pool = p

	opts := dockertest.RunOptions{
		Repository: "mailhog/mailhog",
		Tag: "latest",
		Env: []string{},
		ExposedPorts: []string{"1025", "8025"},
		PortBindings: map[docker.Port][]docker.PortBinding{
			"1025": {
				{HostIP: "0.0.0.0", HostPort: "1026"},
			},
			"8025": {
				{HostIP: "0.0.0.0", HostPort: "8026"},
			},
		},
	}

	resource, err = pool.RunWithOptions(&opts)

	if err != nil {
		fmt.Println(err)
		_ = pool.Purge(resource)
		log.Fatal("Could not start resource")
	}

	// wait for mailhog to start
	time.Sleep(3 * time.Second)

	go mailer.ListenForMail()

	code := m.Run()

	err = pool.Purge(resource)

	if err != nil {
		log.Fatalf("Could not purge resource: %s", err)
	}

	os.Exit(code)
}