package generator

import (
	"fmt"
	"path/filepath"
)

func scaffoldDB(opts ProjectOptions) (ErrorBlock, error) {
	root := opts.ProjectPath

	switch opts.DB {
	case DbMongoose:
		return scaffoldMongoose(root, opts)
	case DbPrisma:
		return scaffoldPrisma(root, opts)
	case DbDrizzle:
		return scaffoldDrizzle(root, opts)
	default:
		return scaffoldMongoose(root, opts)
	}
}

func scaffoldMongoose(root string, opts ProjectOptions) (ErrorBlock, error) {
	envContent := fmt.Sprintf(`PORT=5000
NODE_ENV=development
DATABASE_URL=mongodb://localhost:27017/%s
CORS_ORIGIN=http://localhost:3000
BCRYPT_SALT_ROUNDS=12
JWT_ACCESS_SECRET=your_super_secret_access_key
JWT_ACCESS_EXPIRES_IN=1d
JWT_REFRESH_SECRET=your_jwt_refresh_secret
JWT_REFRESH_EXPIRES_IN=365d
`, opts.ProjectName)

	_ = writeFile(filepath.Join(root, ".env"), envContent)
	_ = writeFile(filepath.Join(root, ".env.example"), `PORT=5000
NODE_ENV=development
DATABASE_URL=
CORS_ORIGIN=http://localhost:3000
BCRYPT_SALT_ROUNDS=12
JWT_ACCESS_SECRET=
JWT_ACCESS_EXPIRES_IN=1d
JWT_REFRESH_SECRET=
JWT_REFRESH_EXPIRES_IN=365d
`)

	_ = writeFile(filepath.Join(root, "src/app/config/index.ts"), `import dotenv from 'dotenv';
import path from 'path';

dotenv.config({ path: path.join(process.cwd(), '.env') });

const requiredEnvVars = ['DATABASE_URL'] as const;
for (const key of requiredEnvVars) {
  if (!process.env[key]) {
    throw new Error(
      ___BT___Missing required environment variable: ${key}. ___BT___ +
        'Check your .env file or .env.example for reference.',
    );
  }
}

export default {
  NODE_ENV: process.env.NODE_ENV ?? 'development',
  port: process.env.PORT ?? 5000,
  databaseUrl: process.env.DATABASE_URL as string,
  cors_origin: process.env.CORS_ORIGIN?.split(',') || ['http://localhost:3000'],
  bcrypt_salt_rounds: Number(process.env.BCRYPT_SALT_ROUNDS || 12),
  jwt_access_secret: process.env.JWT_ACCESS_SECRET,
  jwt_access_expires_in: process.env.JWT_ACCESS_EXPIRES_IN,
  jwt_refresh_secret: process.env.JWT_REFRESH_SECRET,
  jwt_refresh_expires_in: process.env.JWT_REFRESH_EXPIRES_IN,
};
`)

	_ = writeFile(filepath.Join(root, "src/server.ts"), `import { Server } from 'http';
import app from './app';
import config from './app/config';
import logger from './app/utils/logger';
import { connectDB } from './app/utils/connectDB';

let server: Server;

async function bootstrap() {
  try {
    await connectDB();

    server = app.listen(config.port, () => {
      logger.info(___BT___Server running on http://localhost:${config.port}___BT___);
    });

    server.keepAliveTimeout = 65000;
    server.headersTimeout = 66000;
  } catch (error) {
    logger.error('Failed to start server:', error);
    process.exit(1);
  }
}

bootstrap();

process.on('unhandledRejection', (reason) => {
  logger.error('Unhandled Rejection detected, shutting down server...', reason);
  if (server) {
    server.close(() => process.exit(1));
  } else {
    process.exit(1);
  }
});

process.on('uncaughtException', (error) => {
  logger.error('Uncaught Exception detected, shutting down server...', error);
  process.exit(1);
});
`)

	_ = writeFile(filepath.Join(root, "src/app/utils/connectDB.ts"), `import mongoose from 'mongoose';
import config from '../config';
import logger from './logger';

export async function connectDB(): Promise<void> {
  try {
    await mongoose.connect(config.databaseUrl, {
      maxPoolSize: 10,
      serverSelectionTimeoutMS: 5000,
    });
    logger.success('MongoDB connected successfully');
  } catch (error) {
    logger.error('MongoDB connection error:', error);
    throw error;
  }
}
`)

	_ = writeFile(filepath.Join(root, "src/app/utils/QueryBuilder.ts"), `import { FilterQuery, Query } from 'mongoose';

class QueryBuilder<T> {
  public modelQuery: Query<T[], T>;
  public query: Record<string, unknown>;

  constructor(modelQuery: Query<T[], T>, query: Record<string, unknown>) {
    this.modelQuery = modelQuery;
    this.query = query;
  }

  search(searchableFields: string[]) {
    const searchTerm = this.query.searchTerm as string;
    if (searchTerm) {
      this.modelQuery = this.modelQuery.find({
        $or: searchableFields.map(
          (field) =>
            ({
              [field]: { $regex: searchTerm, $options: 'i' },
            }) as FilterQuery<T>,
        ),
      });
    }
    return this;
  }

  filter() {
    const queryObj = { ...this.query };
    const excludeFields = ['searchTerm', 'sort', 'limit', 'page', 'fields'];
    excludeFields.forEach((el) => delete queryObj[el]);

    this.modelQuery = this.modelQuery.find(queryObj as FilterQuery<T>);
    return this;
  }

  sort() {
    const sort = ((this.query.sort as string) || '-createdAt').split(',').join(' ');
    this.modelQuery = this.modelQuery.sort(sort);
    return this;
  }

  paginate() {
    const page = Number(this.query.page) || 1;
    const limit = Number(this.query.limit) || 10;
    const skip = (page - 1) * limit;

    this.modelQuery = this.modelQuery.skip(skip).limit(limit);
    return this;
  }

  fields() {
    const fields = ((this.query.fields as string) || '-__v').split(',').join(' ');
    this.modelQuery = this.modelQuery.select(fields);
    return this;
  }

  async countTotal() {
    const totalQueries = this.modelQuery.getFilter();
    const total = await this.modelQuery.model.countDocuments(totalQueries);
    const page = Number(this.query.page) || 1;
    const limit = Number(this.query.limit) || 10;
    const totalPages = Math.ceil(total / limit);

    return { page, limit, total, totalPages };
  }
}

export default QueryBuilder;
`)

	_ = writeFile(filepath.Join(root, "src/app/errors/handleCastError.ts"), `import mongoose from 'mongoose';
import { TGenericErrorResponse } from '../interfaces/error';

const handleCastError = (err: mongoose.Error.CastError): TGenericErrorResponse => {
  return {
    statusCode: 400,
    message: 'Invalid ID',
    errorSources: [{ path: err.path, message: err.message }],
  };
};

export default handleCastError;
`)

	_ = writeFile(filepath.Join(root, "src/app/errors/handleValidationError.ts"), `import mongoose from 'mongoose';
import { TErrorSources, TGenericErrorResponse } from '../interfaces/error';

const handleValidationError = (err: mongoose.Error.ValidationError): TGenericErrorResponse => {
  const errorSources: TErrorSources = Object.values(err.errors).map((val) => ({
    path: val.path,
    message: val.message,
  }));

  return {
    statusCode: 400,
    message: 'Validation Error',
    errorSources,
  };
};

export default handleValidationError;
`)

	_ = writeFile(filepath.Join(root, "src/app/errors/handleDuplicateError.ts"), `/* eslint-disable @typescript-eslint/no-explicit-any */
import { TGenericErrorResponse } from '../interfaces/error';

const handleDuplicateError = (err: any): TGenericErrorResponse => {
  const match = err.message.match(/"([^"]*)"/);
  const extracted = match && match[1];

  return {
    statusCode: 409,
    message: 'Duplicate Entry',
    errorSources: [
      {
        path: '',
        message: ___BT___${extracted} already exists___BT___,
      },
    ],
  };
};

export default handleDuplicateError;
`)

	return ErrorBlock{
		Imports: `import handleCastError from '../errors/handleCastError';
import handleValidationError from '../errors/handleValidationError';
import handleDuplicateError from '../errors/handleDuplicateError';`,
		Handler: `if (err?.name === 'ValidationError') {
    const simplified = handleValidationError(err);
    statusCode = simplified.statusCode;
    message = simplified.message;
    errorSources = simplified.errorSources;
  } else if (err?.name === 'CastError') {
    const simplified = handleCastError(err);
    statusCode = simplified.statusCode;
    message = simplified.message;
    errorSources = simplified.errorSources;
  } else if (err?.code === 11000) {
    const simplified = handleDuplicateError(err);
    statusCode = simplified.statusCode;
    message = simplified.message;
    errorSources = simplified.errorSources;
  } else `,
	}, nil
}

