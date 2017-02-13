package dotweb

import (
	"golang.org/x/net/websocket"
)

type WebSocket struct {
	Conn *websocket.Conn
}

//send message from websocket.conn
func (ws *WebSocket) SendMessage(msg string) {
	websocket.Message.Send(ws.Conn, msg)
}

//read message from websocket.conn
func (ws *WebSocket) ReadMessage() (string, error) {
	str := ""
	err := websocket.Message.Receive(ws.Conn, &str)
	return str, err
}
