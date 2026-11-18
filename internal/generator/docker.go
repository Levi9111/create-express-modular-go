package generator

import (
	"fmt"
	"path/filepath"
	"regexp"
	"strings"

	"github.com/Levi9111/create-express-modular-go/internal/pm"
)

func scaffoldDocker(opts ProjectOptions) error {
	root := opts.ProjectPath
	serviceName := regexp.MustCompile(`[^a-z0-9-]`).ReplaceAllString(strings.ToLower(opts.ProjectName), "-")

	if err := writeFile(filepath.Join(root, "Dockerfile"), buildDockerfile(opts.PM)); err != nil {
		return err
	}

	if err := writeFile(filepath.Join(root, ".dockerignore"), buildDockerignore(opts.PM)); err != nil {
		return err
	}

	return writeFile(filepath.Join(root, "docker-compose.yml"), buildCompose(serviceName, opts.DB))
}

func buildDockerfile(p pm.PackageManager) string {
	pmStr := string(p)
	lockfile := map[string]string{
		"npm":  "package-lock.json",
		"yarn": "yarn.lock",
		"pnpm": "pnpm-lock.yaml",
		"bun":  "bun.lockb",
	}[pmStr]
	if lockfile == "" {
		lockfile = "package-lock.json"
	}

	copyLock := fmt.Sprintf("COPY package.json %s ./", lockfile)
	if pmStr == "npm" || pmStr == "" {
		copyLock = "COPY package*.json ./"
	}

	pnpmSetup := ""
	if pmStr == "pnpm" {
		pnpmSetup = "\nRUN corepack enable\n"
	}

	buildCmd := map[string]string{
		"npm":  "npm run build",
		"yarn": "yarn build",
		"pnpm": "pnpm build",
		"bun":  "bun run build",
	}[pmStr]
	if buildCmd == "" {
		buildCmd = "npm run build"
	}

	installAll := map[string]string{
		"npm":  "RUN npm ci",
		"yarn": "RUN yarn install --frozen-lockfile",
		"pnpm": "RUN pnpm install --frozen-lockfile",
		"bun":  "RUN bun install --frozen-lockfile",
	}[pmStr]
	if installAll == "" {
		installAll = "RUN npm ci"
	}

	installProd := map[string]string{
		"npm":  "RUN npm ci --omit=dev && npm cache clean --force",
		"yarn": "RUN yarn install --frozen-lockfile --production && yarn cache clean",
		"pnpm": "RUN pnpm install --frozen-lockfile --prod",
		"bun":  "RUN bun install --frozen-lockfile --production",
	}[pmStr]
	if installProd == "" {
		installProd = "RUN npm ci --omit=dev && npm cache clean --force"
	}

	return fmt.Sprintf(`# ── Build stage ───────────────────────────────────────────────
FROM node:22-alpine AS builder

WORKDIR /app
%s
# Install dependencies first (layer cache)
%s
%s

# Copy source and compile
COPY . .
RUN %s

# ── Production stage ───────────────────────────────────────────
FROM node:22-alpine AS production

WORKDIR /app

ENV NODE_ENV=production
%s
# Install only production deps
%s
%s

# Copy compiled output
COPY --from=builder /app/dist ./dist

EXPOSE 5000

# Run as a non-root user for security
RUN addgroup -S appgroup && adduser -S appuser -G appgroup
USER appuser

CMD ["node", "dist/server.js"]
`, pnpmSetup, copyLock, installAll, buildCmd, pnpmSetup, copyLock, installProd)
}

func buildDockerignore(p pm.PackageManager) string {
	currentLock := map[string]string{
		"npm":  "package-lock.json",
		"yarn": "yarn.lock",
		"pnpm": "pnpm-lock.yaml",
		"bun":  "bun.lockb",
	}[string(p)]

	allLocks := []string{"package-lock.json", "yarn.lock", "pnpm-lock.yaml", "bun.lockb"}
	var unused []string
	for _, lf := range allLocks {
		if lf != currentLock {
			unused = append(unused, lf)
		}
	}

	return fmt.Sprintf(`# Dependencies
node_modules/
npm-debug.log*

# Unused lock files
%s

# TypeScript build info
*.tsbuildinfo

# Environment files — never bake secrets into the image
.env
.env.*
!.env.example

# Git
.git/
.gitignore

# Editor
.vscode/
.idea/

# OS
.DS_Store
Thumbs.db

# Tests
coverage/
*.test.ts
*.spec.ts
`, strings.Join(unused, "\n"))
}

func buildCompose(serviceName string, db DbChoice) string {
	dbService := ""
	switch db {
	case DbMongoose:
		dbService = `
  db:
    image: mongo:7
    restart: unless-stopped
    environment:
      MONGO_INITDB_ROOT_USERNAME: ${MONGO_USER:-root}
      MONGO_INITDB_ROOT_PASSWORD: ${MONGO_PASSWORD:-secret}
      MONGO_INITDB_DATABASE: ${MONGO_DB:-appdb}
    ports:
      - '27017:27017'
    volumes:
      - mongo-data:/data/db
    networks:
      - app-network

volumes:
  mongo-data:
`
	case DbPrisma, DbDrizzle:
		dbService = `
  db:
    image: postgres:16-alpine
    restart: unless-stopped
    environment:
      POSTGRES_USER: ${POSTGRES_USER:-postgres}
      POSTGRES_PASSWORD: ${POSTGRES_PASSWORD:-secret}
      POSTGRES_DB: ${POSTGRES_DB:-appdb}
    ports:
      - '5432:5432'
    volumes:
      - pg-data:/var/lib/postgresql/data
    networks:
      - app-network

volumes:
  pg-data:
`
	}

	dependsOn := ""
	if dbService != "" {
		dependsOn = "\n    depends_on:\n      - db"
	}

	return fmt.Sprintf(`version: '3.9'

services:
  %s:
    build:
      context: .
      target: production
    ports:
      - '${PORT:-5000}:5000'
    env_file:
      - .env
    environment:
      NODE_ENV: production
    restart: unless-stopped%s
    networks:
      - app-network
%s
networks:
  app-network:
    driver: bridge
`, serviceName, dependsOn, dbService)
}