func scaffoldPrisma(root string, opts ProjectOptions) (ErrorBlock, error) {
	_ = writeFile(filepath.Join(root, ".env"), `PORT=5000
NODE_ENV=development
DATABASE_URL="postgresql://user:password@localhost:5432/mydb?schema=public"
CORS_ORIGIN=http://localhost:3000
BCRYPT_SALT_ROUNDS=12
JWT_ACCESS_SECRET=your_super_secret_access_key
JWT_ACCESS_EXPIRES_IN=1d
JWT_REFRESH_SECRET=your_jwt_refresh_secret
JWT_REFRESH_EXPIRES_IN=365d
`)

	_ = writeFile(filepath.Join(root, ".env.example"), `PORT=5000
NODE_ENV=development
DATABASE_URL="postgresql://user:password@localhost:5432/mydb?schema=public"
CORS_ORIGIN=http://localhost:3000
BCRYPT_SALT_ROUNDS=12
JWT_ACCESS_SECRET=
JWT_ACCESS_EXPIRES_IN=1d
JWT_REFRESH_SECRET=
JWT_REFRESH_EXPIRES_IN=365d
`)

	_ = writeFile(filepath.Join(root, "src/app/config/index.ts"), `import dotenv from 'dotenv';
import path from 'path';

dotenv.config({ path: path.join(process.cwd(), '.env') });

const requiredEnvVars = ['DATABASE_URL'] as const;
for (const key of requiredEnvVars) {
  if (!process.env[key]) {
    throw new Error(___BT___Missing required environment variable: ${key}.___BT___);
  }
}

export default {
  NODE_ENV: process.env.NODE_ENV ?? 'development',
  port: process.env.PORT ?? 5000,
  databaseUrl: process.env.DATABASE_URL as string,
  cors_origin: process.env.CORS_ORIGIN?.split(',') || ['http://localhost:3000'],
  bcrypt_salt_rounds: Number(process.env.BCRYPT_SALT_ROUNDS || 12),
  jwt_access_secret: process.env.JWT_ACCESS_SECRET,
  jwt_access_expires_in: process.env.JWT_ACCESS_EXPIRES_IN,
  jwt_refresh_secret: process.env.JWT_REFRESH_SECRET,
  jwt_refresh_expires_in: process.env.JWT_REFRESH_EXPIRES_IN,
};
`)

	_ = writeFile(filepath.Join(root, "src/server.ts"), `import { Server } from 'http';
import app from './app';
import config from './app/config';
import logger from './app/utils/logger';
import { prisma } from './app/utils/prisma';

let server: Server;

async function bootstrap() {
  try {
    await prisma.$connect();
    logger.success('Prisma connected to database');

    server = app.listen(config.port, () => {
      logger.info(___BT___Server running on http://localhost:${config.port}___BT___);
    });
    server.keepAliveTimeout = 65000;
    server.headersTimeout = 66000;
  } catch (error) {
    logger.error('Failed to start server:', error);
    process.exit(1);
  }
}

bootstrap();
`)

	_ = writeFile(filepath.Join(root, "prisma/schema.prisma"), `datasource db {
  provider = "postgresql"
  url      = env("DATABASE_URL")
}

generator client {
  provider = "prisma-client-js"
}

model User {
  id        String   @id @default(uuid())
  email     String   @unique
  password  String
  name      String?
  role      String   @default("user")
  createdAt DateTime @default(now())
  updatedAt DateTime @updatedAt
}
`)

	_ = writeFile(filepath.Join(root, "src/app/utils/prisma.ts"), `import { PrismaClient } from '@prisma/client';

export const prisma = new PrismaClient({
  log: process.env.NODE_ENV === 'development' ? ['query', 'error', 'warn'] : ['error'],
});
`)

	return ErrorBlock{
		Imports: "import { Prisma } from '@prisma/client';",
		Handler: `if (err instanceof Prisma.PrismaClientKnownRequestError) {
    if (err.code === 'P2002') {
      statusCode = 409;
      message = 'Duplicate field value entered';
      errorSources = [{ path: String(err.meta?.target ?? ''), message }];
    } else if (err.code === 'P2025') {
      statusCode = 404;
      message = 'Record not found';
      errorSources = [{ path: '', message }];
    } else {
      statusCode = 400;
      message = err.message;
      errorSources = [{ path: '', message: err.message }];
    }
  } else if (err instanceof Prisma.PrismaClientValidationError) {
    statusCode = 400;
    message = 'Prisma validation error';
    errorSources = [{ path: '', message: err.message }];
  } else `,
	}, nil
}

