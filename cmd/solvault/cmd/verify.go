package cmd

import (
	"crypto/sha256"
	"encoding/json"
	"fmt"
	"io"
	"os"
	"path/filepath"
	"time"

	"github.com/spf13/cobra"
)

// verifyCmd represents the verify command
var verifyCmd = &cobra.Command{
	Use:   "verify <mint-address-or-name>",
	Short: "Verify NFT authenticity and optionally publish proof JSON",
	Long: `Verify the authenticity of a backed-up NFT by comparing hashes and 
generating or updating proof documentation.

This command will:
• Recalculate image and metadata hashes
• Compare against stored hash values
• Generate or update proof.json with verification results
• Optionally publish proof to web endpoint

Example:
  solvault verify "Cool Cat #1234"
  solvault verify 7xKXtg2CW87d97TXJSDpbD5jBkheTqA83TZRuJosgAsU --publish
  solvault verify "Midnight Lion #01" --force-recompute`,
	Args: cobra.ExactArgs(1),
	RunE: runVerify,
}

var (
	publish        bool
	forceRecompute bool
	skipOnChain    bool
)

func runVerify(cmd *cobra.Command, args []string) error {
	identifier := args[0]
	fmt.Printf("🔍 Verifying NFT: %s\n", identifier)

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

	// Perform verification
	result, err := performVerification(nftPath)
	if err != nil {
		return err
	}

	// Display results
	if err := displayVerificationResults(result); err != nil {
		return err
	}

	// Generate/update proof
	if err := generateProof(nftPath, result); err != nil {
		return err
	}

	// Publish if requested
	if publish {
		if err := publishProof(nftPath, result); err != nil {
			fmt.Printf("⚠️  Failed to publish proof: %v\n", err)
		}
	}

	return nil
}

type VerificationResult struct {
	NFTName      string
	NFTPath      string
	Status       string
	ImageHash    string
	StoredHash   string
	MetadataHash string
	HashMatch    bool
	HasImage     bool
	HasMetadata  bool
	VerifiedAt   time.Time
	Errors       []string
}

func performVerification(nftPath string) (*VerificationResult, error) {
	result := &VerificationResult{
		NFTName:    filepath.Base(nftPath),
		NFTPath:    nftPath,
		VerifiedAt: time.Now(),
	}

	fmt.Println("🔐 Computing hashes...")

	// Check for required files
	result.HasMetadata = fileExists(filepath.Join(nftPath, "metadata.json"))
	result.HasImage = findImageFile(nftPath) != ""

	if !result.HasImage {
		result.Errors = append(result.Errors, "No image file found")
		result.Status = "incomplete"
		return result, nil
	}

	// Compute image hash
	imageFile := findImageFile(nftPath)
	if imageFile != "" {
		hash, err := computeFileHash(imageFile)
		if err != nil {
			result.Errors = append(result.Errors, fmt.Sprintf("Failed to compute image hash: %v", err))
		} else {
			result.ImageHash = hash
		}
	}

	// Compute metadata hash
	if result.HasMetadata {
		metadataFile := filepath.Join(nftPath, "metadata.json")
		hash, err := computeFileHash(metadataFile)
		if err != nil {
			result.Errors = append(result.Errors, fmt.Sprintf("Failed to compute metadata hash: %v", err))
		} else {
			result.MetadataHash = hash
		}
	}

	// Compare with stored hash
	hashFile := filepath.Join(nftPath, "hash.txt")
	if fileExists(hashFile) {
		if storedHashBytes, err := os.ReadFile(hashFile); err == nil {
			result.StoredHash = string(storedHashBytes)
			result.HashMatch = result.ImageHash == result.StoredHash
		}
	}

	// Determine overall status
	if len(result.Errors) > 0 {
		result.Status = "error"
	} else if result.HashMatch || result.StoredHash == "" {
		result.Status = "authentic"
	} else {
		result.Status = "tampered"
	}

	// Store new hash if none exists or force recompute
	if result.StoredHash == "" || forceRecompute {
		if result.ImageHash != "" {
			if err := os.WriteFile(hashFile, []byte(result.ImageHash), 0644); err != nil {
				result.Errors = append(result.Errors, fmt.Sprintf("Failed to save hash: %v", err))
			} else {
				result.StoredHash = result.ImageHash
				result.HashMatch = true
			}
		}
	}

	return result, nil
}

