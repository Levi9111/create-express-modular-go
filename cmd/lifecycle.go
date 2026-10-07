package cmd

import (
	"os"

	"github.com/Levi9111/create-express-modular-go/internal/modules"
	"github.com/Levi9111/create-express-modular-go/internal/ui"
	"github.com/spf13/cobra"
)

var listCmd = &cobra.Command{
	Use:     "list",
	Aliases: []string{"ls"},
	Short:   "List modules, middlewares, and environment variables",
	Run: func(cmd *cobra.Command, args []string) {
		cwd, err := os.Getwd()
		if err != nil {
			ui.Abort(err.Error())
		}
		if err := modules.ListFeatures(cwd); err != nil {
			ui.Abort(err.Error())
		}
	},
}

var devCmd = &cobra.Command{
	Use:   "dev",
	Short: "Start development server with live reload",
	Run: func(cmd *cobra.Command, args []string) {
		cwd, err := os.Getwd()
		if err != nil {
			ui.Abort(err.Error())
		}
		if err := modules.RunDev(cwd); err != nil {
			ui.Abort(err.Error())
		}
	},
}

var buildCmd = &cobra.Command{
	Use:   "build",
	Short: "Compile TypeScript to dist/",
	Run: func(cmd *cobra.Command, args []string) {
		cwd, err := os.Getwd()
		if err != nil {
			ui.Abort(err.Error())
		}
		if err := modules.RunBuild(cwd); err != nil {
			ui.Abort(err.Error())
		}
	},
}

var startCmd = &cobra.Command{
	Use:   "start",
	Short: "Start production server",
	Run: func(cmd *cobra.Command, args []string) {
		cwd, err := os.Getwd()
		if err != nil {
			ui.Abort(err.Error())
		}
		if err := modules.RunStart(cwd); err != nil {
			ui.Abort(err.Error())
		}
	},
}

var flagCheckFix bool

var checkCmd = &cobra.Command{
	Use:   "check",
	Short: "Run type-check, lint, and format checks",
	Run: func(cmd *cobra.Command, args []string) {
		cwd, err := os.Getwd()
		if err != nil {
			ui.Abort(err.Error())
		}
		if err := modules.RunCheck(cwd, flagCheckFix); err != nil {
			ui.Abort(err.Error())
		}
	},
}

var fixCmd = &cobra.Command{
	Use:   "fix",
	Short: "Auto-fix lint and formatting issues",
	Run: func(cmd *cobra.Command, args []string) {
		cwd, err := os.Getwd()
		if err != nil {
			ui.Abort(err.Error())
		}
		if err := modules.RunFix(cwd); err != nil {
			ui.Abort(err.Error())
		}
	},
}

var ejectCmd = &cobra.Command{
	Use:   "eject",
	Short: "Eject project from CEM CLI dependency",
	Run: func(cmd *cobra.Command, args []string) {
		cwd, err := os.Getwd()
		if err != nil {
			ui.Abort(err.Error())
		}
		if err := modules.RunEject(cwd); err != nil {
			ui.Abort(err.Error())
		}
	},
}

func init() {
	checkCmd.Flags().BoolVar(&flagCheckFix, "fix", false, "Auto-fix lint and formatting issues")
	RootCmd.AddCommand(listCmd)
	RootCmd.AddCommand(devCmd)
	RootCmd.AddCommand(buildCmd)
	RootCmd.AddCommand(startCmd)
	RootCmd.AddCommand(checkCmd)
	RootCmd.AddCommand(fixCmd)
	RootCmd.AddCommand(ejectCmd)
}