func scaffoldDrizzle(root string, opts ProjectOptions) (ErrorBlock, error) {
	_ = writeFile(filepath.Join(root, ".env"), `PORT=5000
NODE_ENV=development
DATABASE_URL="postgres://user:password@localhost:5432/mydb"
CORS_ORIGIN=http://localhost:3000
BCRYPT_SALT_ROUNDS=12
JWT_ACCESS_SECRET=your_super_secret_access_key
JWT_ACCESS_EXPIRES_IN=1d
JWT_REFRESH_SECRET=your_jwt_refresh_secret
JWT_REFRESH_EXPIRES_IN=365d
`)

	_ = writeFile(filepath.Join(root, ".env.example"), `PORT=5000
NODE_ENV=development
DATABASE_URL="postgres://user:password@localhost:5432/mydb"
CORS_ORIGIN=http://localhost:3000
BCRYPT_SALT_ROUNDS=12
JWT_ACCESS_SECRET=
JWT_ACCESS_EXPIRES_IN=1d
JWT_REFRESH_SECRET=
JWT_REFRESH_EXPIRES_IN=365d
`)

	_ = writeFile(filepath.Join(root, "src/app/config/index.ts"), `import dotenv from 'dotenv';
import path from 'path';

dotenv.config({ path: path.join(process.cwd(), '.env') });

const requiredEnvVars = ['DATABASE_URL'] as const;
for (const key of requiredEnvVars) {
  if (!process.env[key]) {
    throw new Error(___BT___Missing required environment variable: ${key}.___BT___);
  }
}

export default {
  NODE_ENV: process.env.NODE_ENV ?? 'development',
  port: process.env.PORT ?? 5000,
  databaseUrl: process.env.DATABASE_URL as string,
  cors_origin: process.env.CORS_ORIGIN?.split(',') || ['http://localhost:3000'],
  bcrypt_salt_rounds: Number(process.env.BCRYPT_SALT_ROUNDS || 12),
  jwt_access_secret: process.env.JWT_ACCESS_SECRET,
  jwt_access_expires_in: process.env.JWT_ACCESS_EXPIRES_IN,
  jwt_refresh_secret: process.env.JWT_REFRESH_SECRET,
  jwt_refresh_expires_in: process.env.JWT_REFRESH_EXPIRES_IN,
};
`)

	_ = writeFile(filepath.Join(root, "drizzle.config.ts"), `import { defineConfig } from 'drizzle-kit';
import config from './src/app/config';

export default defineConfig({
  schema: './src/app/db/schema.ts',
  out: './drizzle',
  dialect: 'postgresql',
  dbCredentials: {
    url: config.databaseUrl,
  },
});
`)

	_ = writeFile(filepath.Join(root, "src/app/db/index.ts"), `import { drizzle } from 'drizzle-orm/node-postgres';
import { Pool } from 'pg';
import config from '../config';
import * as schema from './schema';

export const pool = new Pool({
  connectionString: config.databaseUrl,
  max: 10,
  idleTimeoutMillis: 30000,
  connectionTimeoutMillis: 5000,
});

export const db = drizzle(pool, { schema });
`)

	_ = writeFile(filepath.Join(root, "src/app/db/schema.ts"), `import { pgTable, text, timestamp, uuid } from 'drizzle-orm/pg-core';

export const users = pgTable('users', {
  id: uuid('id').defaultRandom().primaryKey(),
  email: text('email').notNull().unique(),
  password: text('password').notNull(),
  name: text('name'),
  role: text('role').default('user').notNull(),
  createdAt: timestamp('created_at').defaultNow().notNull(),
  updatedAt: timestamp('updated_at').defaultNow().notNull(),
});
`)

	_ = writeFile(filepath.Join(root, "src/server.ts"), `import { Server } from 'http';
import app from './app';
import config from './app/config';
import logger from './app/utils/logger';
import { pool } from './app/db';

let server: Server;

async function bootstrap() {
  try {
    await pool.query('SELECT 1');
    logger.success('PostgreSQL connected via Drizzle');

    server = app.listen(config.port, () => {
      logger.info(___BT___Server running on http://localhost:${config.port}___BT___);
    });
    server.keepAliveTimeout = 65000;
    server.headersTimeout = 66000;
  } catch (error) {
    logger.error('Failed to start server:', error);
    process.exit(1);
  }
}

bootstrap();
`)

	return ErrorBlock{
		Imports: "",
		Handler: `if (err?.code === '23505') {
    statusCode = 409;
    message = 'Duplicate key value violates unique constraint';
    errorSources = [{ path: String(err?.detail ?? ''), message }];
  } else if (err?.code === '23503') {
    statusCode = 400;
    message = 'Foreign key violation';
    errorSources = [{ path: String(err?.detail ?? ''), message }];
  } else `,
	}, nil
}
