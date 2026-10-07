# create-express-modular-go (Experimental)

> High-performance Golang CLI engine for scaffolding production-ready modular Express.js + TypeScript backends in **< 10ms**.

[![Go Version](https://img.shields.io/github/go-mod/go-version/Levi9111/create-express-modular-go)](https://golang.org)
[![License: MIT](https://img.shields.io/badge/License-MIT-blue.svg)](LICENSE)

---

## ⚡ Performance Benchmark

| CLI Engine | Scaffolding Latency | Memory Footprint | Runtime Dependencies |
| :--- | :--- | :--- | :--- |
| **Golang Native Binary (`go-cem`)** | **~8 ms** | **~12 MB** | **0 (Standalone binary)** |
| Node.js Engine (`ts-node`/`node`) | ~130 ms | ~70 MB | Node.js + npm/npx |

*Benchmark measured with `-y --no-install` on Linux x86_64.*

---

## 🚀 Installation & Quick Start

### Build from Source

```bash
git clone https://github.com/Levi9111/create-express-modular-go.git
cd create-express-modular-go
make build

# The compiled binary is located at bin/cem
./bin/cem --help
```

### Or install to `$GOPATH/bin`:

```bash
make install
cem my-api
```

---

## 🛠 Usage & Command Parity

`create-express-modular-go` has **100% feature parity** with the Node.js `create-express-modular` CLI:

### 1. Interactive Scaffolding
```bash
cem my-backend
```

### 2. Zero-Prompt Fast Scaffolding (`-y` / `--yes`)
```bash
cem my-backend -y
cem my-backend -y --db prisma --validator zod --auth --docker --no-install
```

#### Flags:
- `-y, --yes`: Non-interactive mode using defaults
- `--db <mongoose|prisma|drizzle>`: Database ORM/ODM
- `--validator <zod|joi>`: Schema validator
- `--auth`: Include JWT auth module
- `--cookie`: Store tokens in HTTP-only cookies
- `--header`: Return tokens in JSON response headers
- `--docker`: Generate Dockerfile & docker-compose.yml
- `--swagger`: Include OpenAPI 3.0 / Swagger UI documentation
- `--no-install`: Skip `npm install` for instantaneous generation
- `-v, --version`: Print version information

### 3. Module & Middleware Generators
```bash
# Add a complete modular domain (controller, service, route, validation, interface)
cem add module Product

# Add a custom middleware
cem add middleware requestLogger

# Add environment variable to .env and src/app/config/index.ts
cem add env STRIPE_SECRET_KEY

# Remove/unwire module, middleware, or env
cem remove module Product
cem remove middleware requestLogger
cem remove env STRIPE_SECRET_KEY
```

### 4. Lifecycle & Project Commands
```bash
# List all registered modules and route endpoints
cem list

# Run dev server using project package manager (npm/pnpm/yarn/bun)
cem dev

# Build TypeScript to dist/
cem build

# Start production server
cem start

# Architecture health check & lint fix
cem check
cem fix
```

---

## 🧪 Testing

Run all unit and integration tests:

```bash
make test
```

## 📦 Cross-Platform Compilation

Generate standalone binaries for Linux, macOS (Apple Silicon & Intel), and Windows:

```bash
make build-all
```
Outputs:
- `bin/cem-linux-amd64`
- `bin/cem-linux-arm64`
- `bin/cem-darwin-arm64`
- `bin/cem-darwin-amd64`
- `bin/cem-windows-amd64.exe`

---

## 📄 License

MIT © [Levi9111](https://github.com/Levi9111)
