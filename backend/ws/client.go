package ws

import (
	"bytes"
	"log"
	"net/http"
	"time"

	"github.com/gorilla/websocket"
)

const (
	// Time allowed to write a message to the peer.
	WRITE_WAIT = 10 * time.Second

	// Time allowed to read the next pong message from the peer.
	PONG_WAIT = 60 * time.Second

	// Send pings to peer with this period. Must be less than pongWait.
	PING_PERIOD = (PONG_WAIT * 9) / 10

	// Maximum message size allowed from peer.
	MAX_MESSAGE_SIZE = 512
)

var (
	newline = []byte{'\n'}
	space   = []byte{' '}
)

type Client struct {
	hub  *Hub
	conn *websocket.Conn
	send chan []byte
}

var upgrader = websocket.Upgrader{
	ReadBufferSize:  1024,
	WriteBufferSize: 1024,
	CheckOrigin: func(r *http.Request) bool {
		return true
	},
}

func ServeWebSocket(hub *Hub, w http.ResponseWriter, r *http.Request) {
	conn, err := upgrader.Upgrade(w, r, nil)
	if err != nil {
		log.Println(err.Error())
		log.Println("failed to upgrade client connectiono from ", r.RemoteAddr, "to websocket")
		w.WriteHeader(http.StatusInternalServerError)
		w.Write([]byte("failed to upgrade connection to websocket"))
		return
	}

	client := &Client{
		hub:  hub,
		conn: conn,
		send: make(chan []byte),
	}

	hub.register <- client

	go client.write()
	go client.read()
}

func (c *Client) read() {
	// defer the de-registration of the client
	defer func() {
		c.hub.deregister <- c
		c.conn.Close()
	}()

	c.conn.SetReadLimit(MAX_MESSAGE_SIZE)
	c.conn.SetReadDeadline(time.Now().Add(PONG_WAIT))
	c.conn.SetPongHandler(func(string) error {
		c.conn.SetReadDeadline(time.Now().Add(PONG_WAIT))
		return nil
	})

	for {
		_, msg, err := c.conn.ReadMessage()
		if err != nil {
			if websocket.IsUnexpectedCloseError(err, websocket.CloseGoingAway, websocket.CloseAbnormalClosure) {
				log.Printf("error: %v", err)
			}
			break
		}

		message := bytes.TrimSpace(bytes.Replace(msg, newline, space, -1))
		c.hub.broadcast <- message
	}

}

func (c *Client) write() {
	t := time.NewTicker(PING_PERIOD)
	defer func() {
		t.Stop()
		c.conn.Close()
	}()

	for {
		select {
		case message, ok := <-c.send:
			c.conn.SetWriteDeadline(time.Now().Add(WRITE_WAIT))
			if !ok {
				c.conn.WriteMessage(websocket.CloseMessage, []byte{})
				return
			}

			w, err := c.conn.NextWriter(websocket.TextMessage)
			if err != nil {
				return
			}

			w.Write(message)
			n := len(c.send)
			for i := 0; i < n; i++ {
				w.Write(newline)
				w.Write(<-c.send)
			}

			if err := w.Close(); err != nil {
				return
			}
		case <-t.C:
			c.conn.SetWriteDeadline(time.Now().Add(WRITE_WAIT))
			if err := c.conn.WriteMessage(websocket.PingMessage, nil); err != nil {
				return
			}
		}
	}
}
