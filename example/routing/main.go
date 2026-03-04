// Package main demonstrates routing patterns in DotWeb.
// Run: go run main.go
// Test routes listed in the output
package main

import (
	"fmt"
	"github.com/devfeel/dotweb"
)

func main() {
	// Create DotWeb app
	app := dotweb.New()
	app.SetDevelopmentMode()
	
	// ===== Basic Routes =====
	
	// GET request
	app.HttpServer.GET("/", func(ctx dotweb.Context) error {
		return ctx.WriteString("GET / - Home page")
	})
	
	// POST request
	app.HttpServer.POST("/users", func(ctx dotweb.Context) error {
		return ctx.WriteString("POST /users - Create user")
	})
	
	// PUT request
	app.HttpServer.PUT("/users/:id", func(ctx dotweb.Context) error {
		id := ctx.GetRouterName("id")
		return ctx.WriteString("PUT /users/" + id + " - Update user")
	})
	
	// DELETE request
	app.HttpServer.DELETE("/users/:id", func(ctx dotweb.Context) error {
		id := ctx.GetRouterName("id")
		return ctx.WriteString("DELETE /users/" + id + " - Delete user")
	})
	
	// Any method
	app.HttpServer.Any("/any", func(ctx dotweb.Context) error {
		return ctx.WriteString("ANY /any - Method: " + ctx.Request().Method)
	})
	
	// ===== Path Parameters =====
	
	// Single parameter
	app.HttpServer.GET("/users/:id", func(ctx dotweb.Context) error {
		id := ctx.GetRouterName("id")
		return ctx.WriteString("User ID: " + id)
	})
	
	// Multiple parameters
	app.HttpServer.GET("/users/:userId/posts/:postId", func(ctx dotweb.Context) error {
		userId := ctx.GetRouterName("userId")
		postId := ctx.GetRouterName("postId")
		return ctx.WriteString(fmt.Sprintf("User: %s, Post: %s", userId, postId))
	})
	
	// Wildcard (catch-all)
	app.HttpServer.GET("/files/*filepath", func(ctx dotweb.Context) error {
		filepath := ctx.GetRouterName("filepath")
		return ctx.WriteString("File path: " + filepath)
	})
	
	// ===== Route Groups =====
	
	// API group
	api := app.HttpServer.Group("/api")
	api.GET("/health", func(ctx dotweb.Context) error {
		return ctx.WriteString(`{"status": "ok"}`)
	})
	api.GET("/version", func(ctx dotweb.Context) error {
		return ctx.WriteString(`{"version": "1.0.0"}`)
	})
	
	// API v1 group
	v1 := app.HttpServer.Group("/api/v1")
	v1.GET("/users", func(ctx dotweb.Context) error {
		return ctx.WriteString(`{"users": ["Alice", "Bob"]}`)
	})
	v1.POST("/users", func(ctx dotweb.Context) error {
		return ctx.WriteString(`{"created": true}`)
	})
	
	// ===== Print test routes =====
	fmt.Println("🚀 Routing example running at http://localhost:8080")
	fmt.Println("\nBasic routes:")
	fmt.Println("  curl http://localhost:8080/")
	fmt.Println("  curl -X POST http://localhost:8080/users")
	fmt.Println("  curl -X PUT http://localhost:8080/users/123")
	fmt.Println("  curl -X DELETE http://localhost:8080/users/123")
	fmt.Println("  curl -X POST http://localhost:8080/any")
	fmt.Println("\nPath parameters:")
	fmt.Println("  curl http://localhost:8080/users/42")
	fmt.Println("  curl http://localhost:8080/users/42/posts/100")
	fmt.Println("  curl http://localhost:8080/files/path/to/file.txt")
	fmt.Println("\nRoute groups:")
	fmt.Println("  curl http://localhost:8080/api/health")
	fmt.Println("  curl http://localhost:8080/api/version")
	fmt.Println("  curl http://localhost:8080/api/v1/users")
	fmt.Println("  curl -X POST http://localhost:8080/api/v1/users")
	
	if err := app.StartServer(8080); err != nil {
		fmt.Printf("Server error: %v\n", err)
	}
}
