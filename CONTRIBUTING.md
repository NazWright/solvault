# 🤝 Contributing to SolVault

Thank you for your interest in contributing to **SolVault** 💜  
Together, we can build the most trusted NFT backup and proof tool in the Solana ecosystem.

---

## 🧰 Getting Started

### Prerequisites
- [Go 1.22+](https://go.dev/dl/)
- A Solana RPC endpoint (Helius, QuickNode, or public)
- Git + CLI environment (macOS, Linux, or Windows)
- Optional: Notion, IPFS, or Google Drive API keys for integration testing

---

## ⚙️ Local Setup

```bash
git clone https://github.com/NazWright/solvault.git
cd solvault
go mod tidy
go run cmd/solvault/main.go
To build:

bash
Copy code
go build -o solvault cmd/solvault/main.go
To run tests:

bash
Copy code
go test ./...
🌱 Branch Workflow
We use the GitHub Flow model with a few naming conventions.

Branch	Purpose
main	Stable production releases
dev	Active development
feature/<name>	New feature or experiment
fix/<name>	Bug fix or patch
docs/<name>	Documentation updates

Create a new branch for any change:

bash
Copy code
git checkout -b feature/add-ipfs-sync
Push your branch and open a Pull Request (PR) against dev.

🧩 Commit Style
Keep commit messages clear and conventional:

Type	Example
feat:	feat(verifier): add proof JSON output
fix:	fix(listener): correct retry logic for Solana RPC
docs:	docs: update README for v1.0.0
chore:	chore: bump dependencies

🔢 Versioning & Releases
SolVault follows Semantic Versioning (semver.org)

Copy code
vMAJOR.MINOR.PATCH
Example:

bash
Copy code
git checkout main
git pull
git tag -a v1.1.0 -m "Add NFT proof publishing"
git push origin v1.1.0
This triggers an automated GitHub Action to:

Build binaries for macOS, Linux, and Windows

Attach them to the new release page

Update badges in the README automatically

🧠 Code Guidelines
Use go fmt ./... before every commit

Keep public functions documented with GoDoc comments

Avoid committing large binary files — use /build/ or /backups/ for local testing only

Store shared logic in /pkg/solvaultsdk

Write tests near implementation files (*_test.go)

📄 Pull Request Checklist
Before submitting your PR:

 Code builds locally

 All tests pass (go test ./...)

 README or docs updated if applicable

 No .env, backups, or logs committed

 Branch up to date with dev

💬 Communication
🐞 Report bugs → GitHub Issues

💡 Share ideas → GitHub Discussions

📬 Stay connected → @daredevtech

💜 Thank You
Every PR, feature, or idea helps SolVault grow into a secure and transparent verification layer for NFTs.
Built with pride by Naz Wright / DareDevTech 🦾

yaml
Copy code

---

✅ **Placement:**
- `.gitignore` → root of your repo  
- `CONTRIBUTING.md` → root of your repo  

When you push these, GitHub will automatically detect the license and show a **"MIT" badge** and **"Contributing" link** in your repository sidebar.  

Would you like me to generate a **CHANGELOG.md** next — prefilled with your v1.0.0 milestones (CLI core, verify command, daemon prep, etc.) so you can tag your first release cleanly?




