package cmd

import (
	"fmt"

	"github.com/NazWright/solvault/internal/solana"
	"github.com/spf13/cobra"
)

// backupCmd represents the backup command
var backupCmd = &cobra.Command{
	Use:   "backup",
	Short: "Backup an NFT from your wallet",
	Long: `Interactively select a wallet, collection, and NFT to back up.

This command will:
• Prompt for wallet address
• Fetch collections and NFTs
• Let you select which NFT to back up
• Initiate the backup workflow
`,
	RunE: runBackup,
}

func runBackup(cmd *cobra.Command, args []string) error {
	// Load configuration using centralized loader
	config, err := solana.LoadConfig()
	if err != nil {
		fmt.Println("❌ Could not load config. Please run 'solvault init' first.")
		return nil
	}

	walletAddr := config.WalletAddress.String()
	if walletAddr == "" {
		fmt.Println("❌ Wallet address not found in config. Please run 'solvault init' and enter your wallet address.")
		return nil
	}

	// TODO: Fetch collections for walletAddr
	fmt.Printf("Fetching collections for wallet %s...\n", walletAddr)
	// collections := fetchCollections(walletAddr)
	// TODO: Fetch NFTs in collection
	// TODO: Initiate backup workflow

	fmt.Println("✅ (Stub) Backup command initialized. Next: integrate collection/NFT selection and backup logic.")
	return nil
}

func init() {
	rootCmd.AddCommand(backupCmd)
}
