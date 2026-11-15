package modules

import (
	"fmt"
	"os"
	"path/filepath"
	"regexp"
	"strings"

	"github.com/Levi9111/create-express-modular-go/internal/config"
	"github.com/Levi9111/create-express-modular-go/internal/ui"
)

func capitalize(s string) string {
	if s == "" {
		return ""
	}
	return strings.ToUpper(s[:1]) + s[1:]
}

func pluralize(s string) string {
	lower := strings.ToLower(s)
	if strings.HasSuffix(lower, "s") || strings.HasSuffix(lower, "x") || strings.HasSuffix(lower, "z") || strings.HasSuffix(lower, "ch") || strings.HasSuffix(lower, "sh") {
		return lower + "es"
	}
	if strings.HasSuffix(lower, "y") && len(lower) > 1 && !strings.ContainsAny(string(lower[len(lower)-2]), "aeiou") {
		return lower[:len(lower)-1] + "ies"
	}
	return lower + "s"
}

func writeFile(path, content string) error {
	dir := filepath.Dir(path)
	if err := os.MkdirAll(dir, 0755); err != nil {
		return err
	}
	return os.WriteFile(path, []byte(content), 0644)
}

// AddModule scaffolds one or more complete feature modules.
func AddModule(projectRoot string, moduleNames []string) ([]string, error) {
	cfg, err := config.LoadConfig(projectRoot)
	isZod := true
	if err == nil && cfg.Project.Validator == "joi" {
		isZod = false
	}

	added := []string{}
	for _, rawName := range moduleNames {
		if err := ValidateIdentifier(rawName); err != nil {
			return added, err
		}
		moduleName := capitalize(strings.TrimSpace(rawName))
		fileName := strings.ToLower(moduleName)
		moduleDir := filepath.Join(projectRoot, "src/app/modules", moduleName)

		if _, err := os.Stat(moduleDir); err == nil {
			ui.Warn(fmt.Sprintf("Module '%s' already exists, skipping.", moduleName))
			continue
		}

		// 1. Controller
		if err := writeFile(filepath.Join(moduleDir, fileName+".controller.ts"), fmt.Sprintf(`import { Request, Response } from 'express';
import { StatusCodes } from 'http-status-codes';
import { catchAsync } from '../../utils/catchAsync';
import sendResponse from '../../utils/sendResponse';
import { %sService } from './%s.service';

const create%s = catchAsync(async (req: Request, res: Response) => {
  const result = await %sService.create%s(req.body);
  sendResponse(res, {
    statusCode: StatusCodes.CREATED,
    success: true,
    message: '%s created successfully',
    data: result,
  });
});

const getAll%ss = catchAsync(async (req: Request, res: Response) => {
  const result = await %sService.getAll%ss(req.query);
  sendResponse(res, {
    statusCode: StatusCodes.OK,
    success: true,
    message: '%ss retrieved successfully',
    data: result,
  });
});

const getSingle%s = catchAsync(async (req: Request, res: Response) => {
  const result = await %sService.getSingle%s(req.params.id as string);
  sendResponse(res, {
    statusCode: StatusCodes.OK,
    success: true,
    message: '%s retrieved successfully',
    data: result,
  });
});

const update%s = catchAsync(async (req: Request, res: Response) => {
  const result = await %sService.update%s(req.params.id as string, req.body);
  sendResponse(res, {
    statusCode: StatusCodes.OK,
    success: true,
    message: '%s updated successfully',
    data: result,
  });
});

const delete%s = catchAsync(async (req: Request, res: Response) => {
  const result = await %sService.delete%s(req.params.id as string);
  sendResponse(res, {
    statusCode: StatusCodes.OK,
    success: true,
    message: '%s deleted successfully',
    data: result,
  });
});

export const %sController = {
  create%s,
  getAll%ss,
  getSingle%s,
  update%s,
  delete%s,
};
`, moduleName, fileName, moduleName, moduleName, moduleName, moduleName, moduleName, moduleName, moduleName, moduleName, moduleName, moduleName, moduleName, moduleName, moduleName, moduleName, moduleName, moduleName, moduleName, moduleName, moduleName, moduleName, moduleName, moduleName, moduleName, moduleName, moduleName, moduleName)); err != nil {
			return added, err
		}

		// 2. Service
		if err := writeFile(filepath.Join(moduleDir, fileName+".service.ts"), fmt.Sprintf(`import { T%s } from './%s.interface';

const create%s = async (payload: T%s) => {
  // TODO: Implement database insert logic
  return payload;
};

const getAll%ss = async (_query: Record<string, unknown>) => {
  // TODO: Implement database query logic
  return [];
};

const getSingle%s = async (id: string) => {
  // TODO: Implement database findById logic
  return { id };
};

const update%s = async (id: string, payload: Partial<T%s>) => {
  // TODO: Implement database update logic
  return { id, ...payload };
};

const delete%s = async (id: string) => {
  // TODO: Implement database delete logic
  return { id };
};

export const %sService = {
  create%s,
  getAll%ss,
  getSingle%s,
  update%s,
  delete%s,
};
`, moduleName, fileName, moduleName, moduleName, moduleName, moduleName, moduleName, moduleName, moduleName, moduleName, moduleName, moduleName, moduleName, moduleName, moduleName)); err != nil {
			return added, err
		}

		// 3. Interface
		if err := writeFile(filepath.Join(moduleDir, fileName+".interface.ts"), fmt.Sprintf(`export type T%s = {
  id?: string;
  name: string;
  createdAt?: Date;
  updatedAt?: Date;
};
`, moduleName)); err != nil {
			return added, err
		}

		// 4. Validation
		if isZod {
			if err := writeFile(filepath.Join(moduleDir, fileName+".validation.ts"), fmt.Sprintf(`import { z } from 'zod';

const create%sValidationSchema = z.object({
  body: z.object({
    name: z.string({ required_error: 'Name is required' }),
  }),
});

const update%sValidationSchema = z.object({
  body: z.object({
    name: z.string().optional(),
  }),
});

export const %sValidation = {
  create%sValidationSchema,
  update%sValidationSchema,
};
`, moduleName, moduleName, moduleName, moduleName, moduleName)); err != nil {
				return added, err
			}
		} else {
			if err := writeFile(filepath.Join(moduleDir, fileName+".validation.ts"), fmt.Sprintf(`import Joi from 'joi';

const create%sValidationSchema = Joi.object({
  body: Joi.object({
    name: Joi.string().required(),
  }),
});

const update%sValidationSchema = Joi.object({
  body: Joi.object({
    name: Joi.string().optional(),
  }),
});

export const %sValidation = {
  create%sValidationSchema,
  update%sValidationSchema,
};
`, moduleName, moduleName, moduleName, moduleName, moduleName)); err != nil {
				return added, err
			}
		}

		// 5. Constants
		if err := writeFile(filepath.Join(moduleDir, fileName+".constant.ts"), fmt.Sprintf(`export const %sSearchableFields = ['name'];
`, moduleName)); err != nil {
			return added, err
		}

		// 6. Route
		routePath := "/" + pluralize(moduleName)
		if err := writeFile(filepath.Join(moduleDir, fileName+".route.ts"), fmt.Sprintf(`import { Router } from 'express';
import { %sController } from './%s.controller';
import validateRequest from '../../utils/validateRequest';
import { %sValidation } from './%s.validation';

const router = Router();

router.post(
  '/',
  validateRequest(%sValidation.create%sValidationSchema),
  %sController.create%s,
);

router.get('/', %sController.getAll%ss);
router.get('/:id', %sController.getSingle%s);

router.patch(
  '/:id',
  validateRequest(%sValidation.update%sValidationSchema),
  %sController.update%s,
);

router.delete('/:id', %sController.delete%s);

export const %sRoutes = router;
`, moduleName, fileName, moduleName, fileName, moduleName, moduleName, moduleName, moduleName, moduleName, moduleName, moduleName, moduleName, moduleName, moduleName, moduleName, moduleName, moduleName, moduleName, moduleName)); err != nil {
			return added, err
		}

		// 7. Inject Route into src/app/routes/index.ts
		importLine := fmt.Sprintf("import { %sRoutes } from '../modules/%s/%s.route';", moduleName, moduleName, fileName)
		routeLine := fmt.Sprintf("  { path: '%s', route: %sRoutes },", routePath, moduleName)
		if err := InjectRoute(projectRoot, importLine, routeLine); err != nil {
			return added, err
		}

		ui.Success(fmt.Sprintf("Module '%s' created and route wired to '%s'", moduleName, routePath))
		added = append(added, moduleName)
	}
	config.UpdateState(projectRoot, "modules", "add", added)
	return added, nil
}

