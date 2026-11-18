package generator

import (
	"fmt"
	"path/filepath"
)

func scaffoldAuth(opts ProjectOptions) error {
	root := opts.ProjectPath
	authDir := filepath.Join(root, "src/app/modules/Auth")
	utilsDir := filepath.Join(root, "src/app/utils")
	mwDir := filepath.Join(root, "src/app/middlewares")

	// 1. jwt.utils.ts
	err := writeFile(filepath.Join(utilsDir, "jwt.utils.ts"), `import jwt, { JwtPayload, SignOptions } from 'jsonwebtoken';

export const createToken = (
  jwtPayload: { userId: string; role: string },
  secret: string,
  expiresIn: string,
): string => {
  return jwt.sign(jwtPayload, secret, {
    expiresIn,
  } as SignOptions);
};

export const verifyToken = (token: string, secret: string): JwtPayload => {
  return jwt.verify(token, secret) as JwtPayload;
};
`)
	if err != nil {
		return err
	}

	// 2. auth.middleware.ts
	tokenExtraction := ""
	if opts.TokenDelivery == TokenCookie {
		tokenExtraction = `
      const token = req.cookies?.accessToken;

      if (!token) {
        return next(
          new AppError(StatusCodes.UNAUTHORIZED, 'You are not authorized!'),
        );
      }`
	} else {
		tokenExtraction = `
      const authHeader = req.headers.authorization;

      if (!authHeader?.startsWith('Bearer ')) {
        return next(
          new AppError(StatusCodes.UNAUTHORIZED, 'You are not authorized!'),
        );
      }

      const token = authHeader.split(' ')[1];`
	}

	err = writeFile(filepath.Join(mwDir, "auth.middleware.ts"), fmt.Sprintf(`import { NextFunction, Request, Response } from 'express';
import { StatusCodes } from 'http-status-codes';
import { JwtPayload } from 'jsonwebtoken';
import AppError from '../errors/AppError';
import { verifyToken } from '../utils/jwt.utils';
import config from '../config';

declare global {
  // eslint-disable-next-line @typescript-eslint/no-namespace
  namespace Express {
    interface Request {
      user: JwtPayload & { userId: string; role: string };
    }
  }
}

const auth = (...requiredRoles: string[]) => {
  return (req: Request, _res: Response, next: NextFunction): void => {
    try {%s

      const decoded = verifyToken(token, config.jwt_access_secret as string);

      if (requiredRoles.length && !requiredRoles.includes(decoded.role)) {
        return next(
          new AppError(
            StatusCodes.FORBIDDEN,
            'You do not have the required permissions!',
          ),
        );
      }

      req.user = decoded as JwtPayload & { userId: string; role: string };
      next();
    } catch {
      next(new AppError(StatusCodes.UNAUTHORIZED, 'You are not authorized!'));
    }
  };
};

export default auth;
`, tokenExtraction))
	if err != nil {
		return err
	}

	// 3. rateLimiter.middleware.ts
	err = writeFile(filepath.Join(mwDir, "rateLimiter.middleware.ts"), `import rateLimit from 'express-rate-limit';

export const globalRateLimiter = rateLimit({
  windowMs: 15 * 60 * 1000,
  max: 100,
  message: {
    success: false,
    message:
      'Too many requests from this IP, please try again after 15 minutes',
  },
  standardHeaders: true,
  legacyHeaders: false,
});

export const authRateLimiter = rateLimit({
  windowMs: 15 * 60 * 1000,
  max: 5,
  message: {
    success: false,
    message:
      'Too many login attempts from this IP, please try again after 15 minutes',
  },
  standardHeaders: true,
  legacyHeaders: false,
});
`)
	if err != nil {
		return err
	}

	// 4. auth.interface.ts
	err = writeFile(filepath.Join(authDir, "auth.interface.ts"), `export type TUserRole = 'user' | 'admin';

export type TUser = {
  id: string;
  email: string;
  password?: string;
  role: TUserRole;
  createdAt: Date;
  updatedAt: Date;
};

export type TLoginUser = {
  email: string;
  password?: string;
};
`)
	if err != nil {
		return err
	}

	// 5. auth.controller.ts
	cookieLogic := ""
	if opts.TokenDelivery == TokenCookie {
		cookieLogic = `
  const { refreshToken, accessToken } = result;
  const isProduction = config.NODE_ENV === 'production';
  const cookieOptions = {
    httpOnly: true,
    secure: isProduction,
    sameSite: 'strict' as const,
    path: '/',
  };

  res.cookie('accessToken', accessToken, {
    ...cookieOptions,
    maxAge: 15 * 60 * 1000,
  });

  res.cookie('refreshToken', refreshToken, {
    ...cookieOptions,
    maxAge: 7 * 24 * 60 * 60 * 1000,
  });
`
	}

	err = writeFile(filepath.Join(authDir, "auth.controller.ts"), fmt.Sprintf(`import { Request, Response } from 'express';
import { StatusCodes } from 'http-status-codes';
import { catchAsync } from '../../utils/catchAsync';
import sendResponse from '../../utils/sendResponse';
import { AuthService } from './auth.service';
import config from '../../config';

const registerUser = catchAsync(async (req: Request, res: Response) => {
  const result = await AuthService.registerUser(req.body);
  sendResponse(res, {
    statusCode: StatusCodes.CREATED,
    success: true,
    message: 'User registered successfully',
    data: result,
  });
});

const loginUser = catchAsync(async (req: Request, res: Response) => {
  const result = await AuthService.loginUser(req.body);%s
  sendResponse(res, {
    statusCode: StatusCodes.OK,
    success: true,
    message: 'User logged in successfully',
    data: result,
  });
});

export const AuthControllers = {
  registerUser,
  loginUser,
};
`, cookieLogic))
	if err != nil {
		return err
	}

	// 6. auth.model.ts
	modelContent := ""
	switch opts.DB {
	case DbMongoose:
		modelContent = `import { Schema, model } from 'mongoose';
import bcrypt from 'bcrypt';
import { TUserRole } from './auth.interface';

export interface IUser {
  _id?: string;
  email: string;
  password: string;
  role: TUserRole;
  createdAt?: Date;
  updatedAt?: Date;
  comparePassword(plainPassword: string): Promise<boolean>;
}

const userSchema = new Schema<IUser>(
  {
    email: {
      type: String,
      required: true,
      unique: true,
      lowercase: true,
      trim: true,
      match: [/^\S+@\S+\.\S+$/, 'Please provide a valid email'],
    },
    password: {
      type: String,
      required: true,
      select: false,
      minlength: 6,
    },
    role: {
      type: String,
      enum: ['admin', 'user'],
      default: 'user',
    },
  },
  { timestamps: true },
);

userSchema.pre('save', async function () {
  if (!this.isModified('password')) {
    return;
  }
  this.password = await bcrypt.hash(this.password, 10);
});

userSchema.methods.comparePassword = async function (
  plainPassword: string,
): Promise<boolean> {
  return bcrypt.compare(plainPassword, this.password);
};

export const UserModel = model<IUser>('User', userSchema);
`
	case DbPrisma:
		modelContent = `// Add this to prisma/schema.prisma:
//
// model User {
//   id        String   @id @default(cuid())
//   email     String   @unique
//   password  String
//   role      Role     @default(user)
//   createdAt DateTime @default(now())
//   updatedAt DateTime @updatedAt
// }
//
// enum Role {
//   admin
//   user
// }
`
	case DbDrizzle:
		modelContent = `// Add the users table to src/app/db/schema.ts:
//
// import { pgTable, text, varchar, timestamp } from 'drizzle-orm/pg-core';
//
// export const users = pgTable('users', {
//   id:        text('id').primaryKey().$defaultFn(() => crypto.randomUUID()),
//   email:     varchar('email', { length: 255 }).notNull().unique(),
//   password:  varchar('password', { length: 255 }).notNull(),
//   role:      varchar('role', { length: 50 }).notNull().default('user'),
//   createdAt: timestamp('created_at').defaultNow(),
//   updatedAt: timestamp('updated_at').defaultNow(),
// });
`
	}
	if err = writeFile(filepath.Join(authDir, "auth.model.ts"), modelContent); err != nil {
		return err
	}

	// 7. auth.service.ts
	dbImports := ""
	registerLogic := ""
	findUserLogic := ""
	switch opts.DB {
	case DbMongoose:
		dbImports = "import { UserModel } from './auth.model';"
		registerLogic = `const existingUser = await UserModel.findOne({ email: payload.email });
  if (existingUser) {
    throw new AppError(StatusCodes.CONFLICT, 'User already exists');
  }
  const created = await UserModel.create({
    email: payload.email,
    password: payload.password,
    role: 'user',
  });
  return {
    id: String(created._id),
    email: created.email,
    role: created.role,
  };`
		findUserLogic = `const user = await UserModel.findOne({ email: payload.email }).select(
    '+password',
  );
  if (!user) {
    throw new AppError(StatusCodes.NOT_FOUND, 'User not found');
  }`
	case DbPrisma:
		dbImports = "import { PrismaClient } from '@prisma/client';\n\nconst prisma = new PrismaClient();"
		registerLogic = `const hashedPassword = await bcrypt.hash(
    payload.password as string,
    config.bcrypt_salt_rounds,
  );
  const created = await prisma.user.create({
    data: {
      email: payload.email,
      password: hashedPassword,
      role: 'user',
    },
    select: { id: true, email: true, role: true },
  });
  return created;`
		findUserLogic = `const user = await prisma.user.findUnique({
    where: { email: payload.email },
  });
  if (!user) {
    throw new AppError(StatusCodes.NOT_FOUND, 'User not found');
  }`
	case DbDrizzle:
		dbImports = "import { db } from '../../db';\nimport { users } from '../../db/schema';\nimport { eq } from 'drizzle-orm';"
		registerLogic = `const hashedPassword = await bcrypt.hash(
    payload.password as string,
    config.bcrypt_salt_rounds,
  );
  const [created] = await db
    .insert(users)
    .values({
      email: payload.email,
      password: hashedPassword,
      role: 'user',
    })
    .returning({ id: users.id, email: users.email, role: users.role });
  return created;`
		findUserLogic = `const result = await db
    .select()
    .from(users)
    .where(eq(users.email, payload.email))
    .limit(1);
  const user = result[0];
  if (!user) {
    throw new AppError(StatusCodes.NOT_FOUND, 'User not found');
  }`
	}

	err = writeFile(filepath.Join(authDir, "auth.service.ts"), fmt.Sprintf(`import { StatusCodes } from 'http-status-codes';
import bcrypt from 'bcrypt';
import AppError from '../../errors/AppError';
import { TLoginUser, TUser } from './auth.interface';
import { createToken } from '../../utils/jwt.utils';
import config from '../../config';
%s

const registerUser = async (payload: TUser) => {
  %s
};

const loginUser = async (payload: TLoginUser) => {
  %s

  const isPasswordMatch = await bcrypt.compare(
    payload.password ?? '',
    user.password,
  );
  if (!isPasswordMatch) {
    throw new AppError(StatusCodes.UNAUTHORIZED, 'Invalid credentials');
  }

  const jwtPayload = {
    userId: String(
      (user as { id?: string; _id?: unknown }).id ||
        (user as { _id?: unknown })._id,
    ),
    role: user.role,
  };

  const accessToken = createToken(
    jwtPayload,
    config.jwt_access_secret as string,
    config.jwt_access_expires_in as string,
  );

  const refreshToken = createToken(
    jwtPayload,
    config.jwt_refresh_secret as string,
    config.jwt_refresh_expires_in as string,
  );

  return {
    accessToken,
    refreshToken,
  };
};

export const AuthService = {
  registerUser,
  loginUser,
};
`, dbImports, registerLogic, findUserLogic))
	if err != nil {
		return err
	}

	// 7. auth.validation.ts
	if opts.Validator == ValidatorZod {
		err = writeFile(filepath.Join(authDir, "auth.validation.ts"), `import { z } from 'zod';

const registerValidationSchema = z.object({
  body: z.object({
    email: z.string().email(),
    password: z.string().min(6),
  }),
});

const loginValidationSchema = z.object({
  body: z.object({
    email: z.string().email(),
    password: z.string().min(6),
  }),
});

export const AuthValidation = {
  registerValidationSchema,
  loginValidationSchema,
};
`)
	} else {
		err = writeFile(filepath.Join(authDir, "auth.validation.ts"), `import Joi from 'joi';

const registerValidationSchema = Joi.object({
  email: Joi.string().email().required(),
  password: Joi.string().min(6).required(),
});

const loginValidationSchema = Joi.object({
  email: Joi.string().email().required(),
  password: Joi.string().min(6).required(),
});

export const AuthValidation = {
  registerValidationSchema,
  loginValidationSchema,
};
`)
	}
	if err != nil {
		return err
	}

	// 8. auth.route.ts
	validationImport := "import { AuthValidation } from './auth.validation';"
	validationMiddlewareReg := "validateRequest(AuthValidation.registerValidationSchema),"
	validationMiddlewareLog := "validateRequest(AuthValidation.loginValidationSchema),"

	return writeFile(filepath.Join(authDir, "auth.route.ts"), fmt.Sprintf(`import { Router } from 'express';
import { AuthControllers } from './auth.controller';
import validateRequest from '../../utils/validateRequest';
import { authRateLimiter } from '../../middlewares/rateLimiter.middleware';
%s

const router = Router();

router.post(
  '/register',
  %s
  AuthControllers.registerUser,
);

router.post(
  '/login',
  authRateLimiter,
  %s
  AuthControllers.loginUser,
);

export const AuthRoutes = router;
`, validationImport, validationMiddlewareReg, validationMiddlewareLog))
}
