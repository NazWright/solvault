package cmd

import (
	"fmt"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/NazWright/solvault/internal/solana"
	"github.com/spf13/cobra"
)

// watchCmd represents the watch command
var watchCmd = &cobra.Command{
	Use:   "watch",
	Short: "Monitor wallet for new NFT mints and back them up automatically",
	Long: `Watch mode monitors your Solana wallet for new NFT mint events and 
automatically backs up metadata, images, and generates verification hashes.

This command will:
• Connect to Solana RPC endpoint
• Monitor your wallet address for new transactions
• Detect NFT mint events in real-time
• Automatically download and backup NFT data
• Generate proof hashes and metadata

Example:
  solvault watch
  solvault watch --daemon
  solvault watch --poll-interval 15`,
	RunE: runWatch,
}

var (
	daemon       bool
	pollInterval int
)

func runWatch(cmd *cobra.Command, args []string) error {

	fmt.Println("👀 Starting SolVault watcher...")

	// Load configuration using centralized loader
	config, err := solana.LoadConfig()
	if err != nil {
		return fmt.Errorf("failed to load config: %w", err)
	}
	fmt.Printf("✅ Configuration loaded for wallet: %s\n", config.WalletAddress.String())

	if daemon {
		fmt.Println("🔄 Running in daemon mode...")
		// TODO: Implement daemon mode in future version
		fmt.Println("⚠️  Daemon mode not yet implemented. Running in foreground mode.")
	} else {
		fmt.Println("🖥️  Running in foreground mode. Press Ctrl+C to stop.")
	}

	// Set up graceful shutdown
	sigChan := make(chan os.Signal, 1)
	signal.Notify(sigChan, syscall.SIGINT, syscall.SIGTERM)

	// Start monitoring loop
	fmt.Printf("🔍 Monitoring wallet with %d second intervals...\n", pollInterval)
	ticker := time.NewTicker(time.Duration(pollInterval) * time.Second)
	defer ticker.Stop()

	for {
		select {
		case <-ticker.C:
			if err := checkForNewNFTs(); err != nil {
				fmt.Printf("❌ Error checking for NFTs: %v\n", err)
			}
		case <-sigChan:
			fmt.Println("\n🛑 Shutting down SolVault watcher...")
			return nil
		}
	}
}

func validateConfig() error {
	// TODO: Implement configuration validation
	// Check if .env exists and contains required values
	if _, err := os.Stat(".env"); err != nil {
		return fmt.Errorf("configuration file not found. Run 'solvault init' first")
	}

	fmt.Println("✅ Configuration validated")
	return nil
}

func checkForNewNFTs() error {
	// TODO: Implement actual NFT monitoring logic
	// This is a placeholder that will be implemented in the listener module
	fmt.Printf("⏰ [%s] Checking for new NFTs...\n", time.Now().Format("15:04:05"))
	return nil
}

func init() {
	rootCmd.AddCommand(watchCmd)

	watchCmd.Flags().BoolVar(&daemon, "daemon", false, "run in background daemon mode")
	watchCmd.Flags().IntVar(&pollInterval, "poll-interval", 30, "polling interval in seconds")
}
