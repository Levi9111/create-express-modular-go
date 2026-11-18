package cmd

import (
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
	"strings"
	"time"

	"github.com/AlecAivazis/survey/v2"
	"github.com/Levi9111/create-express-modular-go/internal/generator"
	"github.com/Levi9111/create-express-modular-go/internal/modules"
	"github.com/Levi9111/create-express-modular-go/internal/pm"
	"github.com/Levi9111/create-express-modular-go/internal/telemetry"
	"github.com/Levi9111/create-express-modular-go/internal/ui"
	"github.com/Levi9111/create-express-modular-go/internal/updatenotifier"
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
	flagTesting     bool
	flagNoTesting   bool
	flagPm          string
	flagSpec        string
	flagJSON        bool
	flagQuiet       bool
)

type specModule struct {
	Name   string `json:"name"`
	Fields string `json:"fields,omitempty"`
}

func (s *specModule) UnmarshalJSON(data []byte) error {
	var str string
	if err := json.Unmarshal(data, &str); err == nil {
		s.Name = str
		return nil
	}
	type alias specModule
	var obj alias
	if err := json.Unmarshal(data, &obj); err != nil {
		return err
	}
	*s = specModule(obj)
	return nil
}

type projectSpec struct {
	Name          string       `json:"name"`
	DB            string       `json:"db"`
	Validator     string       `json:"validator"`
	Auth          *bool        `json:"auth"`
	TokenDelivery string       `json:"tokenDelivery"`
	Docker        *bool        `json:"docker"`
	Swagger       *bool        `json:"swagger"`
	Testing       *bool        `json:"testing"`
	Modules       []specModule `json:"modules"`
	Middlewares   []string     `json:"middlewares"`
	Env           []string     `json:"env"`
}

func loadProjectSpec(specArg string) (*projectSpec, error) {
	trimmed := strings.TrimSpace(specArg)
	var raw []byte
	if strings.HasPrefix(trimmed, "{") {
		raw = []byte(trimmed)
	} else {
		data, err := os.ReadFile(trimmed)
		if err != nil {
			return nil, err
		}
		raw = data
	}
	var spec projectSpec
	if err := json.Unmarshal(raw, &spec); err != nil {
		return nil, err
	}
	return &spec, nil
}

