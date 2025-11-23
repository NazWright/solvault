// Backup Handler - orchestrates NFT fetching and backup writing
import { fetchNFTsForOwner, fetchMetadataJson } from '../services/nftFetcher';
import { writeBackup, BackupItem } from '../services/backupWriter';
import { logger } from '../utils/logger';

export interface BackupSummary {
  owner: string;
  outDir: string;
  totalFound: number;
  backedUp: number;
  items: BackupItem[];
}

export async function runBackup(
  ownerPubkey: string,
  outDir?: string,
  rpcUrl?: string
): Promise<BackupSummary> {
  logger.info('Starting backup', { owner: ownerPubkey, outDir, rpcUrl });
  
  const nfts = await fetchNFTsForOwner(ownerPubkey, rpcUrl);
  const items: BackupItem[] = [];
  
  for (const nft of nfts) {
    try {
      let metadata = { name: nft.name, mint: nft.mint };
      
      if (nft.uri) {
        const fetchedMetadata = await fetchMetadataJson(nft.uri);
        if (fetchedMetadata) {
          metadata = { ...metadata, ...fetchedMetadata };
        }
      }
      
      const item = await writeBackup(nft.mint, nft.uri, metadata, outDir || './backups');
      items.push(item);
    } catch (err) {
      logger.error('Failed to backup NFT', { mint: nft.mint, error: String(err) });
      items.push({
        mint: nft.mint,
        metadataUri: nft.uri,
        imageSaved: undefined,
        metadata: undefined
      });
    }
  }
  
  const summary: BackupSummary = {
    owner: ownerPubkey,
    outDir: outDir || './backups',
    totalFound: nfts.length,
    backedUp: items.filter(i => i.metadata).length,
    items
  };
  
  logger.info('Backup complete', { 
    totalFound: summary.totalFound, 
    backedUp: summary.backedUp 
  });
  
  return summary;
}
