package generator

import (
	"fmt"
	"path/filepath"
	"strings"
)

func scaffoldDocker(opts ProjectOptions) error {
	root := opts.ProjectPath
	serviceName := strings.ToLower(opts.ProjectName)

	dockerfile := `FROM node:20-alpine AS builder

WORKDIR /app

COPY package*.json ./
RUN npm ci

COPY . .
RUN npm run build

FROM node:20-alpine AS runner

WORKDIR /app
ENV NODE_ENV=production

COPY package*.json ./
RUN npm ci --only=production

COPY --from=builder /app/dist ./dist

EXPOSE 5000

CMD ["node", "dist/server.js"]
`

	dockerignore := `# Dependencies
node_modules/
npm-debug.log*

# TypeScript build info
*.tsbuildinfo

# Environment files
.env
.env.*
!.env.example

# Git
.git/
.gitignore

# OS & Tests
.DS_Store
coverage/
`

	dbService := ""
	dbEnv := "      - DATABASE_URL=mongodb://db:27017/" + opts.ProjectName
	dbVolume := ""

	switch opts.DB {
	case DbMongoose:
		dbService = `  db:
    image: mongo:7-jammy
    restart: always
    ports:
      - '27017:27017'
    volumes:
      - mongo-data:/data/db
`
		dbVolume = `volumes:
  mongo-data:
`
	case DbPrisma, DbDrizzle:
		dbEnv = "      - DATABASE_URL=postgresql://postgres:postgres@db:5432/" + opts.ProjectName + "?schema=public"
		dbService = `  db:
    image: postgres:16-alpine
    restart: always
    environment:
      POSTGRES_USER: postgres
      POSTGRES_PASSWORD: postgres
      POSTGRES_DB: ` + opts.ProjectName + `
    ports:
      - '5432:5432'
    volumes:
      - postgres-data:/var/lib/postgresql/data
`
		dbVolume = `volumes:
  postgres-data:
`
	}

	compose := fmt.Sprintf(`version: '3.9'

services:
  app:
    build:
      context: .
      dockerfile: Dockerfile
    container_name: %s-api
    restart: always
    ports:
      - '5000:5000'
    environment:
      - NODE_ENV=production
      - PORT=5000
%s
    depends_on:
      - db

%s
%s`, serviceName, dbEnv, dbService, dbVolume)

	if err := writeFile(filepath.Join(root, "Dockerfile"), dockerfile); err != nil {
		return err
	}
	if err := writeFile(filepath.Join(root, ".dockerignore"), dockerignore); err != nil {
		return err
	}
	return writeFile(filepath.Join(root, "docker-compose.yml"), compose)
}
