package storage

import (
	"context"
	"time"

	"github.com/NazWright/solvault/internal/fetcher"
	solanago "github.com/gagliardetto/solana-go"
)

// StorageBackend defines the interface for NFT data storage
// This allows us to support different storage backends (local files, databases, etc.)
type StorageBackend interface {
	// SaveNFT stores NFT information with metadata
	SaveNFT(ctx context.Context, nftInfo *fetcher.NFTInfo) error

	// GetNFT retrieves stored NFT information
	GetNFT(ctx context.Context, walletAddr, mintAddr solanago.PublicKey) (*StoredNFT, error)

	// ListNFTs returns all NFTs for a wallet
	ListNFTs(ctx context.Context, walletAddr solanago.PublicKey) ([]*StoredNFT, error)

	// DeleteNFT removes stored NFT data
	DeleteNFT(ctx context.Context, walletAddr, mintAddr solanago.PublicKey) error

	// Close cleans up storage resources
	Close() error
}

// StoredNFT represents NFT data as stored on disk
// This includes the original fetched data plus storage metadata
type StoredNFT struct {
	// Original NFT information
	NFTInfo *fetcher.NFTInfo `json:"nft_info"`

	// Storage metadata
	StoredAt  time.Time `json:"stored_at"`  // When this was saved
	UpdatedAt time.Time `json:"updated_at"` // Last update time
	Version   int       `json:"version"`    // Data version for migrations
	Checksum  string    `json:"checksum"`   // Data integrity check

	// Backup metadata
	BackupPath string    `json:"backup_path"` // Path to image/media backup
	Verified   bool      `json:"verified"`    // Has been verified against blockchain
	LastCheck  time.Time `json:"last_check"`  // Last verification check
}

// BackupStats provides statistics about stored NFT data
type BackupStats struct {
	TotalNFTs       int       `json:"total_nfts"`
	LastBackup      time.Time `json:"last_backup"`
	TotalSize       int64     `json:"total_size_bytes"`
	VerifiedCount   int       `json:"verified_count"`
	UnverifiedCount int       `json:"unverified_count"`
}
