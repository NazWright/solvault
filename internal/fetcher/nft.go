package fetcher

import (
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"strings"
	"time"

	"github.com/NazWright/solvault/internal/solana"
	solanago "github.com/gagliardetto/solana-go"
	"github.com/gagliardetto/solana-go/rpc"
)

// NFTMetadata represents the standard Metaplex NFT metadata structure
type NFTMetadata struct {
	Name                 string      `json:"name"`
	Symbol               string      `json:"symbol"`
	Description          string      `json:"description"`
	Image                string      `json:"image"`
	ExternalURL          string      `json:"external_url,omitempty"`
	AnimationURL         string      `json:"animation_url,omitempty"`
	Attributes           []Attribute `json:"attributes,omitempty"`
	Properties           Properties  `json:"properties,omitempty"`
	SellerFeeBasisPoints int         `json:"seller_fee_basis_points,omitempty"`
	Collection           Collection  `json:"collection,omitempty"`
}

// Attribute represents NFT trait attributes
type Attribute struct {
	TraitType string      `json:"trait_type"`
	Value     interface{} `json:"value"`
}

// Properties contains additional NFT properties
type Properties struct {
	Files     []File                 `json:"files,omitempty"`
	Category  string                 `json:"category,omitempty"`
	Creators  []Creator              `json:"creators,omitempty"`
	MaxSupply int                    `json:"maxSupply,omitempty"`
	Uses      map[string]interface{} `json:"uses,omitempty"`
}

// File represents a file associated with the NFT
type File struct {
	URI  string `json:"uri"`
	Type string `json:"type"`
}

// Creator represents an NFT creator
type Creator struct {
	Address  string `json:"address"`
	Share    int    `json:"share"`
	Verified bool   `json:"verified"`
}

// Collection represents NFT collection information
type Collection struct {
	Name   string `json:"name"`
	Family string `json:"family"`
}

// NFTInfo contains comprehensive information about an NFT
type NFTInfo struct {
	MintAddress  solanago.PublicKey `json:"mint_address"`
	TokenAccount solanago.PublicKey `json:"token_account"`
	Owner        solanago.PublicKey `json:"owner"`
	Metadata     *NFTMetadata       `json:"metadata"`
	MetadataURI  string             `json:"metadata_uri"`
	OnChainData  interface{}        `json:"on_chain_data"`
	FetchedAt    time.Time          `json:"fetched_at"`
	Supply       uint64             `json:"supply"`
	Decimals     uint8              `json:"decimals"`
	MediaFiles   []*MediaFile       `json:"media_files,omitempty"` // Downloaded media files
}

// Fetcher handles fetching NFT metadata from various sources
type Fetcher struct {
	client          *solana.Client
	httpClient      *http.Client
	mediaDownloader *MediaDownloader
}

// NewFetcher creates a new NFT metadata fetcher
func NewFetcher(client *solana.Client) *Fetcher {
	return &Fetcher{
		client: client,
		httpClient: &http.Client{
			Timeout: 30 * time.Second,
		},
		mediaDownloader: NewMediaDownloader(),
	}
}

