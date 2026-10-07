package generator

import "github.com/Levi9111/create-express-modular-go/internal/pm"

type DbChoice string
type ValidatorChoice string
type TokenDelivery string

const (
	DbMongoose DbChoice = "mongoose"
	DbPrisma   DbChoice = "prisma"
	DbDrizzle  DbChoice = "drizzle"

	ValidatorZod ValidatorChoice = "zod"
	ValidatorJoi ValidatorChoice = "joi"

	TokenCookie TokenDelivery = "cookie"
	TokenHeader TokenDelivery = "header"
)

type ProjectOptions struct {
	ProjectName   string
	ProjectPath   string
	DB            DbChoice
	Validator     ValidatorChoice
	UseAuth       bool
	TokenDelivery TokenDelivery
	UseDocker     bool
	UseSwagger    bool
	PM            pm.PackageManager
	SkipInstall   bool
}
