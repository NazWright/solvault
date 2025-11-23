import path from 'path';
import { fetchNFTsForOwner } from '../services/nftFetcher';
import { writeNFTBackup } from '../services/backupWriter';
import { logger } from '../utils/logger';

export type BackupSummary = {
  owner: string;
  outDir: string;
  totalFound: number;
  backedUp: number;
  items: Array<{ mint: string; tokenAccount: string; metadataUri?: string | null; imageSaved?: string | null }>;
};

export async function runBackup(ownerAddress: string, outDirRoot?: string, rpcUrl?: string) : Promise<BackupSummary> {
  const outDirBase = outDirRoot || path.join(process.cwd(), 'solvault_backups');
  logger.info('Backup started', { owner: ownerAddress, outDir: outDirBase, rpc: rpcUrl });

  const nfts = await fetchNFTsForOwner(ownerAddress, rpcUrl);
  const items: BackupSummary['items'] = [];
  let backedUp = 0;

  // For now run sequentially to keep logs simple — can add concurrency controls later
  for (const nft of nfts) {
    try {
      const nftOutDir = path.join(outDirBase, ownerAddress);
      const res = await writeNFTBackup(nftOutDir, nft);
      items.push({ mint: nft.mint, tokenAccount: nft.tokenAccount, metadataUri: nft.metadataUri ?? null, imageSaved: res.imageSaved ?? null });
      backedUp += 1;
    } catch (err) {
      logger.error('Failed to write backup for NFT', { mint: nft.mint, error: String(err) });
    }
  }

  const summary: BackupSummary = {
    owner: ownerAddress,
    outDir: outDirBase,
    totalFound: nfts.length,
    backedUp,
    items,
  };

  logger.info('Backup completed', summary);
  return summary;
}