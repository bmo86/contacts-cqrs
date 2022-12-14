package main

import (
	"contacts_cqrs/database"
	"contacts_cqrs/events"
	"contacts_cqrs/repository"
	"contacts_cqrs/search"
	"fmt"
	"log"
	"net/http"

	"github.com/gorilla/mux"
	"github.com/kelseyhightower/envconfig"
)

type Config struct {
	PostgresDB           string `envconfig:"POSTGRES_DB"`
	PostgresUser         string `envconfig:"POSTGRES_USER"`
	PostgresPassword     string `envconfig:"POSTGRES_PASSWORD"`
	NatsAddress          string `envconfig:"NATS_ADDRESS"`
	ElasticsearchAddress string `envconfig:"ELASTICSEARCH_ADDRESS"`
}

func newRouter() (router *mux.Router) {
	router = mux.NewRouter()
	router.HandleFunc("/contacts", handlerList).Methods(http.MethodGet)
	router.HandleFunc("/search", handlerSearch).Methods(http.MethodGet)
	return
}

func main() {

	var cfg Config
	if err := envconfig.Process("", &cfg); err != nil {
		log.Fatal(err)
	}
	//conn to database-Postgres
	addrPost := fmt.Sprintf("postgres://%s:%s@postgres/%s?sslmode=disable", cfg.PostgresUser, cfg.PostgresPassword, cfg.PostgresDB)
	repo, err := database.CNDatabase(addrPost)
	if err != nil {
		log.Fatal(err)
	}
	repository.SetRepositoriDB(repo)

	//conn to elastic
	es, err := search.NewCNElastic(fmt.Sprintf("http://%s", cfg.ElasticsearchAddress))
	if err != nil {
		log.Fatal(err)
	}
	search.SetRepoSearch(es)
	defer search.Close()

	//conn to nats
	n, err := events.NewNats(fmt.Sprintf("nats://%s", cfg.NatsAddress))
	if err != nil {
		log.Fatal(err)
	}
	err = n.OnCreatedContact(onCreatedContacts)
	if err != nil {
		log.Fatal(err)
	}
	events.SetEvent(n)
	defer events.Close()

	//conn to router
	r := newRouter()
	if err := http.ListenAndServe(":8080", r); err != nil {
		log.Fatal(err)
	}
}
