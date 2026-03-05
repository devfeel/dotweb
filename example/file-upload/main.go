// Package main demonstrates file upload and download in DotWeb.
// Run: go run main.go
// Test: curl -F "file=@test.txt" http://localhost:8080/upload
package main

import (
	"fmt"
	"os"
	"path/filepath"

	"github.com/devfeel/dotweb"
)

func main() {
	// Create DotWeb app
	app := dotweb.New()
	app.SetDevelopmentMode()
	
	// Set max body size (10MB)
	app.HttpServer.SetMaxBodySize(10 * 1024 * 1024)
	
	// Upload single file
	app.HttpServer.POST("/upload", func(ctx dotweb.Context) error {
		// Get uploaded file
		file, err := ctx.Request().FormFile("file")
		if err != nil {
			return ctx.WriteString("❌ Error getting file: " + err.Error())
		}
		
		// Create upload directory
		uploadDir := "./uploads"
		os.MkdirAll(uploadDir, 0755)
		
		// Save file using built-in method
		dst := filepath.Join(uploadDir, file.FileName())
		size, err := file.SaveFile(dst)
		if err != nil {
			return ctx.WriteString("❌ Error saving file: " + err.Error())
		}
		
		return ctx.WriteString(fmt.Sprintf(
			"✅ File uploaded!\n📁 Name: %s\n📊 Size: %d bytes\n📍 Path: %s",
			file.FileName(), size, dst,
		))
	})
	
	// Upload multiple files
	app.HttpServer.POST("/upload/multiple", func(ctx dotweb.Context) error {
		// Get all files
		files, err := ctx.Request().FormFiles()
		if err != nil {
			return ctx.WriteString("❌ Error parsing form: " + err.Error())
		}
		
		uploadDir := "./uploads"
		os.MkdirAll(uploadDir, 0755)
		
		var results []string
		for name, file := range files {
			dst := filepath.Join(uploadDir, file.FileName())
			_, err := file.SaveFile(dst)
			if err != nil {
				results = append(results, fmt.Sprintf("❌ %s: failed to save", name))
				continue
			}
			
			results = append(results, fmt.Sprintf("✅ %s (%s)", name, file.FileName()))
		}
		
		return ctx.WriteString(fmt.Sprintf("Uploaded %d files:\n%s", len(files), 
			fmt.Sprintf("%v", results)))
	})
	
	// Download file
	app.HttpServer.GET("/download/:filename", func(ctx dotweb.Context) error {
		filename := ctx.GetRouterName("filename")
		uploadDir := "./uploads"
		filePath := filepath.Join(uploadDir, filename)
		
		// Check file exists
		if _, err := os.Stat(filePath); os.IsNotExist(err) {
			return ctx.WriteString("❌ File not found: " + filename)
		}
		
		// Set response headers
		ctx.Response().Header().Set("Content-Disposition", "attachment; filename="+filename)
		ctx.Response().Header().Set("Content-Type", "application/octet-stream")
		
		// Read and send file
		data, err := os.ReadFile(filePath)
		if err != nil {
			return ctx.WriteString("❌ Error reading file: " + err.Error())
		}
		
		ctx.Write(200, data)
		return nil
	})
	
	// List uploaded files
	app.HttpServer.GET("/files", func(ctx dotweb.Context) error {
		uploadDir := "./uploads"
		files, err := os.ReadDir(uploadDir)
		if err != nil {
			ctx.WriteString("❌ Error reading directory: " + err.Error())
			return nil
		}
		
		var result string
		for _, file := range files {
			info, _ := file.Info()
			result += fmt.Sprintf("📁 %s (%d bytes)\n", file.Name(), info.Size())
		}
		
		if result == "" {
			ctx.WriteString("📂 No files uploaded yet")
			return nil
		}
		ctx.WriteString("📂 Uploaded files:\n" + result)
		return nil
	})
	
	// Delete file
	app.HttpServer.DELETE("/files/:filename", func(ctx dotweb.Context) error {
		filename := ctx.GetRouterName("filename")
		uploadDir := "./uploads"
		filePath := filepath.Join(uploadDir, filename)
		
		if err := os.Remove(filePath); err != nil {
			ctx.WriteString("❌ Error deleting file: " + err.Error())
			return nil
		}
		ctx.WriteString("✅ File deleted: " + filename)
		return nil
	})
	
	fmt.Println("🚀 File upload example running at http://localhost:8080")
	fmt.Println("\nTest routes:")
	fmt.Println("  # Upload single file")
	fmt.Println("  curl -F 'file=@test.txt' http://localhost:8080/upload")
	fmt.Println("")
	fmt.Println("  # Upload multiple files")
	fmt.Println("  curl -F 'files=@file1.txt' -F 'files=@file2.txt' http://localhost:8080/upload/multiple")
	fmt.Println("")
	fmt.Println("  # List files")
	fmt.Println("  curl http://localhost:8080/files")
	fmt.Println("")
	fmt.Println("  # Download file")
	fmt.Println("  curl http://localhost:8080/download/test.txt -o downloaded.txt")
	fmt.Println("")
	fmt.Println("  # Delete file")
	fmt.Println("  curl -X DELETE http://localhost:8080/files/test.txt")
	
	if err := app.StartServer(8080); err != nil {
		fmt.Printf("Server error: %v\n", err)
	}
}
