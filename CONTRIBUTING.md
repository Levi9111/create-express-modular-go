# Contributing to create-express-modular-go

Thank you for your interest in contributing to **create-express-modular-go**! This project aims to provide sub-10ms, clean-architecture project scaffolding for Express.js + TypeScript backends.

---

## 🛠 Development Setup

### Prerequisites

- **Go:** Version `1.22+`
- **Make:** For build and test automation
- **Git**

### Clone & Build

```bash
# Clone repository
git clone https://github.com/Levi9111/create-express-modular-go.git
cd create-express-modular-go

# Compile binary
make build

# Run local binary
./bin/cem --help
```

---

## 🧪 Testing & Verification

Before submitting any Pull Request, ensure all tests pass and code is formatted:

```bash
# Run unit test suite
make test

# Run tests with race condition detector
go test -race ./...

# Format code
go fmt ./...

# Vet code for common pitfalls
go vet ./...
```

---

## 📝 Commit Convention

We follow [Conventional Commits](https://www.conventionalcommits.org/):

- `feat:` A new feature or generator capability
- `fix:` A bug fix in templates, route injection, or CLI flags
- `docs:` Documentation improvements
- `test:` Adding or refactoring tests
- `refactor:` Code refactoring with no functional change
- `ci:` GitHub Actions or automation changes
- `chore:` Dependency bumps, toolchain adjustments

**Example:**
```bash
git commit -m "feat(generator): add support for redis cache middleware"
```

---

## 🔀 Pull Request Process

1. Fork the repository and create your branch from `main`:
   ```bash
   git checkout -b feat/my-new-feature
   ```
2. Make your changes with focused, logical commits.
3. Add corresponding unit tests in `internal/<package>/<package>_test.go`.
4. Ensure all CI checks and tests pass locally (`make test`).
5. Open a Pull Request referencing any related issues.
