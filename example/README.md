# DotWeb Examples

This directory contains examples demonstrating DotWeb features.

## Quick Start (5 minutes)

```bash
cd quickstart
go run main.go
# Visit http://localhost:8080
```

## Examples Index

| Example | Description | Complexity |
|---------|-------------|------------|
| [quickstart](./quickstart) | Minimal "Hello World" | ★☆☆ |
| [routing](./routing) | Route patterns, params, groups | ★★☆ |
| [middleware](./middleware) | Logging, auth, CORS | ★★☆ |
| [session](./session) | Session management | ★★☆ |
| [group](./group) | Route grouping with 404 handlers | ★★☆ |

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

## Running Examples

```bash
# Run any example
cd example/group
go run main.go

# With hot reload (using air)
air
```

## Common Patterns

### JSON API
```go
app.HttpServer.GET("/api/users", func(ctx dotweb.Context) error {
    ctx.Response().Header().Set("Content-Type", "application/json")
    return ctx.WriteString(`{"users": ["Alice", "Bob"]}`)
})
```

### Error Handling
```go
app.SetExceptionHandle(func(ctx dotweb.Context, err error) {
    ctx.Response().SetContentType(dotweb.MIMEApplicationJSONCharsetUTF8)
    ctx.WriteJsonC(500, map[string]string{"error": err.Error()})
})
```

### File Upload
```go
app.HttpServer.POST("/upload", func(ctx dotweb.Context) error {
    file := ctx.Request().FormFile("file")
    // Save file...
    return ctx.WriteString("Uploaded!")
})
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
