package cmd

import (
	"context"
	"fmt"
	"time"

	"github.com/NazWright/solvault/internal/fetcher"
	"github.com/NazWright/solvault/internal/solana"
	solanago "github.com/gagliardetto/solana-go"
	"github.com/spf13/cobra"
)

// testCmd represents the test command
var testCmd = &cobra.Command{
	Use:   "test [mint-address]",
	Short: "Test NFT fetcher with a specific mint address",
	Long: `Test the NFT fetcher functionality by providing a Solana NFT mint address.

This command will attempt to:
1. Connect to Solana RPC
2. Fetch NFT account information
3. Find associated token accounts
4. Retrieve metadata URI
5. Fetch off-chain metadata

Example:
  solvault test 7pFkKJvNyLwXXGEiP7Xbs8A1r7gVsHkWRu9vH5JnYtEP`,
	Args: cobra.ExactArgs(1),
	RunE: func(cmd *cobra.Command, args []string) error {
		mintAddressStr := args[0]

		// Parse the mint address
		mintAddress, err := solanago.PublicKeyFromBase58(mintAddressStr)
		if err != nil {
			return fmt.Errorf("❌ Invalid mint address format: %w", err)
		}

		fmt.Printf("🧪 Testing NFT fetcher with mint: %s\n\n", mintAddress.String())

		// Load configuration
		fmt.Println("📋 Loading configuration...")
		config, err := solana.LoadConfig()
		if err != nil {
			return fmt.Errorf("❌ Failed to load config: %w", err)
		}

		// Create Solana client
		fmt.Println("🔗 Connecting to Solana...")
		client, err := solana.NewClient(config)
		if err != nil {
			return fmt.Errorf("❌ Failed to create Solana client: %w", err)
		}
		defer client.Close()

		// Test connection
		ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
		defer cancel()

		if err := client.TestConnection(ctx); err != nil {
			return fmt.Errorf("❌ Failed to connect to Solana: %w", err)
		}
		fmt.Println("✅ Connected to Solana RPC")

		// Create NFT fetcher
		fmt.Println("🚀 Creating NFT fetcher...")
		nftFetcher := fetcher.NewFetcher(client)
		defer nftFetcher.Close()

		// Fetch NFT info
		fmt.Println("🔍 Fetching NFT information...")
		ctx2, cancel2 := context.WithTimeout(context.Background(), 30*time.Second)
		defer cancel2()

		nftInfo, err := nftFetcher.FetchNFTInfo(ctx2, mintAddress)
		if err != nil {
			return fmt.Errorf("❌ Failed to fetch NFT info: %w", err)
		}

		// Display results
		fmt.Println("\n🎉 Successfully fetched NFT information!")
		fmt.Println("==================================================")
		fmt.Printf("Mint Address:    %s\n", nftInfo.MintAddress.String())
		fmt.Printf("Token Account:   %s\n", nftInfo.TokenAccount.String())
		fmt.Printf("Owner:           %s\n", nftInfo.Owner.String())
		fmt.Printf("Supply:          %d\n", nftInfo.Supply)
		fmt.Printf("Decimals:        %d\n", nftInfo.Decimals)
		fmt.Printf("Fetched At:      %s\n", nftInfo.FetchedAt.Format(time.RFC3339))

		if nftInfo.MetadataURI != "" {
			fmt.Printf("Metadata URI:    %s\n", nftInfo.MetadataURI)
		} else {
			fmt.Println("Metadata URI:    ⚠️  Not found")
		}

		if nftInfo.Metadata != nil {
			fmt.Println("\n📋 Metadata:")
			fmt.Printf("  Name:          %s\n", nftInfo.Metadata.Name)
			fmt.Printf("  Symbol:        %s\n", nftInfo.Metadata.Symbol)
			fmt.Printf("  Description:   %s\n", nftInfo.Metadata.Description)
			if nftInfo.Metadata.Image != "" {
				fmt.Printf("  Image:         %s\n", nftInfo.Metadata.Image)
			}
			if len(nftInfo.Metadata.Attributes) > 0 {
				fmt.Printf("  Attributes:    %d traits\n", len(nftInfo.Metadata.Attributes))
				for i, attr := range nftInfo.Metadata.Attributes {
					if i < 3 { // Show first 3 attributes
						fmt.Printf("    - %s: %v\n", attr.TraitType, attr.Value)
					}
				}
				if len(nftInfo.Metadata.Attributes) > 3 {
					fmt.Printf("    ... and %d more\n", len(nftInfo.Metadata.Attributes)-3)
				}
			}
		} else {
			fmt.Println("\n📋 Metadata:      ⚠️  Not available")
		}

		fmt.Println("\n✅ Test completed successfully!")
		return nil
	},
}

func init() {
	rootCmd.AddCommand(testCmd)
}