// FetchNFTInfo retrieves comprehensive NFT information including metadata
func (f *Fetcher) FetchNFTInfo(ctx context.Context, mintAddress solanago.PublicKey) (*NFTInfo, error) {
	info := &NFTInfo{
		MintAddress: mintAddress,
		FetchedAt:   time.Now(),
	}

	// Get mint account info
	mintAccount, err := f.client.GetAccountInfo(ctx, mintAddress)
	if err != nil {
		return nil, fmt.Errorf("failed to get mint account info: %w", err)
	}

	// Parse mint data to get supply and decimals
	if len(mintAccount.Data.GetBinary()) >= 44 {
		// Basic mint account structure parsing
		// This is a simplified version - in production you'd want proper mint account parsing
		data := mintAccount.Data.GetBinary()

		// Extract decimals from mint account (byte 44)
		if len(data) > 44 {
			info.Decimals = data[44]
		}

		// For now, assume supply of 1 for NFTs - in production you'd properly parse this
		info.Supply = 1

		// Validate this looks like an NFT (0 decimals is a strong indicator)
		if info.Decimals != 0 {
			return nil, fmt.Errorf("this token has %d decimals - NFTs should have 0 decimals", info.Decimals)
		}
	}

	// Find token accounts for this mint owned by our wallet
	tokenAccounts, err := f.client.GetTokenAccountsByOwner(ctx)
	if err != nil {
		return nil, fmt.Errorf("failed to get token accounts: %w", err)
	}

	var tokenAccount *rpc.TokenAccount
	for _, account := range tokenAccounts {
		// Check if we have parsed data
		rawJSON := account.Account.Data.GetRawJSON()
		if len(rawJSON) > 0 {
			var parsed map[string]interface{}
			if err := json.Unmarshal(rawJSON, &parsed); err == nil {
				// Check if data is under "parsed" key
				var tokenInfo map[string]interface{}
				var ok bool

				if parsedData, exists := parsed["parsed"].(map[string]interface{}); exists {
					tokenInfo, ok = parsedData["info"].(map[string]interface{})
				} else {
					tokenInfo, ok = parsed["info"].(map[string]interface{})
				}

				if ok {
					if mint, ok := tokenInfo["mint"].(string); ok {
						mintPubkey, err := solanago.PublicKeyFromBase58(mint)
						if err == nil && mintPubkey.Equals(mintAddress) {
							tokenAccount = account
							info.TokenAccount = account.Pubkey
							info.Owner = f.client.Config().WalletAddress
							break
						}
					}
				}
			}
		}
	}

	if tokenAccount == nil {
		return nil, fmt.Errorf("token account not found for mint %s", mintAddress.String())
	}

	// Try to find and fetch metadata
	metadataURI, err := f.findMetadataURI(ctx, mintAddress)
	if err != nil {
		// Log warning but continue - some NFTs might not have standard metadata
		fmt.Printf("⚠️  Could not find metadata URI for %s: %v\n", mintAddress.String(), err)
	} else if metadataURI != "" {
		info.MetadataURI = metadataURI
		metadata, err := f.fetchOffChainMetadata(ctx, metadataURI)
		if err != nil {
			fmt.Printf("⚠️  Could not fetch off-chain metadata: %v\n", err)
		} else {
			info.Metadata = metadata
		}
	}

	return info, nil
}

// findMetadataURI attempts to find the metadata URI for an NFT
func (f *Fetcher) findMetadataURI(ctx context.Context, mintAddress solanago.PublicKey) (string, error) {
	// This is a simplified approach. In a full implementation, you would:
	// 1. Derive the metadata account address using Metaplex program
	// 2. Fetch the metadata account data
	// 3. Parse the metadata account to extract the URI

	// For now, we'll use a placeholder approach that checks common patterns
	// In practice, you'd want to implement proper Metaplex metadata parsing

	// Try to find metadata account (simplified version)
	// The actual implementation would use proper PDA derivation
	metadataPubkey, err := f.deriveMetadataAddress(mintAddress)
	if err != nil {
		return "", fmt.Errorf("failed to derive metadata address: %w", err)
	}

	account, err := f.client.GetAccountInfo(ctx, metadataPubkey)
	if err != nil {
		return "", fmt.Errorf("metadata account not found: %w", err)
	}

	// Parse metadata account data (simplified)
	// In practice, you'd use proper Metaplex metadata deserialization
	uri, err := f.parseMetadataURI(account.Data.GetBinary())
	if err != nil {
		return "", fmt.Errorf("failed to parse metadata URI: %w", err)
	}

	return uri, nil
}

