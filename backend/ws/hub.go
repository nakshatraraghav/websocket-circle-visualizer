package ws

import "github.com/nakshatraraghav/hashed-tokens-assignment/backend/lib"

type Hub struct {
	clients map[*Client]bool

	broadcast  chan []byte
	register   chan *Client
	deregister chan *Client
}

func NewHub() *Hub {
	return &Hub{
		clients:    make(map[*Client]bool),
		broadcast:  make(chan []byte),
		register:   make(chan *Client),
		deregister: make(chan *Client),
	}
}

func (hub *Hub) BroadcaseSinCosSamples(sinc, cosc chan float64) {
	for {
		var msg []byte

		select {
		case sin := <-sinc:
			msg = []byte("sin:" + lib.FloatToString(sin))
		case cos := <-cosc:
			msg = []byte("cos:" + lib.FloatToString(cos))
		}

		for client := range hub.clients {
			select {
			case client.send <- msg:
			default:
				close(client.send)
				delete(hub.clients, client)
			}
		}
	}
}

func (hub *Hub) RunHub(sinc, cosc chan float64) {
	go hub.BroadcaseSinCosSamples(sinc, cosc)

	for {
		select {
		case client := <-hub.register:
			hub.clients[client] = true
		case client := <-hub.deregister:
			_, ok := hub.clients[client]
			if ok {
				delete(hub.clients, client)
				close(client.send)
			}
		case message := <-hub.broadcast:

			for client := range hub.clients {
				select {
				case client.send <- message:
				default:
					close(client.send)
					delete(hub.clients, client)
				}
			}
		}
	}
}
