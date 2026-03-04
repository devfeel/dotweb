# File Upload Example

This example demonstrates file upload and download in DotWeb.

## Features

- Single file upload
- Multiple file upload
- File download
- List uploaded files
- Delete files
- File size limits

## Running

```bash
cd example/file-upload
go run main.go
```

## Testing

### Upload Single File

```bash
# Create a test file
echo "Hello, DotWeb!" > test.txt

# Upload file
curl -F 'file=@test.txt' http://localhost:8080/upload
# Output:
# ✅ File uploaded!
# 📁 Name: test.txt
# 📊 Size: 14 bytes
# 📍 Path: ./uploads/test.txt
```

### Upload Multiple Files

```bash
curl -F 'files=@file1.txt' -F 'files=@file2.txt' http://localhost:8080/upload/multiple
# Output:
# Uploaded 2 files:
# [✅ file1.txt ✅ file2.txt]
```

### List Files

```bash
curl http://localhost:8080/files
# Output:
# 📂 Uploaded files:
# 📁 test.txt (14 bytes)
```

### Download File

```bash
curl http://localhost:8080/download/test.txt -o downloaded.txt
```

### Delete File

```bash
curl -X DELETE http://localhost:8080/files/test.txt
# Output: ✅ File deleted: test.txt
```

## API Reference

### Upload File

```go
// Get single file
file, header, err := ctx.Request().FormFile("file")

// Get file content
data, err := io.ReadAll(file)

// Get file name
filename := header.Filename
```

### Upload Multiple Files

```go
// Parse multipart form
err := ctx.Request().ParseMultipartForm(32 << 20) // 32MB

// Get all files
files := ctx.Request().MultipartForm.File["files"]
```

### Download File

```go
// Set headers for download
ctx.Response().Header().Set("Content-Disposition", "attachment; filename="+filename)
ctx.Response().Header().Set("Content-Type", "application/octet-stream")

// Send file data
data, _ := os.ReadFile(filePath)
ctx.Write(200, data)
```

## Configuration

### Set Max Body Size

```go
// 10MB limit
app.HttpServer.SetMaxBodySize(10 * 1024 * 1024)

// Unlimited
app.HttpServer.SetMaxBodySize(-1)
```

## File Upload Helper

DotWeb provides a built-in upload file helper:

```go
// Using UploadFile helper
uploadFile := ctx.Request().UploadFile("file")
if uploadFile != nil {
    filename := uploadFile.Filename
    data := uploadFile.Data  // []byte
    size := len(data)
}
```

## Common Patterns

### Validate File Type

```go
func isValidFileType(filename string) bool {
    ext := strings.ToLower(filepath.Ext(filename))
    allowed := []string{".jpg", ".jpeg", ".png", ".gif", ".pdf"}
    for _, a := range allowed {
        if ext == a {
            return true
        }
    }
    return false
}
```

### Generate Unique Filename

```go
import "github.com/google/uuid"

func uniqueFilename(filename string) string {
    ext := filepath.Ext(filename)
    return uuid.New().String() + ext
}
```

### Check File Size

```go
func checkFileSize(size int64, maxSize int64) bool {
    return size <= maxSize
}
```

## Notes

- Always validate uploaded files
- Set appropriate max body size
- Use unique filenames to avoid conflicts
- Check file types for security
- Clean up old files periodically
