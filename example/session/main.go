// Package main demonstrates session management in DotWeb.
// Run: go run main.go
// Test: curl http://localhost:8080/login  -> sets session
//       curl http://localhost:8080/user   -> get user from session
//       curl http://localhost:8080/logout -> destroy session
package main

import (
	"fmt"
	"github.com/devfeel/dotweb"
	"github.com/devfeel/dotweb/session"
)

func main() {
	// Create DotWeb app
	app := dotweb.New()
	
	// Enable session with default runtime config
	app.HttpServer.SetEnabledSession(true)
	app.HttpServer.SetSessionConfig(session.NewDefaultRuntimeConfig())
	
	// Login - set session
	app.HttpServer.GET("/login", func(ctx dotweb.Context) error {
		ctx.SetSession("user", "Alice")
		ctx.SetSession("role", "admin")
		return ctx.WriteString("✅ Logged in as Alice (admin)")
	})
	
	// Get user info from session
	app.HttpServer.GET("/user", func(ctx dotweb.Context) error {
		user := ctx.GetSession("user")
		role := ctx.GetSession("role")
		
		if user == nil {
			return ctx.WriteString("❌ Not logged in. Visit /login first.")
		}
		
		return ctx.WriteString(fmt.Sprintf("👤 User: %v\n🔑 Role: %v", user, role))
	})
	
	// Logout - destroy session
	app.HttpServer.GET("/logout", func(ctx dotweb.Context) error {
		err := ctx.DestorySession()
		if err != nil {
			return ctx.WriteString("❌ Logout failed: " + err.Error())
		}
		return ctx.WriteString("✅ Logged out successfully")
	})
	
	// Check session exists
	app.HttpServer.GET("/check", func(ctx dotweb.Context) error {
		if ctx.HasSession("user") {
			return ctx.WriteString("✅ Session exists")
		}
		return ctx.WriteString("❌ No session found")
	})
	
	fmt.Println("🚀 Session example running at http://localhost:8080")
	fmt.Println("\nTest routes:")
	fmt.Println("  curl http://localhost:8080/login   -> Set session")
	fmt.Println("  curl http://localhost:8080/user    -> Get session data")
	fmt.Println("  curl http://localhost:8080/check   -> Check session")
	fmt.Println("  curl http://localhost:8080/logout  -> Destroy session")
	
	if err := app.StartServer(8080); err != nil {
		fmt.Printf("Server error: %v\n", err)
	}
}