// RemoveModule deletes module folder(s) and unwires routes.
func RemoveModule(projectRoot string, moduleNames []string) ([]string, error) {
	removed := []string{}
	for _, rawName := range moduleNames {
		if err := ValidateIdentifier(rawName); err != nil {
			return removed, err
		}
		moduleName := capitalize(strings.TrimSpace(rawName))
		moduleDir := filepath.Join(projectRoot, "src/app/modules", moduleName)

		if _, err := os.Stat(moduleDir); os.IsNotExist(err) {
			ui.Warn(fmt.Sprintf("Module '%s' not found under src/app/modules/", moduleName))
			continue
		}

		if err := os.RemoveAll(moduleDir); err != nil {
			return removed, fmt.Errorf("failed to delete module directory: %w", err)
		}

		if err := RemoveRoute(projectRoot, moduleName); err != nil {
			return removed, fmt.Errorf("failed to unwire route: %w", err)
		}
		ui.Success(fmt.Sprintf("Module '%s' removed and route unwired", moduleName))
		removed = append(removed, moduleName)
	}
	config.UpdateState(projectRoot, "modules", "remove", removed)
	return removed, nil
}

// AddMiddleware scaffolds a new custom middleware file.
func AddMiddleware(projectRoot string, middlewareNames []string) ([]string, error) {
	mwDir := filepath.Join(projectRoot, "src/app/middlewares")
	added := []string{}
	for _, name := range middlewareNames {
		cleanName := strings.TrimSuffix(strings.TrimSpace(name), ".middleware")
		cleanName = strings.TrimSuffix(cleanName, ".ts")
		if err := ValidateIdentifier(cleanName); err != nil {
			return added, err
		}
		filePath := filepath.Join(mwDir, cleanName+".middleware.ts")

		if _, err := os.Stat(filePath); err == nil {
			ui.Warn(fmt.Sprintf("Middleware '%s' already exists.", cleanName))
			continue
		}

		content := fmt.Sprintf(`import { NextFunction, Request, Response } from 'express';

export const %s = (req: Request, _res: Response, next: NextFunction): void => {
  // TODO: Add custom middleware logic here
  next();
};
`, cleanName)

		if err := writeFile(filePath, content); err != nil {
			return added, err
		}
		ui.Success(fmt.Sprintf("Middleware '%s' created at src/app/middlewares/%s.middleware.ts", cleanName, cleanName))
		added = append(added, cleanName)
	}
	config.UpdateState(projectRoot, "middlewares", "add", added)
	return added, nil
}

