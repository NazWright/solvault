# SolVault CLI Demo Output

## Test Run Information

**Date:** 2025-11-23
**Command:** `npm run backup -- <WALLET_ADDRESS> ./demo-backups`
**Environment:** Testing environment with limited RPC access

## Expected Output

When running the backup command with a valid wallet address and RPC access, the output should look like this:

```bash
$ npm run backup -- DRpbCBMxVnDK7maPM5tGv6MvB3v1sRMC86PZ8okm21hy ./demo-backups

> solvault@0.1.0 backup
> ts-node cli/index.ts backup DRpbCBMxVnDK7maPM5tGv6MvB3v1sRMC86PZ8okm21hy ./demo-backups

[INFO] Starting backup {"owner":"DRpbCBMxVnDK7maPM5tGv6MvB3v1sRMC86PZ8okm21hy","outDir":"./demo-backups"}
[INFO] Fetching NFTs for owner {"owner":"DRpbCBMxVnDK7maPM5tGv6MvB3v1sRMC86PZ8okm21hy","rpc":"https://api.mainnet-beta.solana.com"}
[INFO] Found 5 NFTs {"count":5}
[INFO] Saved metadata {"mint":"ABC123...","path":"./demo-backups/ABC123.../metadata.json"}
[INFO] Saved image {"mint":"ABC123...","path":"./demo-backups/ABC123.../image.png"}
[INFO] Saved metadata {"mint":"DEF456...","path":"./demo-backups/DEF456.../metadata.json"}
[INFO] Saved image {"mint":"DEF456...","path":"./demo-backups/DEF456.../image.png"}
[INFO] Backup complete {"totalFound":5,"backedUp":5}

--- SolVault Backup Summary ---
Owner: DRpbCBMxVnDK7maPM5tGv6MvB3v1sRMC86PZ8okm21hy
OutDir: ./demo-backups
Total NFTs Found: 5
Backed Up: 5
Items:
 - ABC123def456... (metadataUri: https://arweave.net/...) imageSaved: ./demo-backups/ABC123def456.../image.png
 - DEF456ghi789... (metadataUri: https://arweave.net/...) imageSaved: ./demo-backups/DEF456ghi789.../image.png
 - GHI789jkl012... (metadataUri: https://arweave.net/...) imageSaved: ./demo-backups/GHI789jkl012.../image.png
 - JKL012mno345... (metadataUri: https://arweave.net/...) imageSaved: ./demo-backups/JKL012mno345.../image.png
 - MNO345pqr678... (metadataUri: https://arweave.net/...) imageSaved: ./demo-backups/MNO345pqr678.../image.png
--- End ---
```

## Backup Directory Structure

After a successful backup, the directory structure would look like:

```
demo-backups/
├── ABC123def456.../
│   ├── metadata.json
│   └── image.png
├── DEF456ghi789.../
│   ├── metadata.json
│   └── image.png
├── GHI789jkl012.../
│   ├── metadata.json
│   └── image.png
├── JKL012mno345.../
│   ├── metadata.json
│   └── image.png
└── MNO345pqr678.../
    ├── metadata.json
    └── image.png
```

## Sample metadata.json

```json
{
  "name": "DeGod #1234",
  "mint": "ABC123def456...",
  "symbol": "DEGOD",
  "image": "https://metadata.degods.com/images/1234.png",
  "attributes": [
    {
      "trait_type": "Background",
      "value": "Blue"
    },
    {
      "trait_type": "Skin",
      "value": "Purple"
    }
  ]
}
```

## Running the Demo Locally

To run this demo with your own wallet:

1. Install dependencies:
   ```bash
   npm install
   ```

2. Run the backup command:
   ```bash
   npm run backup -- <YOUR_WALLET_ADDRESS> ./demo-backups
   ```

3. Optional: Use a custom RPC endpoint:
   ```bash
   npm run backup -- <YOUR_WALLET_ADDRESS> ./demo-backups https://your-rpc-endpoint.com
   ```

## Notes

- The demo requires network access to Solana RPC endpoints
- Rate limiting may occur with public RPC endpoints
- For production use, consider using a paid RPC service (Helius, QuickNode, etc.)
- This environment has limited external network access, so live demos may not work

## CLI Help

```bash
$ npm start

> solvault@0.1.0 start
> ts-node cli/index.ts

solvault CLI
Usage: solvault <command> [args]
Commands: backup
```

## Error Handling

If the wallet has no NFTs or if there's an RPC connection error, the CLI will display appropriate error messages:

```bash
[ERROR] Failed to fetch NFTs {"error":"TypeError: fetch failed"}
[ERROR] Backup failed {"error":"TypeError: fetch failed"}
```

This is expected in environments with limited network access or when RPC endpoints are unavailable.
