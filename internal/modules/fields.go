package modules

import (
	"fmt"
	"regexp"
	"strings"
)

var enumValueRegex = regexp.MustCompile(`^[A-Za-z0-9_-]+$`)

func camelCase(s string) string {
	if s == "" {
		return ""
	}
	return strings.ToLower(s[:1]) + s[1:]
}

// BaseFieldType represents a supported --fields type.
type BaseFieldType string

const (
	FieldString  BaseFieldType = "string"
	FieldNumber  BaseFieldType = "number"
	FieldBoolean BaseFieldType = "boolean"
	FieldDate    BaseFieldType = "date"
	FieldEnum    BaseFieldType = "enum"
	FieldRef     BaseFieldType = "ref"
)

// ParsedField represents one parsed field from --fields.
type ParsedField struct {
	Name       string
	BaseType   BaseFieldType
	IsArray    bool
	Optional   bool
	EnumValues []string
	RefModel   string
}

// ParseFields parses a comma-separated --fields DSL string into validated ParsedFields.
func ParseFields(raw string) ([]ParsedField, error) {
	trimmed := strings.TrimSpace(raw)
	if trimmed == "" {
		return nil, nil
	}

	tokens := strings.Split(trimmed, ",")
	var fields []ParsedField

	for _, tok := range tokens {
		tok = strings.TrimSpace(tok)
		if tok == "" {
			continue
		}

		colonIdx := strings.IndexByte(tok, ':')
		if colonIdx == -1 {
			return nil, fmt.Errorf("invalid field %q: expected format \"name:type\" (e.g. \"price:number\")", tok)
		}

		name := strings.TrimSpace(tok[:colonIdx])
		if err := ValidateIdentifier(name); err != nil {
			return nil, fmt.Errorf("invalid field name %q: %w", name, err)
		}

		typePart := strings.TrimSpace(tok[colonIdx+1:])
		if typePart == "" {
			return nil, fmt.Errorf("missing type for field %q", name)
		}

		optional := false
		if strings.HasSuffix(typePart, "?") {
			optional = true
			typePart = strings.TrimSpace(strings.TrimSuffix(typePart, "?"))
		}

		isArray := false
		if strings.HasSuffix(typePart, "[]") {
			isArray = true
			typePart = strings.TrimSpace(strings.TrimSuffix(typePart, "[]"))
		}

		lower := strings.ToLower(typePart)
		if strings.HasPrefix(lower, "enum(") && strings.HasSuffix(typePart, ")") {
			inner := typePart[5 : len(typePart)-1]
			var vals []string
			for _, v := range strings.Split(inner, "|") {
				v = strings.TrimSpace(v)
				if v != "" {
					vals = append(vals, v)
				}
			}
			if len(vals) == 0 {
				return nil, fmt.Errorf("enum field %q must specify at least one value, e.g. enum(a|b)", name)
			}
			for _, v := range vals {
				if !enumValueRegex.MatchString(v) {
					return nil, fmt.Errorf("invalid enum value %q on field %q: use only letters, digits, underscores, or hyphens", v, name)
				}
			}
			fields = append(fields, ParsedField{
				Name:       name,
				BaseType:   FieldEnum,
				IsArray:    isArray,
				Optional:   optional,
				EnumValues: vals,
			})
			continue
		}

		if strings.HasPrefix(lower, "ref(") && strings.HasSuffix(typePart, ")") {
			refModel := strings.TrimSpace(typePart[4 : len(typePart)-1])
			if err := ValidateIdentifier(refModel); err != nil {
				return nil, fmt.Errorf("invalid ref model %q on field %q: %w", refModel, name, err)
			}
			fields = append(fields, ParsedField{
				Name:     name,
				BaseType: FieldRef,
				IsArray:  isArray,
				Optional: optional,
				RefModel: refModel,
			})
			continue
		}

		switch lower {
		case "string":
			fields = append(fields, ParsedField{Name: name, BaseType: FieldString, IsArray: isArray, Optional: optional})
		case "number":
			fields = append(fields, ParsedField{Name: name, BaseType: FieldNumber, IsArray: isArray, Optional: optional})
		case "boolean", "bool":
			fields = append(fields, ParsedField{Name: name, BaseType: FieldBoolean, IsArray: isArray, Optional: optional})
		case "date":
			fields = append(fields, ParsedField{Name: name, BaseType: FieldDate, IsArray: isArray, Optional: optional})
		default:
			return nil, fmt.Errorf("unsupported type %q on field %q; supported: string, number, boolean, date, enum(a|b), ref(Model)", typePart, name)
		}
	}

	return fields, nil
}

