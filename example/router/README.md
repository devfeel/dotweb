# Router Example

This example demonstrates advanced routing features in DotWeb.

## Features

- Basic route registration
- Auto HEAD method
- Method not allowed handler
- Path matching with parameters
- MatchPath helper

## Running

```bash
cd example/router
go run main.go
```

## Testing

```bash
# Basic GET
curl http://localhost:8080/
# Output: index - GET - /

# Path with parameter
curl http://localhost:8080/d/test/y
# Output: index - GET - /d/:x/y - true

# Path with trailing slash
curl http://localhost:8080/x/
# Output: index - GET - /x/

# POST request
curl -X POST http://localhost:8080/post
# Output: index - POST - /post

# Any method
curl -X POST http://localhost:8080/any
curl -X GET http://localhost:8080/any
# Output: any - [METHOD] - /any

# Raw http.HandlerFunc
curl http://localhost:8080/h/func
# Output: go raw http func
```

## Route Registration

### Using HttpServer

```go
app.HttpServer.GET("/", handler)
app.HttpServer.POST("/users", handler)
app.HttpServer.PUT("/users/:id", handler)
app.HttpServer.DELETE("/users/:id", handler)
app.HttpServer.Any("/any", handler)
```

### Using Router

```go
app.HttpServer.Router().GET("/", handler)
app.HttpServer.Router().POST("/users", handler)
```

### Register http.HandlerFunc

```go
func HandlerFunc(w http.ResponseWriter, r *http.Request) {
    w.Write([]byte("go raw http func"))
}

app.HttpServer.RegisterHandlerFunc("GET", "/h/func", HandlerFunc)
```

## Auto Methods

### Auto HEAD

```go
// Automatically handles HEAD requests for GET routes
app.HttpServer.SetEnabledAutoHEAD(true)
```

### Auto OPTIONS

```go
// Automatically handles OPTIONS requests for CORS
app.HttpServer.SetEnabledAutoOPTIONS(true)
```

## Method Not Allowed

```go
app.SetMethodNotAllowedHandle(func(ctx dotweb.Context) {
    ctx.Redirect(301, "/")
    // Or return custom error
    // ctx.WriteString("Method not allowed")
})
```

## Path Matching

```go
func handler(ctx dotweb.Context) error {
    // Get path pattern
    path := ctx.RouterNode().Path()
    // e.g., "/users/:id"
    
    // Check if path matches pattern
    matches := ctx.HttpServer().Router().MatchPath(ctx, "/d/:x/y")
    // returns true if current path matches
    
    return nil
}
```

## API Reference

| Method | Description |
|--------|-------------|
| `server.GET(path, handler)` | Register GET route |
| `server.POST(path, handler)` | Register POST route |
| `server.Any(path, handler)` | Match all methods |
| `server.RegisterHandlerFunc(method, path, handler)` | Register http.HandlerFunc |
| `server.SetEnabledAutoHEAD(bool)` | Auto handle HEAD |
| `server.SetEnabledAutoOPTIONS(bool)` | Auto handle OPTIONS |
| `ctx.RouterNode().Path()` | Get route pattern |
| `router.MatchPath(ctx, pattern)` | Check path match |

## Notes

- Routes are matched in order of registration
- More specific routes should be registered first
- Use `:param` for path parameters
- Use `SetMethodNotAllowedHandle()` for custom 405 responses