// deriveMetadataAddress derives the metadata account address for a mint
func (f *Fetcher) deriveMetadataAddress(mintAddress solanago.PublicKey) (solanago.PublicKey, error) {
	// Metaplex metadata program ID
	metaplexProgramID := solanago.MustPublicKeyFromBase58("metaqbxxUerdq28cj1RbAWkYQm3ybzjb6a8bt518x1s")

	seeds := [][]byte{
		[]byte("metadata"),
		metaplexProgramID.Bytes(),
		mintAddress.Bytes(),
	}

	pda, _, err := solanago.FindProgramAddress(seeds, metaplexProgramID)
	if err != nil {
		return solanago.PublicKey{}, fmt.Errorf("failed to find metadata PDA: %w", err)
	}

	return pda, nil
} // parseMetadataURI extracts the metadata URI from metadata account data
func (f *Fetcher) parseMetadataURI(data []byte) (string, error) {
	// Enhanced parser for Metaplex metadata accounts
	// Based on the Metaplex Token Metadata standard

	if len(data) < 100 {
		return "", fmt.Errorf("metadata account data too short: %d bytes", len(data))
	}

	fmt.Println("\n🔬 Analyzing Metaplex Metadata Account:")
	fmt.Printf("   📊 Size: %d bytes\n", len(data))
	fmt.Printf("   � Account Key: %d", data[0])

	if data[0] == 4 {
		fmt.Println(" ✅ (Valid Metadata Account)")
	} else {
		fmt.Printf(" ❌ (Expected 4, got %d)\n", data[0])
		return "", fmt.Errorf("not a valid metadata account (key = %d, expected 4)", data[0])
	}

	// Skip update authority (32 bytes) and mint (32 bytes)
	offset := 65

	if offset+4 > len(data) {
		return "", fmt.Errorf("data too short for name length")
	}

	// Read name length (little endian u32)
	nameLength := uint32(data[offset]) | uint32(data[offset+1])<<8 |
		uint32(data[offset+2])<<16 | uint32(data[offset+3])<<24
	offset += 4

	if nameLength > 200 {
		return "", fmt.Errorf("name length too large: %d", nameLength)
	}

	// Skip name
	if offset+int(nameLength) > len(data) {
		return "", fmt.Errorf("data too short for name")
	}
	name := string(data[offset : offset+int(nameLength)])
	fmt.Printf("   🏷️  Name: '%s'\n", name)
	offset += int(nameLength)

	// Read symbol length
	if offset+4 > len(data) {
		return "", fmt.Errorf("data too short for symbol length")
	}
	symbolLength := uint32(data[offset]) | uint32(data[offset+1])<<8 |
		uint32(data[offset+2])<<16 | uint32(data[offset+3])<<24
	offset += 4

	if symbolLength > 200 {
		return "", fmt.Errorf("symbol length too large: %d", symbolLength)
	}

	// Skip symbol
	if offset+int(symbolLength) > len(data) {
		return "", fmt.Errorf("data too short for symbol")
	}
	symbol := string(data[offset : offset+int(symbolLength)])
	fmt.Printf("   🔖 Symbol: '%s'\n", symbol)
	offset += int(symbolLength)

	// Read URI length
	if offset+4 > len(data) {
		return "", fmt.Errorf("data too short for URI length")
	}
	uriLength := uint32(data[offset]) | uint32(data[offset+1])<<8 |
		uint32(data[offset+2])<<16 | uint32(data[offset+3])<<24
	offset += 4

	if uriLength > 1000 {
		return "", fmt.Errorf("URI length too large: %d", uriLength)
	}

	// Extract URI
	if offset+int(uriLength) > len(data) {
		return "", fmt.Errorf("data too short for URI")
	}

	uri := string(data[offset : offset+int(uriLength)])

	// Remove null bytes and whitespace padding (common in Metaplex metadata)
	uri = strings.TrimRight(uri, "\x00")
	uri = strings.TrimSpace(uri)

	displayURI := uri
	if len(uri) > 60 {
		displayURI = uri[:57] + "..."
	}
	fmt.Printf("   🌐 Metadata URI: %s\n", displayURI)
	fmt.Println("   ✅ Metadata parsing complete!")

	// Validate URI format
	if len(uri) < 5 {
		return "", fmt.Errorf("URI too short: '%s'", uri)
	}

	// Check for common URI prefixes
	if uri[:4] == "http" || uri[:2] == "ar" || uri[:4] == "ipfs" {
		return uri, nil
	}

	return "", fmt.Errorf("URI format not recognized: '%s'", uri)
}

