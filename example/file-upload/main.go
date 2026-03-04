// Package main demonstrates file upload and download in DotWeb.
// Run: go run main.go
// Test: curl -F "file=@test.txt" http://localhost:8080/upload
package main

import (
	"fmt"
	"io"
	"os"
	"path/filepath"
	"strconv"

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
		file, header, err := ctx.Request().FormFile("file")
		if err != nil {
			return ctx.WriteString("❌ Error getting file: " + err.Error())
		}
		defer file.Close()
		
		// Create upload directory
		uploadDir := "./uploads"
		os.MkdirAll(uploadDir, 0755)
		
		// Create destination file
		dst := filepath.Join(uploadDir, header.Filename)
		dstFile, err := os.Create(dst)
		if err != nil {
			return ctx.WriteString("❌ Error creating file: " + err.Error())
		}
		defer dstFile.Close()
		
		// Copy file content
		written, err := io.Copy(dstFile, file)
		if err != nil {
			return ctx.WriteString("❌ Error saving file: " + err.Error())
		}
		
		return ctx.WriteString(fmt.Sprintf(
			"✅ File uploaded!\n📁 Name: %s\n📊 Size: %d bytes\n📍 Path: %s",
			header.Filename, written, dst,
		))
	})
	
	// Upload multiple files
	app.HttpServer.POST("/upload/multiple", func(ctx dotweb.Context) error {
		// Parse multipart form
		err := ctx.Request().ParseMultipartForm(32 << 20) // 32MB
		if err != nil {
			return ctx.WriteString("❌ Error parsing form: " + err.Error())
		}
		
		// Get all files
		form := ctx.Request().MultipartForm
		files := form.File["files"]
		
		uploadDir := "./uploads"
		os.MkdirAll(uploadDir, 0755)
		
		var results []string
		for _, fileHeader := range files {
			file, err := fileHeader.Open()
			if err != nil {
				results = append(results, fmt.Sprintf("❌ %s: failed to open", fileHeader.Filename))
				continue
			}
			
			dst := filepath.Join(uploadDir, fileHeader.Filename)
			dstFile, err := os.Create(dst)
			if err != nil {
				file.Close()
				results = append(results, fmt.Sprintf("❌ %s: failed to create", fileHeader.Filename))
				continue
			}
			
			io.Copy(dstFile, file)
			dstFile.Close()
			file.Close()
			
			results = append(results, fmt.Sprintf("✅ %s", fileHeader.Filename))
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
		
		return ctx.Write(200, data)
	})
	
	// List uploaded files
	app.HttpServer.GET("/files", func(ctx dotweb.Context) error {
		uploadDir := "./uploads"
		files, err := os.ReadDir(uploadDir)
		if err != nil {
			return ctx.WriteString("❌ Error reading directory: " + err.Error())
		}
		
		var result string
		for _, file := range files {
			info, _ := file.Info()
			result += fmt.Sprintf("📁 %s (%d bytes)\n", file.Name(), info.Size())
		}
		
		if result == "" {
			return ctx.WriteString("📂 No files uploaded yet")
		}
		return ctx.WriteString("📂 Uploaded files:\n" + result)
	})
	
	// Delete file
	app.HttpServer.DELETE("/files/:filename", func(ctx dotweb.Context) error {
		filename := ctx.GetRouterName("filename")
		uploadDir := "./uploads"
		filePath := filepath.Join(uploadDir, filename)
		
		if err := os.Remove(filePath); err != nil {
			return ctx.WriteString("❌ Error deleting file: " + err.Error())
		}
		return ctx.WriteString("✅ File deleted: " + filename)
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
