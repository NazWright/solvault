import fs from 'fs';
import path from 'path';
import fetch from 'node-fetch';
import { logger } from '../utils/logger';

export type NFTBackup = {
  mint: string;
  tokenAccount: string;
  metadataUri?: string | null;
  metadataJson?: any | null;
  imageSaved?: string | null;
};

function ensureDir(dir: string) {
  if (!fs.existsSync(dir)) fs.mkdirSync(dir, { recursive: true });
}

export async function writeNFTBackup(outDir: string, nft: { mint: string; tokenAccount: string; metadataUri?: string | null }) : Promise<NFTBackup> {
  ensureDir(outDir);
  const nftDir = path.join(outDir, nft.mint);
  ensureDir(nftDir);

  const metaPath = path.join(nftDir, 'metadata.json');
  const metaBackup: NFTBackup = { mint: nft.mint, tokenAccount: nft.tokenAccount, metadataUri: nft.metadataUri || null };

  // fetch metadata JSON if URI exists
  if (nft.metadataUri) {
    try {
      const res = await fetch(nft.metadataUri, { timeout: 10000 });
      if (res.ok) {
        const contentType = res.headers.get('content-type') || '';
        if (contentType.includes('application/json') || contentType.includes('json')) {
          const json = await res.json();
          metaBackup.metadataJson = json;
          fs.writeFileSync(metaPath, JSON.stringify(json, null, 2), 'utf8');
          logger.info('Saved metadata JSON', { mint: nft.mint, path: metaPath });
          // attempt to download image if present
          const imgUrl = json.image || json.image_url || json.imageUrl;
          if (imgUrl && typeof imgUrl === 'string') {
            try {
              const imgRes = await fetch(imgUrl, { timeout: 15000 });
              if (imgRes.ok) {
                const ext = path.extname(new URL(imgUrl).pathname) || '.bin';
                const imgPath = path.join(nftDir, 'image' + ext);
                const buffer = Buffer.from(await imgRes.arrayBuffer());
                fs.writeFileSync(imgPath, buffer);
                metaBackup.imageSaved = imgPath;
                logger.info('Saved image file', { mint: nft.mint, image: imgPath });
              } else {
                logger.warn('Image URL fetch returned non-ok', { mint: nft.mint, url: imgUrl, status: imgRes.status });
              }
            } catch (err) {
              logger.warn('Failed to download image', { mint: nft.mint, url: imgUrl, error: String(err) });
            }
          }
        } else {
          // save the raw response
          const text = await res.text();
          fs.writeFileSync(metaPath, text, 'utf8');
          metaBackup.metadataJson = { raw: text };
          logger.info('Saved raw metadata response (non-json content-type)', { mint: nft.mint, path: metaPath });
        }
      } else {
        fs.writeFileSync(metaPath, JSON.stringify({ error: `fetch_status_${res.status}` }), 'utf8');
        logger.warn('Metadata fetch failed', { mint: nft.mint, status: res.status });
      }
    } catch (err) {
      fs.writeFileSync(metaPath, JSON.stringify({ error: String(err) }), 'utf8');
      logger.error('Error fetching metadata URI', { mint: nft.mint, error: String(err) });
    }
  } else {
    // No metadata URI — write minimal file indicating what we have
    fs.writeFileSync(metaPath, JSON.stringify({ note: 'no-metadata-uri-onchain' }, null, 2), 'utf8');
    logger.info('No on-chain metadata URI found', { mint: nft.mint });
  }

  return metaBackup;
}