// fetchOffChainMetadata retrieves and parses metadata from a URI (Arweave, IPFS, HTTP)
func (f *Fetcher) fetchOffChainMetadata(ctx context.Context, uri string) (*NFTMetadata, error) {
	fmt.Printf("   📡 Fetching off-chain metadata from: %s\n", f.getTruncatedURI(uri))

	req, err := http.NewRequestWithContext(ctx, "GET", uri, nil)
	if err != nil {
		return nil, fmt.Errorf("failed to create request: %w", err)
	}

	// Add headers for better compatibility with Arweave and IPFS gateways
	req.Header.Set("User-Agent", "SolVault/1.0 NFT-Backup-Tool")
	req.Header.Set("Accept", "application/json, text/plain, */*")

	resp, err := f.httpClient.Do(req)
	if err != nil {
		return nil, fmt.Errorf("failed to fetch metadata: %w", err)
	}
	defer resp.Body.Close()

	fmt.Printf("   📊 Response: %d %s\n", resp.StatusCode, resp.Status)

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("HTTP error %d fetching metadata", resp.StatusCode)
	}

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("failed to read response body: %w", err)
	}

	fmt.Printf("   📄 Metadata size: %d bytes\n", len(body))

	// Try to parse as standard NFT metadata first
	var metadata NFTMetadata
	if err := json.Unmarshal(body, &metadata); err != nil {
		// If standard parsing fails, try flexible parsing
		fmt.Printf("   🔧 Standard parsing failed, trying flexible parsing...\n")

		flexibleMetadata, flexErr := f.parseFlexibleMetadata(body)
		if flexErr != nil {
			return nil, fmt.Errorf("failed to parse metadata JSON (standard: %v, flexible: %v)", err, flexErr)
		}
		metadata = *flexibleMetadata
	}

	fmt.Printf("   ✅ Successfully parsed metadata for: '%s'\n", metadata.Name)
	return &metadata, nil
}

// getTruncatedURI returns a truncated version of URI for display
func (f *Fetcher) getTruncatedURI(uri string) string {
	if len(uri) <= 60 {
		return uri
	}
	return uri[:57] + "..."
}

// parseFlexibleMetadata handles non-standard metadata formats common in older NFTs
func (f *Fetcher) parseFlexibleMetadata(body []byte) (*NFTMetadata, error) {
	// Parse into a generic map first
	var rawData map[string]interface{}
	if err := json.Unmarshal(body, &rawData); err != nil {
		return nil, fmt.Errorf("failed to parse as JSON object: %w", err)
	}

	metadata := &NFTMetadata{}

	// Extract fields with flexible typing
	if name, ok := rawData["name"].(string); ok {
		metadata.Name = name
	}

	if symbol, ok := rawData["symbol"].(string); ok {
		metadata.Symbol = symbol
	}

	if description, ok := rawData["description"].(string); ok {
		metadata.Description = description
	}

	if image, ok := rawData["image"].(string); ok {
		metadata.Image = image
	}

	if animationURL, ok := rawData["animation_url"].(string); ok {
		metadata.AnimationURL = animationURL
	}

	if externalURL, ok := rawData["external_url"].(string); ok {
		metadata.ExternalURL = externalURL
	}

	// Handle attributes array with flexible typing
	if attrs, ok := rawData["attributes"].([]interface{}); ok {
		for _, attr := range attrs {
			if attrMap, ok := attr.(map[string]interface{}); ok {
				attribute := Attribute{}

				if traitType, ok := attrMap["trait_type"].(string); ok {
					attribute.TraitType = traitType
				}

				// Handle value as any type (string, number, bool)
				if value, exists := attrMap["value"]; exists {
					attribute.Value = value
				}

				metadata.Attributes = append(metadata.Attributes, attribute)
			}
		}
	}

	// Handle properties with flexible creator verification (number vs bool)
	if props, ok := rawData["properties"].(map[string]interface{}); ok {
		metadata.Properties = Properties{}

		if creators, ok := props["creators"].([]interface{}); ok {
			for _, creator := range creators {
				if creatorMap, ok := creator.(map[string]interface{}); ok {
					c := Creator{}

					if address, ok := creatorMap["address"].(string); ok {
						c.Address = address
					}

					if share, ok := creatorMap["share"].(float64); ok {
						c.Share = int(share)
					}

					// Handle verified field as number or boolean
					if verified, ok := creatorMap["verified"]; ok {
						switch v := verified.(type) {
						case bool:
							c.Verified = v
						case float64:
							c.Verified = v != 0 // Convert number to boolean
						case string:
							c.Verified = v == "true" || v == "1"
						}
					}

					metadata.Properties.Creators = append(metadata.Properties.Creators, c)
				}
			}
		}

		if files, ok := props["files"].([]interface{}); ok {
			for _, file := range files {
				if fileMap, ok := file.(map[string]interface{}); ok {
					f := File{}

					if uri, ok := fileMap["uri"].(string); ok {
						f.URI = uri
					}

					if fileType, ok := fileMap["type"].(string); ok {
						f.Type = fileType
					}

					metadata.Properties.Files = append(metadata.Properties.Files, f)
				}
			}
		}

		if category, ok := props["category"].(string); ok {
			metadata.Properties.Category = category
		}
	}

	// Handle seller fee basis points
	if sellerFee, ok := rawData["seller_fee_basis_points"].(float64); ok {
		metadata.SellerFeeBasisPoints = int(sellerFee)
	}

	// Handle collection info
	if collection, ok := rawData["collection"].(map[string]interface{}); ok {
		if name, ok := collection["name"].(string); ok {
			metadata.Collection.Name = name
		}
		if family, ok := collection["family"].(string); ok {
			metadata.Collection.Family = family
		}
	}

	return metadata, nil
}

