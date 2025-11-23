// NFT Fetcher Service - fetches NFT metadata from Solana
import { Connection, PublicKey } from '@solana/web3.js';
import fetch from 'node-fetch';
import { logger } from '../utils/logger';

export interface NFTMetadata {
  mint: string;
  name?: string;
  symbol?: string;
  uri?: string;
  image?: string;
}

// Default RPC endpoint - can be overridden via function parameter or environment variable
const DEFAULT_RPC_ENDPOINT = process.env.SOLANA_RPC_URL || 'https://api.mainnet-beta.solana.com';

export async function fetchNFTsForOwner(
  ownerPubkey: string,
  rpcUrl?: string
): Promise<NFTMetadata[]> {
  const endpoint = rpcUrl || DEFAULT_RPC_ENDPOINT;
  const connection = new Connection(endpoint, 'confirmed');
  
  try {
    logger.info('Fetching NFTs for owner', { owner: ownerPubkey, rpc: endpoint });
    
    const ownerKey = new PublicKey(ownerPubkey);
    const tokenAccounts = await connection.getParsedTokenAccountsByOwner(ownerKey, {
      programId: new PublicKey('TokenkegQfeZyiNwAJbNbGKPFXCWuBvf9Ss623VQ5DA')
    });

    const nfts: NFTMetadata[] = [];
    
    for (const account of tokenAccounts.value) {
      const amount = account.account.data.parsed.info.tokenAmount.uiAmount;
      if (amount === 1) {
        const mintAddress = account.account.data.parsed.info.mint;
        nfts.push({
          mint: mintAddress,
          name: `NFT-${mintAddress.slice(0, 8)}`,
          uri: undefined,
          image: undefined
        });
      }
    }
    
    logger.info(`Found ${nfts.length} NFTs`, { count: nfts.length });
    return nfts;
  } catch (error) {
    logger.error('Failed to fetch NFTs', { error: String(error) });
    throw error;
  }
}

export async function fetchMetadataJson(uri: string): Promise<any> {
  try {
    const response = await fetch(uri);
    if (!response.ok) {
      throw new Error(`HTTP ${response.status}`);
    }
    return await response.json();
  } catch (error) {
    logger.error('Failed to fetch metadata JSON', { uri, error: String(error) });
    return null;
  }
}
