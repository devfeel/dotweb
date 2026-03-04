// Package main demonstrates RESTful JSON API in DotWeb.
// Run: go run main.go
// Test: See README.md for curl examples
package main

import (
	"encoding/json"
	"fmt"
	"strconv"
	"sync"

	"github.com/devfeel/dotweb"
)

// User represents a user entity
type User struct {
	ID    int    `json:"id"`
	Name  string `json:"name"`
	Email string `json:"email"`
}

// ErrorResponse represents an error response
type ErrorResponse struct {
	Error string `json:"error"`
}

// SuccessResponse represents a success response
type SuccessResponse struct {
	Message string `json:"message"`
	Data    interface{} `json:"data,omitempty"`
}

// In-memory database
var (
	users  = make(map[int]*User)
	nextID = 1
	mu     sync.RWMutex
)

func main() {
	// Initialize sample data
	users[1] = &User{ID: 1, Name: "Alice", Email: "alice@example.com"}
	users[2] = &User{ID: 2, Name: "Bob", Email: "bob@example.com"}
	nextID = 3
	
	// Create DotWeb app
	app := dotweb.New()
	app.SetDevelopmentMode()
	
	// Set JSON content type for all responses
	app.HttpServer.Use(func(ctx dotweb.Context) error {
		ctx.Response().Header().Set("Content-Type", "application/json")
		return ctx.NextHandler()
	})
	
	// Global error handler
	app.SetExceptionHandle(func(ctx dotweb.Context, err error) {
		ctx.Response().SetContentType(dotweb.MIMEApplicationJSONCharsetUTF8)
		ctx.WriteJsonC(500, ErrorResponse{Error: err.Error()})
	})
	
	// 404 handler
	app.SetNotFoundHandle(func(ctx dotweb.Context) {
		ctx.Response().SetContentType(dotweb.MIMEApplicationJSONCharsetUTF8)
		ctx.WriteJsonC(404, ErrorResponse{Error: "Not found"})
	})
	
	// API group
	api := app.HttpServer.Group("/api")
	
	// ===== User CRUD =====
	
	// GET /api/users - List all users
	api.GET("/users", listUsers)
	
	// GET /api/users/:id - Get user by ID
	api.GET("/users/:id", getUser)
	
	// POST /api/users - Create user
	api.POST("/users", createUser)
	
	// PUT /api/users/:id - Update user
	api.PUT("/users/:id", updateUser)
	
	// DELETE /api/users/:id - Delete user
	api.DELETE("/users/:id", deleteUser)
	
	// Health check
	api.GET("/health", func(ctx dotweb.Context) error {
		return ctx.WriteJsonC(200, map[string]string{"status": "ok"})
	})
	
	fmt.Println("🚀 JSON API running at http://localhost:8080")
	fmt.Println("\nAPI Endpoints:")
	fmt.Println("  GET    /api/health       - Health check")
	fmt.Println("  GET    /api/users        - List all users")
	fmt.Println("  GET    /api/users/:id    - Get user by ID")
	fmt.Println("  POST   /api/users        - Create user")
	fmt.Println("  PUT    /api/users/:id    - Update user")
	fmt.Println("  DELETE /api/users/:id    - Delete user")
	
	if err := app.StartServer(8080); err != nil {
		fmt.Printf("Server error: %v\n", err)
	}
}

// listUsers returns all users
func listUsers(ctx dotweb.Context) error {
	mu.RLock()
	defer mu.RUnlock()
	
	list := make([]*User, 0, len(users))
	for _, u := range users {
		list = append(list, u)
	}
	
	return ctx.WriteJsonC(200, SuccessResponse{
		Message: "success",
		Data:    list,
	})
}

// getUser returns a user by ID
func getUser(ctx dotweb.Context) error {
	idStr := ctx.GetRouterName("id")
	id, err := strconv.Atoi(idStr)
	if err != nil {
		return ctx.WriteJsonC(400, ErrorResponse{Error: "Invalid user ID"})
	}
	
	mu.RLock()
	user, ok := users[id]
	mu.RUnlock()
	
	if !ok {
		return ctx.WriteJsonC(404, ErrorResponse{Error: "User not found"})
	}
	
	return ctx.WriteJsonC(200, SuccessResponse{
		Message: "success",
		Data:    user,
	})
}

// createUser creates a new user
func createUser(ctx dotweb.Context) error {
	var user User
	if err := json.Unmarshal(ctx.Request().PostBody(), &user); err != nil {
		return ctx.WriteJsonC(400, ErrorResponse{Error: "Invalid JSON"})
	}
	
	if user.Name == "" || user.Email == "" {
		return ctx.WriteJsonC(400, ErrorResponse{Error: "Name and email required"})
	}
	
	mu.Lock()
	user.ID = nextID
	nextID++
	users[user.ID] = &user
	mu.Unlock()
	
	return ctx.WriteJsonC(201, SuccessResponse{
		Message: "User created",
		Data:    &user,
	})
}

// updateUser updates a user
func updateUser(ctx dotweb.Context) error {
	idStr := ctx.GetRouterName("id")
	id, err := strconv.Atoi(idStr)
	if err != nil {
		return ctx.WriteJsonC(400, ErrorResponse{Error: "Invalid user ID"})
	}
	
	mu.RLock()
	user, ok := users[id]
	mu.RUnlock()
	
	if !ok {
		return ctx.WriteJsonC(404, ErrorResponse{Error: "User not found"})
	}
	
	var update User
	if err := json.Unmarshal(ctx.Request().PostBody(), &update); err != nil {
		return ctx.WriteJsonC(400, ErrorResponse{Error: "Invalid JSON"})
	}
	
	mu.Lock()
	if update.Name != "" {
		user.Name = update.Name
	}
	if update.Email != "" {
		user.Email = update.Email
	}
	mu.Unlock()
	
	return ctx.WriteJsonC(200, SuccessResponse{
		Message: "User updated",
		Data:    user,
	})
}

// deleteUser deletes a user
func deleteUser(ctx dotweb.Context) error {
	idStr := ctx.GetRouterName("id")
	id, err := strconv.Atoi(idStr)
	if err != nil {
		return ctx.WriteJsonC(400, ErrorResponse{Error: "Invalid user ID"})
	}
	
	mu.Lock()
	defer mu.Unlock()
	
	if _, ok := users[id]; !ok {
		return ctx.WriteJsonC(404, ErrorResponse{Error: "User not found"})
	}
	
	delete(users, id)
	
	return ctx.WriteJsonC(200, SuccessResponse{
		Message: "User deleted",
	})
}
