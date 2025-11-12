```markdown
# 🧠 SolVault — Solana NFT Backup, Verification & Proof System

![Version](https://img.shields.io/github/v/release/daredevtech/solvault)
![Build](https://github.com/daredevtech/solvault/actions/workflows/release.yml/badge.svg)
![License](https://img.shields.io/github/license/daredevtech/solvault)
![Go](https://img.shields.io/badge/Built%20with-Go-00ADD8?logo=go)

> **Back up, verify, and prove ownership of your NFTs — all from one binary.**

SolVault is a cross-platform app that watches your Solana wallet for new NFT mints, automatically downloads metadata and images, verifies authenticity through on-chain hashes, and (optionally) publishes a public proof page to the web.  

It’s designed to run silently in the background or be opened as a GUI viewer — all powered by the same backend engine.

---

## 🧩 Project Vision

| Phase | Focus | Description |
|:------|:------|:-------------|
| **Phase 1** | CLI Core | Build the command-line utility that watches a wallet, downloads NFT metadata/images, and saves them locally. |
| **Phase 2** | Daemon Mode | Enable background operation with lightweight scheduling and logging. |
| **Phase 3** | GUI Visualization | Add a desktop interface that reads the same backups and displays NFT cards + metadata interactively. |

---

## ⚙️ Phase 1 — CLI Utility (`v1.0.0` target)

### 🔧 Core Behavior
- Start the binary: `./solvault watch`
- Reads config from `.env`
- Monitors wallet for new NFT mints
- Saves each NFT’s:
  - image file  
  - metadata JSON  
  - verification hash  
  - log entry (`backups/log.json`)

### 🧱 Folder Layout
```

cmd/solvault/
├── main.go           # entrypoint
├── root.go           # base command
├── watch.go          # solvault watch
├── verify.go         # solvault verify [--publish]
├── list.go           # solvault list
└── info.go           # solvault info <mint>
internal/
├── listener/         # wallet monitor
├── fetcher/          # metadata/image download
├── verifier/         # hash + proof generator
├── storage/          # file system & logs
└── utils/            # helpers, config, logging
pkg/
└── solvaultsdk/      # reusable developer SDK

````

---

## 🧩 CLI Commands

| Command | Description |
|:---------|:-------------|
| `solvault init` | Initializes `.env` and backup folder. |
| `solvault watch` | Starts watching your wallet for new NFTs. |
| `solvault verify <mint>` | Verifies NFT authenticity and saves proof. |
| `solvault verify <mint> --publish` | Verifies and uploads a public proof page. |
| `solvault list` | Lists all backed-up NFTs. |
| `solvault info <mint>` | Displays detailed metadata for an NFT. |
| `solvault sync` | Pushes data to cloud integrations (optional). |

**Example**
```bash
> solvault init
✅ Initialized SolVault configuration.

> solvault watch
👀 Watching wallet: 5QfQ...ZsLk
🆕 New NFT detected: “Midnight Lion #01”
✅ Saved image + metadata to ~/SolVaultBackups/MidnightLion01/

> solvault verify 9sdfe1xA3s...JKX1L --publish
✅ Authentic NFT verified!
Proof available at: https://proofs.solvault.app/9sdfe1xA3s...JKX1L
````

---

## 🔒 Verification + Proof System

SolVault uses on-chain metadata (Metaplex) and local hashing to verify your NFT image authenticity.

**Verification Steps**

1. Fetch NFT metadata and image URI
2. Compute SHA-256 hash of downloaded image
3. Compare to stored or canonical hash
4. Generate a local proof JSON
5. (Optional) Publish to SolVault web portal

**Proof JSON Example**

```json
{
  "nft_name": "Midnight Lion #01",
  "mint_address": "9sdfe1xA3s...JKX1L",
  "verified_by": "5QfQ...ZsLk",
  "verified_at": "2025-11-12T18:30:00Z",
  "image_hash": "e3b0c44298fc1c149...",
  "status": "authentic",
  "proof_link": "https://proofs.solvault.app/9sdfe1xA3s...JKX1L"
}
```

---

## 🌐 Proof Portal (Phase 2 → 3)

When you publish a proof, it will be viewable at:

```
https://proofs.solvault.app/<mint>
```

Visitors will see:

* NFT image + metadata
* Image hash
* Verification timestamp
* Wallet that performed verification
* “Verified by SolVault” footer

Share-ready for X (Twitter):

```
🧠 Verified my #NFT authenticity with #SolVault  
✅ Midnight Lion #01  
🔗 https://proofs.solvault.app/9sdfe1xA3s...JKX1L  
#Solana #NFTVerification #DareDevTech
```

---

## 🖥️ Phase 2 — Daemon Mode

Run SolVault continuously in the background:

```bash
solvault watch --daemon
```

* Logs stored under `~/.solvault/logs/`
* Managed by system service:

  * macOS → `launchd`
  * Windows → `Task Scheduler`
  * Linux → `systemd`

---

## 💠 Phase 3 — GUI Visualization

When launched normally, SolVault opens a desktop dashboard that reads your existing backups:

* Displays thumbnails & metadata
* Filters NFTs by collection or tag
* Shows proof status (verified / unverified)
* Syncs cloud integrations manually

**Tech Options**

* 🟦 **Fyne (Go-native)** — lightweight single-binary GUI
* 🟣 **Tauri (Rust + TypeScript)** — modern animated interface

---

## 🔢 Versioning

SolVault follows **[Semantic Versioning (semver.org)](https://semver.org)**

```
vMAJOR.MINOR.PATCH
```

| Version  | Feature       | Description                              |
| :------- | :------------ | :--------------------------------------- |
| `v1.0.0` | CLI Utility   | NFT backup + verification (local proofs) |
| `v1.1.0` | Daemon Mode   | Background watcher + logging             |
| `v2.0.0` | GUI Dashboard | NFT visual browser & proof viewer        |
| `v3.0.0` | Cloud Sync    | Notion, Drive, IPFS integrations         |
| `v4.0.0` | Multi-Wallet  | Manage multiple wallets & UI profiles    |

---

## 🧰 Developer Setup

```bash
git clone https://github.com/daredevtech/solvault.git
cd solvault
go mod tidy
go run cmd/solvault/main.go
```

Run tests:

```bash
go test ./...
```

---

## 🤝 Contributing

Pull requests & ideas welcome!
Check [CONTRIBUTING.md](./CONTRIBUTING.md) for guidelines and branch conventions.

---

## 💬 Community & Support

* 🐞 Issues → [GitHub Issues](https://github.com/daredevtech/solvault/issues)
* 💡 Ideas → [Discussions](https://github.com/daredevtech/solvault/discussions)
* 🧑‍💻 Updates → [@daredevtech](https://x.com/daredevtech)

---

## 🛡️ License

**MIT License** © 2025 [DareDevTech](https://x.com/daredevtech)
Free to use · Credit appreciated · Build something remarkable.

---

### 💜 Why SolVault Exists

Your art deserves permanence —
on-chain, in your hands, and verified by truth.

```

---

```
