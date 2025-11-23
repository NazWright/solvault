// Small utility to generate X/Twitter post text in two styles:
// - devStyle(): the original "Gm 𝕏 family" builder voice
// - cryptoStyle(): broader crypto/Twitter friendly copy tuned for collectors/DAOs
//
// Use generateBoth() to get a pair for publishing or PR descriptions (we won't include posts in PRs per instruction).

export type PostPair = {
  dev: string;
  crypto: string;
};

function devStyle(title: string, what: string, why: string, tieIn: string, cta?: string) {
  const ctaLine = cta ? `\n\n${cta}` : '';
  return [
    'Gm 𝕏 family ☀️',
    '',
    title,
    '',
    what,
    '',
    why,
    '',
    tieIn,
    ctaLine,
    '',
    '#Solana #BuildInPublic'
  ].join('\n');
}

function cryptoStyle(title: string, what: string, why: string, value: string, cta?: string) {
  const ctaLine = cta ? `\n\n${cta}` : '';
  return [
    title,
    '',
    what,
    '',
    why,
    '',
    `Why it matters: ${value}`,
    ctaLine,
    '',
    '#Solana #Web3 #NFTs'
  ].join('\n');
}

export function generateBoth(opts: {
  shortTitle: string;
  devWhat: string;
  devWhy: string;
  devTieIn: string;
  cryptoWhat?: string;
  cryptoWhy?: string;
  cryptoValue?: string;
  cta?: string;
}): PostPair {
  const dev = devStyle(
    opts.shortTitle,
    opts.devWhat,
    opts.devWhy,
    opts.devTieIn,
    opts.cta
  );

  const crypto = cryptoStyle(
    opts.shortTitle,
    opts.cryptoWhat ?? opts.devWhat,
    opts.cryptoWhy ?? opts.devWhy,
    opts.cryptoValue ?? 'Local-first backups keep your assets safe — verifiable & portable.',
    opts.cta
  );

  return { dev, crypto };
}
