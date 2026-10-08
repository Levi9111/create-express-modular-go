package cmd

import (
	"fmt"
	"os"
	"path/filepath"

	"github.com/AlecAivazis/survey/v2"
	"github.com/Levi9111/create-express-modular-go/internal/generator"
	"github.com/Levi9111/create-express-modular-go/internal/pm"
	"github.com/Levi9111/create-express-modular-go/internal/telemetry"
	"github.com/Levi9111/create-express-modular-go/internal/ui"
	"github.com/Levi9111/create-express-modular-go/internal/version"
	"github.com/spf13/cobra"
)

var (
	flagYes         bool
	flagDb          string
	flagValidator   string
	flagAuth        bool
	flagNoAuth      bool
	flagCookie      bool
	flagHeader      bool
	flagDocker      bool
	flagNoDocker    bool
	flagSwagger     bool
	flagNoSwagger   bool
	flagNoInstall   bool
	flagSkipInstall bool
	flagVersion     bool
)

var RootCmd = &cobra.Command{
	Use:   "cem [project-name]",
	Short: "Scaffold a modular, domain-driven Express + TypeScript backend",
	Long:  `CEM CLI (Go Engine) — Fast, domain-driven Express + TypeScript project scaffolder.`,
	Args:  cobra.ArbitraryArgs,
	Run: func(cmd *cobra.Command, args []string) {
		if flagVersion {
			fmt.Println(version.Short())
			return
		}

		ui.PrintBanner(version.Version)

		projectName := ""
		if len(args) > 0 {
			projectName = args[0]
		}

		cwd, err := os.Getwd()
		if err != nil {
			ui.Abort("Failed to get current directory: " + err.Error())
		}

		detectedPM := pm.DetectPM(cwd)

		// Defaults for -y / --yes
		dbChoice := generator.DbMongoose
		validatorChoice := generator.ValidatorZod
		useAuth := true
		tokenDelivery := generator.TokenCookie
		useDocker := true
		useSwagger := true
		skipInstall := flagNoInstall || flagSkipInstall

		if flagYes {
			if projectName == "" {
				projectName = "my-api"
			}
			// Apply explicit flag overrides
			if flagDb != "" {
				dbChoice = generator.DbChoice(flagDb)
			}
			if flagValidator != "" {
				validatorChoice = generator.ValidatorChoice(flagValidator)
			}
			if flagNoAuth {
				useAuth = false
			}
			if flagHeader {
				tokenDelivery = generator.TokenHeader
			}
			if flagNoDocker {
				useDocker = false
			}
			if flagNoSwagger {
				useSwagger = false
			}
		} else {
			// Interactive prompts
			if projectName == "" {
				prompt := &survey.Input{
					Message: "Project name:",
					Default: "my-api",
				}
				_ = survey.AskOne(prompt, &projectName)
			}

			// Check directory collision
			targetPath := filepath.Join(cwd, projectName)
			if _, err := os.Stat(targetPath); err == nil {
				ui.Abort(fmt.Sprintf("Directory '%s' already exists.", projectName))
			}

			// Database selection
			dbPrompt := &survey.Select{
				Message: "Select database / ORM:",
				Options: []string{"Mongoose (MongoDB)", "Prisma (PostgreSQL)", "Drizzle (PostgreSQL)"},
				Default: "Mongoose (MongoDB)",
			}
			var dbAns string
			_ = survey.AskOne(dbPrompt, &dbAns)
			switch dbAns {
			case "Prisma (PostgreSQL)":
				dbChoice = generator.DbPrisma
			case "Drizzle (PostgreSQL)":
				dbChoice = generator.DbDrizzle
			default:
				dbChoice = generator.DbMongoose
			}

			// Validator selection
			valPrompt := &survey.Select{
				Message: "Select schema validator:",
				Options: []string{"Zod", "Joi"},
				Default: "Zod",
			}
			var valAns string
			_ = survey.AskOne(valPrompt, &valAns)
			if valAns == "Joi" {
				validatorChoice = generator.ValidatorJoi
			} else {
				validatorChoice = generator.ValidatorZod
			}

			// Auth prompt
			authPrompt := &survey.Confirm{
				Message: "Include JWT Authentication setup (bcrypt, rate limiting)?",
				Default: true,
			}
			_ = survey.AskOne(authPrompt, &useAuth)

			if useAuth {
				deliveryPrompt := &survey.Select{
					Message: "Token delivery strategy:",
					Options: []string{"HTTP-only cookies (recommended)", "Authorization header"},
					Default: "HTTP-only cookies (recommended)",
				}
				var delAns string
				_ = survey.AskOne(deliveryPrompt, &delAns)
				if delAns == "Authorization header" {
					tokenDelivery = generator.TokenHeader
				} else {
					tokenDelivery = generator.TokenCookie
				}
			}

			// Docker prompt
			dockerPrompt := &survey.Confirm{
				Message: "Generate Docker configuration (Dockerfile + compose)?",
				Default: true,
			}
			_ = survey.AskOne(dockerPrompt, &useDocker)

			// Swagger prompt
			swaggerPrompt := &survey.Confirm{
				Message: "Include Swagger (OpenAPI 3.0) API documentation at /docs?",
				Default: true,
			}
			_ = survey.AskOne(swaggerPrompt, &useSwagger)
		}

		projectPath := filepath.Join(cwd, projectName)
		opts := generator.ProjectOptions{
			ProjectName:   projectName,
			ProjectPath:   projectPath,
			DB:            dbChoice,
			Validator:     validatorChoice,
			UseAuth:       useAuth,
			TokenDelivery: tokenDelivery,
			UseDocker:     useDocker,
			UseSwagger:    useSwagger,
			PM:            detectedPM,
			SkipInstall:   skipInstall,
		}

		// Run generator
		if err := generator.GenerateProject(opts); err != nil {
			ui.Abort(err.Error())
		}

		// Dependency installation
		ui.SectionHeader("Installing dependencies")
		if !skipInstall {
			installSpin := ui.NewSpinner(fmt.Sprintf("Installing dependencies via %s...", detectedPM))
			installSpin.Start()
			cmdName, cmdArgs := pm.InitialInstallCmd(detectedPM)
			_, err := pm.RunSilentCommand(projectPath, cmdName, cmdArgs...)
			installSpin.Stop()
			if err != nil {
				ui.Err("Dependency installation failed:")
				fmt.Println(err.Error())
			} else {
				ui.Success(fmt.Sprintf("Dependencies installed via %s", detectedPM))
			}
		} else {
			ui.Warn("Skipping dependencies installation (--no-install).")
		}

		// Fire-and-forget anonymous telemetry ping
		go telemetry.ReportInstall("create", version.Version, string(detectedPM))

		// Summary and completion
		ui.PrintSummary(projectName, string(dbChoice), string(validatorChoice), useAuth, useDocker, useSwagger)
		ui.PrintNextSteps(projectName)
	},
}

