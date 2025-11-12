# 🧠 SolVault Pull Request

Thank you for contributing to **SolVault** — where we back up, verify, and prove authenticity on-chain 🔒  

Please make sure your PR follows the conventions below 👇

---

## 🔗 Related Issue

> Example:  
> Fixes #12  
> *(Only reference issue numbers here — never inside commit messages)*

---

## ✨ Summary

Briefly describe what this PR does:

- What feature, fix, or refactor was introduced?
- Any notable design or security considerations?
- Did it add or modify commands (e.g., `solvault verify`, `solvault sync`)?

---

## 💻 Changes

List the key changes introduced in this PR:

- [ ] Added new CLI command
- [ ] Improved daemon or GUI handling
- [ ] Updated documentation
- [ ] Fixed bug or regression
- [ ] Added tests
- [ ] Refactored existing code

---

## 🧩 Verification

Steps to test this PR locally:

```bash
git fetch origin <your-branch>
go build -o solvault cmd/solvault/main.go
./solvault <command>
