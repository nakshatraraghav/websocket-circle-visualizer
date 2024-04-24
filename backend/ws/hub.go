package ws

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

func (hub *Hub) RunHub() {
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
