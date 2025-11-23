# X/Twitter Post Guidelines

This document outlines the two posting styles used for SolVault communications on X (Twitter).

## Overview

SolVault uses **two distinct voices** for different audiences:

1. **Dev Style** — for builders, technical followers, and #BuildInPublic community
2. **Crypto Style** — for collectors, DAOs, NFT enthusiasts, and broader crypto audience

Both styles maintain authenticity while serving different reader expectations.

---

## Dev Style (Builder Voice)

### Structure

```
Gm 𝕏 family ☀️

[Short Title]

[What was built/shipped]

[Why it matters technically]

[How it ties to the bigger vision]

[Optional CTA]

#Solana #BuildInPublic
```

### Characteristics

- **Tone:** Conversational, transparent, human
- **Greeting:** Always "Gm 𝕏 family ☀️"
- **Focus:** Technical implementation, challenges overcome, learning
- **Hashtags:** #Solana #BuildInPublic (consistent)
- **Length:** Medium (4-6 lines of substance)

### Example Template

```
Gm 𝕏 family ☀️

🚀 Just shipped [feature name]

Built [what] to solve [problem]. Used [tech/approach] because [reason].

Why this matters: [technical value prop]

Next up: [what's coming]

#Solana #BuildInPublic
```

---

## Crypto Style (Collector/DAO Voice)

### Structure

```
[Strong Title/Hook]

[What it does for users]

[Why they should care]

Why it matters: [Value proposition]

[Optional CTA]

#Solana #Web3 #NFTs
```

### Characteristics

- **Tone:** Professional, value-focused, benefit-driven
- **Opening:** Direct hook or question (no "Gm")
- **Focus:** User benefits, security, ownership, trust
- **Hashtags:** #Solana #Web3 #NFTs (broader reach)
- **Length:** Concise (3-4 punchy lines)

### Example Template

```
🔒 [Feature Name] is live

[What users can now do]. [How it works, briefly].

Built for [target audience]: [benefit 1], [benefit 2], [benefit 3].

Why it matters: [Core value — security, ownership, portability]

[CTA if applicable]

#Solana #Web3 #NFTs
```

---

## When to Use Each Style

| Situation | Style | Reason |
|-----------|-------|--------|
| Shipping a new feature | Dev | Show the technical work and learning |
| Major milestone (e.g., v1.0 launch) | Both | Dev for builders, Crypto for users |
| Bug fix or technical improvement | Dev | Technical audience cares about quality |
| New integration (IPFS, Notion, etc.) | Crypto | User-facing benefit is the story |
| Community update / roadmap | Both | Keep both audiences engaged |
| Demo or tutorial | Dev | Building in public = sharing process |

---

## Content Guidelines

### Do's ✅

- **Be authentic**: Real progress, real challenges, real excitement
- **Show, don't just tell**: Screenshots, code snippets, demo links
- **Engage**: Ask questions, respond to replies, build community
- **Credit others**: Tag libraries, contributors, inspiration sources
- **Use emojis sparingly**: 1-2 per post for visual breaks
- **Keep it scannable**: Line breaks, bullet points when needed

### Don'ts ❌

- **Don't over-promise**: Ship first, then post
- **Don't spam hashtags**: 2-3 relevant tags max
- **Don't post just to post**: Quality > frequency
- **Don't ignore replies**: Community engagement matters
- **Don't copy-paste**: Each post should feel fresh

---

## Examples (Reference Only)

### Dev Style Example

```
Gm 𝕏 family ☀️

🧪 Tested the new NFT verification flow today

Metadata hashing + on-chain comparison = bulletproof authenticity checks. Took 3 tries to get the Metaplex integration right, but now it's solid.

Next: Adding IPFS fallback for metadata that's not on Arweave.

#Solana #BuildInPublic
```

### Crypto Style Example

```
🔐 Your NFTs deserve better backups

SolVault now auto-saves metadata + images locally whenever you mint. No cloud dependencies. No permissions. Just your wallet + your files.

Why it matters: If Arweave goes down or a project rug pulls, you still own the proof.

Try it: [link]

#Solana #Web3 #NFTs
```

---

## Integration with xPoster.ts

The `utils/xPoster.ts` utility provides a `generateBoth()` function to create both styles programmatically:

```typescript
import { generateBoth } from './utils/xPoster';

const posts = generateBoth({
  shortTitle: '🚀 Shipped CLI v0.1.0',
  devWhat: 'Built a TypeScript CLI that backs up your NFTs locally.',
  devWhy: 'Solves the "what if the project site goes down" problem.',
  devTieIn: 'Next: daemon mode for auto-backups.',
  cryptoWhat: 'Your NFTs are now backed up locally — no cloud required.',
  cryptoWhy: 'One command, and your collection is yours forever.',
  cryptoValue: 'True ownership means having the files, not just the on-chain pointer.',
  cta: 'Try it: npm install -g solvault'
});

console.log('Dev post:\n', posts.dev);
console.log('\nCrypto post:\n', posts.crypto);
```

---

## Notes

- **Do NOT include live marketing posts in PR bodies** — reference the guidelines, but don't paste full posts into PR descriptions
- Posts are for X/Twitter only — PR descriptions should be technical and factual
- When sharing progress, use Dev Style; when announcing user-facing value, use Crypto Style
- Both styles can coexist: post Dev Style on your personal account, Crypto Style on the project account

---

## Maintenance

Update this guide as the voice evolves. If you find a format that performs better, document it here and update the `xPoster.ts` templates accordingly.

**Last updated:** 2025-11-23
