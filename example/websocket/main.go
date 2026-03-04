// Package main demonstrates WebSocket in DotWeb.
// Run: go run main.go
// Test: Use a WebSocket client (e.g., wscat or browser)
package main

import (
	"fmt"
	"log"

	"github.com/devfeel/dotweb"
	"golang.org/x/net/websocket"
)

// Connected clients
var clients = make(map[*websocket.Conn]bool)
var broadcast = make(chan string)

func main() {
	// Create DotWeb app
	app := dotweb.New()
	app.SetDevelopmentMode()
	
	// WebSocket endpoint - echo server
	app.HttpServer.GET("/ws", func(ctx dotweb.Context) error {
		// Check if WebSocket request
		if !ctx.IsWebSocket() {
			return ctx.WriteString("This endpoint requires WebSocket connection")
		}
		
		// Get WebSocket connection
		ws := ctx.WebSocket()
		
		// Register client
		clients[ws.Conn] = true
		log.Printf("Client connected. Total: %d", len(clients))
		
		// Send welcome message
		ws.SendMessage("Welcome to DotWeb WebSocket!")
		
		// Read messages in loop
		for {
			msg, err := ws.ReadMessage()
			if err != nil {
				log.Printf("Client disconnected: %v", err)
				delete(clients, ws.Conn)
				break
			}
			
			log.Printf("Received: %s", msg)
			
			// Echo back
			ws.SendMessage("Echo: " + msg)
		}
		
		return nil
	})
	
	// WebSocket chat endpoint
	app.HttpServer.GET("/chat", func(ctx dotweb.Context) error {
		if !ctx.IsWebSocket() {
			return ctx.WriteString("This endpoint requires WebSocket connection")
		}
		
		ws := ctx.WebSocket()
		clients[ws.Conn] = true
		
		// Get username from query
		username := ctx.Request().QueryString("name")
		if username == "" {
			username = "Anonymous"
		}
		
		// Announce join
		broadcast <- fmt.Sprintf("🔔 %s joined the chat", username)
		
		// Read messages
		for {
			msg, err := ws.ReadMessage()
			if err != nil {
				delete(clients, ws.Conn)
				broadcast <- fmt.Sprintf("🚪 %s left the chat", username)
				break
			}
			
			broadcast <- fmt.Sprintf("💬 %s: %s", username, msg)
		}
		
		return nil
	})
	
	// HTTP endpoint to check WebSocket status
	app.HttpServer.GET("/status", func(ctx dotweb.Context) error {
		return ctx.WriteString(fmt.Sprintf(
			"WebSocket Server Status\n"+
				"Connected clients: %d\n"+
				"Endpoints:\n"+
				"  - ws://localhost:8080/ws (Echo)\n"+
				"  - ws://localhost:8080/chat?name=YourName (Chat)",
			len(clients),
		))
	})
	
	// Start broadcast goroutine
	go handleBroadcast()
	
	fmt.Println("🚀 WebSocket example running at http://localhost:8080")
	fmt.Println("\nWebSocket endpoints:")
	fmt.Println("  ws://localhost:8080/ws          - Echo server")
	fmt.Println("  ws://localhost:8080/chat?name=X - Chat room")
	fmt.Println("\nHTTP status:")
	fmt.Println("  curl http://localhost:8080/status")
	fmt.Println("\nTest with wscat:")
	fmt.Println("  wscat -c ws://localhost:8080/ws")
	fmt.Println("  wscat -c 'ws://localhost:8080/chat?name=Alice'")
	
	if err := app.StartServer(8080); err != nil {
		log.Fatal(err)
	}
}

// handleBroadcast sends messages to all connected clients
func handleBroadcast() {
	for msg := range broadcast {
		for conn := range clients {
			// Use websocket.Message.Send directly
			err := websocket.Message.Send(conn, msg)
			if err != nil {
				conn.Close()
				delete(clients, conn)
			}
		}
	}
}
