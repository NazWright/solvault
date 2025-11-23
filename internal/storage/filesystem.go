package storage

import (
	"context"
	"crypto/sha256"
	"encoding/json"
	"fmt"
	"io/fs"
	"os"
	"path/filepath"
	"time"

	"github.com/NazWright/solvault/internal/fetcher"
	solanago "github.com/gagliardetto/solana-go"
)

// FileStorage implements StorageBackend using local filesystem
//
// Directory structure:
// backup_dir/
//
//	└── wallets/
//	    └── {wallet_address}/
//	        └── nfts/
//	            └── {mint_address}/
//	                ├── nft_data.json     (StoredNFT struct)
//	                ├── metadata.json     (off-chain metadata)
//	                └── media/            (images, videos, etc.)
type FileStorage struct {
	baseDir     string      // Root directory for all backups
	permissions fs.FileMode // File permissions for created files
}

// NewFileStorage creates a new file-based storage backend
func NewFileStorage(baseDir string) (*FileStorage, error) {
	// Explanation: We create the base directory structure upfront
	// This ensures we have write permissions and the path exists
	if err := os.MkdirAll(baseDir, 0755); err != nil {
		return nil, fmt.Errorf("failed to create base directory %s: %w", baseDir, err)
	}

	return &FileStorage{
		baseDir:     baseDir,
		permissions: 0644, // Read/write for owner, read for others
	}, nil
}

// SaveNFT stores NFT information to the filesystem
func (fs *FileStorage) SaveNFT(ctx context.Context, nftInfo *fetcher.NFTInfo) error {
	// Explanation: We build a path that's organized and human-readable
	// wallet/nfts/mint/ structure makes it easy to browse backups
	nftDir := fs.buildNFTPath(nftInfo.Owner, nftInfo.MintAddress)

	// Create directory structure
	if err := os.MkdirAll(nftDir, 0755); err != nil {
		return fmt.Errorf("failed to create NFT directory %s: %w", nftDir, err)
	}

	// Create stored NFT with metadata
	storedNFT := &StoredNFT{
		NFTInfo:    nftInfo,
		StoredAt:   time.Now(),
		UpdatedAt:  time.Now(),
		Version:    1, // Start with version 1
		BackupPath: nftDir,
		Verified:   false,       // Will be verified later
		LastCheck:  time.Time{}, // Not checked yet
	}

	// Calculate checksum for data integrity
	// Explanation: This helps us detect if files get corrupted
	checksum, err := fs.calculateChecksum(nftInfo)
	if err != nil {
		return fmt.Errorf("failed to calculate checksum: %w", err)
	}
	storedNFT.Checksum = checksum

	// Save main NFT data
	nftDataPath := filepath.Join(nftDir, "nft_data.json")
	if err := fs.saveJSON(nftDataPath, storedNFT); err != nil {
		return fmt.Errorf("failed to save NFT data: %w", err)
	}

	// Save metadata separately if available
	// Explanation: Separate files make it easier to examine metadata
	if nftInfo.Metadata != nil {
		metadataPath := filepath.Join(nftDir, "metadata.json")
		if err := fs.saveJSON(metadataPath, nftInfo.Metadata); err != nil {
			return fmt.Errorf("failed to save metadata: %w", err)
		}
	}

	// Create media directory and save media file info if available
	if len(nftInfo.MediaFiles) > 0 {
		mediaDir := filepath.Join(nftDir, "media")
		if err := os.MkdirAll(mediaDir, 0755); err != nil {
			return fmt.Errorf("failed to create media directory: %w", err)
		}

		// Save media manifest for tracking downloaded files
		mediaManifestPath := filepath.Join(nftDir, "media_manifest.json")
		if err := fs.saveJSON(mediaManifestPath, nftInfo.MediaFiles); err != nil {
			return fmt.Errorf("failed to save media manifest: %w", err)
		}
	}

	return nil
}

