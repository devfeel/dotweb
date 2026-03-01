package main

import (
	"fmt"
	"os"
	"time"

	"github.com/devfeel/dotweb"
)

// File Upload Example
// This demonstrates how to handle file uploads with dotweb

func main() {
	// init DotApp
	app := dotweb.New()

	// Set upload directory
	uploadDir := "./uploads"
	if err := os.MkdirAll(uploadDir, 0755); err != nil {
		fmt.Println("Error creating upload dir:", err)
		return
	}

	// Show upload form
	app.HttpServer.GET("/upload", func(ctx dotweb.Context) error {
		return ctx.WriteString(`
<!DOCTYPE html>
<html>
<head><title>File Upload</title></head>
<body>
<form enctype="multipart/form-data" method="post" action="/upload">
<input type="file" name="file"><br>
<input type="submit" value="Upload">
</form>
</body>
</html>
		`)
	})

	// Handle file upload
	app.HttpServer.POST("/upload", func(ctx dotweb.Context) error {
		file, err := ctx.Request().FormFile("file")
		if err != nil {
			return ctx.WriteString("Error uploading file: " + err.Error())
		}

		// Generate unique filename
		filename := fmt.Sprintf("file_%d%s", time.Now().Unix(), file.GetFileExt())
		serverPath := fmt.Sprintf("%s/%s", uploadDir, filename)

		// Save file
		size, err := file.SaveFile(serverPath)
		if err != nil {
			return ctx.WriteString("Error saving file: " + err.Error())
		}

		return ctx.WriteJson(map[string]interface{}{
			"message":    "File uploaded successfully",
			"filename":   filename,
			"size":       size,
			"contentType": file.Header.Header.Get("Content-Type"),
		})
	})

	// List uploaded files
	app.HttpServer.GET("/files", func(ctx dotweb.Context) error {
		files, err := os.ReadDir(uploadDir)
		if err != nil {
			return ctx.WriteString("Error reading files")
		}

		var fileList []string
		for _, f := range files {
			if !f.IsDir() {
				fileList = append(fileList, f.Name())
			}
		}
		return ctx.WriteJson(fileList)
	})

	fmt.Println("File Upload Example server starting on :8080")
	fmt.Println("Upload form: GET /upload")
	fmt.Println("Upload handler: POST /upload")
	fmt.Println("List files: GET /files")
	err := app.StartServer(8080)
	fmt.Println("server error => ", err)
}
