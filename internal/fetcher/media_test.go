package fetcher

import (
	"context"
	"net/http"
	"net/http/httptest"
	"os"
	"path/filepath"
	"testing"
	"time"
)

func TestMediaDownloader(t *testing.T) {
	// Create test server that serves a small image
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "image/png")
		// Send a minimal PNG (1x1 transparent pixel)
		pngData := []byte{
			0x89, 0x50, 0x4E, 0x47, 0x0D, 0x0A, 0x1A, 0x0A, 0x00, 0x00, 0x00, 0x0D,
			0x49, 0x48, 0x44, 0x52, 0x00, 0x00, 0x00, 0x01, 0x00, 0x00, 0x00, 0x01,
			0x08, 0x06, 0x00, 0x00, 0x00, 0x1F, 0x15, 0xC4, 0x89, 0x00, 0x00, 0x00,
			0x0A, 0x49, 0x44, 0x41, 0x54, 0x78, 0x9C, 0x63, 0x00, 0x01, 0x00, 0x00,
			0x05, 0x00, 0x01, 0x0D, 0x0A, 0x2D, 0xB4, 0x00, 0x00, 0x00, 0x00, 0x49,
			0x45, 0x4E, 0x44, 0xAE, 0x42, 0x60, 0x82,
		}
		w.Write(pngData)
	}))
	defer server.Close()

	// Create temporary directory for test
	tempDir, err := os.MkdirTemp("", "media_test")
	if err != nil {
		t.Fatalf("Failed to create temp dir: %v", err)
	}
	defer os.RemoveAll(tempDir)

	// Test media downloader
	downloader := NewMediaDownloader()
	defer downloader.Close()

	ctx := context.Background()
	mediaFile, err := downloader.DownloadMedia(ctx, server.URL+"/test.png", tempDir)
	if err != nil {
		t.Fatalf("Failed to download media: %v", err)
	}

	// Verify downloaded file
	if mediaFile.MediaType != MediaTypeImage {
		t.Errorf("Expected media type %s, got %s", MediaTypeImage, mediaFile.MediaType)
	}

	if mediaFile.ContentType != "image/png" {
		t.Errorf("Expected content type image/png, got %s", mediaFile.ContentType)
	}

	if mediaFile.Size == 0 {
		t.Error("Expected file size > 0")
	}

	if mediaFile.Checksum == "" {
		t.Error("Expected non-empty checksum")
	}

	// Verify file exists on disk
	if _, err := os.Stat(mediaFile.LocalPath); os.IsNotExist(err) {
		t.Error("Downloaded file does not exist on disk")
	}

	// Verify filename extraction
	expectedPath := filepath.Join(tempDir, "test.png")
	if mediaFile.LocalPath != expectedPath {
		t.Errorf("Expected path %s, got %s", expectedPath, mediaFile.LocalPath)
	}
}

func TestMediaDownloader_LargeFile(t *testing.T) {
	// Test that large files are rejected
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Length", "200000000") // 200MB
		w.Header().Set("Content-Type", "image/jpeg")
		w.WriteHeader(http.StatusOK)
	}))
	defer server.Close()

	tempDir, err := os.MkdirTemp("", "media_test_large")
	if err != nil {
		t.Fatalf("Failed to create temp dir: %v", err)
	}
	defer os.RemoveAll(tempDir)

	downloader := NewMediaDownloader()
	defer downloader.Close()

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	_, err = downloader.DownloadMedia(ctx, server.URL+"/large.jpg", tempDir)
	if err == nil {
		t.Error("Expected error for large file, but got none")
	}
}

func TestMediaType_Detection(t *testing.T) {
	downloader := NewMediaDownloader()

	tests := []struct {
		contentType string
		filename    string
		expected    MediaType
	}{
		{"image/jpeg", "test.jpg", MediaTypeImage},
		{"video/mp4", "test.mp4", MediaTypeVideo},
		{"audio/mpeg", "test.mp3", MediaTypeAudio},
		{"application/octet-stream", "test.png", MediaTypeImage},
		{"text/plain", "test.txt", MediaTypeUnknown},
	}

	for _, test := range tests {
		result := downloader.determineMediaType(test.contentType, test.filename)
		if result != test.expected {
			t.Errorf("For content-type %s and filename %s, expected %s but got %s",
				test.contentType, test.filename, test.expected, result)
		}
	}
}
