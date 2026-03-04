# Configuration Example

This example demonstrates how to configure DotWeb using different config file formats.

## Config Files

| File | Format | Description |
|------|--------|-------------|
| `dotweb.json` | JSON | JSON configuration |
| `dotweb.yaml` | YAML | YAML configuration |
| `dotweb.conf` | INI | INI-style configuration |
| `userconf.xml` | XML | XML configuration |

## Running

```bash
cd example/config
go run main.go
```

## Configuration Methods

### 1. Classic Mode (with config file)

```go
// Load config from directory
app := dotweb.Classic("/path/to/config")

// Or use current directory
app := dotweb.Classic(file.GetCurrentDirectory())
```

Classic mode automatically loads:
- `dotweb.json`
- `dotweb.yaml`
- `dotweb.conf`

### 2. Programmatic Configuration

```go
app := dotweb.New()

// Enable features
app.SetEnabledLog(true)
app.SetDevelopmentMode()

// Server configuration
app.HttpServer.SetEnabledSession(true)
app.HttpServer.SetEnabledGzip(true)
app.HttpServer.SetMaxBodySize(10 * 1024 * 1024) // 10MB
```

## Config File Structure

### JSON (`dotweb.json`)

```json
{
  "App": {
    "EnabledLog": true,
    "LogPath": "./logs"
  },
  "HttpServer": {
    "Port": 8080,
    "EnabledSession": true,
    "EnabledGzip": true,
    "MaxBodySize": 10485760
  }
}
```

### YAML (`dotweb.yaml`)

```yaml
App:
  EnabledLog: true
  LogPath: ./logs

HttpServer:
  Port: 8080
  EnabledSession: true
  EnabledGzip: true
  MaxBodySize: 10485760
```

### INI (`dotweb.conf`)

```ini
[App]
EnabledLog = true
LogPath = ./logs

[HttpServer]
Port = 8080
EnabledSession = true
EnabledGzip = true
MaxBodySize = 10485760
```

## Common Settings

| Setting | Method | Description |
|---------|--------|-------------|
| Log | `app.SetEnabledLog(true)` | Enable logging |
| Log Path | `app.SetLogPath("./logs")` | Log directory |
| Dev Mode | `app.SetDevelopmentMode()` | Development mode |
| Prod Mode | `app.SetProductionMode()` | Production mode |
| Session | `app.HttpServer.SetEnabledSession(true)` | Enable session |
| Gzip | `app.HttpServer.SetEnabledGzip(true)` | Enable gzip compression |
| Max Body | `app.HttpServer.SetMaxBodySize(bytes)` | Max request body size |

## Notes

- Config files are loaded in order: JSON → YAML → INI
- Programmatic config overrides file config
- Use `dotweb.Classic()` for quick setup with defaults
- Use `dotweb.New()` for full control
