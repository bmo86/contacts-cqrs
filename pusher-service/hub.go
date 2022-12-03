package main

import (
	"contacts_cqrs/means"
	"encoding/json"
	"log"
	"net/http"
	"sync"

	"github.com/gorilla/websocket"
)

var upgrader = websocket.Upgrader{
	CheckOrigin: func(r *http.Request) bool { return true },
}

type Hub struct {
	clients    []*Client
	register   chan *Client
	unregister chan *Client
	mutex      *sync.Mutex
}

func NewHub() *Hub {
	return &Hub{
		clients:    make([]*Client, 0),
		register:   make(chan *Client),
		unregister: make(chan *Client),
		mutex:      &sync.Mutex{},
	}
}

func (hub *Hub) Run() {
	for {
		select {
		case cli := <-hub.register:
			hub.onConnect(cli)
		case cli := <-hub.unregister:
			hub.disconnect(cli)
		}
	}
}

func (hub *Hub) Broadcast(msg interface{}, ignore *Client) {
	data, _ := json.Marshal(msg)
	for _, client := range hub.clients {
		if client != ignore {
			client.outbound <- data
		}
	}
}

func (h *Hub) onConnect(cli *Client) {
	log.Println("Connect client : ", cli.socket.RemoteAddr())
	h.mutex.Lock()
	defer h.mutex.Unlock()
	cli.id = cli.socket.RemoteAddr().String()
	h.clients = append(h.clients, cli)
}

func (h *Hub) disconnect(cli *Client) {
	log.Println("Connect Client : ", cli.socket.RemoteAddr())
	cli.Close()
	h.mutex.Lock()
	defer h.mutex.Unlock()

	i := -1
	for j, c := range h.clients {
		if c.id == cli.id {
			i = j
			break
		}
	}

	copy(h.clients[i:], h.clients[i+1:])
	h.clients[len(h.clients)-1] = nil
	h.clients = h.clients[:len(h.clients)-1]
}

// handler to ws conn
func (hub *Hub) HandlerWs(w http.ResponseWriter, r *http.Request) {
	socket, err := upgrader.Upgrade(w, r, nil)
	if err != nil {
		log.Println(err)
		means.MessageErr(http.StatusInternalServerError, err.Error(), w)
		return
	}
	client := NewClient(hub, socket)
	hub.register <- client
	go client.Write()
}
