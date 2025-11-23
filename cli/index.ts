#!/usr/bin/env node
// Small CLI dispatcher. Add new commands under cli/commands.
import { cliBackup } from './commands/backup';

async function main() {
  const argv = process.argv.slice(2);
  const cmd = argv[0];
  if (!cmd) {
    console.log('solvault CLI');
    console.log('Usage: solvault <command> [args]');
    console.log('Commands: backup');
    process.exit(0);
  }

  if (cmd === 'backup') {
    await cliBackup(argv.slice(1));
    return;
  }

  console.log(`Unknown command: ${cmd}`);
  process.exit(1);
}

main().catch((err) => {
  console.error('Unhandled CLI error', err);
  process.exit(1);
});