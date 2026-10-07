package generator

import (
	"path/filepath"
)

type ErrorBlock struct {
	Imports string
	Handler string
}

func scaffoldValidator(opts ProjectOptions) (ErrorBlock, error) {
	root := opts.ProjectPath

	if opts.Validator == ValidatorZod {
		err := writeFile(filepath.Join(root, "src/app/utils/validateRequest.ts"), `import { z } from 'zod';
import { NextFunction, Request, Response } from 'express';
import { catchAsync } from './catchAsync';

const validateRequest = (schema: z.ZodType) => {
  return catchAsync(
    async (req: Request, _res: Response, next: NextFunction) => {
      await schema.parseAsync({
        body: req.body,
        query: req.query,
        params: req.params,
        cookies: req.cookies,
      });
      return next();
    },
  );
};

export default validateRequest;
`)
		if err != nil {
			return ErrorBlock{}, err
		}

		err = writeFile(filepath.Join(root, "src/app/errors/handleZodError.ts"), `import { ZodError } from 'zod';
import { TErrorSources, TGenericErrorResponse } from '../interfaces/error';

const handleZodError = (err: ZodError): TGenericErrorResponse => {
  const errorSources: TErrorSources = err.issues.map((issue) => {
    return {
      path: issue.path[issue.path.length - 1] ?? '',
      message: issue.message,
    };
  });

  return {
    statusCode: 400,
    message: 'Validation Error',
    errorSources,
  };
};

export default handleZodError;
`)
		if err != nil {
			return ErrorBlock{}, err
		}

		return ErrorBlock{
			Imports: `import { ZodError } from 'zod';
import handleZodError from '../errors/handleZodError';`,
			Handler: `if (err instanceof ZodError) {
    const simplified = handleZodError(err);
    statusCode = simplified.statusCode;
    message = simplified.message;
    errorSources = simplified.errorSources;
  } else `,
		}, nil
	}

	// Joi
	err := writeFile(filepath.Join(root, "src/app/utils/validateRequest.ts"), `import Joi from 'joi';
import { NextFunction, Request, Response } from 'express';
import { catchAsync } from './catchAsync';

const validateRequest = (schema: Joi.ObjectSchema) => {
  return catchAsync(async (req: Request, _res: Response, next: NextFunction) => {
    const { error } = schema.validate(req.body, { abortEarly: false });
    if (error) throw error;
    return next();
  });
};

export default validateRequest;
`)
	if err != nil {
		return ErrorBlock{}, err
	}

	return ErrorBlock{
		Imports: "import Joi from 'joi';",
		Handler: `if (err?.isJoi === true || err instanceof Joi.ValidationError) {
    statusCode = 400;
    message = 'Validation Error';
    errorSources = err.details.map((detail: Joi.ValidationErrorItem) => ({
      path: String(detail.path[detail.path.length - 1] ?? ''),
      message: detail.message.replace(/['"]/g, ''),
    }));
  } else `,
	}, nil
}
