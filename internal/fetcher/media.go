package fetcher

import (
	"context"
	"crypto/sha256"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"os"
	"path/filepath"
	"strings"
	"time"
)

// MediaType represents the type of media file
type MediaType string

const (
	MediaTypeImage     MediaType = "image"
	MediaTypeVideo     MediaType = "video"
	MediaTypeAnimation MediaType = "animation"
	MediaTypeAudio     MediaType = "audio"
	MediaTypeUnknown   MediaType = "unknown"
)

// MediaFile represents a downloaded media file
type MediaFile struct {
	URL          string    `json:"url"`
	LocalPath    string    `json:"local_path"`
	Filename     string    `json:"filename"`
	MediaType    MediaType `json:"media_type"`
	ContentType  string    `json:"content_type"`
	Size         int64     `json:"size"`
	Checksum     string    `json:"checksum"`
	DownloadedAt time.Time `json:"downloaded_at"`
}

// MediaDownloader handles downloading and storing NFT media files
type MediaDownloader struct {
	client      *http.Client
	maxFileSize int64 // Maximum file size in bytes (default 100MB)
}

// NewMediaDownloader creates a new media downloader
func NewMediaDownloader() *MediaDownloader {
	return &MediaDownloader{
		client: &http.Client{
			Timeout: 60 * time.Second, // Longer timeout for media downloads
		},
		maxFileSize: 100 * 1024 * 1024, // 100MB default limit
	}
}

// DownloadMedia downloads media from a URL and stores it locally
func (md *MediaDownloader) DownloadMedia(ctx context.Context, mediaURL, targetDir string) (*MediaFile, error) {
	// Parse and validate URL
	parsedURL, err := url.Parse(mediaURL)
	if err != nil {
		return nil, fmt.Errorf("invalid media URL: %w", err)
	}

	// Create target directory
	if err := os.MkdirAll(targetDir, 0755); err != nil {
		return nil, fmt.Errorf("failed to create media directory: %w", err)
	}

	// Determine filename from URL
	filename := md.extractFilename(parsedURL)
	if filename == "" {
		filename = fmt.Sprintf("media_%d", time.Now().Unix())
	}

	// Create request
	req, err := http.NewRequestWithContext(ctx, "GET", mediaURL, nil)
	if err != nil {
		return nil, fmt.Errorf("failed to create request: %w", err)
	}

	// Add user agent to avoid blocking
	req.Header.Set("User-Agent", "SolVault/1.0 NFT-Backup-Tool")

	// Execute request
	resp, err := md.client.Do(req)
	if err != nil {
		return nil, fmt.Errorf("failed to download media: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("HTTP error %d downloading media", resp.StatusCode)
	}

	// Check content length
	if resp.ContentLength > md.maxFileSize {
		return nil, fmt.Errorf("file too large: %d bytes (max %d)", resp.ContentLength, md.maxFileSize)
	}

	// Determine media type and adjust filename if needed
	contentType := resp.Header.Get("Content-Type")
	mediaType := md.determineMediaType(contentType, filename)

	// Add extension if missing
	if !strings.Contains(filename, ".") {
		if ext := md.getExtensionForContentType(contentType); ext != "" {
			filename += ext
		}
	}

	localPath := filepath.Join(targetDir, filename)

	// Create file and download with size limit
	file, err := os.Create(localPath)
	if err != nil {
		return nil, fmt.Errorf("failed to create file %s: %w", localPath, err)
	}
	defer file.Close()

	// Use limited reader to prevent huge downloads
	limitedReader := &io.LimitedReader{
		R: resp.Body,
		N: md.maxFileSize,
	}

	// Copy with checksum calculation
	hash := sha256.New()
	multiWriter := io.MultiWriter(file, hash)

	bytesWritten, err := io.Copy(multiWriter, limitedReader)
	if err != nil {
		os.Remove(localPath) // Cleanup on error
		return nil, fmt.Errorf("failed to write media file: %w", err)
	}

	// Check if we hit the size limit
	if limitedReader.N == 0 && resp.ContentLength == -1 {
		os.Remove(localPath)
		return nil, fmt.Errorf("file too large: exceeded %d bytes", md.maxFileSize)
	}

	// Calculate final checksum
	checksum := fmt.Sprintf("%x", hash.Sum(nil))

	mediaFile := &MediaFile{
		URL:          mediaURL,
		LocalPath:    localPath,
		Filename:     filename,
		MediaType:    mediaType,
		ContentType:  contentType,
		Size:         bytesWritten,
		Checksum:     checksum,
		DownloadedAt: time.Now(),
	}

	return mediaFile, nil
}

// extractFilename extracts a filename from URL path
func (md *MediaDownloader) extractFilename(u *url.URL) string {
	path := u.Path
	if path == "" || path == "/" {
		return ""
	}

	// Get the last part of the path
	filename := filepath.Base(path)

	// Remove query parameters if they got included
	if idx := strings.Index(filename, "?"); idx != -1 {
		filename = filename[:idx]
	}

	return filename
}

// determineMediaType determines the media type from content type and filename
func (md *MediaDownloader) determineMediaType(contentType, filename string) MediaType {
	contentType = strings.ToLower(contentType)
	filename = strings.ToLower(filename)

	// Check content type first
	switch {
	case strings.HasPrefix(contentType, "image/"):
		return MediaTypeImage
	case strings.HasPrefix(contentType, "video/"):
		return MediaTypeVideo
	case strings.HasPrefix(contentType, "audio/"):
		return MediaTypeAudio
	case contentType == "application/octet-stream" && strings.Contains(filename, ".gif"):
		return MediaTypeAnimation
	}

	// Fallback to filename extension
	switch {
	case strings.HasSuffix(filename, ".jpg") || strings.HasSuffix(filename, ".jpeg") ||
		strings.HasSuffix(filename, ".png") || strings.HasSuffix(filename, ".gif") ||
		strings.HasSuffix(filename, ".webp") || strings.HasSuffix(filename, ".svg"):
		return MediaTypeImage
	case strings.HasSuffix(filename, ".mp4") || strings.HasSuffix(filename, ".webm") ||
		strings.HasSuffix(filename, ".mov") || strings.HasSuffix(filename, ".avi"):
		return MediaTypeVideo
	case strings.HasSuffix(filename, ".mp3") || strings.HasSuffix(filename, ".wav") ||
		strings.HasSuffix(filename, ".ogg"):
		return MediaTypeAudio
	}

	return MediaTypeUnknown
}

// getExtensionForContentType returns appropriate file extension for content type
func (md *MediaDownloader) getExtensionForContentType(contentType string) string {
	contentType = strings.ToLower(contentType)

	switch contentType {
	case "image/jpeg":
		return ".jpg"
	case "image/png":
		return ".png"
	case "image/gif":
		return ".gif"
	case "image/webp":
		return ".webp"
	case "image/svg+xml":
		return ".svg"
	case "video/mp4":
		return ".mp4"
	case "video/webm":
		return ".webm"
	case "audio/mpeg":
		return ".mp3"
	case "audio/wav":
		return ".wav"
	case "audio/ogg":
		return ".ogg"
	default:
		return ""
	}
}

// SetMaxFileSize sets the maximum allowed file size for downloads
func (md *MediaDownloader) SetMaxFileSize(maxSize int64) {
	md.maxFileSize = maxSize
}

// Close cleans up the downloader resources
func (md *MediaDownloader) Close() error {
	md.client.CloseIdleConnections()
	return nil
}
