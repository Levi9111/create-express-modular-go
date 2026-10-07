package generator

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"
)

func writeFile(path, content string) error {
	dir := filepath.Dir(path)
	if err := os.MkdirAll(dir, 0755); err != nil {
		return err
	}
	content = strings.ReplaceAll(content, "___BT___", "`")
	return os.WriteFile(path, []byte(content), 0644)
}

func scaffoldCore(opts ProjectOptions) error {
	root := opts.ProjectPath

	// 1. Errors
	err := writeFile(filepath.Join(root, "src/app/errors/AppError.ts"), `class AppError extends Error {
  public statusCode: number;

  constructor(statusCode: number, message: string, stack: string = '') {
    super(message);
    this.statusCode = statusCode;
    if (stack) {
      this.stack = stack;
    } else {
      Error.captureStackTrace(this, this.constructor);
    }
  }
}

export default AppError;
`)
	if err != nil {
		return err
	}

	// 2. Interfaces
	err = writeFile(filepath.Join(root, "src/app/interfaces/error.ts"), `export type TErrorSources = {
  path: string | number;
  message: string;
}[];

export type TGenericErrorResponse = {
  statusCode: number;
  message: string;
  errorSources: TErrorSources;
};
`)
	if err != nil {
		return err
	}

	// 3. Middlewares: notFound
	err = writeFile(filepath.Join(root, "src/app/middlewares/notFound.middleware.ts"), `import { NextFunction, Request, Response } from 'express';
import { StatusCodes } from 'http-status-codes';

const notFound = (_req: Request, res: Response, _next: NextFunction): void => {
  res.status(StatusCodes.NOT_FOUND).json({
    success: false,
    message: 'API Not Found!',
    error: '',
  });
};

export default notFound;
`)
	if err != nil {
		return err
	}

	// 4. Utils: catchAsync, sendResponse, logger, welcomePage
	err = writeFile(filepath.Join(root, "src/app/utils/catchAsync.ts"), `import { NextFunction, Request, RequestHandler, Response } from 'express';

export const catchAsync = (fn: RequestHandler) => {
  return (req: Request, res: Response, next: NextFunction) => {
    Promise.resolve(fn(req, res, next)).catch(next);
  };
};
`)
	if err != nil {
		return err
	}

	err = writeFile(filepath.Join(root, "src/app/utils/sendResponse.ts"), `import { Response } from 'express';

type TResponse<T> = {
  statusCode: number;
  success: boolean;
  message?: string;
  meta?: {
    page: number;
    limit: number;
    total: number;
    totalPages: number;
  };
  data: T | null;
};

const sendResponse = <T>(res: Response, data: TResponse<T>): void => {
  res.status(data.statusCode).json({
    success: data.success,
    message: data.message,
    meta: data.meta,
    data: data.data,
  });
};

export default sendResponse;
`)
	if err != nil {
		return err
	}

	err = writeFile(filepath.Join(root, "src/app/utils/logger.ts"), `const logger = {
  info: (...args: unknown[]) => console.log('\x1b[36m[INFO]\x1b[0m', ...args),
  warn: (...args: unknown[]) => console.warn('\x1b[33m[WARN]\x1b[0m', ...args),
  error: (...args: unknown[]) => console.error('\x1b[31m[ERROR]\x1b[0m', ...args),
  success: (...args: unknown[]) => console.log('\x1b[32m[SUCCESS]\x1b[0m', ...args),
};

export default logger;
`)
	if err != nil {
		return err
	}

	err = writeFile(filepath.Join(root, "src/app/utils/welcomePage.ts"), fmt.Sprintf(`export function cemWelcomePage(): string {
  return ___BT___<!DOCTYPE html>
<html lang="en">
<head>
  <meta charset="UTF-8">
  <title>%s — Modular Express</title>
  <style>
    body { font-family: system-ui, sans-serif; background: #0a0e17; color: #e2e8f0; display: grid; place-items: center; min-height: 100vh; margin: 0; }
    .card { background: #111827; border: 1px solid #1f2937; border-radius: 12px; padding: 2.5rem; max-width: 540px; text-align: center; box-shadow: 0 20px 25px -5px rgba(0, 0, 0, 0.5); }
    h1 { color: #38bdf8; margin: 0 0 0.5rem; font-size: 1.75rem; }
    p { color: #94a3b8; line-height: 1.6; margin: 0 0 1.5rem; }
    .badge { display: inline-block; padding: 0.25rem 0.75rem; border-radius: 9999px; font-size: 0.8rem; font-weight: 600; background: #0284c7; color: white; margin-bottom: 1.5rem; }
    .links a { display: inline-block; margin: 0 0.5rem; padding: 0.5rem 1rem; border-radius: 6px; text-decoration: none; font-weight: 500; font-size: 0.9rem; }
    .btn-primary { background: #0284c7; color: white; }
    .btn-secondary { background: #1f2937; color: #cbd5e1; }
  </style>
</head>
<body>
  <div class="card">
    <div class="badge">SYSTEM ONLINE</div>
    <h1>%s</h1>
    <p>Modular Express + TypeScript Backend scaffolded with Create Express Modular (Go Engine).</p>
    <div class="links">
      <a href="/api/v1/health" class="btn-secondary">Health Check</a>
      <a href="/docs" class="btn-primary">API Documentation</a>
    </div>
  </div>
</body>
</html>___BT___;
}
`, opts.ProjectName, opts.ProjectName))
	if err != nil {
		return err
	}

	// 5. Routes: src/app/routes/index.ts
	authImport := ""
	authRoute := ""
	if opts.UseAuth {
		authImport = "import { AuthRoutes } from '../modules/Auth/auth.route';\n"
		authRoute = "  { path: '/auth', route: AuthRoutes },\n"
	}

	err = writeFile(filepath.Join(root, "src/app/routes/index.ts"), fmt.Sprintf(`import { Router } from 'express';
%s// --- INJECT IMPORTS HERE ---

const router = Router();

const moduleRoutes: { path: string; route: Router }[] = [
%s  // --- INJECT ROUTES HERE ---
];

// Health check endpoint
router.get('/health', (_req, res) => {
  res.json({ status: 'ok', timestamp: new Date().toISOString() });
});

moduleRoutes.forEach((route) => router.use(route.path, route.route));

export default router;
`, authImport, authRoute))
	if err != nil {
		return err
	}

	// 6. Swagger config (if enabled)
	if opts.UseSwagger {
		err = writeFile(filepath.Join(root, "src/app/config/swagger.ts"), `import swaggerJsdoc from 'swagger-jsdoc';
import config from './index';

const options: swaggerJsdoc.Options = {
  definition: {
    openapi: '3.0.0',
    info: {
      title: 'CEM Modular Express API',
      version: '1.0.0',
      description: 'Interactive OpenAPI 3.0 documentation generated by Create Express Modular CLI.',
    },
    servers: [
      {
        url: ___BT___http://localhost:${config.port}/api/v1___BT___,
        description: 'Development Server',
      },
    ],
    components: {
      securitySchemes: {
        bearerAuth: {
          type: 'http',
          scheme: 'bearer',
          bearerFormat: 'JWT',
        },
      },
    },
  },
  apis: [
    './src/app/modules/**/*.route.ts',
    './src/app/modules/**/*.controller.ts',
  ],
};

export const swaggerSpec = swaggerJsdoc(options);
`)
		if err != nil {
			return err
		}
	}

	// 7. App entry: src/app.ts
	return scaffoldAppFile(opts)
}

