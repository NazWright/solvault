package cmd

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"
	"time"

	"github.com/spf13/cobra"
)

// listCmd represents the list command
var listCmd = &cobra.Command{
	Use:   "list",
	Short: "Show all locally backed-up NFTs",
	Long: `List all NFTs that have been backed up locally by SolVault.

This command will:
• Scan the backup directory for NFT folders
• Display NFT names, backup dates, and verification status
• Show summary statistics
• Filter results by collection or status

Example:
  solvault list
  solvault list --collection "Cool Cats"
  solvault list --status verified
  solvault list --format json`,
	RunE: runList,
}

var (
	collection string
	status     string
	format     string
	showHashes bool
)

func runList(cmd *cobra.Command, args []string) error {
	fmt.Println("📋 Listing backed-up NFTs...")

	// Get backup directory from config or default
	backupDir, err := getBackupDirectory()
	if err != nil {
		return err
	}

	// Scan for NFT directories
	nfts, err := scanNFTDirectories(backupDir)
	if err != nil {
		return err
	}

	// Apply filters
	filteredNFTs := filterNFTs(nfts)

	if len(filteredNFTs) == 0 {
		fmt.Println("📭 No NFTs found matching criteria")
		return nil
	}

	// Display results
	switch format {
	case "json":
		return displayJSON(filteredNFTs)
	default:
		return displayTable(filteredNFTs)
	}
}

type NFTInfo struct {
	Name        string
	Path        string
	BackupDate  time.Time
	HasMetadata bool
	HasImage    bool
	HasHash     bool
	HasProof    bool
	Status      string
}

func getBackupDirectory() (string, error) {
	// TODO: Load from .env configuration
	homeDir, err := os.UserHomeDir()
	if err != nil {
		return "", fmt.Errorf("failed to get home directory: %w", err)
	}
	return filepath.Join(homeDir, "SolVaultBackups"), nil
}

func scanNFTDirectories(backupDir string) ([]NFTInfo, error) {
	var nfts []NFTInfo

	// Check if backup directory exists
	if _, err := os.Stat(backupDir); os.IsNotExist(err) {
		return nfts, fmt.Errorf("backup directory not found: %s. Run 'solvault init' first", backupDir)
	}

	// Scan directories
	entries, err := os.ReadDir(backupDir)
	if err != nil {
		return nfts, fmt.Errorf("failed to read backup directory: %w", err)
	}

	for _, entry := range entries {
		if !entry.IsDir() {
			continue
		}

		nftPath := filepath.Join(backupDir, entry.Name())
		nftInfo, err := analyzeNFTDirectory(entry.Name(), nftPath)
		if err != nil {
			fmt.Printf("⚠️  Warning: Failed to analyze %s: %v\n", entry.Name(), err)
			continue
		}

		nfts = append(nfts, nftInfo)
	}

	return nfts, nil
}

func analyzeNFTDirectory(name, path string) (NFTInfo, error) {
	info := NFTInfo{
		Name: name,
		Path: path,
	}

	// Get directory modification time as backup date
	if stat, err := os.Stat(path); err == nil {
		info.BackupDate = stat.ModTime()
	}

	// Check for required files
	info.HasMetadata = fileExists(filepath.Join(path, "metadata.json"))
	info.HasHash = fileExists(filepath.Join(path, "hash.txt"))
	info.HasProof = fileExists(filepath.Join(path, "proof.json"))

	// Check for image files
	imageExtensions := []string{".png", ".jpg", ".jpeg", ".gif", ".svg", ".webp"}
	for _, ext := range imageExtensions {
		if fileExists(filepath.Join(path, "image"+ext)) {
			info.HasImage = true
			break
		}
	}

	// Determine status
	if info.HasMetadata && info.HasImage && info.HasHash {
		if info.HasProof {
			info.Status = "verified"
		} else {
			info.Status = "backed-up"
		}
	} else {
		info.Status = "incomplete"
	}

	return info, nil
}

func fileExists(path string) bool {
	_, err := os.Stat(path)
	return err == nil
}

func filterNFTs(nfts []NFTInfo) []NFTInfo {
	var filtered []NFTInfo

	for _, nft := range nfts {
		// Filter by collection
		if collection != "" && !strings.Contains(strings.ToLower(nft.Name), strings.ToLower(collection)) {
			continue
		}

		// Filter by status
		if status != "" && nft.Status != status {
			continue
		}

		filtered = append(filtered, nft)
	}

	return filtered
}

func displayTable(nfts []NFTInfo) error {
	fmt.Printf("\n📊 Found %d NFTs:\n\n", len(nfts))
	fmt.Printf("%-30s %-12s %-20s %s\n", "NAME", "STATUS", "BACKUP DATE", "FILES")
	fmt.Println(strings.Repeat("-", 80))

	for _, nft := range nfts {
		files := buildFileStatus(nft)
		date := nft.BackupDate.Format("2006-01-02 15:04")
		fmt.Printf("%-30s %-12s %-20s %s\n",
			truncateString(nft.Name, 28),
			nft.Status,
			date,
			files)
	}

	// Summary
	fmt.Printf("\n📈 Summary:\n")
	statusCounts := make(map[string]int)
	for _, nft := range nfts {
		statusCounts[nft.Status]++
	}

	for status, count := range statusCounts {
		fmt.Printf("   %s: %d\n", status, count)
	}

	return nil
}

func displayJSON(nfts []NFTInfo) error {
	// TODO: Implement JSON output
	fmt.Println("📋 JSON output not yet implemented")
	return displayTable(nfts)
}

func buildFileStatus(nft NFTInfo) string {
	var parts []string

	if nft.HasMetadata {
		parts = append(parts, "M")
	}
	if nft.HasImage {
		parts = append(parts, "I")
	}
	if nft.HasHash {
		parts = append(parts, "H")
	}
	if nft.HasProof {
		parts = append(parts, "P")
	}

	return strings.Join(parts, ",")
}

func truncateString(s string, length int) string {
	if len(s) <= length {
		return s
	}
	return s[:length-2] + ".."
}

func init() {
	rootCmd.AddCommand(listCmd)

	listCmd.Flags().StringVar(&collection, "collection", "", "filter by collection name")
	listCmd.Flags().StringVar(&status, "status", "", "filter by status (verified, backed-up, incomplete)")
	listCmd.Flags().StringVar(&format, "format", "table", "output format (table, json)")
	listCmd.Flags().BoolVar(&showHashes, "show-hashes", false, "display file hashes")
}