func findImageFile(nftPath string) string {
	imageExtensions := []string{"image.png", "image.jpg", "image.jpeg", "image.gif", "image.svg", "image.webp"}

	for _, ext := range imageExtensions {
		path := filepath.Join(nftPath, ext)
		if fileExists(path) {
			return path
		}
	}

	// Fallback: look for any image file
	entries, err := os.ReadDir(nftPath)
	if err != nil {
		return ""
	}

	for _, entry := range entries {
		if entry.IsDir() {
			continue
		}

		name := entry.Name()
		ext := filepath.Ext(name)
		switch ext {
		case ".png", ".jpg", ".jpeg", ".gif", ".svg", ".webp":
			return filepath.Join(nftPath, name)
		}
	}

	return ""
}

func computeFileHash(filePath string) (string, error) {
	file, err := os.Open(filePath)
	if err != nil {
		return "", err
	}
	defer file.Close()

	hasher := sha256.New()
	if _, err := io.Copy(hasher, file); err != nil {
		return "", err
	}

	return fmt.Sprintf("sha256:%x", hasher.Sum(nil)), nil
}

func displayVerificationResults(result *VerificationResult) error {
	fmt.Printf("\n🔍 Verification Results\n")
	fmt.Printf("═══════════════════════════════════════════════════════════════════════════════\n")
	fmt.Printf("NFT Name:     %s\n", result.NFTName)
	fmt.Printf("Status:       %s", result.Status)

	// Add status emoji
	switch result.Status {
	case "authentic":
		fmt.Printf(" ✅")
	case "tampered":
		fmt.Printf(" ❌")
	case "incomplete":
		fmt.Printf(" ⚠️")
	case "error":
		fmt.Printf(" 🚫")
	}
	fmt.Println()

	fmt.Printf("Verified At:  %s\n", result.VerifiedAt.Format("2006-01-02 15:04:05"))

	if result.ImageHash != "" {
		fmt.Printf("\n🔐 Hash Information\n")
		fmt.Printf("───────────────────────────────────────────────────────────────────────────────\n")
		fmt.Printf("Current Hash: %s\n", result.ImageHash)
		if result.StoredHash != "" {
			fmt.Printf("Stored Hash:  %s\n", result.StoredHash)
			if result.HashMatch {
				fmt.Printf("Hash Match:   ✅ Verified\n")
			} else {
				fmt.Printf("Hash Match:   ❌ MISMATCH - Possible tampering detected!\n")
			}
		}
	}

	if result.MetadataHash != "" {
		fmt.Printf("Metadata Hash: %s\n", result.MetadataHash)
	}

	// Show errors if any
	if len(result.Errors) > 0 {
		fmt.Printf("\n🚫 Errors\n")
		fmt.Printf("───────────────────────────────────────────────────────────────────────────────\n")
		for _, err := range result.Errors {
			fmt.Printf("• %s\n", err)
		}
	}

	return nil
}

func generateProof(nftPath string, result *VerificationResult) error {
	fmt.Printf("📝 Generating proof document...\n")

	proof := map[string]interface{}{
		"nft_name":            result.NFTName,
		"mint_address":        "", // TODO: Extract from metadata or parameter
		"verified_by":         fmt.Sprintf("SolVault %s", Version),
		"verified_at":         result.VerifiedAt.Format(time.RFC3339),
		"image_hash":          result.ImageHash,
		"metadata_hash":       result.MetadataHash,
		"status":              result.Status,
		"hash_match":          result.HashMatch,
		"verification_method": "local_sha256",
	}

	// Add error information if present
	if len(result.Errors) > 0 {
		proof["errors"] = result.Errors
	}

	// Write proof file
	proofPath := filepath.Join(nftPath, "proof.json")
	proofData, err := json.MarshalIndent(proof, "", "  ")
	if err != nil {
		return fmt.Errorf("failed to marshal proof data: %w", err)
	}

	if err := os.WriteFile(proofPath, proofData, 0644); err != nil {
		return fmt.Errorf("failed to write proof file: %w", err)
	}

	fmt.Printf("✅ Proof saved to: %s\n", proofPath)
	return nil
}

func publishProof(nftPath string, result *VerificationResult) error {
	fmt.Printf("🌐 Publishing proof...\n")

	// TODO: Implement actual proof publishing
	// This would upload the proof.json and potentially the image to a web endpoint
	// and return a shareable URL

	fmt.Printf("⚠️  Proof publishing not yet implemented\n")
	fmt.Printf("   Proof file available locally at: %s/proof.json\n", nftPath)

	return nil
}

func init() {
	rootCmd.AddCommand(verifyCmd)

	verifyCmd.Flags().BoolVar(&publish, "publish", false, "publish proof to web endpoint")
	verifyCmd.Flags().BoolVar(&forceRecompute, "force-recompute", false, "recompute and update stored hashes")
	verifyCmd.Flags().BoolVar(&skipOnChain, "skip-onchain", false, "skip on-chain verification (local only)")
}
