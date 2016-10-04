package handler

import (
	"chat/trace"
	"log"
	"net/http"

	"github.com/gorilla/websocket"
	"os"
	"github.com/stretchr/objx"
)

const (
	socketBufferSize = 1024
	messageBufferSize = 256
)

type room struct {
	forward chan *message
	join    chan *client
	leave   chan *client
	clients map[*client]bool
	t       trace.Tracer
}

var upgrader = &websocket.Upgrader{
	ReadBufferSize:  socketBufferSize,
	WriteBufferSize: socketBufferSize,
}

func (r *room) run() {
	for {
		select {
		case client := <-r.join:
			r.clients[client] = true
			r.t.Trace("New client joined")
		case client := <-r.leave:
			delete(r.clients, client)
			close(client.send)
			r.t.Trace("Client left")
		case msg := <-r.forward:
			for client := range r.clients {
				select {
				case client.send <- msg:
					r.t.Trace(msg.Message, " -- sent to client")
				default:
					delete(r.clients, client)
					close(client.send)
					r.t.Trace(msg.Message, " -- failed to send, cleaned up client")
				}
			}
		}
	}
}

func (r *room) ServeHTTP(w http.ResponseWriter, req *http.Request) {
	socket, err := upgrader.Upgrade(w, req, nil)
	if err != nil {
		log.Fatal("ServeHTTP:", err)
		return
	}
	authCookie, err := req.Cookie("auth")
	if err != nil {
		log.Fatal("Failed to get auth cookie:", err)
		return
	}
	client := &client{
		socket: socket,
		send:   make(chan *message, messageBufferSize),
		room:   r,
		userData: objx.MustFromBase64(authCookie.Value),
	}

	r.join <- client
	defer func() {
		r.leave <- client
	}()

	go client.write()
	client.read()
}

func NewRoom() *room {
	r := &room{
		forward: make(chan *message),
		join:    make(chan *client),
		leave:   make(chan *client),
		clients: make(map[*client]bool),
		t:       trace.New(os.Stdout),
	}

	go r.run()

	return r
}
