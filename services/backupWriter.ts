// Backup Writer Service - writes NFT data to local filesystem
import * as fs from 'fs';
import * as path from 'path';
import fetch from 'node-fetch';
import { logger } from '../utils/logger';

export interface BackupItem {
  mint: string;
  metadataUri?: string;
  imageSaved?: string;
  metadata?: any;
}

export async function writeBackup(
  mint: string,
  metadataUri: string | undefined,
  metadata: any,
  outDir: string
): Promise<BackupItem> {
  const backupDir = path.resolve(outDir || './backups');
  
  // Ensure backup directory exists
  if (!fs.existsSync(backupDir)) {
    fs.mkdirSync(backupDir, { recursive: true });
  }
  
  const nftDir = path.join(backupDir, mint);
  if (!fs.existsSync(nftDir)) {
    fs.mkdirSync(nftDir, { recursive: true });
  }
  
  // Save metadata JSON
  const metadataPath = path.join(nftDir, 'metadata.json');
  fs.writeFileSync(metadataPath, JSON.stringify(metadata, null, 2));
  logger.info('Saved metadata', { mint, path: metadataPath });
  
  let imageSaved: string | undefined;
  
  // Download image if available
  if (metadata?.image) {
    try {
      const imageUrl = metadata.image;
      const response = await fetch(imageUrl);
      if (response.ok) {
        const buffer = await response.buffer();
        const ext = imageUrl.endsWith('.png') ? 'png' : 'jpg';
        const imagePath = path.join(nftDir, `image.${ext}`);
        fs.writeFileSync(imagePath, buffer);
        imageSaved = imagePath;
        logger.info('Saved image', { mint, path: imagePath });
      }
    } catch (err) {
      logger.warn('Failed to download image', { mint, error: String(err) });
    }
  }
  
  return {
    mint,
    metadataUri,
    imageSaved,
    metadata
  };
}
