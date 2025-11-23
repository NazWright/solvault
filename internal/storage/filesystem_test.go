package storage

import (
	"context"
	"os"
	"path/filepath"
	"testing"
	"time"

	"github.com/NazWright/solvault/internal/fetcher"
	solanago "github.com/gagliardetto/solana-go"
)

// TestFileStorage_SaveAndGet tests basic save/retrieve operations
func TestFileStorage_SaveAndGet(t *testing.T) {
	// Explanation: We use a temporary directory for testing
	// This ensures tests don't interfere with real data
	tempDir, err := os.MkdirTemp("", "solvault_test_*")
	if err != nil {
		t.Fatalf("Failed to create temp dir: %v", err)
	}
	defer os.RemoveAll(tempDir) // Clean up

	// Create storage
	storage, err := NewFileStorage(tempDir)
	if err != nil {
		t.Fatalf("Failed to create storage: %v", err)
	}

	// Create test NFT data
	walletAddr := solanago.MustPublicKeyFromBase58("h6VG3SKVfCjFavPC8r5ztnSCJFFPhm6yDmzbZF8fEQP")
	mintAddr := solanago.MustPublicKeyFromBase58("ANg3FsUmzYDzvPffk9sv6EX15Jke13gPCtEBRQm2wL3")

	testNFT := &fetcher.NFTInfo{
		MintAddress:  mintAddr,
		TokenAccount: solanago.MustPublicKeyFromBase58("AZCdUmUV3JLpiL8jmpughB8zMP3sS6VZdbA1ga2Jj2dJ"),
		Owner:        walletAddr,
		Supply:       1,
		Decimals:     0,
		FetchedAt:    time.Now(),
		Metadata: &fetcher.NFTMetadata{
			Name:        "Test NFT",
			Symbol:      "TEST",
			Description: "A test NFT for SolVault",
			Image:       "https://example.com/image.png",
		},
	}

	ctx := context.Background()

	// Test saving
	err = storage.SaveNFT(ctx, testNFT)
	if err != nil {
		t.Fatalf("Failed to save NFT: %v", err)
	}

	// Test retrieving
	storedNFT, err := storage.GetNFT(ctx, walletAddr, mintAddr)
	if err != nil {
		t.Fatalf("Failed to get NFT: %v", err)
	}

	// Verify data integrity
	if storedNFT.NFTInfo.MintAddress != mintAddr {
		t.Errorf("Mint address mismatch: got %v, want %v", storedNFT.NFTInfo.MintAddress, mintAddr)
	}

	if storedNFT.NFTInfo.Metadata.Name != "Test NFT" {
		t.Errorf("Name mismatch: got %v, want %v", storedNFT.NFTInfo.Metadata.Name, "Test NFT")
	}

	if storedNFT.Checksum == "" {
		t.Error("Checksum should not be empty")
	}

	if storedNFT.Version != 1 {
		t.Errorf("Version mismatch: got %v, want %v", storedNFT.Version, 1)
	}
}

// TestFileStorage_ListNFTs tests listing multiple NFTs
func TestFileStorage_ListNFTs(t *testing.T) {
	tempDir, err := os.MkdirTemp("", "solvault_test_*")
	if err != nil {
		t.Fatalf("Failed to create temp dir: %v", err)
	}
	defer os.RemoveAll(tempDir)

	storage, err := NewFileStorage(tempDir)
	if err != nil {
		t.Fatalf("Failed to create storage: %v", err)
	}

	walletAddr := solanago.MustPublicKeyFromBase58("h6VG3SKVfCjFavPC8r5ztnSCJFFPhm6yDmzbZF8fEQP")
	ctx := context.Background()

	// Create multiple test NFTs
	for i := 0; i < 3; i++ {
		// Generate different mint addresses for each NFT
		mintAddr := solanago.NewWallet().PublicKey()

		testNFT := &fetcher.NFTInfo{
			MintAddress:  mintAddr,
			TokenAccount: solanago.NewWallet().PublicKey(),
			Owner:        walletAddr,
			Supply:       1,
			Decimals:     0,
			FetchedAt:    time.Now(),
		}

		err = storage.SaveNFT(ctx, testNFT)
		if err != nil {
			t.Fatalf("Failed to save NFT %d: %v", i, err)
		}
	}

	// List all NFTs
	nfts, err := storage.ListNFTs(ctx, walletAddr)
	if err != nil {
		t.Fatalf("Failed to list NFTs: %v", err)
	}

	if len(nfts) != 3 {
		t.Errorf("Expected 3 NFTs, got %d", len(nfts))
	}
}

// TestFileStorage_DirectoryStructure verifies the directory layout
func TestFileStorage_DirectoryStructure(t *testing.T) {
	tempDir, err := os.MkdirTemp("", "solvault_test_*")
	if err != nil {
		t.Fatalf("Failed to create temp dir: %v", err)
	}
	defer os.RemoveAll(tempDir)

	storage, err := NewFileStorage(tempDir)
	if err != nil {
		t.Fatalf("Failed to create storage: %v", err)
	}

	walletAddr := solanago.MustPublicKeyFromBase58("h6VG3SKVfCjFavPC8r5ztnSCJFFPhm6yDmzbZF8fEQP")
	mintAddr := solanago.MustPublicKeyFromBase58("ANg3FsUmzYDzvPffk9sv6EX15Jke13gPCtEBRQm2wL3")

	testNFT := &fetcher.NFTInfo{
		MintAddress:  mintAddr,
		TokenAccount: solanago.NewWallet().PublicKey(),
		Owner:        walletAddr,
		Supply:       1,
		Decimals:     0,
		FetchedAt:    time.Now(),
		Metadata: &fetcher.NFTMetadata{
			Name: "Directory Test NFT",
		},
	}

	ctx := context.Background()
	err = storage.SaveNFT(ctx, testNFT)
	if err != nil {
		t.Fatalf("Failed to save NFT: %v", err)
	}

	// Verify directory structure
	expectedDir := filepath.Join(tempDir, "wallets", walletAddr.String(), "nfts", mintAddr.String())

	// Check main directory exists
	if _, err := os.Stat(expectedDir); os.IsNotExist(err) {
		t.Errorf("Expected directory does not exist: %s", expectedDir)
	}

	// Check main data file exists
	nftDataFile := filepath.Join(expectedDir, "nft_data.json")
	if _, err := os.Stat(nftDataFile); os.IsNotExist(err) {
		t.Errorf("NFT data file does not exist: %s", nftDataFile)
	}

	// Check metadata file exists (since we provided metadata)
	metadataFile := filepath.Join(expectedDir, "metadata.json")
	if _, err := os.Stat(metadataFile); os.IsNotExist(err) {
		t.Errorf("Metadata file does not exist: %s", metadataFile)
	}
}
