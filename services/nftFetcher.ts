import { Connection, PublicKey } from '@solana/web3.js';
import fetch from 'node-fetch';
import { logger } from '../utils/logger';

// NFT fetcher: MVP logic to find potential NFT mints for an owner and try to extract an on-chain metadata URI.
// Heuristic: parsed token accounts with amount === 1 & decimals === 0, read metadata PDA, extract URI substring.

const TOKEN_PROGRAM_ID = new PublicKey('TokenkegQfeZyiNwAJbNbGKPFXCWuBvf9Ss623VQ5DA');
const METADATA_PROGRAM_ID = new PublicKey('metaqbxxUerdq28cj1RbAWkYQm3ybzjb6a8bt518x1s');

export type FetchedNFT = {
  mint: string;
  tokenAccount: string;
  metadataUri?: string | null;
  metadataRaw?: string | null;
};

function extractUriFromAccountData(buffer: Buffer): string | null {
  try {
    const text = buffer.toString('utf8');
    const httpIndex = text.indexOf('http');
    if (httpIndex === -1) return null;
    const rest = text.slice(httpIndex);
    const nullIdx = rest.indexOf('\x00');
    return nullIdx === -1 ? rest.trim() : rest.slice(0, nullIdx).trim();
  } catch (err) {
    return null;
  }
}

export async function fetchNFTsForOwner(ownerAddress: string, rpcUrl = 'https://api.mainnet-beta.solana.com'): Promise<FetchedNFT[]> {
  const connection = new Connection(rpcUrl, { commitment: 'confirmed' });
  const ownerPubkey = new PublicKey(ownerAddress);
  logger.info('Starting token account fetch', { owner: ownerAddress, rpc: rpcUrl });

  const resp = await connection.getParsedTokenAccountsByOwner(ownerPubkey, { programId: TOKEN_PROGRAM_ID });

  const potentialNFTs: { mint: string; tokenAccount: string }[] = [];

  for (const { pubkey, account } of resp.value) {
    try {
      const parsed = (account.data as any).parsed;
      if (!parsed) continue;
      const info = parsed.info;
      const tokenAmount = info.tokenAmount;
      if (!tokenAmount) continue;
      if (tokenAmount.uiAmount === 1 && tokenAmount.decimals === 0) {
        potentialNFTs.push({ mint: info.mint, tokenAccount: pubkey.toBase58() });
      }
    } catch (err) {
      logger.debug('Skipping token account due to parse error', { pubkey: pubkey.toBase58(), error: String(err) });
    }
  }

  logger.info('Found potential NFT token accounts', { owner: ownerAddress, count: potentialNFTs.length });

  const results: FetchedNFT[] = [];

  for (const { mint, tokenAccount } of potentialNFTs) {
    let metadataUri: string | null = null;
    let metadataRaw: string | null = null;
    try {
      const mintKey = new PublicKey(mint);
      const [metadataPDA] = await PublicKey.findProgramAddress(
        [Buffer.from('metadata'), METADATA_PROGRAM_ID.toBuffer(), mintKey.toBuffer()],
        METADATA_PROGRAM_ID
      );
      const accountInfo = await connection.getAccountInfo(metadataPDA);
      if (accountInfo && accountInfo.data) {
        const buf = Buffer.from(accountInfo.data);
        metadataUri = extractUriFromAccountData(buf);
        metadataRaw = buf.toString('base64');
      }
    } catch (err) {
      logger.debug('Could not read metadata account', { mint, error: String(err) });
    }

    if (metadataUri) {
      try {
        const r = await fetch(metadataUri, { timeout: 8000 });
        if (!r.ok) {
          logger.debug('Metadata URI returned non-ok', { mint, uri: metadataUri, status: r.status });
        }
      } catch (err) {
        logger.debug('Failed to fetch metadata URI proactively', { mint, uri: metadataUri, error: String(err) });
      }
    }

    results.push({ mint, tokenAccount, metadataUri, metadataRaw });
  }

  return results;
}