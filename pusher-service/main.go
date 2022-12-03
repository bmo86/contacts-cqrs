package main

import (
	"contacts_cqrs/events"
	"fmt"
	"log"
	"net/http"

	"github.com/kelseyhightower/envconfig"
)

type Config struct {
	NatsAddr string `envconfig:"NATS_ADRRESS"`
}

func main() {
	var cfg Config
	err := envconfig.Process("", &cfg)
	if err != nil {
		log.Fatal(err)
	}

	hub := NewHub()

	n, err := events.NewNats(fmt.Sprintf("nats://%s", cfg.NatsAddr))
	if err != nil {
		log.Fatal(err)
	}
	err = n.OnCreatedContact(func(m events.CreatedContactMessage) {
		hub.Broadcast(newCreatedContactMessage(m.ID, m.Name, m.Lastname, m.Image, m.Email, m.Phone, m.Status, m.CreatedAt, m.UpdateData), nil)
	})

	if err != nil {
		log.Fatal(err)
	}

	events.SetEvent(n)
	defer events.Close()

	go hub.Run()
	http.HandleFunc("/ws", hub.HandlerWs)
	err = http.ListenAndServe(":8080", nil)
	if err != nil {
		log.Fatal(err)
	}

}