func scaffoldAppFile(opts ProjectOptions) error {
	var imports []string
	imports = append(imports,
		"import express, { Application, Request, Response } from 'express';",
		"import cors from 'cors';",
		"import helmet from 'helmet';",
		"import compression from 'compression';",
		"import config from './app/config';",
		"import logger from './app/utils/logger';",
		"import { cemWelcomePage } from './app/utils/welcomePage';",
	)

	if opts.UseSwagger {
		imports = append(imports,
			"import swaggerUi from 'swagger-ui-express';",
			"import { swaggerSpec } from './app/config/swagger';",
		)
	}

	if opts.UseAuth && opts.TokenDelivery == TokenCookie {
		imports = append(imports, "import cookieParser from 'cookie-parser';")
	}

	if opts.UseAuth {
		imports = append(imports, "import { globalRateLimiter } from './app/middlewares/rateLimiter.middleware';")
	}

	imports = append(imports,
		"import router from './app/routes';",
		"import notFound from './app/middlewares/notFound.middleware';",
		"import globalErrorHandler from './app/middlewares/globalErrorHandler.middleware';",
		"",
		"const app: Application = express();",
		"",
		"app.disable('x-powered-by');",
		"",
		"app.use(helmet());",
		"app.use(compression());",
		"app.use(cors({ origin: config.cors_origin, credentials: true }));",
		"app.use(express.json());",
		"app.use(express.urlencoded({ extended: true }));",
	)

	if opts.UseAuth && opts.TokenDelivery == TokenCookie {
		imports = append(imports, "app.use(cookieParser());")
	}

	if opts.UseAuth {
		imports = append(imports, "app.use(globalRateLimiter);")
	}

	if opts.UseSwagger {
		imports = append(imports,
			"",
			"// Swagger UI Docs",
			"app.use('/docs', swaggerUi.serve, swaggerUi.setup(swaggerSpec));",
		)
	}

	imports = append(imports,
		"",
		"// Root welcome landing page",
		"app.get('/', (_req: Request, res: Response) => {",
		"  res.send(cemWelcomePage());",
		"});",
		"",
		"// API v1 router",
		"app.use('/api/v1', router);",
		"",
		"// Error handling",
		"app.use(notFound);",
		"app.use(globalErrorHandler);",
		"",
		"export default app;",
		"",
	)

	return writeFile(filepath.Join(opts.ProjectPath, "src/app.ts"), strings.Join(imports, "\n"))
}
