# create-express-modular-go (Experimental)

> **High-Performance Native Golang Engine for Scaffolding Clean-Architecture Express.js + TypeScript Backends in Sub-10 Milliseconds.**

[![Go Report Card](https://goreportcard.com/badge/github.com/Levi9111/create-express-modular-go)](https://goreportcard.com/report/github.com/Levi9111/create-express-modular-go)
[![Go Version](https://img.shields.io/github/go-mod/go-version/Levi9111/create-express-modular-go)](https://golang.org)
[![License: MIT](https://img.shields.io/badge/License-MIT-blue.svg)](LICENSE)
[![Original CLI](https://img.shields.io/badge/npm-create--express--modular-CB3837?logo=npm)](https://www.npmjs.com/package/create-express-modular)

---

## 📖 Table of Contents

- [⚡ Why Go? (Sub-10ms Scaffolding)](#-why-go-sub-10ms-scaffolding)
- [📥 Installation](#-installation)
  - [Method 1: Go Install (One-Liner)](#method-1-go-install-recommended)
  - [Method 2: Build From Source (`make install`)](#method-2-build-from-source)
  - [Method 3: System-Wide Install (`/usr/local/bin`)](#method-3-system-wide-install)
  - [Method 4: Local Binary Build](#method-4-local-binary-build)
- [🚀 Quick Start & Usage](#-quick-start--usage)
  - [1. Scaffold a New Project](#1-scaffold-a-new-project)
  - [2. Domain & Module Generator](#2-domain--module-generator)
  - [3. Lifecycle Runners](#3-lifecycle-runners)
- [⚙️ CLI Reference & Flags](#️-cli-reference--flags)
- [🎯 Feature Parity with Original CLI](#-feature-parity-with-original-cli)
- [📊 Benchmarks & Technical Deep Dive](#-benchmarks--technical-deep-dive)
- [🏗 Project Architecture](#-project-architecture)
- [📦 Cross-Platform Compilation](#-cross-platform-compilation)
- [📄 License](#-license)

---

## ⚡ Why Go? (Sub-10ms Scaffolding)

`create-express-modular-go` (`cem`) is a standalone, compiled reimplementation of the [`create-express-modular`](https://www.npmjs.com/package/create-express-modular) CLI.

While traditional Node.js/npx scaffolding takes several seconds to boot V8 and load module dependencies, this Go engine:

- **Starts in ~2ms** (zero V8 boot overhead).
- **Generates complete architectures in ~8ms** using in-memory embedded templates (`//go:embed`).
- **Has zero runtime dependencies:** Standalone binary with no need for Node.js or npm to run the CLI itself.
- **Includes fire-and-forget goroutine telemetry** that never delays terminal output.

---

## 📥 Installation

Choose any of the following installation methods:

### Method 1: Go Install (Recommended)

If you have Go (1.20+) installed:

```bash
go install github.com/Levi9111/create-express-modular-go@latest
```

Make sure your Go bin path is in your `PATH` (if not already):

```bash
# Add to ~/.bashrc or ~/.zshrc
export PATH="$HOME/go/bin:$PATH"
```

Verify:

```bash
cem --version
```

---

### Method 2: Build From Source

Clone and install using the included Makefile:

```bash
# 1. Clone the repository
git clone https://github.com/Levi9111/create-express-modular-go.git
cd create-express-modular-go

# 2. Compile and install to $GOPATH/bin
make install

# 3. Ensure ~/go/bin is in your PATH
export PATH="$HOME/go/bin:$PATH"

# 4. Verify
cem --help
```

---

### Method 3: System-Wide Install

To install globally for all system users without modifying `$PATH`:

```bash
cd create-express-modular-go
make build
sudo cp bin/cem /usr/local/bin/cem

# Verify
cem --help
```

---

### Method 4: Local Binary Build

You can also build and run directly in the repository directory without installing globally:

```bash
make build

# Run the generated local binary
./bin/cem my-api -y
```

---

## 🚀 Quick Start & Usage

### 1. Scaffold a New Project

#### Zero-Prompt Mode (Sub-10ms Scaffolding)

Instantly generate a full production Express + TypeScript project with recommended defaults:

```bash
cem my-api -y
```

#### Interactive Mode

Select your preferred database, validator, authentication strategy, and Docker setup interactively:

```bash
cem my-api
```

#### Skip Dependency Installation

If you want instant file generation without running `npm install`:

```bash
cem my-api -y --no-install
```

#### Custom Flag Scaffolding

```bash
cem my-api -y \
  --db prisma \
  --validator zod \
  --auth \
  --cookie \
  --docker \
  --swagger
```

---

### 2. Domain & Module Generator

The `cem` CLI comes with domain generator commands to rapidly build out clean, modular architectures:

```bash
cd my-api

# Generate a complete domain module (controller, service, route, validation, interface)
# Automatically registers and wires routes in src/app/routes/index.ts
cem add module Product

# Generate a custom middleware in src/app/middlewares/
cem add middleware requestLogger

# Add an environment variable across .env, .env.example, and config/index.ts
cem add env STRIPE_SECRET_KEY

# Cleanly remove resources and auto-unwire routes
cem remove module Product
cem remove middleware requestLogger
cem remove env STRIPE_SECRET_KEY
```

---

### 3. Lifecycle Runners

Run daily development tasks directly with `cem` (auto-detects `bun`, `pnpm`, `yarn`, or `npm`):

```bash
# Start development server with live reload
cem dev

# Run comprehensive checks (type-check, ESLint, Prettier)
cem check

# Auto-fix linting and formatting issues
cem fix

# Compile TypeScript to dist/
cem build

# Start production server
cem start

# Inspect registered modules, middlewares, and env variables
cem list
```

---

## ⚙️ CLI Reference & Flags

### Root Flags (`cem [project-name]`)

| Flag                 | Description                                           | Default    |
| :------------------- | :---------------------------------------------------- | :--------- |
| `-y, --yes`          | Scaffold using recommended defaults non-interactively | `false`    |
| `--db <engine>`      | Database engine (`mongoose`, `prisma`, `drizzle`)     | `mongoose` |
| `--validator <type>` | Schema validator (`zod`, `joi`)                       | `zod`      |
| `--auth`             | Enable JWT authentication module                      | `true`     |
| `--no-auth`          | Disable JWT authentication module                     | `false`    |
| `--cookie`           | Deliver JWT tokens via HTTP-only cookies              | `true`     |
| `--header`           | Deliver JWT tokens via Authorization headers          | `false`    |
| `--docker`           | Generate Dockerfile and docker-compose.yml            | `true`     |
| `--no-docker`        | Skip Docker file generation                           | `false`    |
| `--swagger`          | Enable Swagger (OpenAPI 3.0) documentation at `/docs` | `true`     |
| `--no-swagger`       | Disable Swagger documentation                         | `false`    |
| `--no-install`       | Skip dependency installation                          | `false`    |
| `-v, --version`      | Print CLI version                                     | -          |
| `-h, --help`         | Display help menu                                     | -          |

---

## 🎯 Feature Parity with Original CLI

`create-express-modular-go` maintains **100% feature parity** with the Node.js version:

- [x] **Zero-prompt scaffolding (`-y` / `--yes`)**
- [x] **Interactive survey mode** powered by `survey/v2`
- [x] **Database selection:** Mongoose (with `QueryBuilder`), Prisma, Drizzle ORM
- [x] **Validator selection:** Zod, Joi
- [x] **Turnkey Auth:** Complete JWT auth module with cookie or header token delivery
- [x] **Docker support:** Production `Dockerfile`, `.dockerignore`, `docker-compose.yml`
- [x] **Documentation:** OpenAPI 3.0 / Swagger UI, `README.md`, `AGENTS.md`, `CLAUDE.md`
- [x] **Domain generators:** `cem add/remove module`, `middleware`, `env` (with AST route wiring)
- [x] **Lifecycle runners:** `cem dev`, `cem build`, `cem start`, `cem list`, `cem check`, `cem fix`
- [x] **Cybernetic Terminal HUD:** ANSI gradient styling, spinners, and summary cards

---

## 📊 Benchmarks & Technical Deep Dive

A head-to-head benchmark running non-interactive scaffolding (`-y --no-install`) on Linux x86_64:

| Metric                   | Node.js Engine (`create-express-modular`) | Golang Engine (`create-express-modular-go`) | Difference       |
| :----------------------- | :---------------------------------------- | :------------------------------------------ | :--------------- |
| **CLI Cold Startup**     | `~118 ms`                                 | **`2 ms`**                                  | **59x faster**   |
| **Disk Scaffolding**     | `~132 ms`                                 | **`8 ms` (0.008s)**                         | **16.5x faster** |
| **Peak Memory (RSS)**    | `68 MB`                                   | **`11.8 MB`**                               | **5.7x lighter** |
| **Binary Size**          | `~45 MB` (with `node_modules`)            | **`9.0 MB` (standalone)**                   | **5x smaller**   |
| **Runtime Dependencies** | Node.js 18+, npm/npx                      | **Zero (Native compiled binary)**           | 100% portable    |

### Latency Breakdown: Where Does Time Go?

```text
Total Scaffolding Latency Breakdown:
┌────────────────────────────────────────────────────────────────────────┐
│ [1] Dependency Download & Install (npm install): ~1,500ms - 4,500ms   │  92.5% of total time
├────────────────────────────────────────────────────────────────────────┤
│ [2] Blocking Network Calls (e.g. telemetry): ~0ms (goroutine)          │  0.0% of total time
├────────────────────────────────────────────────────────────────────────┤
│ [3] CLI Runtime Startup: ~2ms                                          │  0.1% of total time
├────────────────────────────────────────────────────────────────────────┤
│ [4] File Generation & Disk I/O (Go engine): ~8ms                       │  0.4% of total time
└────────────────────────────────────────────────────────────────────────┘
```

---

## 🏗 Project Architecture

```
create-express-modular-go/
├── cmd/                          # Cobra CLI commands
│   ├── root.go                   # Main 'cem' root command & flags
│   ├── add.go                    # 'cem add module|middleware|env'
│   ├── remove.go                 # 'cem remove module|middleware|env'
│   └── lifecycle.go              # 'cem dev|build|start|list|check|fix|eject'
├── internal/
│   ├── config/                   # 'cem-cli.json' configuration loader & writer
│   ├── generator/                # Scaffolding orchestrator
│   │   ├── embed.go              # go:embed base template loader
│   │   ├── core.go               # server.ts, app.ts, routes, utils, errors
│   │   ├── db.go                 # Mongoose, Prisma, Drizzle integrations
│   │   ├── validator.go          # Zod & Joi schema validators
│   │   ├── auth.go               # JWT authentication module & middlewares
│   │   ├── docker.go             # Dockerfile & compose generators
│   │   ├── package_json.go       # Dependency resolution & script configuration
│   │   └── templates/base/       # Static base files (tsconfig, eslint, prettier)
│   ├── modules/                  # Resource generators & route AST injector
│   ├── pm/                       # Package manager detection (bun, pnpm, yarn, npm)
│   ├── telemetry/                # Non-blocking goroutine telemetry
│   └── ui/                       # Cybernetic HUD, colors, spinners & banners
├── main.go                       # Application entry point
├── Makefile                      # Build, test, and cross-compile targets
└── README.md
```

---

## 📦 Cross-Platform Compilation

Generate standalone, zero-dependency binaries for all major operating systems in one command:

```bash
make build-all
```

Outputs in `bin/`:

- `cem-linux-amd64` (Linux 64-bit)
- `cem-linux-arm64` (Linux ARM / Raspberry Pi)
- `cem-darwin-arm64` (macOS Apple Silicon)
- `cem-darwin-amd64` (macOS Intel)
- `cem-windows-amd64.exe` (Windows 64-bit)

---

## 🧪 Testing

```bash
make test
```

---

## 📄 License

MIT © [Levi9111](https://github.com/Levi9111)
