package main

import (
	"context"
	"fmt"
	"os"
	"path/filepath"
	"strings"
	"time"

	"github.com/NazWright/solvault/internal/fetcher"
	"github.com/NazWright/solvault/internal/solana"
	"github.com/NazWright/solvault/internal/storage"
	solanago "github.com/gagliardetto/solana-go"
)

func main() {
	fmt.Println("🚀 SolVault Enhanced NFT Demo - Rich Metadata & Media")
	fmt.Println(strings.Repeat("=", 60))
	
	// Test with known NFTs that have rich metadata
	// Using actual NFT mints that should have proper metadata and images
	testNFTs := []string{
		"ANg3FsUmzYDzvPffk9sv6EX15Jke13gPCtEBRQm2wL3", // Known NFT from our previous demo
		"7atgF8KQo4wJrD5ATGX7t1V6BgJyNsLWWH6nxiccKcjb", // Another known NFT
		"5FusRj5CjtQZPfaCu3gYTsE75k9GdxR8q4RmrP7LwRAx", // Magic Eden NFT example
		"8Rt3bXX5PpMUhYRrwVsKoRW4EtsRaKMw7rtx9hYM9hp",  // SMB NFT example
	}
	
	selectedNFT := testNFTs[0] // Default selection
	
	// Allow user to specify NFT mint address
	if len(os.Args) > 1 {
		selectedNFT = os.Args[1]
		fmt.Printf("🎯 Testing with user-provided NFT: %s\n\n", selectedNFT)
	} else {
		fmt.Printf("🎯 Testing with default NFT: %s\n", selectedNFT)
		fmt.Printf("   💡 You can specify your own: go run enhanced_demo.go <mint_address>\n\n")
	}
	
	ctx := context.Background()
	
	// Initialize Solana client
	fmt.Print("🔗 Connecting to Solana mainnet...")
	enhancedLoadingDots(3)
	
	config, err := solana.LoadConfig()
	if err != nil {
		fmt.Printf("\n❌ Failed to load config: %v\n", err)
		fmt.Println("\n💡 Make sure you have a valid Solana config or RPC endpoint")
		return
	}
	
	client, err := solana.NewClient(config)
	if err != nil {
		fmt.Printf("\n❌ Failed to create client: %v\n", err)
		return
	}
	fmt.Printf("✅ Connected to %s\n", config.RPCURL)
	enhancedPause()
	
	// Initialize fetcher with media downloader
	fmt.Print("⚙️  Initializing enhanced NFT fetcher...")
	enhancedLoadingDots(2)
	nftFetcher := fetcher.NewFetcher(client)
	defer nftFetcher.Close()
	fmt.Println("✅ Media downloader ready")
	enhancedPause()
	
	// Parse mint address
	mintPubkey, err := solanago.PublicKeyFromBase58(selectedNFT)
	if err != nil {
		fmt.Printf("❌ Invalid mint address: %v\n", err)
		return
	}
	
	// Fetch NFT info with metadata
	fmt.Print("📡 Fetching comprehensive NFT data...")
	enhancedProgressBar(30)
	
	nftInfo, err := nftFetcher.FetchNFTInfoDemo(ctx, mintPubkey)
	if err != nil {
		fmt.Printf("\n❌ Failed to fetch NFT: %v\n", err)
		fmt.Println("\n💡 This might be a token without NFT metadata or an invalid mint")
		return
	}
	
	fmt.Println("✅ NFT data retrieved successfully!")
	
	// Display comprehensive NFT information
	displayNFTInfo(nftInfo)
	enhancedPause()
	
	// Initialize storage for backup
	fmt.Print("💾 Setting up demo backup storage...")
	enhancedLoadingDots(2)
	
	backupDir := "demo_backups"
	
	// Clear previous demo backups
	fmt.Print("🧹 Clearing previous demo data...")
	os.RemoveAll(backupDir)
	enhancedLoadingDots(1)
	
	fileStorage, err := storage.NewFileStorage(backupDir)
	if err != nil {
		fmt.Printf("\n❌ Failed to create storage: %v\n", err)
		return
	}
	
	fmt.Printf("✅ Clean demo storage ready\n")
	enhancedPause()
	
	// Download media if available
	if nftInfo.Metadata != nil && hasMediaURLs(nftInfo.Metadata) {
		fmt.Println("🖼️  Detected media files in NFT metadata!")
		fmt.Print("📥 Downloading media files...")
		
		// Create media directory  
		mediaDir := filepath.Join(backupDir, "wallets", "demo", "nfts", nftInfo.MintAddress.String(), "media")
		
		enhancedProgressBar(20)
		
		err := nftFetcher.DownloadMediaFiles(ctx, nftInfo, mediaDir)
		if err != nil {
			fmt.Printf("\n⚠️  Media download encountered issues: %v\n", err)
		}
		
		if len(nftInfo.MediaFiles) > 0 {
			fmt.Printf("✅ Downloaded %d media files!\n", len(nftInfo.MediaFiles))
			displayMediaFiles(nftInfo.MediaFiles)
		} else {
			fmt.Println("⚠️  No media files were downloaded")
		}
		enhancedPause()
	} else {
		fmt.Println("ℹ️  No media URLs found in metadata (minimal NFT)")
	}
	
	// Save complete NFT backup
	fmt.Print("💽 Creating comprehensive backup...")
	enhancedProgressBar(15)
	
	// Set demo wallet address for storage
	demoWallet, _ := solanago.PublicKeyFromBase58("11111111111111111111111111111112") // Dummy wallet
	nftInfo.Owner = demoWallet
	
	err = fileStorage.SaveNFT(ctx, nftInfo)
	if err != nil {
		fmt.Printf("\n❌ Failed to save NFT: %v\n", err)
		return
	}
	
	fmt.Println("✅ Backup completed successfully!")
	
	// Display backup summary
	displayBackupSummary(backupDir, nftInfo)
	
	fmt.Println("\n🎉 Enhanced demo completed!")
	fmt.Println("📁 Backup files created in:", backupDir)
	
	if len(nftInfo.MediaFiles) > 0 {
		fmt.Printf("🖼️  Media files saved to: %s\n", 
			filepath.Join(backupDir, "wallets", "demo", "nfts", nftInfo.MintAddress.String(), "media"))
	}
	
	// Clean up demo backups after showing results
	fmt.Print("\n🧹 Cleaning up demo files...")
	enhancedLoadingDots(2)
	os.RemoveAll(backupDir)
	fmt.Println("✅ Demo cleanup complete!")
}