func tsFieldType(f ParsedField, isMongoose bool) string {
	var base string
	switch f.BaseType {
	case FieldString:
		base = "string"
	case FieldNumber:
		base = "number"
	case FieldBoolean:
		base = "boolean"
	case FieldDate:
		base = "Date"
	case FieldEnum:
		quoted := make([]string, len(f.EnumValues))
		for i, v := range f.EnumValues {
			quoted[i] = fmt.Sprintf("'%s'", v)
		}
		base = strings.Join(quoted, " | ")
	case FieldRef:
		if isMongoose {
			base = "Types.ObjectId"
		} else {
			base = "string"
		}
	}
	if f.IsArray {
		return base + "[]"
	}
	return base
}

func mongooseFieldDef(f ParsedField) string {
	req := ""
	if !f.Optional {
		req = ", required: true"
	}
	var inner string
	switch f.BaseType {
	case FieldString:
		inner = fmt.Sprintf("{ type: String%s, trim: true }", req)
	case FieldNumber:
		inner = fmt.Sprintf("{ type: Number%s }", req)
	case FieldBoolean:
		inner = fmt.Sprintf("{ type: Boolean%s }", req)
	case FieldDate:
		inner = fmt.Sprintf("{ type: Date%s }", req)
	case FieldEnum:
		quoted := make([]string, len(f.EnumValues))
		for i, v := range f.EnumValues {
			quoted[i] = fmt.Sprintf("'%s'", v)
		}
		inner = fmt.Sprintf("{ type: String, enum: [%s]%s }", strings.Join(quoted, ", "), req)
	case FieldRef:
		inner = fmt.Sprintf("{ type: Schema.Types.ObjectId, ref: '%s'%s }", f.RefModel, req)
	}
	if f.IsArray {
		return "[" + inner + "]"
	}
	return inner
}

func zodFieldDef(f ParsedField, allOptional bool) string {
	var expr string
	switch f.BaseType {
	case FieldString, FieldRef:
		if f.Optional || allOptional {
			expr = "z.string()"
		} else {
			expr = fmt.Sprintf("z.string().min(1, '%s is required')", f.Name)
		}
	case FieldNumber:
		expr = "z.number()"
	case FieldBoolean:
		expr = "z.boolean()"
	case FieldDate:
		expr = "z.coerce.date()"
	case FieldEnum:
		quoted := make([]string, len(f.EnumValues))
		for i, v := range f.EnumValues {
			quoted[i] = fmt.Sprintf("'%s'", v)
		}
		expr = fmt.Sprintf("z.enum([%s])", strings.Join(quoted, ", "))
	}
	if f.IsArray {
		expr = fmt.Sprintf("z.array(%s)", expr)
	}
	if f.Optional || allOptional {
		expr += ".optional()"
	}
	return expr
}

func joiFieldDef(f ParsedField, allOptional bool) string {
	var expr string
	switch f.BaseType {
	case FieldString, FieldRef:
		expr = "Joi.string()"
	case FieldNumber:
		expr = "Joi.number()"
	case FieldBoolean:
		expr = "Joi.boolean()"
	case FieldDate:
		expr = "Joi.date()"
	case FieldEnum:
		quoted := make([]string, len(f.EnumValues))
		for i, v := range f.EnumValues {
			quoted[i] = fmt.Sprintf("'%s'", v)
		}
		expr = fmt.Sprintf("Joi.string().valid(%s)", strings.Join(quoted, ", "))
	}
	if f.IsArray {
		expr = fmt.Sprintf("Joi.array().items(%s)", expr)
	}
	if !f.Optional && !allOptional {
		expr += ".required()"
	} else {
		expr += ".optional()"
	}
	return expr
}

func searchableFieldNames(fields []ParsedField) []string {
	var out []string
	for _, f := range fields {
		if !f.IsArray && (f.BaseType == FieldString || f.BaseType == FieldEnum) {
			out = append(out, f.Name)
		}
	}
	return out
}

