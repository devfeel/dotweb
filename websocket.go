package dotweb

import (
	"net/http"

	"golang.org/x/net/websocket"
)

type WebSocket struct {
	Conn *websocket.Conn
}

// Request get http request
func (ws *WebSocket) Request() *http.Request {
	return ws.Conn.Request()
}

// SendMessage send message from websocket.conn
func (ws *WebSocket) SendMessage(msg string) error {
	return websocket.Message.Send(ws.Conn, msg)
}

// ReadMessage read message from websocket.conn
func (ws *WebSocket) ReadMessage() (string, error) {
	str := ""
	err := websocket.Message.Receive(ws.Conn, &str)
	return str, err
}
