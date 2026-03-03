# Group SetNotFoundHandle Example

This example demonstrates how to use `Group.SetNotFoundHandle` to set custom 404 handlers for router groups.

## Features

- **Group-level 404 handler**: Set custom 404 response for specific route groups
- **Priority**: Group-level handler takes priority over app-level handler
- **Flexible**: Different groups can have different 404 handlers

## Usage

```bash
# Run the example
go run main.go

# Test routes
curl http://localhost:8080/                    # Welcome page
curl http://localhost:8080/api/users          # API: Users list
curl http://localhost:8080/api/health         # API: Health check
curl http://localhost:8080/api/unknown        # API: 404 (group handler)
curl http://localhost:8080/web/index          # Web: Index page
curl http://localhost:8080/web/unknown        # Web: 404 (global handler)
curl http://localhost:8080/unknown            # Global: 404 (global handler)
```

## Expected Responses

### API Group (custom 404)
```bash
$ curl http://localhost:8080/api/unknown
{"code": 404, "message": "API 404 - Resource not found", "hint": "Check API documentation for available endpoints"}
```

### Web Group (uses global 404)
```bash
$ curl http://localhost:8080/web/unknown
{"code": 404, "message": "Global 404 - Page not found"}
```

### Global 404
```bash
$ curl http://localhost:8080/unknown
{"code": 404, "message": "Global 404 - Page not found"}
```

## Code Explanation

```go
// Set global 404 handler (fallback)
app.SetNotFoundHandle(func(ctx dotweb.Context) error {
    return ctx.WriteString(`{"code": 404, "message": "Global 404"}`)
})

// Create API group with custom 404 handler
apiGroup := app.HttpServer.Group("/api")
apiGroup.SetNotFoundHandle(func(ctx dotweb.Context) error {
    return ctx.WriteString(`{"code": 404, "message": "API 404"}`)
})

// Web group uses global 404 (no SetNotFoundHandle)
webGroup := app.HttpServer.Group("/web")
```

## Use Cases

1. **API vs Web**: Return JSON for API 404s, HTML for Web 404s
2. **Versioned APIs**: Different 404 messages for v1 vs v2 APIs
3. **Multi-tenant**: Custom 404 per tenant group
4. **Internationalization**: Different language 404 messages per group