var RootCmd = &cobra.Command{
	Use:   "cem [project-name]",
	Short: "Scaffold a modular, domain-driven Express + TypeScript backend",
	Long:  `CEM CLI (Go Engine) — Fast, domain-driven Express + TypeScript project scaffolder.`,
	Args:  cobra.ArbitraryArgs,
	PersistentPreRun: func(cmd *cobra.Command, args []string) {
		if flagJSON {
			ui.SetOutputMode(ui.ModeJSON)
		} else if flagQuiet {
			ui.SetOutputMode(ui.ModeQuiet)
		} else {
			ui.SetOutputMode(ui.ModeNormal)
		}
	},
	Run: func(cmd *cobra.Command, args []string) {
		if flagVersion {
			fmt.Println(version.Short())
			return
		}

		ui.PrintBanner(version.Version)

		updateCh := make(chan string, 1)
		go func() {
			updateCh <- updatenotifier.CheckForUpdates()
		}()

		projectName := ""
		if len(args) > 0 {
			projectName = args[0]
		}

		cwd, err := os.Getwd()
		if err != nil {
			ui.Abort("Failed to get current directory: " + err.Error())
		}

		detectedPM := pm.DetectPM(cwd)
		if flagPm != "" {
			switch flagPm {
			case "npm", "yarn", "pnpm", "bun":
				detectedPM = pm.PackageManager(flagPm)
			default:
				ui.Abort(fmt.Sprintf("Invalid --pm choice %q. Allowed: npm, yarn, pnpm, bun", flagPm))
			}
		}

		// Defaults for -y / --yes
		dbChoice := generator.DbMongoose
		validatorChoice := generator.ValidatorZod
		useAuth := true
		tokenDelivery := generator.TokenCookie
		useDocker := true
		useSwagger := true
		useTesting := false
		skipInstall := flagNoInstall || flagSkipInstall

		var spec *projectSpec
		if flagSpec != "" {
			var err error
			spec, err = loadProjectSpec(flagSpec)
			if err != nil {
				ui.Abort("Invalid --spec input: " + err.Error())
			}
			flagYes = true
			if spec.Name != "" {
				projectName = spec.Name
			}
			if spec.DB != "" {
				flagDb = spec.DB
			}
			if spec.Validator != "" {
				flagValidator = spec.Validator
			}
			if spec.Auth != nil {
				useAuth = *spec.Auth
			}
			if spec.TokenDelivery != "" {
				tokenDelivery = generator.TokenDelivery(spec.TokenDelivery)
			}
			if spec.Docker != nil {
				useDocker = *spec.Docker
			}
			if spec.Swagger != nil {
				useSwagger = *spec.Swagger
			}
			if spec.Testing != nil {
				useTesting = *spec.Testing
			}
		}

		if flagYes {
			if projectName == "" {
				projectName = "my-api"
			}
			// Apply explicit flag overrides
			if flagDb != "" {
				switch flagDb {
				case "mongoose", "prisma", "drizzle":
					dbChoice = generator.DbChoice(flagDb)
				default:
					ui.Abort(fmt.Sprintf("Invalid --db choice %q. Allowed: mongoose, prisma, drizzle", flagDb))
				}
			}
			if flagValidator != "" {
				switch flagValidator {
				case "zod", "joi":
					validatorChoice = generator.ValidatorChoice(flagValidator)
				default:
					ui.Abort(fmt.Sprintf("Invalid --validator choice %q. Allowed: zod, joi", flagValidator))
				}
			}
			if spec != nil && spec.TokenDelivery != "" {
				switch spec.TokenDelivery {
				case "cookie", "header":
					tokenDelivery = generator.TokenDelivery(spec.TokenDelivery)
				default:
					ui.Abort(fmt.Sprintf("Invalid tokenDelivery %q. Allowed: cookie, header", spec.TokenDelivery))
				}
			}
			if flagNoAuth {
				useAuth = false
			} else if flagAuth {
				useAuth = true
			}
			if flagHeader {
				tokenDelivery = generator.TokenHeader
			} else if flagCookie {
				tokenDelivery = generator.TokenCookie
			}
			if flagNoDocker {
				useDocker = false
			} else if flagDocker {
				useDocker = true
			}
			if flagNoSwagger {
				useSwagger = false
			} else if flagSwagger {
				useSwagger = true
			}
			if flagNoTesting {
				useTesting = false
			} else if flagTesting {
				useTesting = true
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

		if err := modules.ValidateProjectName(projectName); err != nil {
			ui.Abort(err.Error())
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
			UseTesting:    useTesting,
			PM:            detectedPM,
			SkipInstall:   skipInstall,
		}

		// Run generator
		if err := generator.GenerateProject(opts); err != nil {
			ui.Abort(err.Error())
		}

		// One-shot --spec extras (modules, middlewares, env)
		if spec != nil {
			for _, mod := range spec.Modules {
				if mod.Name == "" {
					continue
				}
				if _, err := modules.AddModuleWithFields(projectPath, []string{mod.Name}, mod.Fields); err != nil {
					ui.Abort(err.Error())
				}
			}
			if len(spec.Middlewares) > 0 {
				if _, err := modules.AddMiddleware(projectPath, spec.Middlewares); err != nil {
					ui.Abort(err.Error())
				}
			}
			if len(spec.Env) > 0 {
				if _, err := modules.AddEnv(projectPath, spec.Env); err != nil {
					ui.Abort(err.Error())
				}
			}
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
				ui.Err("Dependency installation failed: " + err.Error())
			} else {
				ui.Success(fmt.Sprintf("Dependencies installed via %s", detectedPM))
			}
		} else {
			ui.Warn("Skipping dependencies installation (--no-install).")
		}

		// Fire-and-forget anonymous telemetry ping
		go telemetry.ReportInstall("create", version.Version, string(detectedPM))

		if ui.GetOutputMode() == ui.ModeJSON {
			specMods := []string{}
			specMws := []string{}
			specEnv := []string{}
			if spec != nil {
				if spec.Modules != nil {
					for _, m := range spec.Modules {
						if m.Name != "" {
							specMods = append(specMods, m.Name)
						}
					}
				}
				if spec.Middlewares != nil {
					specMws = spec.Middlewares
				}
				if spec.Env != nil {
					specEnv = spec.Env
				}
			}
			ui.PrintJSON(map[string]any{
				"ok":          true,
				"engine":      "go",
				"project":     projectName,
				"path":        projectPath,
				"db":          dbChoice,
				"validator":   validatorChoice,
				"auth":        useAuth,
				"docker":      useDocker,
				"swagger":     useSwagger,
				"testing":     useTesting,
				"pm":          detectedPM,
				"modules":     specMods,
				"middlewares": specMws,
				"env":         specEnv,
			})
			return
		}

		// Summary and completion
		ui.PrintSummary(projectName, string(dbChoice), string(validatorChoice), useAuth, useDocker, useSwagger)
		ui.PrintNextSteps(projectName)

		var latestVersion string
		select {
		case latestVersion = <-updateCh:
		case <-time.After(300 * time.Millisecond):
		}
		ui.PrintGoNotice(version.Version, latestVersion, updatenotifier.IsUpdateAvailable(version.Version, latestVersion))
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
	RootCmd.Flags().BoolVar(&flagTesting, "testing", false, "Enable Vitest + Supertest testing setup")
	RootCmd.Flags().BoolVar(&flagNoTesting, "no-testing", false, "Disable testing setup")
	RootCmd.Flags().StringVar(&flagPm, "pm", "", "Package manager override (npm | yarn | pnpm | bun)")
	RootCmd.Flags().StringVar(&flagSpec, "spec", "", "One-shot JSON spec file path or inline JSON")
	RootCmd.PersistentFlags().BoolVar(&flagJSON, "json", false, "Output compact JSON for AI agents and scripts")
	RootCmd.PersistentFlags().BoolVarP(&flagQuiet, "quiet", "q", false, "Suppress banners, spinners, and decorative output")
	RootCmd.Flags().BoolVar(&flagNoInstall, "no-install", false, "Skip dependency installation")
	RootCmd.Flags().BoolVar(&flagSkipInstall, "skip-install", false, "Alias for --no-install")
	RootCmd.Flags().BoolVarP(&flagVersion, "version", "v", false, "Print CLI version")
}