// FetchNFTInfoDemo fetches NFT information for demo purposes (without ownership check)
func (f *Fetcher) FetchNFTInfoDemo(ctx context.Context, mintAddress solanago.PublicKey) (*NFTInfo, error) {
	info := &NFTInfo{
		MintAddress: mintAddress,
		FetchedAt:   time.Now(),
	}

	// Get mint account info
	mintAccount, err := f.client.GetAccountInfo(ctx, mintAddress)
	if err != nil {
		return nil, fmt.Errorf("failed to get mint account info: %w", err)
	}

	// Parse mint data to get supply and decimals
	if len(mintAccount.Data.GetBinary()) >= 44 {
		data := mintAccount.Data.GetBinary()

		// Extract decimals from mint account (byte 44)
		if len(data) > 44 {
			info.Decimals = data[44]
		}

		// For now, assume supply of 1 for NFTs - in production you'd properly parse this
		info.Supply = 1

		// Validate this looks like an NFT (0 decimals is a strong indicator)
		if info.Decimals != 0 {
			return nil, fmt.Errorf("this token has %d decimals - NFTs should have 0 decimals", info.Decimals)
		}
	}

	// Set demo owner (we don't check actual ownership for demo)
	demoWallet, _ := solanago.PublicKeyFromBase58("11111111111111111111111111111112")
	info.Owner = demoWallet
	info.TokenAccount = demoWallet // Dummy token account for demo

	// Try to find and fetch metadata
	metadataURI, err := f.findMetadataURI(ctx, mintAddress)
	if err != nil {
		fmt.Printf("⚠️  Could not find metadata URI for %s: %v\n", mintAddress.String(), err)
	} else if metadataURI != "" {
		info.MetadataURI = metadataURI
		metadata, err := f.fetchOffChainMetadata(ctx, metadataURI)
		if err != nil {
			fmt.Printf("⚠️  Could not fetch off-chain metadata: %v\n", err)
		} else {
			info.Metadata = metadata
		}
	}

	return info, nil
}

// DownloadMediaFiles downloads all media files associated with an NFT
func (f *Fetcher) DownloadMediaFiles(ctx context.Context, nftInfo *NFTInfo, mediaDir string) error {
	if nftInfo.Metadata == nil {
		return nil // No metadata, no media to download
	}

	var mediaURLs []string

	// Collect media URLs from metadata
	if nftInfo.Metadata.Image != "" {
		mediaURLs = append(mediaURLs, nftInfo.Metadata.Image)
	}
	if nftInfo.Metadata.AnimationURL != "" {
		mediaURLs = append(mediaURLs, nftInfo.Metadata.AnimationURL)
	}

	// Collect URLs from properties.files array
	if nftInfo.Metadata.Properties.Files != nil {
		for _, file := range nftInfo.Metadata.Properties.Files {
			if file.URI != "" {
				mediaURLs = append(mediaURLs, file.URI)
			}
		}
	}

	// Download each media file
	for _, mediaURL := range mediaURLs {
		mediaFile, err := f.mediaDownloader.DownloadMedia(ctx, mediaURL, mediaDir)
		if err != nil {
			fmt.Printf("⚠️  Failed to download media %s: %v\n", mediaURL, err)
			continue // Skip failed downloads but continue with others
		}

		// Add to NFT info
		nftInfo.MediaFiles = append(nftInfo.MediaFiles, mediaFile)
		fmt.Printf("✅ Downloaded media: %s (%s, %d bytes)\n",
			mediaFile.Filename, mediaFile.MediaType, mediaFile.Size)
	}

	return nil
}

// Close cleans up the fetcher resources
func (f *Fetcher) Close() error {
	f.httpClient.CloseIdleConnections()
	if f.mediaDownloader != nil {
		f.mediaDownloader.Close()
	}
	return nil
}