// GetNFT retrieves stored NFT information
func (fs *FileStorage) GetNFT(ctx context.Context, walletAddr, mintAddr solanago.PublicKey) (*StoredNFT, error) {
	nftDataPath := filepath.Join(fs.buildNFTPath(walletAddr, mintAddr), "nft_data.json")

	var storedNFT StoredNFT
	if err := fs.loadJSON(nftDataPath, &storedNFT); err != nil {
		if os.IsNotExist(err) {
			return nil, fmt.Errorf("NFT not found: %s", mintAddr.String())
		}
		return nil, fmt.Errorf("failed to load NFT data: %w", err)
	}

	return &storedNFT, nil
}

// ListNFTs returns all NFTs for a wallet
func (fs *FileStorage) ListNFTs(ctx context.Context, walletAddr solanago.PublicKey) ([]*StoredNFT, error) {
	walletDir := filepath.Join(fs.baseDir, "wallets", walletAddr.String(), "nfts")

	// Check if wallet directory exists
	if _, err := os.Stat(walletDir); os.IsNotExist(err) {
		return []*StoredNFT{}, nil // Empty slice, not an error
	}

	var nfts []*StoredNFT

	// Walk through all NFT directories
	// Explanation: We scan subdirectories to find all stored NFTs
	err := filepath.Walk(walletDir, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return err
		}

		// Look for nft_data.json files
		if info.Name() == "nft_data.json" {
			var storedNFT StoredNFT
			if loadErr := fs.loadJSON(path, &storedNFT); loadErr != nil {
				// Log error but continue with other NFTs
				fmt.Printf("⚠️  Warning: failed to load %s: %v\n", path, loadErr)
				return nil
			}
			nfts = append(nfts, &storedNFT)
		}

		return nil
	})

	if err != nil {
		return nil, fmt.Errorf("failed to scan wallet directory: %w", err)
	}

	return nfts, nil
}

// DeleteNFT removes stored NFT data
func (fs *FileStorage) DeleteNFT(ctx context.Context, walletAddr, mintAddr solanago.PublicKey) error {
	nftDir := fs.buildNFTPath(walletAddr, mintAddr)

	// Check if directory exists
	if _, err := os.Stat(nftDir); os.IsNotExist(err) {
		return fmt.Errorf("NFT not found: %s", mintAddr.String())
	}

	// Remove entire NFT directory
	if err := os.RemoveAll(nftDir); err != nil {
		return fmt.Errorf("failed to delete NFT directory: %w", err)
	}

	return nil
}

// Close cleans up storage resources (no-op for file storage)
func (fs *FileStorage) Close() error {
	return nil
}

// Helper methods

// buildNFTPath constructs the filesystem path for an NFT
func (fs *FileStorage) buildNFTPath(walletAddr, mintAddr solanago.PublicKey) string {
	return filepath.Join(
		fs.baseDir,
		"wallets",
		walletAddr.String(),
		"nfts",
		mintAddr.String(),
	)
}

// saveJSON marshals and saves data as JSON
func (fs *FileStorage) saveJSON(filePath string, data interface{}) error {
	// Pretty-print JSON for human readability
	// Explanation: Indented JSON makes it easier to examine backup files
	jsonData, err := json.MarshalIndent(data, "", "  ")
	if err != nil {
		return fmt.Errorf("failed to marshal JSON: %w", err)
	}

	if err := os.WriteFile(filePath, jsonData, fs.permissions); err != nil {
		return fmt.Errorf("failed to write file: %w", err)
	}

	return nil
}

// loadJSON loads and unmarshals JSON data
func (fs *FileStorage) loadJSON(filePath string, target interface{}) error {
	data, err := os.ReadFile(filePath)
	if err != nil {
		return err
	}

	if err := json.Unmarshal(data, target); err != nil {
		return fmt.Errorf("failed to unmarshal JSON: %w", err)
	}

	return nil
}

// calculateChecksum computes SHA256 hash of NFT data for integrity checking
func (fs *FileStorage) calculateChecksum(nftInfo *fetcher.NFTInfo) (string, error) {
	// Explanation: We hash the core NFT data to detect corruption
	// This helps ensure backup integrity over time
	data, err := json.Marshal(nftInfo)
	if err != nil {
		return "", err
	}

	hash := sha256.Sum256(data)
	return fmt.Sprintf("%x", hash), nil
}
