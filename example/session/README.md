# Session Management Example

This example demonstrates how to use session management in DotWeb.

## Features

- Enable session middleware
- Set/Get session values
- Check session existence
- Destroy session (logout)

## Running

```bash
cd example/session
go run main.go
```

## Testing

```bash
# Login - set session
curl http://localhost:8080/login
# Output: ✅ Logged in as Alice (admin)

# Get user info from session
curl http://localhost:8080/user
# Output: 👤 User: Alice
#         🔑 Role: admin

# Check session exists
curl http://localhost:8080/check
# Output: ✅ Session exists

# Logout - destroy session
curl http://localhost:8080/logout
# Output: ✅ Logged out successfully

# Check session again
curl http://localhost:8080/check
# Output: ❌ No session found
```

## Session Configuration

### Runtime Mode (Default)

```go
app.HttpServer.SetSessionConfig(session.NewDefaultRuntimeConfig())
```

Session data is stored in memory. Good for development and single-instance deployment.

### Redis Mode

```go
// Without auth
app.HttpServer.SetSessionConfig(
    session.NewDefaultRedisConfig("redis://192.168.1.100:6379/0"),
)

// With auth
app.HttpServer.SetSessionConfig(
    session.NewDefaultRedisConfig("redis://:password@192.168.1.100:6379/0"),
)
```

Session data is stored in Redis. Recommended for production with multiple instances.

## API Reference

| Method | Description |
|--------|-------------|
| `ctx.SetSession(key, value)` | Set session value |
| `ctx.GetSession(key)` | Get session value (returns `interface{}`) |
| `ctx.HasSession(key)` | Check if session key exists |
| `ctx.DestorySession()` | Destroy current session |

## Notes

- Sessions are identified by a cookie named `DOTWEB_SESSION_ID` by default
- Session ID is automatically generated and managed by DotWeb
- Always enable session before using: `app.HttpServer.SetEnabledSession(true)`
