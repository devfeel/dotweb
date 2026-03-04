# RESTful JSON API Example

This example demonstrates how to build a RESTful JSON API with DotWeb.

## Features

- RESTful CRUD operations
- JSON request/response handling
- Error handling
- Global middleware
- API versioning (groups)
- Concurrent-safe data storage

## Running

```bash
cd example/json-api
go run main.go
```

## API Endpoints

| Method | Endpoint | Description |
|--------|----------|-------------|
| GET | /api/health | Health check |
| GET | /api/users | List all users |
| GET | /api/users/:id | Get user by ID |
| POST | /api/users | Create user |
| PUT | /api/users/:id | Update user |
| DELETE | /api/users/:id | Delete user |

## Testing

### Health Check

```bash
curl http://localhost:8080/api/health
# Output: {"status":"ok"}
```

### List Users

```bash
curl http://localhost:8080/api/users
# Output: {"message":"success","data":[{"id":1,"name":"Alice","email":"alice@example.com"}...]}
```

### Get User

```bash
curl http://localhost:8080/api/users/1
# Output: {"message":"success","data":{"id":1,"name":"Alice","email":"alice@example.com"}}
```

### Create User

```bash
curl -X POST http://localhost:8080/api/users \
  -H "Content-Type: application/json" \
  -d '{"name":"Charlie","email":"charlie@example.com"}'
# Output: {"message":"User created","data":{"id":3,"name":"Charlie","email":"charlie@example.com"}}
```

### Update User

```bash
curl -X PUT http://localhost:8080/api/users/1 \
  -H "Content-Type: application/json" \
  -d '{"name":"Alice Updated"}'
# Output: {"message":"User updated","data":{"id":1,"name":"Alice Updated","email":"alice@example.com"}}
```

### Delete User

```bash
curl -X DELETE http://localhost:8080/api/users/1
# Output: {"message":"User deleted"}
```

### Error Responses

```bash
# Invalid ID
curl http://localhost:8080/api/users/abc
# Output: {"error":"Invalid user ID"}

# User not found
curl http://localhost:8080/api/users/999
# Output: {"error":"User not found"}

# Missing fields
curl -X POST http://localhost:8080/api/users \
  -H "Content-Type: application/json" \
  -d '{"name":""}'
# Output: {"error":"Name and email required"}
```

## Code Structure

### JSON Response Helper

```go
// Success response
return ctx.WriteJsonC(200, SuccessResponse{
    Message: "success",
    Data:    user,
})

// Error response
return ctx.WriteJsonC(404, ErrorResponse{
    Error: "User not found",
})
```

### JSON Request Parsing

```go
var user User
if err := json.Unmarshal(ctx.Request().PostBody(), &user); err != nil {
    return ctx.WriteJsonC(400, ErrorResponse{Error: "Invalid JSON"})
}
```

### Global Middleware

```go
// Set JSON content type for all responses
app.HttpServer.Use(func(ctx dotweb.Context) error {
    ctx.Response().Header().Set("Content-Type", "application/json")
    return ctx.NextHandler()
})
```

### Error Handling

```go
// Global exception handler
app.SetExceptionHandle(func(ctx dotweb.Context, err error) {
    ctx.Response().SetContentType(dotweb.MIMEApplicationJSONCharsetUTF8)
    ctx.WriteJsonC(500, ErrorResponse{Error: err.Error()})
})

// 404 handler
app.SetNotFoundHandle(func(ctx dotweb.Context) {
    ctx.Response().SetContentType(dotweb.MIMEApplicationJSONCharsetUTF8)
    ctx.WriteJsonC(404, ErrorResponse{Error: "Not found"})
})
```

## RESTful Best Practices

### 1. Use Proper HTTP Methods

```go
GET    /api/users     // List
GET    /api/users/:id // Get
POST   /api/users     // Create
PUT    /api/users/:id // Update
DELETE /api/users/:id // Delete
```

### 2. Use Appropriate Status Codes

```go
200 // OK - Successful GET, PUT, DELETE
201 // Created - Successful POST
400 // Bad Request - Invalid input
404 // Not Found - Resource doesn't exist
500 // Internal Server Error - Server error
```

### 3. Use Consistent Response Format

```go
// Success
{
    "message": "success",
    "data": { ... }
}

// Error
{
    "error": "Error message"
}
```

### 4. Use Route Groups

```go
api := app.HttpServer.Group("/api")
api.GET("/users", listUsers)
api.GET("/users/:id", getUser)
```

## Extending

### Add Authentication

```go
api.Use(func(ctx dotweb.Context) error {
    token := ctx.Request().Header.Get("Authorization")
    if token == "" {
        return ctx.WriteJsonC(401, ErrorResponse{Error: "Unauthorized"})
    }
    return ctx.NextHandler()
})
```

### Add Pagination

```go
func listUsers(ctx dotweb.Context) error {
    page := ctx.QueryValue("page")
    limit := ctx.QueryValue("limit")
    // Implement pagination...
}
```

### Add Validation

```go
func validateUser(user User) error {
    if user.Name == "" {
        return errors.New("name required")
    }
    if !strings.Contains(user.Email, "@") {
        return errors.New("invalid email")
    }
    return nil
}
```

## Notes

- This example uses in-memory storage for simplicity
- For production, use a database (MySQL, PostgreSQL, MongoDB, etc.)
- Always validate input data
- Use proper authentication for production APIs
- Consider rate limiting for public APIs
