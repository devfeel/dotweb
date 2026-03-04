# DotWeb Examples

This directory contains examples demonstrating DotWeb features.

## Quick Start (5 minutes)

```bash
cd quickstart
go run main.go
# Visit http://localhost:8080
```

## Examples Index

### 🚀 Getting Started

| Example | Description | Complexity |
|---------|-------------|------------|
| [quickstart](./quickstart) | Minimal "Hello World" | ★☆☆ |
| [routing](./routing) | Route patterns, params, groups | ★★☆ |
| [group](./group) | Route grouping with 404 handlers | ★★☆ |

### 🔧 Core Features

| Example | Description | Complexity |
|---------|-------------|------------|
| [middleware](./middleware) | Logging, auth, CORS | ★★☆ |
| [session](./session) | Session management | ★★☆ |
| [bind](./bind) | Data binding (form, JSON) | ★★☆ |
| [config](./config) | Configuration files | ★★☆ |
| [router](./router) | Advanced routing | ★★☆ |

### 🌐 Web Features

| Example | Description | Complexity |
|---------|-------------|------------|
| [json-api](./json-api) | RESTful API with CRUD | ★★☆ |
| [file-upload](./file-upload) | File upload/download | ★★☆ |
| [websocket](./websocket) | WebSocket (echo, chat) | ★★★ |

### 🧪 Testing

| Example | Description | Complexity |
|---------|-------------|------------|
| [mock](./mock) | Mock mode for testing | ★★☆ |

## Feature Examples

### 1. Basic Routing
```go
app.HttpServer.GET("/", handler)
app.HttpServer.POST("/users", handler)
app.HttpServer.PUT("/users/:id", handler)
app.HttpServer.DELETE("/users/:id", handler)
```

### 2. Route Parameters
```go
// Path parameter
app.HttpServer.GET("/users/:id", func(ctx dotweb.Context) error {
    id := ctx.GetRouterName("id")
    return ctx.WriteString("User ID: " + id)
})

// Wildcard
app.HttpServer.GET("/files/*filepath", func(ctx dotweb.Context) error {
    path := ctx.GetRouterName("filepath")
    return ctx.WriteString("File: " + path)
})
```

### 3. Route Groups
```go
api := app.HttpServer.Group("/api")
api.GET("/users", listUsers)
api.POST("/users", createUser)
api.GET("/health", healthCheck)

// Group-level 404 handler
api.SetNotFoundHandle(func(ctx dotweb.Context) error {
    return ctx.WriteString(`{"error": "API endpoint not found"}`)
})
```

### 4. Middleware
```go
app.HttpServer.Use(func(ctx dotweb.Context) error {
    // Before handler
    ctx.Items().Set("startTime", time.Now())
    
    err := ctx.NextHandler()  // Call next handler
    
    // After handler
    duration := time.Since(ctx.Items().Get("startTime").(time.Time))
    log.Printf("Request took %v", duration)
    
    return err
})
```

### 5. Session
```go
app.HttpServer.SetEnabledSession(true)
app.HttpServer.SetSessionConfig(session.NewDefaultRuntimeConfig())

app.HttpServer.GET("/login", func(ctx dotweb.Context) error {
    ctx.SetSession("user", "admin")
    return ctx.WriteString("Logged in!")
})
```

### 6. Data Binding
```go
type User struct {
    Name string `json:"name" form:"name"`
    Age  int    `json:"age" form:"age"`
}

app.HttpServer.POST("/users", func(ctx dotweb.Context) error {
    user := new(User)
    if err := ctx.Bind(user); err != nil {
        return err
    }
    return ctx.WriteString(fmt.Sprintf("Created: %s", user.Name))
})
```

### 7. JSON API
```go
app.HttpServer.GET("/api/users", func(ctx dotweb.Context) error {
    ctx.Response().Header().Set("Content-Type", "application/json")
    return ctx.WriteString(`{"users": ["Alice", "Bob"]}`)
})

// Or use WriteJsonC
app.HttpServer.GET("/api/user", func(ctx dotweb.Context) error {
    return ctx.WriteJsonC(200, map[string]string{
        "name": "Alice",
        "email": "alice@example.com",
    })
})
```

### 8. File Upload
```go
app.HttpServer.POST("/upload", func(ctx dotweb.Context) error {
    file, header, err := ctx.Request().FormFile("file")
    if err != nil {
        return err
    }
    defer file.Close()
    
    // Save file...
    return ctx.WriteString("Uploaded: " + header.Filename)
})
```

### 9. WebSocket
```go
app.HttpServer.GET("/ws", func(ctx dotweb.Context) error {
    if !ctx.IsWebSocket() {
        return ctx.WriteString("Requires WebSocket")
    }
    
    ws := ctx.WebSocket()
    
    for {
        msg, err := ws.ReadMessage()
        if err != nil {
            break
        }
        ws.SendMessage("Echo: " + msg)
    }
    
    return nil
})
```

### 10. Error Handling
```go
app.SetExceptionHandle(func(ctx dotweb.Context, err error) {
    ctx.Response().SetContentType(dotweb.MIMEApplicationJSONCharsetUTF8)
    ctx.WriteJsonC(500, map[string]string{"error": err.Error()})
})

app.SetNotFoundHandle(func(ctx dotweb.Context) {
    ctx.Response().SetContentType(dotweb.MIMEApplicationJSONCharsetUTF8)
    ctx.WriteJsonC(404, map[string]string{"error": "Not found"})
})
```

## Running Examples

```bash
# Run any example
cd example/session
go run main.go

# With hot reload (using air)
air
```

## Project Structure

For larger projects, consider this structure:

```
myapp/
├── main.go
├── config/
│   └── config.yaml
├── handlers/
│   ├── user.go
│   └── auth.go
├── middleware/
│   ├── auth.go
│   └── logger.go
├── models/
│   └── user.go
└── routes/
    └── routes.go
```

## Testing

```bash
# Run all tests
go test ./...

# With coverage
go test ./... -coverprofile=coverage.out
go tool cover -html=coverage.out
```

## Documentation

- [DotWeb GitHub](https://github.com/devfeel/dotweb)
- [API Documentation](https://pkg.go.dev/github.com/devfeel/dotweb)
- [Examples Repository](https://github.com/devfeel/dotweb-example)

## Support

- QQ Group: 193409346
- Gitter: [devfeel-dotweb](https://gitter.im/devfeel-dotweb)