func init() {
	RootCmd.Flags().BoolVarP(&flagYes, "yes", "y", false, "Scaffold using recommended defaults non-interactively")
	RootCmd.Flags().StringVar(&flagDb, "db", "", "Database choice (mongoose | prisma | drizzle)")
	RootCmd.Flags().StringVar(&flagValidator, "validator", "", "Validator choice (zod | joi)")
	RootCmd.Flags().BoolVar(&flagAuth, "auth", false, "Enable JWT auth")
	RootCmd.Flags().BoolVar(&flagNoAuth, "no-auth", false, "Disable JWT auth")
	RootCmd.Flags().BoolVar(&flagCookie, "cookie", false, "Use HTTP-only cookies for tokens")
	RootCmd.Flags().BoolVar(&flagHeader, "header", false, "Use Authorization header for tokens")
	RootCmd.Flags().BoolVar(&flagDocker, "docker", false, "Enable Docker files")
	RootCmd.Flags().BoolVar(&flagNoDocker, "no-docker", false, "Disable Docker files")
	RootCmd.Flags().BoolVar(&flagSwagger, "swagger", false, "Enable Swagger API docs")
	RootCmd.Flags().BoolVar(&flagNoSwagger, "no-swagger", false, "Disable Swagger API docs")
	RootCmd.Flags().BoolVar(&flagNoInstall, "no-install", false, "Skip dependency installation")
	RootCmd.Flags().BoolVar(&flagSkipInstall, "skip-install", false, "Alias for --no-install")
	RootCmd.Flags().BoolVarP(&flagVersion, "version", "v", false, "Print CLI version")
}
