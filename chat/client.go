package main

import (
	"github.com/gorilla/websocket"
)

// clientはチャットを行っている1人のユーザーを表します
type client struct {
	socket *websocket.Conn
	send   chan []byte
	room   *room
}

func (c *client) read () {
	for {
		if _, msg, err := c.socket.ReadMessage(); err == nil {
			c.room.forward <- msg
		} else {
			break
		}
	}
	if err := c.socket.Close(); err != nil {
		// TODO エラー処理
	}
}

func (c *client) write() {
	for msg := range c.send {
		if err := c.socket.WriteMessage(websocket.TextMessage, msg); err != nil {
			break
		}
	}
	if err := c.socket.Close(); err != nil {
		// TODO エラー処理
	}
}
