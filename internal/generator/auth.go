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
    message: 'Too many requests from this IP, please try again after 15 minutes',
  },
  standardHeaders: true,
  legacyHeaders: false,
});

export const authRateLimiter = rateLimit({
  windowMs: 15 * 60 * 1000,
  max: 5,
  message: {
    success: false,
    message: 'Too many login attempts from this IP, please try again after 15 minutes',
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

  res.cookie('refreshToken', refreshToken, {
    secure: config.NODE_ENV === 'production',
    httpOnly: true,
  });

  res.cookie('accessToken', accessToken, {
    secure: config.NODE_ENV === 'production',
    httpOnly: true,
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

	// 6. auth.service.ts
	err = writeFile(filepath.Join(authDir, "auth.service.ts"), `import { StatusCodes } from 'http-status-codes';
import bcrypt from 'bcrypt';
import AppError from '../../errors/AppError';
import { TLoginUser, TUser } from './auth.interface';
import { createToken } from '../../utils/jwt.utils';
import config from '../../config';

const registerUser = async (payload: TUser) => {
  const hashedPassword = await bcrypt.hash(
    payload.password as string,
    config.bcrypt_salt_rounds,
  );

  return {
    id: 'user-id-placeholder',
    email: payload.email,
    role: payload.role || 'user',
    password: hashedPassword,
  };
};

const loginUser = async (payload: TLoginUser) => {
  // In production, fetch user from DB and compare passwords using bcrypt.compare
  const jwtPayload = {
    userId: 'user-id-placeholder',
    role: 'user',
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
`)
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
    role: z.enum(['user', 'admin']).optional(),
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
  body: Joi.object({
    email: Joi.string().email().required(),
    password: Joi.string().min(6).required(),
    role: Joi.string().valid('user', 'admin').optional(),
  }),
});

const loginValidationSchema = Joi.object({
  body: Joi.object({
    email: Joi.string().email().required(),
    password: Joi.string().min(6).required(),
  }),
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
