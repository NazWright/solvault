package cmd

import (
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"

	"github.com/spf13/cobra"
)

// infoCmd represents the info command
var infoCmd = &cobra.Command{
	Use:   "info <mint-address-or-name>",
	Short: "Display NFT metadata and hash information",
	Long: `Display detailed information about a backed-up NFT including metadata,
file hashes, verification status, and proof information.

This command will:
• Show NFT metadata (name, description, attributes)
• Display file hashes and verification status
• Show backup location and file sizes
• Display proof information if available

Example:
  solvault info "Cool Cat #1234"
  solvault info 7xKXtg2CW87d97TXJSDpbD5jBkheTqA83TZRuJosgAsU
  solvault info --format json "Midnight Lion #01"`,
	Args: cobra.ExactArgs(1),
	RunE: runInfo,
}

var (
	infoFormat string
	showFiles  bool
)

func runInfo(cmd *cobra.Command, args []string) error {
	identifier := args[0]
	fmt.Printf("🔍 Looking up NFT: %s\n", identifier)

	// Get backup directory
	backupDir, err := getBackupDirectory()
	if err != nil {
		return err
	}

	// Find NFT directory
	nftPath, err := findNFTDirectory(backupDir, identifier)
	if err != nil {
		return err
	}

	// Load NFT information
	nftInfo, err := loadNFTInfo(nftPath)
	if err != nil {
		return err
	}

	// Display information
	switch infoFormat {
	case "json":
		return displayNFTInfoJSON(nftInfo)
	default:
		return displayNFTInfoTable(nftInfo)
	}
}

type DetailedNFTInfo struct {
	NFTInfo
	Metadata  map[string]interface{}
	Hash      string
	ProofData map[string]interface{}
	Files     []FileInfo
	TotalSize int64
}

type FileInfo struct {
	Name string
	Size int64
	Path string
}

func findNFTDirectory(backupDir, identifier string) (string, error) {
	// First try exact match by directory name
	exactPath := filepath.Join(backupDir, identifier)
	if _, err := os.Stat(exactPath); err == nil {
		return exactPath, nil
	}

	// Scan all directories for matches
	entries, err := os.ReadDir(backupDir)
	if err != nil {
		return "", fmt.Errorf("failed to read backup directory: %w", err)
	}

	var matches []string
	for _, entry := range entries {
		if !entry.IsDir() {
			continue
		}

		// Check if directory name contains identifier (case-insensitive)
		name := entry.Name()
		if contains(name, identifier) {
			matches = append(matches, filepath.Join(backupDir, name))
		}
	}

	if len(matches) == 0 {
		return "", fmt.Errorf("NFT not found: %s", identifier)
	}

	if len(matches) > 1 {
		fmt.Printf("⚠️  Multiple matches found:\n")
		for i, match := range matches {
			fmt.Printf("  %d. %s\n", i+1, filepath.Base(match))
		}
		return "", fmt.Errorf("multiple matches found, please be more specific")
	}

	return matches[0], nil
}

func contains(s, substr string) bool {
	// Simple case-insensitive contains check
	s = toLower(s)
	substr = toLower(substr)
	return len(s) >= len(substr) && findSubstring(s, substr)
}

func toLower(s string) string {
	result := make([]byte, len(s))
	for i, b := range []byte(s) {
		if b >= 'A' && b <= 'Z' {
			result[i] = b + 32
		} else {
			result[i] = b
		}
	}
	return string(result)
}

func findSubstring(s, substr string) bool {
	if len(substr) == 0 {
		return true
	}
	for i := 0; i <= len(s)-len(substr); i++ {
		if s[i:i+len(substr)] == substr {
			return true
		}
	}
	return false
}

func loadNFTInfo(nftPath string) (*DetailedNFTInfo, error) {
	name := filepath.Base(nftPath)

	// Get basic NFT info
	basicInfo, err := analyzeNFTDirectory(name, nftPath)
	if err != nil {
		return nil, err
	}

	detailed := &DetailedNFTInfo{
		NFTInfo: basicInfo,
	}

	// Load metadata if available
	if detailed.HasMetadata {
		if metadata, err := loadJSONFile(filepath.Join(nftPath, "metadata.json")); err == nil {
			detailed.Metadata = metadata
		}
	}

	// Load hash if available
	if detailed.HasHash {
		if hashBytes, err := os.ReadFile(filepath.Join(nftPath, "hash.txt")); err == nil {
			detailed.Hash = string(hashBytes)
		}
	}

	// Load proof if available
	if detailed.HasProof {
		if proof, err := loadJSONFile(filepath.Join(nftPath, "proof.json")); err == nil {
			detailed.ProofData = proof
		}
	}

	// Get file information
	detailed.Files, detailed.TotalSize = getFileInfo(nftPath)

	return detailed, nil
}

