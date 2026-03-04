# Mock Example

This example demonstrates how to use mock mode for testing in DotWeb.

## What is Mock Mode?

Mock mode allows you to intercept requests and return pre-defined responses, useful for:
- Development testing
- API prototyping
- Integration testing
- Offline development

## Running

```bash
cd example/mock
go run main.go
```

## Testing

```bash
# Without mock: returns actual handler response
# With mock: returns mock data

curl http://localhost:8080/
# Output: mock data
```

## Using Mock

### 1. Register String Response

```go
func AppMock() dotweb.Mock {
    m := dotweb.NewStandardMock()
    
    // Register mock for specific path
    m.RegisterString("/", "mock data")
    
    return m
}

// Apply mock
app.SetMock(AppMock())
```

### 2. Register JSON Response

```go
m.RegisterJson("/api/users", `{"users": ["Alice", "Bob"]}`)
```

### 3. Register File Response

```go
m.RegisterFile("/download", "./test.pdf")
```

### 4. Register Custom Handler

```go
m.RegisterHandler("/custom", func(ctx dotweb.Context) error {
    return ctx.WriteString("custom mock response")
})
```

## Mock Configuration

```go
// Enable mock mode
app.SetMock(AppMock())

// Mock responses are used instead of actual handlers
// when the path matches a registered mock
```

## Mock Types

| Method | Description |
|--------|-------------|
| `RegisterString(path, data)` | Return string |
| `RegisterJson(path, json)` | Return JSON |
| `RegisterFile(path, filepath)` | Return file |
| `RegisterHandler(path, handler)` | Custom handler |

## Testing Flow

```
Request → Mock Check → Mock Response (if registered)
                      → Actual Handler (if not registered)
```

## Use Cases

### 1. Development

Mock external API responses during development:

```go
m.RegisterJson("/api/weather", `{"temp": 25, "city": "Beijing"}`)
```

### 2. Testing

Mock database responses for unit tests:

```go
m.RegisterJson("/api/users/1", `{"id": 1, "name": "Test User"}`)
```

### 3. Prototyping

Define API responses before implementing:

```go
m.RegisterJson("/api/products", `[
    {"id": 1, "name": "Product A"},
    {"id": 2, "name": "Product B"}
]`)
```

## Notes

- Mock mode is for development/testing only
- Do not use in production
- Mock responses take precedence over actual handlers
- Useful for frontend development before backend is ready
