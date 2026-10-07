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
func AddModule(projectRoot string, moduleNames []string) error {
	cfg, err := config.LoadConfig(projectRoot)
	isZod := true
	if err == nil && cfg.Validator == "joi" {
		isZod = false
	}

	for _, rawName := range moduleNames {
		moduleName := capitalize(rawName)
		fileName := strings.ToLower(moduleName)
		moduleDir := filepath.Join(projectRoot, "src/app/modules", moduleName)

		if _, err := os.Stat(moduleDir); err == nil {
			ui.Warn(fmt.Sprintf("Module '%s' already exists, skipping.", moduleName))
			continue
		}

		// 1. Controller
		_ = writeFile(filepath.Join(moduleDir, fileName+".controller.ts"), fmt.Sprintf(`import { Request, Response } from 'express';
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
`, moduleName, fileName, moduleName, moduleName, moduleName, moduleName, moduleName, moduleName, moduleName, moduleName, moduleName, moduleName, moduleName, moduleName, moduleName, moduleName, moduleName, moduleName, moduleName, moduleName, moduleName, moduleName, moduleName, moduleName, moduleName, moduleName, moduleName, moduleName))

		// 2. Service
		_ = writeFile(filepath.Join(moduleDir, fileName+".service.ts"), fmt.Sprintf(`import { T%s } from './%s.interface';

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
`, moduleName, fileName, moduleName, moduleName, moduleName, moduleName, moduleName, moduleName, moduleName, moduleName, moduleName, moduleName, moduleName, moduleName, moduleName))

		// 3. Interface
		_ = writeFile(filepath.Join(moduleDir, fileName+".interface.ts"), fmt.Sprintf(`export type T%s = {
  id?: string;
  name: string;
  createdAt?: Date;
  updatedAt?: Date;
};
`, moduleName))

		// 4. Validation
		if isZod {
			_ = writeFile(filepath.Join(moduleDir, fileName+".validation.ts"), fmt.Sprintf(`import { z } from 'zod';

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
`, moduleName, moduleName, moduleName, moduleName, moduleName))
		} else {
			_ = writeFile(filepath.Join(moduleDir, fileName+".validation.ts"), fmt.Sprintf(`import Joi from 'joi';

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
`, moduleName, moduleName, moduleName, moduleName, moduleName))
		}

		// 5. Constants
		_ = writeFile(filepath.Join(moduleDir, fileName+".constant.ts"), fmt.Sprintf(`export const %sSearchableFields = ['name'];
`, moduleName))

		// 6. Route
		routePath := "/" + pluralize(moduleName)
		_ = writeFile(filepath.Join(moduleDir, fileName+".route.ts"), fmt.Sprintf(`import { Router } from 'express';
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
`, moduleName, fileName, moduleName, fileName, moduleName, moduleName, moduleName, moduleName, moduleName, moduleName, moduleName, moduleName, moduleName, moduleName, moduleName, moduleName, moduleName, moduleName, moduleName))

		// 7. Inject Route into src/app/routes/index.ts
		importLine := fmt.Sprintf("import { %sRoutes } from '../modules/%s/%s.route';", moduleName, moduleName, fileName)
		routeLine := fmt.Sprintf("  { path: '%s', route: %sRoutes },", routePath, moduleName)
		_ = InjectRoute(projectRoot, importLine, routeLine)

		ui.Success(fmt.Sprintf("Module '%s' created and route wired to '%s'", moduleName, routePath))
	}
	return nil
}

// RemoveModule deletes module folder(s) and unwires routes.
func RemoveModule(projectRoot string, moduleNames []string) error {
	for _, rawName := range moduleNames {
		moduleName := capitalize(rawName)
		moduleDir := filepath.Join(projectRoot, "src/app/modules", moduleName)

		if _, err := os.Stat(moduleDir); os.IsNotExist(err) {
			ui.Warn(fmt.Sprintf("Module '%s' not found under src/app/modules/", moduleName))
			continue
		}

		if err := os.RemoveAll(moduleDir); err != nil {
			return fmt.Errorf("failed to delete module directory: %w", err)
		}

		_ = RemoveRoute(projectRoot, moduleName)
		ui.Success(fmt.Sprintf("Module '%s' removed and route unwired", moduleName))
	}
	return nil
}

// AddMiddleware scaffolds a new custom middleware file.
func AddMiddleware(projectRoot string, middlewareNames []string) error {
	mwDir := filepath.Join(projectRoot, "src/app/middlewares")
	for _, name := range middlewareNames {
		cleanName := strings.TrimSuffix(name, ".middleware")
		cleanName = strings.TrimSuffix(cleanName, ".ts")
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
			return err
		}
		ui.Success(fmt.Sprintf("Middleware '%s' created at src/app/middlewares/%s.middleware.ts", cleanName, cleanName))
	}
	return nil
}

// RemoveMiddleware removes one or more custom middleware files safely.
func RemoveMiddleware(projectRoot string, middlewareNames []string) error {
	protected := map[string]bool{
		"globalErrorHandler": true,
		"notFound":           true,
		"auth":               true,
		"rateLimiter":        true,
	}

	mwDir := filepath.Join(projectRoot, "src/app/middlewares")
	for _, name := range middlewareNames {
		cleanName := strings.TrimSuffix(name, ".middleware")
		cleanName = strings.TrimSuffix(cleanName, ".ts")

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
			return err
		}
		ui.Success(fmt.Sprintf("Middleware '%s' deleted", cleanName))
	}
	return nil
}

// AddEnv adds environment variable(s) to .env, .env.example, and config/index.ts.
func AddEnv(projectRoot string, keys []string) error {
	for _, rawKey := range keys {
		key := strings.ToUpper(strings.TrimSpace(rawKey))
		if key == "" {
			continue
		}

		// Append to .env
		envPath := filepath.Join(projectRoot, ".env")
		if data, err := os.ReadFile(envPath); err == nil {
			content := string(data)
			if !strings.Contains(content, key+"=") {
				content += fmt.Sprintf("%s=\n", key)
				_ = os.WriteFile(envPath, []byte(content), 0644)
			}
		}

		// Append to .env.example
		exPath := filepath.Join(projectRoot, ".env.example")
		if data, err := os.ReadFile(exPath); err == nil {
			content := string(data)
			if !strings.Contains(content, key+"=") {
				content += fmt.Sprintf("%s=\n", key)
				_ = os.WriteFile(exPath, []byte(content), 0644)
			}
		}

		ui.Success(fmt.Sprintf("Added '%s' to .env and .env.example", key))
	}
	return nil
}

// RemoveEnv removes environment variable(s) from .env and .env.example.
func RemoveEnv(projectRoot string, keys []string) error {
	for _, rawKey := range keys {
		key := strings.ToUpper(strings.TrimSpace(rawKey))
		if key == "" {
			continue
		}

		keyRegex := regexp.MustCompile(fmt.Sprintf(`(?m)^%s=.*(?:\r?\n)?`, key))

		for _, fname := range []string{".env", ".env.example"} {
			fpath := filepath.Join(projectRoot, fname)
			if data, err := os.ReadFile(fpath); err == nil {
				cleaned := keyRegex.ReplaceAllString(string(data), "")
				_ = os.WriteFile(fpath, []byte(cleaned), 0644)
			}
		}
		ui.Success(fmt.Sprintf("Removed '%s' from .env and .env.example", key))
	}
	return nil
}

// ListFeatures lists modules, middlewares, and env variables in the project.
func ListFeatures(projectRoot string) error {
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
