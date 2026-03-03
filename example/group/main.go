package main

import (
	"fmt"
	"github.com/devfeel/dotweb"
)

func main() {
	// Create DotWeb app
	app := dotweb.New()
	
	// Set global 404 handler
	app.SetNotFoundHandle(func(ctx dotweb.Context) {
		ctx.Response().Header().Set("Content-Type", "application/json")
		ctx.WriteString(`{"code": 404, "message": "Global 404 - Page not found"}`)
	})
	
	// Create API group
	apiGroup := app.HttpServer.Group("/api")
	
	// Set group-level 404 handler
	apiGroup.SetNotFoundHandle(func(ctx dotweb.Context) {
		ctx.Response().Header().Set("Content-Type", "application/json")
		ctx.WriteString(`{"code": 404, "message": "API 404 - Resource not found", "hint": "Check API documentation for available endpoints"}`)
	})
	
	// Register API routes
	apiGroup.GET("/users", func(ctx dotweb.Context) error {
		return ctx.WriteString(`{"users": ["Alice", "Bob", "Charlie"]}`)
	})
	
	apiGroup.GET("/health", func(ctx dotweb.Context) error {
		return ctx.WriteString(`{"status": "ok"}`)
	})
	
	// Create Web group (no custom 404 handler, will use global)
	webGroup := app.HttpServer.Group("/web")
	
	webGroup.GET("/index", func(ctx dotweb.Context) error {
		return ctx.WriteString("<h1>Welcome to Web</h1>")
	})
	
	// Root route
	app.HttpServer.GET("/", func(ctx dotweb.Context) error {
		return ctx.WriteString("Welcome to DotWeb! Try:\n" +
			"- GET /api/users (exists)\n" +
			"- GET /api/unknown (API 404)\n" +
			"- GET /web/index (exists)\n" +
			"- GET /web/unknown (Global 404)\n" +
			"- GET /unknown (Global 404)")
	})
	
	fmt.Println("Server starting on :8080...")
	fmt.Println("\nTest routes:")
	fmt.Println("  curl http://localhost:8080/                    - Welcome page")
	fmt.Println("  curl http://localhost:8080/api/users         - API: Users list")
	fmt.Println("  curl http://localhost:8080/api/health        - API: Health check")
	fmt.Println("  curl http://localhost:8080/api/unknown       - API: 404 (group handler)")
	fmt.Println("  curl http://localhost:8080/web/index         - Web: Index page")
	fmt.Println("  curl http://localhost:8080/web/unknown       - Web: 404 (global handler)")
	fmt.Println("  curl http://localhost:8080/unknown           - Global: 404 (global handler)")
	
	// Start server
	err := app.StartServer(8080)
	if err != nil {
		fmt.Println("Server error:", err)
	}
}
