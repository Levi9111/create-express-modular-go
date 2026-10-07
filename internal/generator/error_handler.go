package generator

import (
	"fmt"
	"path/filepath"
)

func scaffoldGlobalErrorHandler(root string, dbBlock, valBlock ErrorBlock) error {
	content := fmt.Sprintf(`/* eslint-disable @typescript-eslint/no-explicit-any */

import { Request, Response, NextFunction } from 'express';
import { TErrorSources } from '../interfaces/error';
import AppError from '../errors/AppError';
import config from '../config';
%s
%s

const globalErrorHandler = (
  err: any,
  _req: Request,
  res: Response,
  next: NextFunction,
): void => {
  if (res.headersSent) {
    return next(err);
  }

  let statusCode = 500;
  let message = 'Something went wrong!';
  let errorSources: TErrorSources = [
    { path: '', message: 'Something went wrong' },
  ];

  %s%sif (err instanceof AppError) {
    statusCode = err.statusCode;
    message = err.message;
    errorSources = [{ path: '', message: err.message }];
  } else if (err instanceof Error) {
    message = err.message;
    errorSources = [{ path: '', message: err.message }];
  }

  res.status(statusCode).json({
    success: false,
    message,
    errorSources,
    stack: config.NODE_ENV === 'development' ? err?.stack : undefined,
  });
};

export default globalErrorHandler;
`, valBlock.Imports, dbBlock.Imports, valBlock.Handler, dbBlock.Handler)

	return writeFile(filepath.Join(root, "src/app/middlewares/globalErrorHandler.middleware.ts"), content)
}
