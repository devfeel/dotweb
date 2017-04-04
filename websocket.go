package dotweb

import (
	"golang.org/x/net/websocket"
	"net/http"
)

type WebSocket struct {
	Conn *websocket.Conn
}

//get http request
func (ws *WebSocket) Request() *http.Request {
	return ws.Conn.Request()
}

//send message from websocket.conn
func (ws *WebSocket) SendMessage(msg string) error {
	return websocket.Message.Send(ws.Conn, msg)
}

//read message from websocket.conn
func (ws *WebSocket) ReadMessage() (string, error) {
	str := ""
	err := websocket.Message.Receive(ws.Conn, &str)
	return str, err
}