// RemoveMiddleware removes one or more custom middleware files safely.
func RemoveMiddleware(projectRoot string, middlewareNames []string) ([]string, error) {
	protected := map[string]bool{
		"globalErrorHandler": true,
		"notFound":           true,
		"auth":               true,
		"rateLimiter":        true,
	}

	mwDir := filepath.Join(projectRoot, "src/app/middlewares")
	removed := []string{}
	for _, name := range middlewareNames {
		cleanName := strings.TrimSuffix(strings.TrimSpace(name), ".middleware")
		cleanName = strings.TrimSuffix(cleanName, ".ts")
		if err := ValidateIdentifier(cleanName); err != nil {
			return removed, err
		}

		if protected[cleanName] {
			ui.Err(fmt.Sprintf("Cannot remove protected core middleware '%s'", cleanName))
			continue
		}

		filePath := filepath.Join(mwDir, cleanName+".middleware.ts")
		if _, err := os.Stat(filePath); os.IsNotExist(err) {
			ui.Warn(fmt.Sprintf("Middleware '%s' does not exist.", cleanName))
			continue
		}

		if err := os.Remove(filePath); err != nil {
			return removed, err
		}
		ui.Success(fmt.Sprintf("Middleware '%s' deleted", cleanName))
		removed = append(removed, cleanName)
	}
	config.UpdateState(projectRoot, "middlewares", "remove", removed)
	return removed, nil
}

