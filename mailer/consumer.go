package mailer

import (
	"fmt"
	"github.com/confluentinc/confluent-kafka-go/kafka"
	"github.com/vmihailenco/msgpack/v5"
	"myapp/config"
)

func CreateConsumer() {
	c, err := kafka.NewConsumer(&kafka.ConfigMap{
		"bootstrap.servers": config.Config("BOOTSTRAP_SERVER"),
		"security.protocol": "SASL_SSL",
		"sasl.username":     config.Config("USERNAME"),
		"sasl.password":     config.Config("PASSWORD"),
		"sasl.mechanism":    "PLAIN",
		"group.id":          config.Config("CONSUMER_GROUP_ID"),
		"auto.offset.reset": "earliest",
	})

	if err != nil {
		panic(err)
	}

	err = c.SubscribeTopics([]string{config.Config("CONSUME_TOPIC")}, nil)
	if err != nil {
		panic(err)
	}

	for {
		msg, err := c.ReadMessage(-1)
		if err == nil {
			fmt.Printf("Message on %s: %s\n", msg.TopicPartition, string(msg.Value))
			email := new(Email)
			err = msgpack.Unmarshal(msg.Value, email)
			fmt.Println(email)
			err = EmailRepoImpl{}.Create(email)
			if err != nil {
				fmt.Println(err)
			}
		} else {
			// The client will automatically try to recover from all errors.
			fmt.Printf("Consumer error: %v (%v)\n", err, msg)
		}
	}

	c.Close()
}