func prismaModelBlock(moduleName string, fields []ParsedField) string {
	lines := []string{
		fmt.Sprintf("model %s {", moduleName),
		"  id        String   @id @default(cuid())",
	}
	for _, f := range fields {
		var pType string
		switch f.BaseType {
		case FieldString, FieldEnum, FieldRef:
			pType = "String"
		case FieldNumber:
			pType = "Float"
		case FieldBoolean:
			pType = "Boolean"
		case FieldDate:
			pType = "DateTime"
		}
		if f.IsArray {
			pType += "[]"
		} else if f.Optional {
			pType += "?"
		}
		lines = append(lines, fmt.Sprintf("  %-10s %s", f.Name, pType))
	}
	lines = append(lines, "  createdAt DateTime @default(now())", "  updatedAt DateTime @updatedAt", "}")
	return strings.Join(lines, "\n")
}

func drizzleTableBlock(tableName string, fields []ParsedField) string {
	lines := []string{
		fmt.Sprintf("export const %s = pgTable('%s', {", tableName, tableName),
		"  id: uuid('id').defaultRandom().primaryKey(),",
	}
	for _, f := range fields {
		var col string
		switch f.BaseType {
		case FieldString, FieldEnum, FieldRef:
			col = fmt.Sprintf("text('%s')", f.Name)
		case FieldNumber:
			col = fmt.Sprintf("doublePrecision('%s')", f.Name)
		case FieldBoolean:
			col = fmt.Sprintf("boolean('%s')", f.Name)
		case FieldDate:
			col = fmt.Sprintf("timestamp('%s')", f.Name)
		}
		if f.IsArray {
			col += ".array()"
		}
		if !f.Optional {
			col += ".notNull()"
		}
		lines = append(lines, fmt.Sprintf("  %s: %s,", f.Name, col))
	}
	lines = append(lines,
		"  createdAt: timestamp('created_at').defaultNow().notNull(),",
		"  updatedAt: timestamp('updated_at').defaultNow().notNull(),",
		"});",
	)
	return strings.Join(lines, "\n")
}

func containsStr(slice []string, target string) bool {
	for _, s := range slice {
		if s == target {
			return true
		}
	}
	return false
}

func buildGoInterfaceContent(moduleName string, fields []ParsedField, isMongoose bool) string {
	if len(fields) == 0 {
		return fmt.Sprintf(`export type T%s = {
  id?: string;
  name?: string;
  createdAt?: Date;
  updatedAt?: Date;
};
`, moduleName)
	}

	needsTypes := false
	if isMongoose {
		for _, f := range fields {
			if f.BaseType == FieldRef {
				needsTypes = true
				break
			}
		}
	}

	importLine := ""
	if needsTypes {
		importLine = "import { Types } from 'mongoose';\n\n"
	}

	var lines []string
	for _, f := range fields {
		opt := ""
		if f.Optional {
			opt = "?"
		}
		lines = append(lines, fmt.Sprintf("  %s%s: %s;", f.Name, opt, tsFieldType(f, isMongoose)))
	}

	return fmt.Sprintf(`%sexport type T%s = {
  id?: string;
%s
  createdAt?: Date;
  updatedAt?: Date;
};
`, importLine, moduleName, strings.Join(lines, "\n"))
}

func buildGoMongooseModelContent(moduleName, fileName string, fields []ParsedField) string {
	schemaBody := "    // TODO: Define your fields here\n    // name: { type: String, required: true, trim: true },"
	if len(fields) > 0 {
		var lines []string
		for _, f := range fields {
			lines = append(lines, fmt.Sprintf("    %s: %s,", f.Name, mongooseFieldDef(f)))
		}
		schemaBody = strings.Join(lines, "\n")
	}

	return fmt.Sprintf(`import { Schema, model } from 'mongoose';
import { T%s } from './%s.interface';

const %sSchema = new Schema<T%s>(
  {
%s
  },
  { timestamps: true },
);

export const %sModel = model<T%s>('%s', %sSchema);
`, moduleName, fileName, fileName, moduleName, schemaBody, moduleName, moduleName, moduleName, fileName)
}

