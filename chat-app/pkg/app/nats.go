package app

import (
	"log"
	"os"
	"time"

	"github.com/Ta-SeenJunaid/hybrid-cloud-chat/chat-app/pkg/apis"
	"github.com/nats-io/nats.go"
)

var Sender string
var Receiver string
var NatsConnection *nats.Conn

func InitializeNats() {
	env, exists := os.LookupEnv("SENDER")
	if !exists || len(env) == 0 {
		log.Panic("SENDER env cannot be empty")
	}
	Sender = env

	env, exists = os.LookupEnv("RECEIVER")
	if !exists || len(env) == 0 {
		log.Panic("RECEIVER env cannot be empty")
	}
	Receiver = env

	natsURL, exists := os.LookupEnv("NATS_URL")
	if !exists {
		natsURL = nats.DefaultURL
	}

	var err error
	NatsConnection, err = nats.Connect(natsURL)
	if err != nil {
		log.Fatalln(err)
	}

	_, err = NatsConnection.Subscribe(Sender, func(msg *nats.Msg) {
		Messages = append(Messages, apis.Message{
			Author: Receiver,
			Body:   string(msg.Data),
			Time:   time.Now().Format("2022-06-02 02:23:45"),
		})
	})
	if err != nil {
		log.Fatalln(err)
	}
}
