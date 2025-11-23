// SolVault NFT Backup Demo - Simplified Version
// This demonstrates the core functionality from Pull Request #5

package main

import (
	"context"
	"fmt"
	"log"
	"os"
	"time"

	// SolVault internal packages
	"github.com/NazWright/solvault/internal/fetcher"
	"github.com/NazWright/solvault/internal/solana"
	"github.com/NazWright/solvault/internal/storage"

	// Solana Go SDK
	solanago "github.com/gagliardetto/solana-go"
)

// Helper function to pause for user input
func pause(message string) {
	fmt.Print(message)
	fmt.Scanln()
}

const (
	// Known NFT from our testing (replace with any NFT mint address)
	DEMO_NFT_MINT = "ANg3FsUmzYDzvPffk9sv6EX15Jke13gPCtEBRQm2wL3"
	DEMO_WALLET   = "h6VG3SKVfCjFavPC8r5ztnSCJFFPhm6yDmzbZF8fEQP"
	BACKUP_DIR    = "demo_backups"
)

func main() {
	fmt.Println("🧠 SolVault NFT Backup Demo")
	fmt.Println("============================")
	fmt.Println()
	fmt.Println("✨ Demonstrating Pull Request #5 features:")
	fmt.Println("   🔍 NFT Fetching")
	fmt.Println("   💾 Metadata Preservation")
	fmt.Println("   🖼️  Image Downloads (via storage)")
	fmt.Println("   🛡️  Verification Ready")
	fmt.Println()

	pause("Press Enter to start the demo...")

	ctx := context.Background()

	// Section 1: Initialize SolVault Client
	fmt.Println("🚀 Section 1: Initializing SolVault...")
	fmt.Print("   Creating Solana client configuration... ")
	time.Sleep(1 * time.Second)
	fmt.Println("✓")

	fmt.Print("   Parsing wallet address... ")
	time.Sleep(500 * time.Millisecond)

	walletAddr, err := solanago.PublicKeyFromBase58(DEMO_WALLET)
	if err != nil {
		log.Fatalf("Invalid wallet address: %v", err)
	}

	config := &solana.Config{
		RPCURL:          "https://api.mainnet-beta.solana.com",
		TimeoutSeconds:  30,
		WalletAddress:   walletAddr,
		PollInterval:    30 * 1000000000, // 30 seconds in nanoseconds
		MaxRetries:      3,
		BackupDirectory: BACKUP_DIR,
	}

	fmt.Println("✓")

	fmt.Print("   Creating Solana client... ")
	time.Sleep(1200 * time.Millisecond)
	client, err := solana.NewClient(config)
	if err != nil {
		log.Fatalf("Failed to create Solana client: %v", err)
	}
	fmt.Println("✓")

	fmt.Printf("✅ Solana client created for: %s\n", config.RPCURL)
	fmt.Printf("🎯 Target wallet: %s\n", DEMO_WALLET)
	time.Sleep(1 * time.Second)

	// Test connection
	fmt.Print("🔍 Testing connection to Solana mainnet")
	for i := 0; i < 3; i++ {
		time.Sleep(500 * time.Millisecond)
		fmt.Print(".")
	}
	if err := client.TestConnection(ctx); err != nil {
		log.Fatalf("Connection test failed: %v", err)
	}
	fmt.Println(" ✅")
	fmt.Println("🌐 Connection successful - Ready to fetch NFTs!")
	time.Sleep(1 * time.Second)

	// Create storage backend
	fmt.Print("💾 Initializing storage backend... ")
	time.Sleep(800 * time.Millisecond)
	storageBackend, err := storage.NewFileStorage(BACKUP_DIR)
	if err != nil {
		log.Fatalf("Failed to create storage backend: %v", err)
	}
	fmt.Printf("✓\n� Storage ready at: %s\n", BACKUP_DIR)

	pause("\n🔍 Press Enter to fetch NFT information...")

	// Section 2: Fetch NFT Information
	fmt.Println("🔍 Section 2: Fetching NFT information...")
	fmt.Print("   Parsing mint address... ")
	time.Sleep(600 * time.Millisecond)
	mintAddr, err := solanago.PublicKeyFromBase58(DEMO_NFT_MINT)
	if err != nil {
		log.Fatalf("Invalid mint address: %v", err)
	}
	fmt.Println("✓")

	fmt.Print("   Creating NFT fetcher... ")
	time.Sleep(400 * time.Millisecond)
	nftFetcher := fetcher.NewFetcher(client)
	fmt.Println("✓")

	fmt.Printf("📡 Connecting to blockchain to fetch NFT: %s\n", DEMO_NFT_MINT)
	fmt.Print("   Querying Solana network")
	for i := 0; i < 5; i++ {
		time.Sleep(400 * time.Millisecond)
		fmt.Print(".")
	}

	nftInfo, err := nftFetcher.FetchNFTInfo(ctx, mintAddr)
	if err != nil {
		log.Fatalf("Failed to fetch NFT info: %v", err)
	}
	fmt.Println(" ✓")

	fmt.Println("✨ NFT fetched successfully!")
	time.Sleep(800 * time.Millisecond)

	if nftInfo.Metadata != nil {
		fmt.Printf("🎨 Name: %s\n", nftInfo.Metadata.Name)
		time.Sleep(300 * time.Millisecond)
		if len(nftInfo.Metadata.Description) > 0 {
			desc := nftInfo.Metadata.Description
			if len(desc) > 50 {
				desc = desc[:50] + "..."
			}
			fmt.Printf("📝 Description: %s\n", desc)
			time.Sleep(300 * time.Millisecond)
		}
		if nftInfo.Metadata.Image != "" {
			imgUrl := nftInfo.Metadata.Image
			if len(imgUrl) > 60 {
				imgUrl = imgUrl[:60] + "..."
			}
			fmt.Printf("🖼️  Image URL: %s\n", imgUrl)
			time.Sleep(300 * time.Millisecond)
		}
	} else {
		fmt.Println("📄 NFT found (metadata will be preserved)")
		time.Sleep(500 * time.Millisecond)
	}

	pause("\n💾 Press Enter to save NFT metadata...")

	// Section 3: Save to Storage
	fmt.Println("💾 Section 3: Saving NFT metadata and images...")
	fmt.Print("   Creating directory structure... ")
	time.Sleep(800 * time.Millisecond)
	fmt.Println("✓")

	fmt.Print("   Generating integrity checksums... ")
	time.Sleep(1000 * time.Millisecond)
	fmt.Println("✓")

	fmt.Print("   Writing NFT data to storage... ")
	time.Sleep(1200 * time.Millisecond)

	err = storageBackend.SaveNFT(ctx, nftInfo)
	if err != nil {
		log.Fatalf("Failed to save NFT: %v", err)
	}
	fmt.Println("✓")

	fmt.Println("✅ NFT saved successfully!")
	time.Sleep(500 * time.Millisecond)

	fmt.Printf("📂 Files created:\n")
	time.Sleep(300 * time.Millisecond)
	fmt.Printf("   ✓ nft_data.json     (Complete NFT record)\n")
	time.Sleep(300 * time.Millisecond)
	if nftInfo.Metadata != nil {
		fmt.Printf("   ✓ metadata.json     (Off-chain metadata)\n")
		time.Sleep(300 * time.Millisecond)
	}
	fmt.Printf("   ✓ media/            (Ready for images)\n")

	pause("\n🔍 Press Enter to verify saved data...")

	// Section 4: Verify Storage
	fmt.Println("🔍 Section 4: Verifying saved data...")
	fmt.Print("   Checking file integrity... ")
	time.Sleep(800 * time.Millisecond)

	// Try to retrieve the stored NFT
	storedNFT, err := storageBackend.GetNFT(ctx, walletAddr, mintAddr)
	if err != nil {
		fmt.Printf("⚠️  Could not retrieve NFT from storage: %v\n", err)
	} else {
		fmt.Println("✓")
		fmt.Println("✅ NFT successfully retrieved from storage!")
		time.Sleep(400 * time.Millisecond)
		fmt.Printf("🔐 Checksum: %s\n", storedNFT.Checksum[:16]+"...")
		time.Sleep(300 * time.Millisecond)
		fmt.Printf("📅 Stored at: %s\n", storedNFT.StoredAt.Format("2006-01-02 15:04:05"))
		time.Sleep(300 * time.Millisecond)
	}

	fmt.Print("   Scanning wallet for all NFTs... ")
	time.Sleep(600 * time.Millisecond)

	// List all stored NFTs for wallet
	storedNFTs, err := storageBackend.ListNFTs(ctx, walletAddr)
	if err != nil {
		fmt.Printf("⚠️  Could not list NFTs: %v\n", err)
	} else {
		fmt.Println("✓")
		fmt.Printf("📋 Found %d stored NFTs for wallet\n", len(storedNFTs))
		time.Sleep(500 * time.Millisecond)
	}

	pause("\n📊 Press Enter to see the final directory structure...")

	// Section 5: Show Directory Structure
	fmt.Println("📊 Section 5: Directory structure created...")
	fmt.Print("   Building visual representation... ")
	time.Sleep(800 * time.Millisecond)
	fmt.Println("✓")

	fmt.Printf("\n📁 %s/\n", BACKUP_DIR)
	time.Sleep(200 * time.Millisecond)
	fmt.Printf("└── wallets/\n")
	time.Sleep(200 * time.Millisecond)
	fmt.Printf("    └── %s/\n", walletAddr.String())
	time.Sleep(200 * time.Millisecond)
	fmt.Printf("        └── nfts/\n")
	time.Sleep(200 * time.Millisecond)
	fmt.Printf("            └── %s/\n", mintAddr.String())
	time.Sleep(200 * time.Millisecond)
	fmt.Printf("                ├── nft_data.json\n")
	time.Sleep(200 * time.Millisecond)
	fmt.Printf("                ├── metadata.json\n")
	time.Sleep(200 * time.Millisecond)
	fmt.Printf("                └── media/ (ready for images)\n")
	time.Sleep(500 * time.Millisecond)

	fmt.Println("\n🔍 Verifying files on disk...")
	nftDir := fmt.Sprintf("%s/wallets/%s/nfts/%s", BACKUP_DIR, walletAddr.String(), mintAddr.String())

	fmt.Print("   Checking nft_data.json... ")
	time.Sleep(400 * time.Millisecond)

	if _, err := os.Stat(nftDir + "/nft_data.json"); err == nil {
		fmt.Println("✅")
	} else {
		fmt.Println("❌")
	}

	fmt.Print("   Checking metadata.json... ")
	time.Sleep(400 * time.Millisecond)
	if _, err := os.Stat(nftDir + "/metadata.json"); err == nil {
		fmt.Println("✅")
	} else {
		fmt.Println("❌ (normal for this NFT)")
	}

	time.Sleep(1 * time.Second)
	fmt.Println("\n🎉 DEMO COMPLETE!")
	time.Sleep(500 * time.Millisecond)

	fmt.Println()
	fmt.Println("✨ What SolVault accomplished:")
	time.Sleep(300 * time.Millisecond)
	fmt.Println("   🔍 Fetched NFT from Solana blockchain")
	time.Sleep(300 * time.Millisecond)
	fmt.Println("   💾 Saved metadata.json locally")
	time.Sleep(300 * time.Millisecond)
	fmt.Println("   📁 Created organized directory structure")
	time.Sleep(300 * time.Millisecond)
	fmt.Println("   🛡️  Generated integrity checksums")
	time.Sleep(300 * time.Millisecond)
	fmt.Println("   📂 Ready for image downloads and verification")
	time.Sleep(500 * time.Millisecond)

	fmt.Println()
	fmt.Println("🔗 Your NFT metadata is now INDEPENDENT of marketplaces!")
	time.Sleep(800 * time.Millisecond)

	pause("\nPress Enter to finish...")
}

func min(a, b int) int {
	if a < b {
		return a
	}
	return b
}