func displayNFTInfo(nftInfo *fetcher.NFTInfo) {
	fmt.Println("\n📋 Comprehensive NFT Information:")
	fmt.Printf("   • Mint Address: %s\n", nftInfo.MintAddress)
	fmt.Printf("   • Token Supply: %d\n", nftInfo.Supply)
	fmt.Printf("   • Decimals: %d\n", nftInfo.Decimals)
	fmt.Printf("   • Fetched At: %s\n", nftInfo.FetchedAt.Format("2006-01-02 15:04:05"))
	
	if nftInfo.MetadataURI != "" {
		fmt.Printf("   • Metadata URI: %s\n", truncateString(nftInfo.MetadataURI, 60))
	}
	
	if nftInfo.Metadata != nil {
		fmt.Println("\n🏷️  Metadata Details:")
		fmt.Printf("   • Name: %s\n", nftInfo.Metadata.Name)
		fmt.Printf("   • Symbol: %s\n", nftInfo.Metadata.Symbol)
		
		if nftInfo.Metadata.Description != "" {
			fmt.Printf("   • Description: %s\n", truncateString(nftInfo.Metadata.Description, 80))
		}
		
		if nftInfo.Metadata.Image != "" {
			fmt.Printf("   • Image URL: %s\n", truncateString(nftInfo.Metadata.Image, 60))
		}
		
		if nftInfo.Metadata.AnimationURL != "" {
			fmt.Printf("   • Animation URL: %s\n", truncateString(nftInfo.Metadata.AnimationURL, 60))
		}
		
		if len(nftInfo.Metadata.Attributes) > 0 {
			fmt.Printf("   • Attributes: %d traits found\n", len(nftInfo.Metadata.Attributes))
			
			// Show first few attributes
			for i, attr := range nftInfo.Metadata.Attributes {
				if i >= 3 {
					fmt.Printf("     ... and %d more\n", len(nftInfo.Metadata.Attributes)-3)
					break
				}
				fmt.Printf("     - %s: %v\n", attr.TraitType, attr.Value)
			}
		}
		
		if nftInfo.Metadata.Collection.Name != "" {
			fmt.Printf("   • Collection: %s\n", nftInfo.Metadata.Collection.Name)
		}
	} else {
		fmt.Println("\n⚠️  No off-chain metadata found")
		fmt.Println("   This might be a minimal token or the metadata URI is inaccessible")
	}
}

