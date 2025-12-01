package cmd

import (
	"context"
	"encoding/json"
	"fmt"
	"time"

	"github.com/NazWright/solvault/internal/fetcher"
	"github.com/NazWright/solvault/internal/solana"
	solanago "github.com/gagliardetto/solana-go"
	"github.com/spf13/cobra"
)

// listTokensCmd represents the list-tokens command
var prettyOutput bool
var listTokensCmd = &cobra.Command{
	Use:   "list-tokens",
	Short: "List all NFTs in your wallet",
	Long: `List all NFTs in your configured wallet.

This will show you only the NFTs (tokens with supply=1 and decimals=0) that your wallet owns,
along with their mint addresses that you can use for testing.`,
	RunE: func(cmd *cobra.Command, args []string) error {
		fmt.Println("🔍 Loading your token accounts...")

		// Load configuration
		config, err := solana.LoadConfig()
		if err != nil {
			return fmt.Errorf("❌ Failed to load config: %w", err)
		}

		fmt.Printf("📋 Wallet: %s\n", config.WalletAddress.String())
		fmt.Printf("🌐 RPC: %s\n\n", config.RPCURL)

		// Create Solana client
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

		// Get token accounts
		fmt.Println("🔗 Fetching token accounts...")
		tokenAccounts, err := client.GetTokenAccountsByOwner(ctx)
		if err != nil {
			return fmt.Errorf("❌ Failed to get token accounts: %w", err)
		}

		if len(tokenAccounts) == 0 {
			fmt.Println("📭 No token accounts found in this wallet.")
			return nil
		}

		fmt.Printf("🔍 Found %d token account(s), filtering for NFTs...\n\n", len(tokenAccounts))

		nftCount := 0
		fetcherObj := fetcher.NewFetcher(client)

		for _, account := range tokenAccounts {
			rawJSON := account.Account.Data.GetRawJSON()
			if len(rawJSON) > 0 {
				var parsed map[string]interface{}
				if err := json.Unmarshal(rawJSON, &parsed); err == nil {
					var tokenInfo map[string]interface{}
					var ok bool
					if parsedData, exists := parsed["parsed"].(map[string]interface{}); exists {
						tokenInfo, ok = parsedData["info"].(map[string]interface{})
					} else {
						tokenInfo, ok = parsed["info"].(map[string]interface{})
					}
					if ok {
						var mint string
						var decimals float64
						var amount string
						var uiAmount float64
						if m, ok := tokenInfo["mint"].(string); ok {
							mint = m
						}
						if tokenAmount, ok := tokenInfo["tokenAmount"].(map[string]interface{}); ok {
							if a, ok := tokenAmount["amount"].(string); ok {
								amount = a
							}
							if d, ok := tokenAmount["decimals"].(float64); ok {
								decimals = d
							}
							if ua, ok := tokenAmount["uiAmount"].(float64); ok {
								uiAmount = ua
							}
						}
						if decimals == 0 && amount == "1" && uiAmount == 1 {
							nftCount++
							mintPubkey, err := solanago.PublicKeyFromBase58(mint)
							if err == nil {
								ctxMeta, cancelMeta := context.WithTimeout(context.Background(), 10*time.Second)
								defer cancelMeta()
								nftInfo, err := fetcherObj.FetchNFTInfo(ctxMeta, mintPubkey)
								if prettyOutput {
									fmt.Println("━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━")
									fmt.Printf("🖼️  NFT #%d\n", nftCount)
									if nftInfo.Metadata != nil {
										if nftInfo.Metadata.Name != "" {
											fmt.Printf("🏷️  Name: %s\n", nftInfo.Metadata.Name)
											fmt.Println("   The name of your NFT.")
										}
										if nftInfo.Metadata.Collection.Name != "" {
											fmt.Printf("📚 Collection: %s\n", nftInfo.Metadata.Collection.Name)
											fmt.Println("   The collection or series this NFT belongs to.")
										}
										if nftInfo.Metadata.Description != "" {
											fmt.Printf("📝 Description: %s\n", nftInfo.Metadata.Description)
											fmt.Println("   What this NFT is about.")
										}
										if nftInfo.Metadata.Image != "" {
											fmt.Printf("🖼️  Image URL: %s\n", nftInfo.Metadata.Image)
											fmt.Println("   Link to the NFT's image.")
										}
										fmt.Printf("🆔 NFT ID: %s\n", mint)
										fmt.Println("   Unique identifier for this NFT.")
										if len(nftInfo.Metadata.Attributes) > 0 {
											fmt.Printf("🔖 Attributes: ")
											for _, attr := range nftInfo.Metadata.Attributes {
												fmt.Printf("[%s: %v] ", attr.TraitType, attr.Value)
											}
											fmt.Println()
											fmt.Println("   Special traits or properties.")
										}
										fmt.Printf("🔗 Metadata URI: %s\n", nftInfo.MetadataURI)
										fmt.Println("   Link to full NFT details.")
									} else {
										fmt.Printf("🆔 NFT ID: %s\n", mint)
										fmt.Printf("🔗 Metadata URI: %s\n", nftInfo.MetadataURI)
										fmt.Printf("⚠️  Metadata not found\n")
									}
									fmt.Println("━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━\n")
								} else {
									// Default technical output
									fmt.Printf("NFT #%d:\n", nftCount)
									fmt.Printf("  Account Address: %s\n", account.Pubkey.String())
									fmt.Printf("  Mint Address:    %s\n", mint)
									if err == nil && nftInfo.Metadata != nil {
										fmt.Printf("  Name:            %s\n", nftInfo.Metadata.Name)
										fmt.Printf("  Symbol:          %s\n", nftInfo.Metadata.Symbol)
										fmt.Printf("  Description:     %s\n", nftInfo.Metadata.Description)
										fmt.Printf("  Image:           %s\n", nftInfo.Metadata.Image)
										if nftInfo.Metadata.Collection.Name != "" {
											fmt.Printf("  Collection:      %s\n", nftInfo.Metadata.Collection.Name)
										}
										if len(nftInfo.Metadata.Attributes) > 0 {
											fmt.Printf("  Attributes:      ")
											for _, attr := range nftInfo.Metadata.Attributes {
												fmt.Printf("[%s: %v] ", attr.TraitType, attr.Value)
											}
											fmt.Println()
										}
										fmt.Printf("  Metadata URI:    %s\n", nftInfo.MetadataURI)
									} else if err == nil {
										fmt.Printf("  Metadata URI:    %s\n", nftInfo.MetadataURI)
									} else {
										fmt.Printf("  Metadata:        (not found)\n")
									}
									fmt.Printf("  Amount:          %s (Supply: 1)\n", amount)
									fmt.Printf("  Decimals:        %.0f (NFT characteristic)\n", decimals)
									if state, ok := tokenInfo["state"].(string); ok {
										fmt.Printf("  State:           %s\n", state)
									}
									fmt.Println()
								}
							} else {
								fmt.Printf("  Metadata:        (invalid mint pubkey)\n")
							}
						}
					}
				}
			}
		}

		if nftCount == 0 {
			fmt.Println("📭 No NFTs found in this wallet.")
			fmt.Println("💡 NFTs are tokens with exactly 1 supply and 0 decimals.")
		} else {
			fmt.Printf("✅ Found %d NFT(s) in your wallet!\n\n", nftCount)
			fmt.Println("💡 To test the NFT fetcher, use any of the mint addresses above:")
			fmt.Println("   solvault test <mint-address>")
		}

		return nil
	},
}

// getKeys returns the keys of a map for debugging
func getKeys(m map[string]interface{}) []string {
	keys := make([]string, 0, len(m))
	for k := range m {
		keys = append(keys, k)
	}
	return keys
}

func init() {
	rootCmd.AddCommand(listTokensCmd)
	listTokensCmd.Flags().BoolVar(&prettyOutput, "pretty", false, "Show NFTs in a visually friendly format")
}
