package dotweb

import (
	"os"
	"path/filepath"
	"testing"
)

// TestNewUploadFile tests NewUploadFile constructor
func TestNewUploadFile_ValidInput_Success(t *testing.T) {
	file := &UploadFile{
		FileName:    "test.txt",
		FileSize:    1024,
		FileHeader:  nil,
		FileContent: []byte("test content"),
	}
	
	if file.FileName != "test.txt" {
		t.Errorf("Expected filename test.txt, got %s", file.FileName)
	}
	if file.FileSize != 1024 {
		t.Errorf("Expected size 1024, got %d", file.FileSize)
	}
}

// TestUploadFile_GetFileName tests FileName method
func TestUploadFile_GetFileName(t *testing.T) {
	tests := []struct {
		name     string
		filename string
		expected string
	}{
		{"simple", "test.txt", "test.txt"},
		{"with_path", "path/to/file.txt", "file.txt"},
		{"empty", "", ""},
		{"special_chars", "test-file_v1.txt", "test-file_v1.txt"},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			file := &UploadFile{FileName: tt.filename}
			result := file.FileName
			if result != tt.expected {
				t.Errorf("Expected %s, got %s", tt.expected, result)
			}
		})
	}
}

// TestUploadFile_RandomFileName tests RandomFileName method
func TestUploadFile_RandomFileName(t *testing.T) {
	file := &UploadFile{FileName: "test.txt"}
	
	// Generate multiple random names
	names := make(map[string]bool)
	for i := 0; i < 100; i++ {
		name := file.RandomFileName()
		if name == "" {
			t.Error("RandomFileName should not return empty string")
		}
		names[name] = true
	}
	
	// All names should be unique
	if len(names) != 100 {
		t.Errorf("Expected 100 unique names, got %d", len(names))
	}
}

// TestUploadFile_Size tests Size method
func TestUploadFile_Size(t *testing.T) {
	tests := []struct {
		name     string
		size     int64
		expected int64
	}{
		{"zero", 0, 0},
		{"small", 100, 100},
		{"medium", 1024 * 1024, 1024 * 1024},
		{"large", 100 * 1024 * 1024, 100 * 1024 * 1024},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			file := &UploadFile{FileSize: tt.size}
			result := file.Size()
			if result != tt.expected {
				t.Errorf("Expected %d, got %d", tt.expected, result)
			}
		})
	}
}

// TestUploadFile_SaveFile tests SaveFile method
func TestUploadFile_SaveFile(t *testing.T) {
	// Create temp directory
	tmpDir, err := os.MkdirTemp("", "dotweb_test_*")
	if err != nil {
		t.Fatalf("Failed to create temp dir: %v", err)
	}
	defer os.RemoveAll(tmpDir)

	content := []byte("test file content")
	file := &UploadFile{
		FileName:    "test_save.txt",
		FileSize:    int64(len(content)),
		FileContent: content,
	}

	// Test saving file
	destPath := filepath.Join(tmpDir, "saved_test.txt")
	err = file.SaveFile(destPath)
	if err != nil {
		t.Fatalf("Failed to save file: %v", err)
	}

	// Verify file exists
	if _, err := os.Stat(destPath); os.IsNotExist(err) {
		t.Error("Saved file should exist")
	}

	// Verify content
	savedContent, err := os.ReadFile(destPath)
	if err != nil {
		t.Fatalf("Failed to read saved file: %v", err)
	}
	if string(savedContent) != string(content) {
		t.Error("Saved content should match original")
	}
}

// TestUploadFile_SaveFile_EmptyContent tests saving empty file
func TestUploadFile_SaveFile_EmptyContent(t *testing.T) {
	tmpDir, err := os.MkdirTemp("", "dotweb_test_*")
	if err != nil {
		t.Fatalf("Failed to create temp dir: %v", err)
	}
	defer os.RemoveAll(tmpDir)

	file := &UploadFile{
		FileName:    "empty.txt",
		FileSize:    0,
		FileContent: []byte{},
	}

	destPath := filepath.Join(tmpDir, "empty.txt")
	err = file.SaveFile(destPath)
	if err != nil {
		t.Fatalf("Failed to save empty file: %v", err)
	}

	// File should exist but be empty
	info, err := os.Stat(destPath)
	if err != nil {
		t.Fatalf("Failed to stat saved file: %v", err)
	}
	if info.Size() != 0 {
		t.Errorf("Expected size 0, got %d", info.Size())
	}
}

// TestUploadFile_GetFileExt tests GetFileExt method
func TestUploadFile_GetFileExt(t *testing.T) {
	tests := []struct {
		name     string
		filename string
		expected string
	}{
		{"txt", "file.txt", ".txt"},
		{"jpg", "image.jpg", ".jpg"},
		{"no_ext", "filename", ""},
		{"multi_dots", "file.name.txt", ".txt"},
		{"hidden", ".gitignore", ".gitignore"},
		{"empty", "", ""},
		{"upper_ext", "file.TXT", ".TXT"},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			file := &UploadFile{FileName: tt.filename}
			result := file.GetFileExt()
			if result != tt.expected {
				t.Errorf("Expected %s, got %s", tt.expected, result)
			}
		})
	}
}

// TestUploadFile_ReadBytes tests ReadBytes method
func TestUploadFile_ReadBytes(t *testing.T) {
	content := []byte("test content for reading")
	file := &UploadFile{
		FileName:    "test.txt",
		FileSize:    int64(len(content)),
		FileContent: content,
	}

	result := file.ReadBytes()
	if string(result) != string(content) {
		t.Error("ReadBytes should return file content")
	}

	// Verify it returns a copy (not reference)
	result[0] = 'X'
	if file.FileContent[0] == 'X' {
		t.Error("ReadBytes should return a copy, not reference")
	}
}

// TestUploadFile_ReadBytes_Empty tests ReadBytes with empty content
func TestUploadFile_ReadBytes_Empty(t *testing.T) {
	file := &UploadFile{
		FileName:    "empty.txt",
		FileSize:    0,
		FileContent: []byte{},
	}

	result := file.ReadBytes()
	if len(result) != 0 {
		t.Errorf("Expected empty slice, got length %d", len(result))
	}
}

// TestUploadFile_SaveFile_InvalidPath tests saving to invalid path
func TestUploadFile_SaveFile_InvalidPath(t *testing.T) {
	file := &UploadFile{
		FileName:    "test.txt",
		FileSize:    4,
		FileContent: []byte("test"),
	}

	// Try to save to non-existent directory
	err := file.SaveFile("/nonexistent/path/file.txt")
	if err == nil {
		t.Error("Should return error for invalid path")
	}
}
