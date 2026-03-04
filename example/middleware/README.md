# Middleware Example

This example demonstrates how to use middleware in DotWeb.

## Features

- Global middleware
- Route-level middleware
- Group-level middleware
- Exclude specific routes
- Custom middleware implementation

## Running

```bash
cd example/middleware
go run main.go
```

## Middleware Types

### 1. Global Middleware

Applied to all routes:

```go
app.Use(NewAccessFmtLog("app"))
```

### 2. Route-level Middleware

Applied to specific routes:

```go
server.Router().GET("/use", Index).Use(NewAccessFmtLog("Router-use"))
```

### 3. Group-level Middleware

Applied to all routes in a group:

```go
g := server.Group("/api").Use(NewAuthMiddleware("secret"))
g.GET("/users", listUsers)
```

### 4. Exclude Routes

Skip middleware for specific routes:

```go
middleware := NewAccessFmtLog("appex")
middleware.Exclude("/index")
middleware.Exclude("/v1/machines/queryIP/:IP")
app.Use(middleware)

// Or exclude from first middleware
app.ExcludeUse(NewAccessFmtLog("appex1"), "/")
```

## Custom Middleware

```go
func NewAccessFmtLog(name string) dotweb.HandlerFunc {
    return func(ctx dotweb.Context) error {
        // Before handler
        start := time.Now()
        log.Printf("[%s] %s %s", name, ctx.Request().Method, ctx.Request().Url())
        
        // Call next handler
        err := ctx.NextHandler()
        
        // After handler
        duration := time.Since(start)
        log.Printf("[%s] Request took %v", name, duration)
        
        return err
    }
}
```

## Testing

```bash
# All routes go through middleware
curl http://localhost:8080/
curl http://localhost:8080/index

# Check middleware chain in group routes
curl http://localhost:8080/A/
curl http://localhost:8080/A/B/
curl http://localhost:8080/A/C/
```

## Middleware Chain

When using groups with middleware, they form a chain:

```
Global Middleware
    ↓
Group Middleware (A)
    ↓
Group Middleware (B)
    ↓
Route Middleware
    ↓
Handler
```

Use `ctx.RouterNode().GroupMiddlewares()` to inspect the chain.

## API Reference

| Method | Description |
|--------|-------------|
| `app.Use(middleware...)` | Add global middleware |
| `route.Use(middleware)` | Add route-level middleware |
| `group.Use(middleware)` | Add group-level middleware |
| `app.ExcludeUse(middleware, path)` | Exclude path from middleware |
| `middleware.Exclude(path)` | Exclude path from middleware |
| `ctx.NextHandler()` | Call next middleware/handler |
| `ctx.RouterNode().GroupMiddlewares()` | Get middleware chain |

## Notes

- Middleware is executed in the order they are added
- Call `ctx.NextHandler()` to pass control to the next middleware
- Without `ctx.NextHandler()`, the middleware chain stops
- Use `Exclude()` to skip middleware for certain routes
