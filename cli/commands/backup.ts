#!/usr/bin/env node
import { runBackup } from '../../handlers/backupHandler';
import { logger } from '../../utils/logger';

export async function cliBackup(args: string[]) {
  // Very small arg parser: node index.js <owner> [outDir] [rpcUrl]
  if (args.length < 1) {
    console.error('Usage: solvault backup <ownerWalletPubkey> [outDir] [rpcUrl]');
    process.exit(1);
  }
  const owner = args[0];
  const outDir = args[1];
  const rpcUrl = args[2];

  try {
    const summary = await runBackup(owner, outDir, rpcUrl);
    // human-friendly summary print
    console.log('--- SolVault Backup Summary ---');
    console.log(`Owner: ${summary.owner}`);
    console.log(`OutDir: ${summary.outDir}`);
    console.log(`Total NFTs Found: ${summary.totalFound}`);
    console.log(`Backed Up: ${summary.backedUp}`);
    console.log('Items:');
    for (const it of summary.items) {
      console.log(` - ${it.mint} (metadataUri: ${it.metadataUri ?? 'none'}) imageSaved: ${it.imageSaved ?? 'none'}`);
    }
    console.log('--- End ---');
    process.exit(0);
  } catch (err) {
    logger.error('Backup failed', { error: String(err) });
    process.exit(2);
  }
}
