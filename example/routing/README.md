# Routing Example

This example demonstrates various routing patterns in DotWeb.

## Features

- HTTP methods (GET, POST, PUT, DELETE, ANY)
- Path parameters (`:id`, `:userId/:postId`)
- Wildcard routes (`*filepath`)
- Route groups (`/api`, `/api/v1`)

## Running

```bash
cd example/routing
go run main.go
```

## Testing

### Basic Routes

```bash
# GET request
curl http://localhost:8080/
# Output: GET / - Home page

# POST request
curl -X POST http://localhost:8080/users
# Output: POST /users - Create user

# PUT request
curl -X PUT http://localhost:8080/users/123
# Output: PUT /users/123 - Update user

# DELETE request
curl -X DELETE http://localhost:8080/users/123
# Output: DELETE /users/123 - Delete user

# Any method
curl -X POST http://localhost:8080/any
# Output: ANY /any - Method: POST
```

### Path Parameters

```bash
# Single parameter
curl http://localhost:8080/users/42
# Output: User ID: 42

# Multiple parameters
curl http://localhost:8080/users/42/posts/100
# Output: User: 42, Post: 100

# Wildcard (catch-all)
curl http://localhost:8080/files/path/to/file.txt
# Output: File path: /path/to/file.txt
```

### Route Groups

```bash
# API group
curl http://localhost:8080/api/health
# Output: {"status": "ok"}

curl http://localhost:8080/api/version
# Output: {"version": "1.0.0"}

# API v1 group
curl http://localhost:8080/api/v1/users
# Output: {"users": ["Alice", "Bob"]}

curl -X POST http://localhost:8080/api/v1/users
# Output: {"created": true}
```

## Routing Patterns

### 1. Named Parameters

Use `:name` to capture path segments:

```go
// /users/123 -> id = "123"
app.HttpServer.GET("/users/:id", handler)

// /users/42/posts/100 -> userId = "42", postId = "100"
app.HttpServer.GET("/users/:userId/posts/:postId", handler)
```

Get parameter value:

```go
id := ctx.GetRouterName("id")
```

### 2. Wildcard Routes

Use `*name` to capture everything after the prefix:

```go
// /files/path/to/file.txt -> filepath = "/path/to/file.txt"
app.HttpServer.GET("/files/*filepath", handler)
```

### 3. Route Groups

Organize routes with common prefix:

```go
// All routes under /api
api := app.HttpServer.Group("/api")
api.GET("/health", healthHandler)
api.GET("/users", listUsersHandler)

// Nested groups
v1 := app.HttpServer.Group("/api/v1")
v1.GET("/users", listUsersV1Handler)
```

### 4. Group-level Middleware

Apply middleware to a group:

```go
api := app.HttpServer.Group("/api")
api.Use(authMiddleware)  // Apply to all /api/* routes
api.GET("/users", listUsersHandler)
```

## API Reference

| Method | Description |
|--------|-------------|
| `app.HttpServer.GET(path, handler)` | Register GET route |
| `app.HttpServer.POST(path, handler)` | Register POST route |
| `app.HttpServer.PUT(path, handler)` | Register PUT route |
| `app.HttpServer.DELETE(path, handler)` | Register DELETE route |
| `app.HttpServer.ANY(path, handler)` | Match all HTTP methods |
| `app.HttpServer.Group(prefix)` | Create route group |
| `ctx.GetRouterName(name)` | Get path parameter value |

## Notes

- Parameters are extracted from the path and can be accessed via `ctx.GetRouterName()`
- Wildcard captures the rest of the URL including slashes
- Route groups can be nested
- Use `app.SetNotFoundHandle()` for custom 404 handling