func buildGoServiceContent(moduleName, fileName, dbChoice string, hasFields bool) string {
	plural := pluralize(moduleName)
	switch dbChoice {
	case "mongoose":
		return fmt.Sprintf(`import { StatusCodes } from 'http-status-codes';
import AppError from '../../errors/AppError';
import QueryBuilder from '../../utils/QueryBuilder';
import { T%s } from './%s.interface';
import { %sModel } from './%s.model';
import { %sSearchableFields } from './%s.constant';

const create%s = async (payload: T%s) => {
  return %sModel.create(payload);
};

const getAll%ss = async (query: Record<string, unknown>) => {
  const qb = new QueryBuilder(%sModel.find(), query)
    .search(%sSearchableFields)
    .filter()
    .sort()
    .paginate()
    .fields();
  const [meta, data] = await Promise.all([qb.countTotal(), qb.modelQuery]);
  return { meta, data };
};

const getSingle%s = async (id: string) => {
  const result = await %sModel.findById(id);
  if (!result) {
    throw new AppError(StatusCodes.NOT_FOUND, '%s not found');
  }
  return result;
};

const update%s = async (id: string, payload: Partial<T%s>) => {
  const result = await %sModel.findByIdAndUpdate(id, payload, {
    new: true,
    runValidators: true,
  });
  if (!result) {
    throw new AppError(StatusCodes.NOT_FOUND, '%s not found');
  }
  return result;
};

const delete%s = async (id: string) => {
  const result = await %sModel.findByIdAndDelete(id);
  if (!result) {
    throw new AppError(StatusCodes.NOT_FOUND, '%s not found');
  }
  return result;
};

export const %sService = {
  create%s,
  getAll%ss,
  getSingle%s,
  update%s,
  delete%s,
};
`, moduleName, fileName, moduleName, fileName, moduleName, fileName,
			moduleName, moduleName, moduleName,
			moduleName, moduleName, moduleName,
			moduleName, moduleName, moduleName,
			moduleName, moduleName, moduleName, moduleName,
			moduleName, moduleName, moduleName,
			moduleName, moduleName, moduleName, moduleName, moduleName, moduleName)

	case "prisma":
		camelName := camelCase(moduleName)
		return fmt.Sprintf(`import { StatusCodes } from 'http-status-codes';
import { PrismaClient } from '@prisma/client';
import AppError from '../../errors/AppError';
import { T%s } from './%s.interface';

const prisma = new PrismaClient();
// eslint-disable-next-line @typescript-eslint/no-explicit-any
const modelDelegate = (prisma as unknown as Record<string, any>).%s;

const create%s = async (payload: T%s) => {
  return modelDelegate.create({ data: payload });
};

const getAll%ss = async (query: Record<string, unknown>) => {
  const page = Number(query.page) || 1;
  const limit = Number(query.limit) || 10;
  const skip = (page - 1) * limit;
  const [data, total] = await Promise.all([
    modelDelegate.findMany({ skip, take: limit }),
    modelDelegate.count(),
  ]);
  return {
    meta: { page, limit, total, totalPages: Math.ceil(total / limit) },
    data,
  };
};

const getSingle%s = async (id: string) => {
  const result = await modelDelegate.findUnique({ where: { id } });
  if (!result) {
    throw new AppError(StatusCodes.NOT_FOUND, '%s not found');
  }
  return result;
};

const update%s = async (id: string, payload: Partial<T%s>) => {
  await getSingle%s(id);
  return modelDelegate.update({ where: { id }, data: payload });
};

const delete%s = async (id: string) => {
  await getSingle%s(id);
  return modelDelegate.delete({ where: { id } });
};

export const %sService = {
  create%s,
  getAll%ss,
  getSingle%s,
  update%s,
  delete%s,
};
`, moduleName, fileName, camelName,
			moduleName, moduleName,
			moduleName,
			moduleName, moduleName,
			moduleName, moduleName, moduleName,
			moduleName, moduleName,
			moduleName, moduleName, moduleName, moduleName, moduleName, moduleName)

	case "drizzle":
		if hasFields {
			return fmt.Sprintf(`import { StatusCodes } from 'http-status-codes';
import { eq } from 'drizzle-orm';
import AppError from '../../errors/AppError';
import { db } from '../../db';
import { %s } from '../../db/schema';
import { T%s } from './%s.interface';

const create%s = async (payload: T%s) => {
  const [created] = await db.insert(%s).values(payload).returning();
  return created;
};

const getAll%ss = async (query: Record<string, unknown>) => {
  const page = Number(query.page) || 1;
  const limit = Number(query.limit) || 10;
  const offset = (page - 1) * limit;
  const data = await db.select().from(%s).limit(limit).offset(offset);
  return { meta: { page, limit }, data };
};

const getSingle%s = async (id: string) => {
  const [result] = await db
    .select()
    .from(%s)
    .where(eq(%s.id, id))
    .limit(1);
  if (!result) {
    throw new AppError(StatusCodes.NOT_FOUND, '%s not found');
  }
  return result;
};

const update%s = async (id: string, payload: Partial<T%s>) => {
  const [updated] = await db
    .update(%s)
    .set({ ...payload, updatedAt: new Date() })
    .where(eq(%s.id, id))
    .returning();
  if (!updated) {
    throw new AppError(StatusCodes.NOT_FOUND, '%s not found');
  }
  return updated;
};

const delete%s = async (id: string) => {
  const [deleted] = await db
    .delete(%s)
    .where(eq(%s.id, id))
    .returning();
  if (!deleted) {
    throw new AppError(StatusCodes.NOT_FOUND, '%s not found');
  }
  return deleted;
};

export const %sService = {
  create%s,
  getAll%ss,
  getSingle%s,
  update%s,
  delete%s,
};
`, plural, moduleName, fileName,
				moduleName, moduleName, plural,
				moduleName, plural,
				moduleName, plural, plural, moduleName,
				moduleName, moduleName, plural, plural, moduleName,
				moduleName, plural, plural, moduleName,
				moduleName, moduleName, moduleName, moduleName, moduleName, moduleName)
		}
	}

	return fmt.Sprintf(`import { T%s } from './%s.interface';

const create%s = async (payload: T%s) => {
  // TODO: Define %s table in src/app/db/schema.ts (or pass --fields when generating)
  return payload;
};

const getAll%ss = async (_query: Record<string, unknown>) => {
  return [];
};

const getSingle%s = async (id: string) => {
  return { id };
};

const update%s = async (id: string, payload: Partial<T%s>) => {
  return { id, ...payload };
};

const delete%s = async (id: string) => {
  return { id };
};

export const %sService = {
  create%s,
  getAll%ss,
  getSingle%s,
  update%s,
  delete%s,
};
`, moduleName, fileName, moduleName, moduleName, plural, moduleName, moduleName, moduleName, moduleName, moduleName, moduleName, moduleName, moduleName, moduleName, moduleName, moduleName)
}

