// Package main demonstrates the simplest DotWeb application.
// Run: go run main.go
// Test: curl http://localhost:8080/
package main

import (
	"fmt"
	"github.com/devfeel/dotweb"
)

func main() {
	// Create a new DotWeb application
	app := dotweb.New()
	
	// Register a simple route
	app.HttpServer.GET("/", func(ctx dotweb.Context) error {
		return ctx.WriteString("Hello, DotWeb! 🐾")
	})
	
	// Start the server
	fmt.Println("🚀 Server running at http://localhost:8080")
	fmt.Println("Press Ctrl+C to stop")
	
	if err := app.StartServer(8080); err != nil {
		fmt.Printf("Server error: %v\n", err)
	}
}