func displayMediaFiles(mediaFiles []*fetcher.MediaFile) {
	fmt.Println("\n🖼️  Downloaded Media Files:")
	for i, media := range mediaFiles {
		fmt.Printf("   %d. %s\n", i+1, media.Filename)
		fmt.Printf("      • Type: %s (%s)\n", media.MediaType, media.ContentType)
		fmt.Printf("      • Size: %s\n", formatBytes(media.Size))
		fmt.Printf("      • Checksum: %s\n", media.Checksum[:16]+"...")
		fmt.Printf("      • Downloaded: %s\n", media.DownloadedAt.Format("15:04:05"))
	}
}

func displayBackupSummary(backupDir string, nftInfo *fetcher.NFTInfo) {
	fmt.Println("\n📦 Backup Summary:")
	
	nftPath := filepath.Join(backupDir, "wallets", "demo", "nfts", nftInfo.MintAddress.String())
	
	// Check what files were created
	files := []string{"nft_data.json", "metadata.json", "media_manifest.json"}
	for _, file := range files {
		fullPath := filepath.Join(nftPath, file)
		if _, err := os.Stat(fullPath); err == nil {
			if stat, err := os.Stat(fullPath); err == nil {
				fmt.Printf("   ✅ %s (%s)\n", file, formatBytes(stat.Size()))
			}
		}
	}
	
	// Check media directory
	mediaDir := filepath.Join(nftPath, "media")
	if _, err := os.Stat(mediaDir); err == nil {
		fmt.Printf("   📁 media/ directory with %d files\n", len(nftInfo.MediaFiles))
	}
}

func hasMediaURLs(metadata *fetcher.NFTMetadata) bool {
	return metadata.Image != "" || metadata.AnimationURL != "" || len(metadata.Properties.Files) > 0
}

func truncateString(s string, maxLen int) string {
	if len(s) <= maxLen {
		return s
	}
	return s[:maxLen-3] + "..."
}

func formatBytes(bytes int64) string {
	const unit = 1024
	if bytes < unit {
		return fmt.Sprintf("%d B", bytes)
	}
	div, exp := int64(unit), 0
	for n := bytes / unit; n >= unit; n /= unit {
		div *= unit
		exp++
	}
	return fmt.Sprintf("%.1f %cB", float64(bytes)/float64(div), "KMGTPE"[exp])
}

// Helper functions for enhanced demo effects
func enhancedPause() {
	fmt.Print("\n⏸️  Press Enter to continue...")
	fmt.Scanln()
	fmt.Println()
}

func enhancedLoadingDots(count int) {
	for i := 0; i < count; i++ {
		time.Sleep(500 * time.Millisecond)
		fmt.Print(".")
	}
	time.Sleep(200 * time.Millisecond)
}

func enhancedProgressBar(steps int) {
	fmt.Print("\n   [")
	for i := 0; i < steps; i++ {
		time.Sleep(100 * time.Millisecond)
		fmt.Print("█")
	}
	fmt.Print("] ")
	time.Sleep(200 * time.Millisecond)
}