// AddEnv adds environment variable(s) to .env, .env.example, and config/index.ts.
func AddEnv(projectRoot string, keys []string) ([]string, error) {
	added := []string{}
	for _, rawKey := range keys {
		key := strings.ToUpper(strings.TrimSpace(rawKey))
		if err := ValidateEnvKey(key); err != nil {
			return added, err
		}

		// Append to .env
		envPath := filepath.Join(projectRoot, ".env")
		if data, err := os.ReadFile(envPath); err == nil {
			content := string(data)
			if !strings.Contains(content, key+"=") {
				content += fmt.Sprintf("%s=\n", key)
				if err := os.WriteFile(envPath, []byte(content), 0644); err != nil {
					return added, err
				}
			}
		}

		// Append to .env.example
		exPath := filepath.Join(projectRoot, ".env.example")
		if data, err := os.ReadFile(exPath); err == nil {
			content := string(data)
			if !strings.Contains(content, key+"=") {
				content += fmt.Sprintf("%s=\n", key)
				if err := os.WriteFile(exPath, []byte(content), 0644); err != nil {
					return added, err
				}
			}
		}

		// Append to src/app/config/index.ts
		cfgPath := filepath.Join(projectRoot, "src/app/config/index.ts")
		if data, err := os.ReadFile(cfgPath); err == nil {
			content := string(data)
			propName := strings.ToLower(key)
			if !strings.Contains(content, "process.env."+key) {
				entry := fmt.Sprintf("  %s: process.env.%s,\n", propName, key)
				if idx := strings.LastIndex(content, "};"); idx != -1 {
					updated := content[:idx] + entry + content[idx:]
					if err := os.WriteFile(cfgPath, []byte(updated), 0644); err != nil {
						return added, err
					}
				}
			}
		}

		ui.Success(fmt.Sprintf("Added '%s' to .env, .env.example, and src/app/config/index.ts", key))
		added = append(added, key)
	}
	config.UpdateState(projectRoot, "env", "add", added)
	return added, nil
}

// RemoveEnv removes environment variable(s) from .env and .env.example.
func RemoveEnv(projectRoot string, keys []string) ([]string, error) {
	removed := []string{}
	for _, rawKey := range keys {
		key := strings.ToUpper(strings.TrimSpace(rawKey))
		if err := ValidateEnvKey(key); err != nil {
			return removed, err
		}

		keyRegex := regexp.MustCompile(fmt.Sprintf(`(?m)^%s=.*(?:\r?\n)?`, regexp.QuoteMeta(key)))

		for _, fname := range []string{".env", ".env.example"} {
			fpath := filepath.Join(projectRoot, fname)
			if data, err := os.ReadFile(fpath); err == nil {
				cleaned := keyRegex.ReplaceAllString(string(data), "")
				if err := os.WriteFile(fpath, []byte(cleaned), 0644); err != nil {
					return removed, err
				}
			}
		}

		// Remove from src/app/config/index.ts
		cfgPath := filepath.Join(projectRoot, "src/app/config/index.ts")
		if data, err := os.ReadFile(cfgPath); err == nil {
			cfgRegex := regexp.MustCompile(fmt.Sprintf(`(?m)^\s*%s:\s*process\.env\.%s,?\r?\n?`, regexp.QuoteMeta(strings.ToLower(key)), regexp.QuoteMeta(key)))
			cleaned := cfgRegex.ReplaceAllString(string(data), "")
			if err := os.WriteFile(cfgPath, []byte(cleaned), 0644); err != nil {
				return removed, err
			}
		}

		ui.Success(fmt.Sprintf("Removed '%s' from .env, .env.example, and src/app/config/index.ts", key))
		removed = append(removed, key)
	}
	config.UpdateState(projectRoot, "env", "remove", removed)
	return removed, nil
}

