package solana

import (
	"context"
	"fmt"
	"time"

	"github.com/gagliardetto/solana-go"
	"github.com/gagliardetto/solana-go/rpc"
)

// Client wraps the Solana RPC client with our configuration
type Client struct {
	rpc    *rpc.Client
	config *Config
}

// NewClient creates a new Solana client with the given configuration
func NewClient(config *Config) (*Client, error) {
	if err := config.Validate(); err != nil {
		return nil, fmt.Errorf("invalid configuration: %w", err)
	}

	rpcClient := rpc.New(config.RPCURL)

	client := &Client{
		rpc:    rpcClient,
		config: config,
	}

	return client, nil
}

// TestConnection verifies that we can connect to the Solana RPC endpoint
func (c *Client) TestConnection(ctx context.Context) error {
	ctx, cancel := context.WithTimeout(ctx, time.Duration(c.config.TimeoutSeconds)*time.Second)
	defer cancel()

	_, err := c.rpc.GetVersion(ctx)
	if err != nil {
		return fmt.Errorf("failed to connect to Solana RPC: %w", err)
	}

	return nil
}

// GetTokenAccountsByOwner retrieves all token accounts owned by the configured wallet
func (c *Client) GetTokenAccountsByOwner(ctx context.Context) ([]*rpc.TokenAccount, error) {
	ctx, cancel := context.WithTimeout(ctx, time.Duration(c.config.TimeoutSeconds)*time.Second)
	defer cancel()

	// Get all token accounts for the wallet
	result, err := c.rpc.GetTokenAccountsByOwner(
		ctx,
		c.config.WalletAddress,
		&rpc.GetTokenAccountsConfig{
			ProgramId: &solana.TokenProgramID,
		},
		&rpc.GetTokenAccountsOpts{
			Encoding: solana.EncodingJSONParsed,
		},
	)
	if err != nil {
		return nil, fmt.Errorf("failed to get token accounts: %w", err)
	}

	return result.Value, nil
}

// GetAccountInfo retrieves account information for a given public key
func (c *Client) GetAccountInfo(ctx context.Context, pubkey solana.PublicKey) (*rpc.Account, error) {
	ctx, cancel := context.WithTimeout(ctx, time.Duration(c.config.TimeoutSeconds)*time.Second)
	defer cancel()

	result, err := c.rpc.GetAccountInfo(ctx, pubkey)
	if err != nil {
		return nil, fmt.Errorf("failed to get account info for %s: %w", pubkey.String(), err)
	}

	if result.Value == nil {
		return nil, fmt.Errorf("account not found: %s", pubkey.String())
	}

	return result.Value, nil
}

// GetTransaction retrieves transaction details by signature
func (c *Client) GetTransaction(ctx context.Context, signature solana.Signature) (*rpc.GetTransactionResult, error) {
	ctx, cancel := context.WithTimeout(ctx, time.Duration(c.config.TimeoutSeconds)*time.Second)
	defer cancel()

	result, err := c.rpc.GetTransaction(
		ctx,
		signature,
		&rpc.GetTransactionOpts{
			Encoding:   solana.EncodingJSONParsed,
			Commitment: rpc.CommitmentFinalized,
		},
	)
	if err != nil {
		return nil, fmt.Errorf("failed to get transaction %s: %w", signature.String(), err)
	}

	return result, nil
}

// GetSignaturesForAddress retrieves recent transaction signatures for an address
func (c *Client) GetSignaturesForAddress(ctx context.Context, address solana.PublicKey, limit int) ([]*rpc.TransactionSignature, error) {
	ctx, cancel := context.WithTimeout(ctx, time.Duration(c.config.TimeoutSeconds)*time.Second)
	defer cancel()

	limitUint := uint64(limit)
	result, err := c.rpc.GetConfirmedSignaturesForAddress2(
		ctx,
		address,
		&rpc.GetConfirmedSignaturesForAddress2Opts{
			Limit:      &limitUint,
			Commitment: rpc.CommitmentFinalized,
		},
	)
	if err != nil {
		return nil, fmt.Errorf("failed to get signatures for address %s: %w", address.String(), err)
	}

	return result, nil
}

// Config returns the client's configuration
func (c *Client) Config() *Config {
	return c.config
}

// Close cleans up the client resources
func (c *Client) Close() error {
	// The gagliardetto client doesn't require explicit closing
	return nil
}
