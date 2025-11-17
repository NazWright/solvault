package solana

import (
	"fmt"
	"os"
	"strconv"
	"time"

	"github.com/gagliardetto/solana-go"
	"github.com/joho/godotenv"
)

// Config holds all Solana-related configuration
type Config struct {
	RPCURL          string
	WebSocketURL    string
	WalletAddress   solana.PublicKey
	PollInterval    time.Duration
	MaxRetries      int
	TimeoutSeconds  int
	BackupDirectory string
	PublishEndpoint string
	PublishAPIKey   string
}

// LoadConfig loads configuration from environment variables
func LoadConfig() (*Config, error) {
	// Try to load .env file, but don't fail if it doesn't exist
	_ = godotenv.Load()

	config := &Config{}

	// Required fields
	config.RPCURL = os.Getenv("SOLANA_RPC_URL")
	if config.RPCURL == "" {
		return nil, fmt.Errorf("SOLANA_RPC_URL environment variable is required")
	}

	config.WebSocketURL = os.Getenv("SOLANA_WEBSOCKET_URL")
	if config.WebSocketURL == "" {
		return nil, fmt.Errorf("SOLANA_WEBSOCKET_URL environment variable is required")
	}

	walletAddr := os.Getenv("WALLET_ADDRESS")
	if walletAddr == "" || walletAddr == "your_wallet_address_here" {
		return nil, fmt.Errorf("WALLET_ADDRESS environment variable is required and must be set to a valid Solana address")
	}

	var err error
	config.WalletAddress, err = solana.PublicKeyFromBase58(walletAddr)
	if err != nil {
		return nil, fmt.Errorf("invalid wallet address format: %w", err)
	}

	config.BackupDirectory = os.Getenv("BACKUP_DIRECTORY")
	if config.BackupDirectory == "" {
		homeDir, err := os.UserHomeDir()
		if err != nil {
			return nil, fmt.Errorf("failed to get home directory: %w", err)
		}
		config.BackupDirectory = fmt.Sprintf("%s/SolVaultBackups", homeDir)
	}

	// Optional fields with defaults
	config.PublishEndpoint = os.Getenv("PUBLISH_ENDPOINT")
	config.PublishAPIKey = os.Getenv("PUBLISH_API_KEY")

	// Parse numeric fields with defaults
	pollInterval := os.Getenv("POLL_INTERVAL_SECONDS")
	if pollInterval == "" {
		config.PollInterval = 30 * time.Second
	} else {
		seconds, err := strconv.Atoi(pollInterval)
		if err != nil {
			return nil, fmt.Errorf("invalid POLL_INTERVAL_SECONDS: %w", err)
		}
		config.PollInterval = time.Duration(seconds) * time.Second
	}

	maxRetries := os.Getenv("MAX_RETRIES")
	if maxRetries == "" {
		config.MaxRetries = 3
	} else {
		config.MaxRetries, err = strconv.Atoi(maxRetries)
		if err != nil {
			return nil, fmt.Errorf("invalid MAX_RETRIES: %w", err)
		}
	}

	timeoutSeconds := os.Getenv("TIMEOUT_SECONDS")
	if timeoutSeconds == "" {
		config.TimeoutSeconds = 60
	} else {
		config.TimeoutSeconds, err = strconv.Atoi(timeoutSeconds)
		if err != nil {
			return nil, fmt.Errorf("invalid TIMEOUT_SECONDS: %w", err)
		}
	}

	return config, nil
}

// Validate checks if the configuration is valid
func (c *Config) Validate() error {
	if c.RPCURL == "" {
		return fmt.Errorf("RPC URL is required")
	}

	if c.WalletAddress.IsZero() {
		return fmt.Errorf("wallet address is required")
	}

	if c.PollInterval <= 0 {
		return fmt.Errorf("poll interval must be positive")
	}

	if c.MaxRetries < 0 {
		return fmt.Errorf("max retries cannot be negative")
	}

	if c.TimeoutSeconds <= 0 {
		return fmt.Errorf("timeout seconds must be positive")
	}

	return nil
}