func buildGoValidationContent(moduleName string, isZod bool, fields []ParsedField) string {
	if isZod {
		if len(fields) == 0 {
			return fmt.Sprintf(`import { z } from 'zod';

const create%sValidationSchema = z.object({
  body: z.object({
    // TODO: Define fields
    // name: z.string().min(1, 'Name is required'),
  }),
});

const update%sValidationSchema = z.object({
  body: z.object({
    // TODO: Define update fields (all optional)
  }),
});

export const %sValidation = {
  create%sValidationSchema,
  update%sValidationSchema,
};
`, moduleName, moduleName, moduleName, moduleName, moduleName)
		}

		var lines []string
		for _, f := range fields {
			lines = append(lines, fmt.Sprintf("    %s: %s,", f.Name, zodFieldDef(f, false)))
		}
		return fmt.Sprintf(`import { z } from 'zod';

const create%sValidationSchema = z.object({
  body: z.object({
%s
  }),
});

const update%sValidationSchema = z.object({
  body: create%sValidationSchema.shape.body.partial(),
});

export const %sValidation = {
  create%sValidationSchema,
  update%sValidationSchema,
};
`, moduleName, strings.Join(lines, "\n"), moduleName, moduleName, moduleName, moduleName, moduleName)
	}

	if len(fields) == 0 {
		return fmt.Sprintf(`import Joi from 'joi';

const create%sValidationSchema = Joi.object({
  // TODO: Define fields
  // name: Joi.string().required(),
});

const update%sValidationSchema = Joi.object({
  // TODO: Define update fields
});

export const %sValidation = {
  create%sValidationSchema,
  update%sValidationSchema,
};
`, moduleName, moduleName, moduleName, moduleName, moduleName)
	}

	var createLines, updateLines []string
	for _, f := range fields {
		createLines = append(createLines, fmt.Sprintf("  %s: %s,", f.Name, joiFieldDef(f, false)))
		updateLines = append(updateLines, fmt.Sprintf("  %s: %s,", f.Name, joiFieldDef(f, true)))
	}
	return fmt.Sprintf(`import Joi from 'joi';

const create%sValidationSchema = Joi.object({
%s
});

const update%sValidationSchema = Joi.object({
%s
});

export const %sValidation = {
  create%sValidationSchema,
  update%sValidationSchema,
};
`, moduleName, strings.Join(createLines, "\n"), moduleName, strings.Join(updateLines, "\n"), moduleName, moduleName, moduleName)
}
