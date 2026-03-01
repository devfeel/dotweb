package main

import (
	"fmt"
	"time"

	"github.com/devfeel/dotweb"
)

// JWT Middleware Example
// This demonstrates how to use JWT authentication with dotweb

func main() {
	// init DotApp
	app := dotweb.New()

	// Set route
	app.HttpServer.GET("/api/public", func(ctx dotweb.Context) error {
		return ctx.WriteString("This is a public endpoint")
	})

	// Protected route - in real app, add JWT validation middleware
	app.HttpServer.GET("/api/private", func(ctx dotweb.Context) error {
		return ctx.WriteJson(map[string]interface{}{
			"message": "This is a private endpoint",
			"user":    "admin",
		})
	})

	// Login endpoint - returns simple token (use proper JWT in production)
	app.HttpServer.POST("/api/login", func(ctx dotweb.Context) error {
		// In production, validate credentials from database
		username := ctx.PostFormValue("username")
		password := ctx.PostFormValue("password")

		// Simple validation (replace with real auth)
		if username == "admin" && password == "password" {
			// Generate simple token (use proper JWT library in production)
			token := generateSimpleToken(username)
			return ctx.WriteJson(map[string]string{
				"token": token,
			})
		}
		return ctx.WriteString("Invalid credentials")
	})

	fmt.Println("JWT Example server starting on :8080")
	fmt.Println("Public endpoint: GET /api/public")
	fmt.Println("Private endpoint: GET /api/private")
	fmt.Println("Login: POST /api/login?username=admin&password=password")
	err := app.StartServer(8080)
	fmt.Println("server error => ", err)
}

// generateSimpleToken creates a simple token (for demonstration only)
// In production, use: github.com/golang-jwt/jwt
func generateSimpleToken(username string) string {
	exp := time.Now().Add(time.Hour).Unix()
	return fmt.Sprintf("%s.%d.%s", username, exp, "signature")
}
