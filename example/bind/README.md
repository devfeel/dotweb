# Data Binding Example

This example demonstrates how to bind request data to Go structs in DotWeb.

## Features

- Bind form data to struct
- Bind JSON body to struct
- Custom binder implementation
- Using JSON tags for binding

## Running

```bash
cd example/bind
go run main.go
```

## Testing

### Bind Form Data (POST)

```bash
# Using JSON tag for binding
curl -X POST http://localhost:8080/ \
  -d 'UserName=Alice&Sex=1'
# Output: TestBind [no error] &{Alice 1}
```

### Bind Query Parameters (GET)

```bash
curl "http://localhost:8080/getbind?user=Bob&sex=2"
# Output: GetBind [no error] &{Bob 2}
```

### Bind JSON Body (POST)

```bash
curl -X POST http://localhost:8080/jsonbind \
  -H "Content-Type: application/json" \
  -d '{"user":"Charlie","sex":1}'
# Output: PostBind [no error] &{Charlie 1}
```

## Binding Methods

### 1. Auto Bind (Form/JSON)

```go
type UserInfo struct {
    UserName string `json:"user" form:"user"`
    Sex      int    `json:"sex" form:"sex"`
}

func handler(ctx dotweb.Context) error {
    user := new(UserInfo)
    if err := ctx.Bind(user); err != nil {
        return err
    }
    // user.UserName, user.Sex are populated
    return nil
}
```

### 2. Bind JSON Body

```go
func handler(ctx dotweb.Context) error {
    user := new(UserInfo)
    if err := ctx.BindJsonBody(user); err != nil {
        return err
    }
    return nil
}
```

### 3. Custom Binder

```go
// Implement dotweb.Binder interface
type userBinder struct{}

func (b *userBinder) Bind(i interface{}, ctx dotweb.Context) error {
    // Custom binding logic
    return nil
}

func (b *userBinder) BindJsonBody(i interface{}, ctx dotweb.Context) error {
    // Custom JSON binding logic
    return nil
}

// Register custom binder
app.HttpServer.SetBinder(newUserBinder())
```

## Configuration

### Enable JSON Tag

```go
// Use JSON tags instead of form tags
app.HttpServer.SetEnabledBindUseJsonTag(true)
```

## API Reference

| Method | Description |
|--------|-------------|
| `ctx.Bind(struct)` | Auto bind from form/JSON |
| `ctx.BindJsonBody(struct)` | Bind from JSON body |
| `app.HttpServer.SetBinder(binder)` | Set custom binder |
| `app.HttpServer.SetEnabledBindUseJsonTag(bool)` | Use JSON tags |

## Notes

- Default tag name is `form`
- Enable `SetEnabledBindUseJsonTag(true)` to use JSON tags
- Custom binder allows implementing your own binding logic
- Supports JSON, XML, and form data content types