func loadJSONFile(path string) (map[string]interface{}, error) {
	data, err := os.ReadFile(path)
	if err != nil {
		return nil, err
	}

	var result map[string]interface{}
	if err := json.Unmarshal(data, &result); err != nil {
		return nil, err
	}

	return result, nil
}

func getFileInfo(dirPath string) ([]FileInfo, int64) {
	var files []FileInfo
	var totalSize int64

	entries, err := os.ReadDir(dirPath)
	if err != nil {
		return files, 0
	}

	for _, entry := range entries {
		if entry.IsDir() {
			continue
		}

		filePath := filepath.Join(dirPath, entry.Name())
		if stat, err := os.Stat(filePath); err == nil {
			files = append(files, FileInfo{
				Name: entry.Name(),
				Size: stat.Size(),
				Path: filePath,
			})
			totalSize += stat.Size()
		}
	}

	return files, totalSize
}

func displayNFTInfoTable(info *DetailedNFTInfo) error {
	fmt.Printf("\n📋 NFT Information\n")
	fmt.Printf("═══════════════════════════════════════════════════════════════════════════════\n")
	fmt.Printf("Name:         %s\n", info.Name)
	fmt.Printf("Status:       %s\n", info.Status)
	fmt.Printf("Backup Date:  %s\n", info.BackupDate.Format("2006-01-02 15:04:05"))
	fmt.Printf("Location:     %s\n", info.Path)
	fmt.Printf("Total Size:   %s\n", formatBytes(info.TotalSize))

	// Metadata section
	if info.Metadata != nil {
		fmt.Printf("\n📝 Metadata\n")
		fmt.Printf("───────────────────────────────────────────────────────────────────────────────\n")

		if name, ok := info.Metadata["name"].(string); ok {
			fmt.Printf("NFT Name:     %s\n", name)
		}
		if desc, ok := info.Metadata["description"].(string); ok && desc != "" {
			fmt.Printf("Description:  %s\n", truncateString(desc, 60))
		}
		if image, ok := info.Metadata["image"].(string); ok {
			fmt.Printf("Image URI:    %s\n", image)
		}
	}

	// Hash section
	if info.Hash != "" {
		fmt.Printf("\n🔐 Verification\n")
		fmt.Printf("───────────────────────────────────────────────────────────────────────────────\n")
		fmt.Printf("Hash:         %s\n", info.Hash)
	}

	// Files section
	if showFiles && len(info.Files) > 0 {
		fmt.Printf("\n📁 Files\n")
		fmt.Printf("───────────────────────────────────────────────────────────────────────────────\n")
		for _, file := range info.Files {
			fmt.Printf("%-20s %10s\n", file.Name, formatBytes(file.Size))
		}
	}

	// Proof section
	if info.ProofData != nil {
		fmt.Printf("\n✅ Proof Information\n")
		fmt.Printf("───────────────────────────────────────────────────────────────────────────────\n")
		if verifiedAt, ok := info.ProofData["verified_at"].(string); ok {
			fmt.Printf("Verified At:  %s\n", verifiedAt)
		}
		if verifiedBy, ok := info.ProofData["verified_by"].(string); ok {
			fmt.Printf("Verified By:  %s\n", verifiedBy)
		}
	}

	return nil
}

func displayNFTInfoJSON(info *DetailedNFTInfo) error {
	// TODO: Implement JSON output
	fmt.Println("📋 JSON output not yet implemented")
	return displayNFTInfoTable(info)
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

func init() {
	rootCmd.AddCommand(infoCmd)

	infoCmd.Flags().StringVar(&infoFormat, "format", "table", "output format (table, json)")
	infoCmd.Flags().BoolVar(&showFiles, "show-files", false, "show detailed file information")
}
