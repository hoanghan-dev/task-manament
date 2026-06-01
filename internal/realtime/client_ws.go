package realtime

import (
	"log"

	"github.com/gorilla/websocket"
)

type Client struct {
	conn *websocket.Conn
	send chan []byte
}

func NewClient(conn *websocket.Conn) *Client {
	return &Client{
		conn: conn,
		send: make(chan []byte, 256),
	}
}

func (c *Client) Send(message []byte) {
	c.send <- message
}

// gửi dữ liệu cho client
func (c *Client) WriteLoop() {
	defer c.conn.Close()

	for message := range c.send {
		err := c.conn.WriteMessage(websocket.TextMessage, message)
		if err != nil {
			log.Println("Write websocket message failed:", err)
			return
		}
	}
}

// Đọc dữ liệu từ client
func (c *Client) ReadLoop() {
	defer c.conn.Close()

	for {
		_, _, err := c.conn.ReadMessage()

		if err != nil {
			log.Println("Read websocket message failed:", err)
			return
		}
	}
}
