package main

import (
	"fmt"

	"github.com/devfeel/dotweb"
)

// CORS Example
// This demonstrates how to handle Cross-Origin Resource Sharing (CORS)

func main() {
	// init DotApp
	app := dotweb.New()

	// API endpoint with CORS headers
	app.HttpServer.GET("/api/data", func(ctx dotweb.Context) error {
		// Set CORS headers
		ctx.Response().Header().Set("Access-Control-Allow-Origin", "*")
		ctx.Response().Header().Set("Access-Control-Allow-Methods", "GET, POST, PUT, DELETE, OPTIONS")
		ctx.Response().Header().Set("Access-Control-Allow-Headers", "Content-Type, Authorization")

		return ctx.WriteJson(map[string]interface{}{
			"message": "Hello from API",
			"data":     []string{"item1", "item2", "item3"},
		})
	})

	// POST endpoint
	app.HttpServer.POST("/api/data", func(ctx dotweb.Context) error {
		// Set CORS headers
		ctx.Response().Header().Set("Access-Control-Allow-Origin", "*")
		ctx.Response().Header().Set("Access-Control-Allow-Methods", "GET, POST, PUT, DELETE, OPTIONS")
		ctx.Response().Header().Set("Access-Control-Allow-Headers", "Content-Type, Authorization")

		return ctx.WriteJson(map[string]string{
			"status":  "success",
			"message": "Data received",
		})
	})

	fmt.Println("CORS Example server starting on :8080")
	fmt.Println("API endpoint: GET /api/data")
	fmt.Println("API endpoint: POST /api/data")
	fmt.Println("Test with: curl -H 'Origin: *' http://localhost:8080/api/data")
	err := app.StartServer(8080)
	fmt.Println("server error => ", err)
}
