package cmd

import (
	"fmt"
	"os"
	"path/filepath"

	"github.com/spf13/cobra"
)

// initCmd represents the init command
var initCmd = &cobra.Command{
	Use:   "init",
	Short: "Initialize SolVault configuration and backup directories",
	Long: `Initialize SolVault by creating configuration files and backup directories.

This command will:
• Create a .env configuration file with default settings
• Set up the backup directory structure (~/SolVaultBackups)
• Validate configuration and connectivity
• Guide you through the setup process

Example:
  solvault init
  solvault init --backup-dir /custom/backup/path`,
	RunE: runInit,
}

var (
	backupDir string
	force     bool
)

func runInit(cmd *cobra.Command, args []string) error {
	fmt.Println("🚀 Initializing SolVault...")

	// Set default backup directory if not specified
	if backupDir == "" {
		homeDir, err := os.UserHomeDir()
		if err != nil {
			return fmt.Errorf("failed to get home directory: %w", err)
		}
		backupDir = filepath.Join(homeDir, "SolVaultBackups")
	}

	// Create backup directory
	if err := createBackupDirectory(); err != nil {
		return err
	}

	// Create .env file
	if err := createEnvFile(); err != nil {
		return err
	}

	fmt.Println("✅ SolVault initialized successfully!")
	fmt.Printf("   Backup directory: %s\n", backupDir)
	fmt.Println("   Configuration: .env")
	fmt.Println("")
	fmt.Println("Next steps:")
	fmt.Println("1. Edit .env with your Solana RPC endpoint and wallet address")
	fmt.Println("2. Run 'solvault watch' to start monitoring for new NFTs")

	return nil
}

func createBackupDirectory() error {
	fmt.Printf("📁 Creating backup directory: %s\n", backupDir)

	if err := os.MkdirAll(backupDir, 0755); err != nil {
		return fmt.Errorf("failed to create backup directory: %w", err)
	}

	return nil
}

func createEnvFile() error {
	envPath := ".env"

	// Check if .env already exists
	if _, err := os.Stat(envPath); err == nil && !force {
		fmt.Printf("⚠️  .env file already exists. Use --force to overwrite\n")
		return nil
	}

	fmt.Printf("📝 Creating configuration file: %s\n", envPath)

	envContent := fmt.Sprintf(`# SolVault Configuration
# Edit these values according to your setup

# Solana RPC Configuration
SOLANA_RPC_URL=https://api.mainnet-beta.solana.com
SOLANA_WEBSOCKET_URL=wss://api.mainnet-beta.solana.com

# Your Solana wallet address to monitor
WALLET_ADDRESS=your_wallet_address_here

# Backup Settings
BACKUP_DIRECTORY=%s

# Optional: Proof Publishing (leave empty to disable)
PUBLISH_ENDPOINT=
PUBLISH_API_KEY=

# Monitoring Settings
POLL_INTERVAL_SECONDS=30
MAX_RETRIES=3
TIMEOUT_SECONDS=60
`, backupDir)

	if err := os.WriteFile(envPath, []byte(envContent), 0644); err != nil {
		return fmt.Errorf("failed to create .env file: %w", err)
	}

	return nil
}

func init() {
	rootCmd.AddCommand(initCmd)

	initCmd.Flags().StringVar(&backupDir, "backup-dir", "", "custom backup directory path")
	initCmd.Flags().BoolVar(&force, "force", false, "overwrite existing .env file")
}
