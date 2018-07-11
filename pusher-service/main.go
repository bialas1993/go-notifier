package main

import (
	"fmt"
	"log"
	"net/http"
	"time"

	"github.com/kelseyhightower/envconfig"
	"github.com/bialas1993/go-notifier/event"
	"github.com/tinrab/retry"
)

type Config struct {
	NatsAddress string `envconfig:"NATS_ADDRESS"`
}

func main() {
	var cfg Config
	err := envconfig.Process("", &cfg)
	if err != nil {
		log.Fatal(err)
	}

	// Connect to Nats
	hub := newHub()
	retry.ForeverSleep(2*time.Second, func(_ int) error {
		es, err := event.NewNats(fmt.Sprintf("nats://%s", cfg.NatsAddress))
		if err != nil {
			log.Println(err)
			return err
		}

		// Push messages to clients
		err = es.OnNotifyCreated(func(m event.NotifyCreatedMessage) {
			log.Printf("Notify received: %v\n", m)
			hub.broadcast(
				newNotifyCreatedMessage(m.ID, m.Title, m.Body, m.Service, m.CreatedAt),
				nil)
		})
		if err != nil {
			log.Println(err)
			return err
		}

		event.SetEventStore(es)
		return nil
	})
	defer event.Close()

	// Run WebSocket server
	go hub.run()
	http.HandleFunc("/pusher", hub.handleWebSocket)
	err = http.ListenAndServe(":8080", nil)
	if err != nil {
		log.Fatal(err)
	}
}
