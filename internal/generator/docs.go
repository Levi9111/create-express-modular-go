package generator

import (
	"fmt"
	"path/filepath"
)

func scaffoldDocs(opts ProjectOptions) error {
	root := opts.ProjectPath

	readme := fmt.Sprintf(`# %s

Modular, domain-driven Express + TypeScript backend.

> Scaffolded with [create-express-modular](https://create-express-modular.lovable.app/) (cem)

---

## Stack
- **Runtime**: Node.js + Express
- **Language**: TypeScript
- **Database**: %s
- **Validation**: %s
- **Authentication**: %t
- **Docker**: %t
- **Swagger Docs**: %t

---

## Getting Started

### 1. Install dependencies
`+"```bash"+`
npm install
`+"```"+`

### 2. Configure environment
Check `+"`.env`"+` and adjust database credentials as necessary.

### 3. Start development server
`+"```bash"+`
npm run dev
`+"```"+`

Server runs at `+"http://localhost:5000"+`.
API documentation is available at `+"http://localhost:5000/docs"+`.

---

## Available Commands
- `+"`npm run dev`"+` — Start development server with hot reload
- `+"`npm run build`"+` — Compile TypeScript to dist/
- `+"`npm run start`"+` — Start production server
- `+"`npm run check`"+` — Run type-check, lint, and format checks

---

Built with [create-express-modular](https://create-express-modular.lovable.app/).
`, opts.ProjectName, opts.DB, opts.Validator, opts.UseAuth, opts.UseDocker, opts.UseSwagger)

	agentsMd := fmt.Sprintf(`# AGENTS.md — %s

## Token-Saving Workflow
1. Read `+"`cem-cli.json`"+` first (stack, live `+"`state.modules`"+`, `+"`state.middlewares`"+`, `+"`state.env`"+`) instead of scanning directories.
2. Use `+"`cem`"+` with `+"`--json`"+` or `+"`--quiet`"+` for compact, zero-noise output.
3. Never hand-write boilerplate files — generate them with `+"`cem add`"+`.

## Stack
- Runtime: Node.js ≥ 22 | Express + TypeScript | PM: %s
- DB: %s | Validator: %s | Auth: %t | Docker: %t | Swagger: %t | Testing: %t

## CLI Commands
- `+"`cem list --json`"+` — inspect modules, middlewares, env
- `+"`cem add module <Name...> [--fields \"name:string,price:number,inStock:boolean?,category:enum(a|b),owner:ref(User),tags:string[]?\"] --json`"+` — scaffold module(s) with typed schema, validation, wired routes, & CRUD service
- `+"`cem add middleware <name...> --json`"+` — scaffold `+"`<name>.middleware.ts`"+`
- `+"`cem add env <KEY...> --json`"+` — add env key(s) to `+"`.env`"+`, `+"`.env.example`"+`, `+"`src/app/config/index.ts`"+`
- `+"`cem remove module|middleware|env <name...> --json`"+` — remove & unwire
- `+"`cem check --json`"+` / `+"`cem check --fix --json`"+` — typecheck, lint, format

## Rules
- Wrap handlers with `+"`catchAsync()`"+`; respond with `+"`sendResponse()`"+`; validate with `+"`validateRequest()`"+`.
- Middlewares MUST be named `+"`<name>.middleware.ts`"+` in `+"`src/app/middlewares/`"+`.
- Never read `+"`process.env`"+` directly — import `+"`config`"+` from `+"`src/app/config`"+`.
- Keep `+"`INJECT`"+` markers in `+"`src/app/routes/index.ts`"+` intact.
`, opts.ProjectName, opts.PM, opts.DB, opts.Validator, opts.UseAuth, opts.UseDocker, opts.UseSwagger, opts.UseTesting)

	claudeMd := fmt.Sprintf(`# CLAUDE.md — %s

See `+"`AGENTS.md`"+` for full rules and `+"`cem-cli.json`"+` for live project state (`+"`modules`"+`, `+"`middlewares`"+`, `+"`env`"+`).

## Stack
- %s | %s | %s

## Commands (Prefer `+"`--json`"+` / `+"`--quiet`"+` to minimize output)
- `+"`cem list --json`"+` — inspect project state
- `+"`cem add module <Name...> [--fields \"name:string,price:number\"] --json`"+` — generate typed module, CRUD service, & wired route
- `+"`cem add middleware <name...> --json`"+` — generate `+"`<name>.middleware.ts`"+`
- `+"`cem add env <KEY...> --json`"+` — add env var
- `+"`cem check --json`"+` — typecheck, lint, format
`, opts.ProjectName, opts.DB, opts.Validator, opts.PM)

	if err := writeFile(filepath.Join(root, "README.md"), readme); err != nil {
		return err
	}
	if err := writeFile(filepath.Join(root, "AGENTS.md"), agentsMd); err != nil {
		return err
	}
	return writeFile(filepath.Join(root, "CLAUDE.md"), claudeMd)
}
