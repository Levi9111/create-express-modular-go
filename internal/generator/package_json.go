package generator

import (
	"encoding/json"
	"os"
	"path/filepath"
)

type PackageJson struct {
	Name            string            `json:"name"`
	Version         string            `json:"version"`
	Description     string            `json:"description,omitempty"`
	Main            string            `json:"main,omitempty"`
	Scripts         map[string]string `json:"scripts"`
	Dependencies    map[string]string `json:"dependencies"`
	DevDependencies map[string]string `json:"devDependencies"`
	Overrides       map[string]string `json:"overrides,omitempty"`
	Resolutions     map[string]string `json:"resolutions,omitempty"`
}

var versionMap = map[string]string{
	// Core
	"express":           "^4.21.2",
	"cors":              "^2.8.5",
	"dotenv":            "^16.4.7",
	"http-status-codes": "^2.3.0",
	"helmet":            "^8.0.0",
	"compression":       "^1.8.0",
	// DB
	"mongoose":       "^8.9.5",
	"prisma":         "^6.2.1",
	"@prisma/client": "^6.2.1",
	"drizzle-orm":    "^0.38.4",
	"drizzle-kit":    "^0.30.2",
	"pg":             "^8.13.1",
	// Validator
	"zod": "^3.24.1",
	"joi": "^17.13.3",
	// Auth
	"bcrypt":             "^5.1.1",
	"jsonwebtoken":       "^9.0.2",
	"express-rate-limit": "^7.5.0",
	"cookie-parser":      "^1.4.7",
	// Swagger
	"swagger-ui-express": "^5.0.1",
	"swagger-jsdoc":      "^6.2.8",
	// Dev
	"@types/express":            "^5.0.0",
	"@types/cors":               "^2.8.17",
	"@types/compression":        "^1.7.5",
	"@types/node":               "^22.10.7",
	"@types/pg":                 "^8.11.10",
	"@types/joi":                "^17.2.3",
	"@types/bcrypt":             "^5.0.2",
	"@types/jsonwebtoken":       "^9.0.8",
	"@types/cookie-parser":      "^1.4.8",
	"@types/swagger-ui-express": "^4.1.7",
	"@types/swagger-jsdoc":      "^6.0.4",
	"typescript":                "^5.7.3",
	"tsx":                       "^4.19.2",
	"eslint":                    "^10.0.0",
	"@eslint/js":                "^10.0.0",
	"typescript-eslint":         "^8.20.0",
	"eslint-config-prettier":    "^10.0.1",
	"prettier":                  "^3.4.2",
	"create-express-modular":    "^3.3.10",
}

func configurePackageJson(opts ProjectOptions) error {
	pkgPath := filepath.Join(opts.ProjectPath, "package.json")
	data, err := os.ReadFile(pkgPath)
	if err != nil {
		return err
	}

	var pkg PackageJson
	if err := json.Unmarshal(data, &pkg); err != nil {
		return err
	}

	pkg.Name = opts.ProjectName
	if pkg.Dependencies == nil {
		pkg.Dependencies = make(map[string]string)
	}
	if pkg.DevDependencies == nil {
		pkg.DevDependencies = make(map[string]string)
	}

	// Core
	prodDeps := []string{"dotenv", "http-status-codes", "express", "cors", "helmet", "compression"}
	devDeps := []string{"@types/express", "@types/cors", "@types/compression", "@types/node", "typescript", "tsx", "eslint", "@eslint/js", "typescript-eslint", "eslint-config-prettier", "prettier"}

	// DB
	switch opts.DB {
	case DbMongoose:
		prodDeps = append(prodDeps, "mongoose")
	case DbPrisma:
		prodDeps = append(prodDeps, "@prisma/client")
		devDeps = append(devDeps, "prisma")
	case DbDrizzle:
		prodDeps = append(prodDeps, "drizzle-orm", "pg")
		devDeps = append(devDeps, "drizzle-kit", "@types/pg")
	}

	// Validator
	if opts.Validator == ValidatorZod {
		prodDeps = append(prodDeps, "zod")
	} else {
		prodDeps = append(prodDeps, "joi")
		devDeps = append(devDeps, "@types/joi")
	}

	// Auth
	if opts.UseAuth {
		prodDeps = append(prodDeps, "bcrypt", "jsonwebtoken", "express-rate-limit")
		devDeps = append(devDeps, "@types/bcrypt", "@types/jsonwebtoken")
		if opts.TokenDelivery == TokenCookie {
			prodDeps = append(prodDeps, "cookie-parser")
			devDeps = append(devDeps, "@types/cookie-parser")
		}
	}

	// Swagger
	if opts.UseSwagger {
		prodDeps = append(prodDeps, "swagger-ui-express", "swagger-jsdoc")
		devDeps = append(devDeps, "@types/swagger-ui-express", "@types/swagger-jsdoc")
	}

	// CEM CLI
	devDeps = append(devDeps, "create-express-modular")

	for _, dep := range prodDeps {
		if ver, ok := versionMap[dep]; ok {
			pkg.Dependencies[dep] = ver
		} else {
			pkg.Dependencies[dep] = "latest"
		}
	}

	for _, dep := range devDeps {
		if ver, ok := versionMap[dep]; ok {
			pkg.DevDependencies[dep] = ver
		} else {
			pkg.DevDependencies[dep] = "latest"
		}
	}

	// Overrides for glob vulnerability
	if pkg.Overrides == nil {
		pkg.Overrides = make(map[string]string)
	}
	pkg.Overrides["glob"] = "^13.0.6"

	if pkg.Resolutions == nil {
		pkg.Resolutions = make(map[string]string)
	}
	pkg.Resolutions["glob"] = "^13.0.6"

	out, err := json.MarshalIndent(pkg, "", "  ")
	if err != nil {
		return err
	}
	return os.WriteFile(pkgPath, out, 0644)
}
