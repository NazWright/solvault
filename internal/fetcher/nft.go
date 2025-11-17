package fetcher

import (
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
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
}

// Fetcher handles fetching NFT metadata from various sources
type Fetcher struct {
	client     *solana.Client
	httpClient *http.Client
}

// NewFetcher creates a new NFT metadata fetcher
func NewFetcher(client *solana.Client) *Fetcher {
	return &Fetcher{
		client: client,
		httpClient: &http.Client{
			Timeout: 30 * time.Second,
		},
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
		info.Supply = 1   // Most NFTs have supply of 1
		info.Decimals = 0 // Most NFTs have 0 decimals
	}

	// Find token accounts for this mint owned by our wallet
	tokenAccounts, err := f.client.GetTokenAccountsByOwner(ctx)
	if err != nil {
		return nil, fmt.Errorf("failed to get token accounts: %w", err)
	}

	var tokenAccount *rpc.TokenAccount
	for _, account := range tokenAccounts {
		if account.Account.Data.GetParsed() != nil {
			parsed := account.Account.Data.GetParsed()
			if tokenInfo, ok := parsed["info"].(map[string]interface{}); ok {
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
	// This is a simplified parser. In practice, you'd use proper borsh deserialization
	// or the Metaplex Go SDK for parsing metadata accounts

	if len(data) < 100 {
		return "", fmt.Errorf("metadata account data too short")
	}

	// Skip to URI section (this is a rough approximation)
	// Real implementation would properly deserialize the metadata struct
	for i := 50; i < len(data)-4; i++ {
		if i+4 < len(data) {
			// Look for URI length prefix and extract URI
			// This is simplified - actual implementation would follow Metaplex format
			if data[i] == 0 && data[i+1] == 0 && data[i+2] > 0 && data[i+2] < 200 {
				length := int(data[i+2])
				if i+3+length < len(data) {
					uri := string(data[i+3 : i+3+length])
					if len(uri) > 10 && (uri[:4] == "http" || uri[:2] == "ar://") {
						return uri, nil
					}
				}
			}
		}
	}

	return "", fmt.Errorf("URI not found in metadata")
}

// fetchOffChainMetadata retrieves and parses metadata from a URI
func (f *Fetcher) fetchOffChainMetadata(ctx context.Context, uri string) (*NFTMetadata, error) {
	req, err := http.NewRequestWithContext(ctx, "GET", uri, nil)
	if err != nil {
		return nil, fmt.Errorf("failed to create request: %w", err)
	}

	resp, err := f.httpClient.Do(req)
	if err != nil {
		return nil, fmt.Errorf("failed to fetch metadata: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("HTTP error %d fetching metadata", resp.StatusCode)
	}

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("failed to read response body: %w", err)
	}

	var metadata NFTMetadata
	if err := json.Unmarshal(body, &metadata); err != nil {
		return nil, fmt.Errorf("failed to parse metadata JSON: %w", err)
	}

	return &metadata, nil
}

// Close cleans up the fetcher resources
func (f *Fetcher) Close() error {
	f.httpClient.CloseIdleConnections()
	return nil
}
