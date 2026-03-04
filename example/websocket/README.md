# WebSocket Example

This example demonstrates WebSocket support in DotWeb.

## Features

- WebSocket echo server
- Chat room with broadcast
- Connection management
- HTTP status endpoint

## Running

```bash
cd example/websocket
go run main.go
```

## Testing

### Using wscat (recommended)

Install wscat:
```bash
npm install -g wscat
```

### Echo Server

```bash
wscat -c ws://localhost:8080/ws

# Send message
> Hello, DotWeb!
< Echo: Hello, DotWeb!
```

### Chat Room

Open multiple terminals:

```bash
# Terminal 1
wscat -c 'ws://localhost:8080/chat?name=Alice'
> Hi everyone!
< 🔔 Alice joined the chat
< 💬 Alice: Hi everyone!

# Terminal 2
wscat -c 'ws://localhost:8080/chat?name=Bob'
< 🔔 Alice joined the chat
< 🔔 Bob joined the chat
< 💬 Alice: Hi everyone!
```

### Check Status

```bash
curl http://localhost:8080/status
# Output:
# WebSocket Server Status
# Connected clients: 2
# Endpoints:
#   - ws://localhost:8080/ws (Echo)
#   - ws://localhost:8080/chat?name=YourName (Chat)
```

### Using Browser

```html
<script>
const ws = new WebSocket('ws://localhost:8080/ws');

ws.onopen = () => {
    console.log('Connected');
    ws.send('Hello from browser!');
};

ws.onmessage = (event) => {
    console.log('Received:', event.data);
};

ws.onerror = (error) => {
    console.error('Error:', error);
};
</script>
```

## WebSocket API

### Check WebSocket Request

```go
func handler(ctx dotweb.Context) error {
    if !ctx.IsWebSocket() {
        return ctx.WriteString("Requires WebSocket")
    }
    
    ws := ctx.WebSocket()
    // ...
}
```

### Send Message

```go
ws := ctx.WebSocket()
err := ws.SendMessage("Hello, client!")
```

### Read Message

```go
ws := ctx.WebSocket()
msg, err := ws.ReadMessage()
if err != nil {
    // Client disconnected
    return err
}
```

### Get Underlying Connection

```go
ws := ctx.WebSocket()
conn := ws.Conn  // *websocket.Conn
req := ws.Request()  // *http.Request
```

## Common Patterns

### Echo Server

```go
app.HttpServer.GET("/ws", func(ctx dotweb.Context) error {
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

### Chat Room with Broadcast

```go
var clients = make(map[*websocket.Conn]bool)
var broadcast = make(chan string)

func handler(ctx dotweb.Context) error {
    ws := ctx.WebSocket()
    clients[ws.Conn] = true
    
    for {
        msg, err := ws.ReadMessage()
        if err != nil {
            delete(clients, ws.Conn)
            break
        }
        broadcast <- msg
    }
    
    return nil
}

func broadcaster() {
    for msg := range broadcast {
        for conn := range clients {
            websocket.Message.Send(conn, msg)
        }
    }
}
```

### Heartbeat/Ping

```go
func heartbeat(ws *dotweb.WebSocket) {
    ticker := time.NewTicker(30 * time.Second)
    defer ticker.Stop()
    
    for range ticker.C {
        if err := ws.SendMessage("ping"); err != nil {
            return
        }
    }
}

// Start in handler
go heartbeat(ctx.WebSocket())
```

## Notes

- WebSocket uses `golang.org/x/net/websocket` package
- Always check `ctx.IsWebSocket()` before using `ctx.WebSocket()`
- Handle connection errors (client disconnect)
- Use goroutines for concurrent message handling
- Consider adding heartbeat/ping for long connections
