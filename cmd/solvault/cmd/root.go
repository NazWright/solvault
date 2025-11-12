package cmd

import (
	"fmt"

	"github.com/spf13/cobra"
)

var (
	// Version information
	Version   = "1.0.0"
	BuildTime = "dev"
	GitCommit = "dev"
)

// rootCmd represents the base command when called without any subcommands
var rootCmd = &cobra.Command{
	Use:   "solvault",
	Short: "SolVault - Solana NFT backup and verification tool",
	Long: `🔒 SolVault - Where we back up, verify, and prove authenticity on-chain

SolVault is a self-contained Solana NFT backup and verification system that runs 
locally, monitors your wallet for new NFT mints, downloads metadata and images, 
verifies authenticity through on-chain hashes, and optionally publishes proof pages.

Built with clarity. Verified with truth. Leave nothing unbacked.`,
	Version: fmt.Sprintf("%s (built %s, commit %s)", Version, BuildTime, GitCommit),
}

// Execute adds all child commands to the root command and sets flags appropriately.
// This is called by main.main(). It only needs to happen once to the rootCmd.
func Execute() error {
	return rootCmd.Execute()
}

func init() {
	// Global flags can be added here
	rootCmd.PersistentFlags().BoolP("verbose", "v", false, "enable verbose output")
	rootCmd.PersistentFlags().String("config", "", "config file (default is $HOME/.solvault.env)")
}