// ListFeatures lists modules, middlewares, and env variables in the project.
func ListFeatures(projectRoot string) error {
	if ui.GetOutputMode() == ui.ModeJSON {
		cfg, _ := config.LoadConfig(projectRoot)
		modList := []string{}
		if entries, err := os.ReadDir(filepath.Join(projectRoot, "src/app/modules")); err == nil {
			for _, e := range entries {
				if e.IsDir() {
					modList = append(modList, e.Name())
				}
			}
		}
		mwList := []string{}
		if entries, err := os.ReadDir(filepath.Join(projectRoot, "src/app/middlewares")); err == nil {
			for _, e := range entries {
				if !e.IsDir() && strings.HasSuffix(e.Name(), ".ts") {
					mwList = append(mwList, strings.TrimSuffix(e.Name(), ".middleware.ts"))
				}
			}
		}
		envList := []string{}
		if data, err := os.ReadFile(filepath.Join(projectRoot, ".env.example")); err == nil {
			for _, line := range strings.Split(string(data), "\n") {
				if t := strings.TrimSpace(line); t != "" && !strings.HasPrefix(t, "#") {
					if parts := strings.SplitN(t, "=", 2); len(parts) > 0 {
						envList = append(envList, strings.TrimSpace(parts[0]))
					}
				}
			}
		}
		ui.PrintJSON(map[string]any{
			"ok":          true,
			"config":      cfg,
			"modules":     modList,
			"middlewares": mwList,
			"env":         envList,
		})
		return nil
	}

	ui.SectionHeader("CEM Project Snapshot")

	// Modules
	fmt.Println(ui.Bold("Feature Modules:"))
	modulesDir := filepath.Join(projectRoot, "src/app/modules")
	if entries, err := os.ReadDir(modulesDir); err == nil {
		for _, e := range entries {
			if e.IsDir() {
				fmt.Printf("  %s %s\n", ui.Cyan("•"), e.Name())
			}
		}
	} else {
		fmt.Println("  (No modules directory found)")
	}
	fmt.Println()

	// Middlewares
	fmt.Println(ui.Bold("Middlewares:"))
	mwDir := filepath.Join(projectRoot, "src/app/middlewares")
	if entries, err := os.ReadDir(mwDir); err == nil {
		for _, e := range entries {
			if !e.IsDir() && strings.HasSuffix(e.Name(), ".ts") {
				name := strings.TrimSuffix(e.Name(), ".middleware.ts")
				fmt.Printf("  %s %s\n", ui.Green("•"), name)
			}
		}
	}
	fmt.Println()

	// Env vars
	fmt.Println(ui.Bold("Environment Variables (.env.example):"))
	exPath := filepath.Join(projectRoot, ".env.example")
	if data, err := os.ReadFile(exPath); err == nil {
		lines := strings.Split(string(data), "\n")
		for _, line := range lines {
			line = strings.TrimSpace(line)
			if line != "" && !strings.HasPrefix(line, "#") {
				parts := strings.SplitN(line, "=", 2)
				fmt.Printf("  %s %s\n", ui.Yellow("•"), parts[0])
			}
		}
	}
	fmt.Println()
	return nil
}
