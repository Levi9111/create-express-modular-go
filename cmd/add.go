package cmd

import (
	"os"

	"github.com/Levi9111/create-express-modular-go/internal/modules"
	"github.com/Levi9111/create-express-modular-go/internal/ui"
	"github.com/spf13/cobra"
)

var addCmd = &cobra.Command{
	Use:   "add",
	Short: "Add modules, middlewares, or environment variables",
}

var addModuleCmd = &cobra.Command{
	Use:   "module [names...]",
	Short: "Scaffold one or more feature modules",
	Args:  cobra.MinimumNArgs(1),
	Run: func(cmd *cobra.Command, args []string) {
		cwd, err := os.Getwd()
		if err != nil {
			ui.Abort(err.Error())
		}
		if err := modules.AddModule(cwd, args); err != nil {
			ui.Abort(err.Error())
		}
	},
}

var addMiddlewareCmd = &cobra.Command{
	Use:   "middleware [names...]",
	Short: "Create one or more custom middleware files",
	Args:  cobra.MinimumNArgs(1),
	Run: func(cmd *cobra.Command, args []string) {
		cwd, err := os.Getwd()
		if err != nil {
			ui.Abort(err.Error())
		}
		if err := modules.AddMiddleware(cwd, args); err != nil {
			ui.Abort(err.Error())
		}
	},
}

var addEnvCmd = &cobra.Command{
	Use:   "env [keys...]",
	Short: "Add one or more environment variables",
	Args:  cobra.MinimumNArgs(1),
	Run: func(cmd *cobra.Command, args []string) {
		cwd, err := os.Getwd()
		if err != nil {
			ui.Abort(err.Error())
		}
		if err := modules.AddEnv(cwd, args); err != nil {
			ui.Abort(err.Error())
		}
	},
}

func init() {
	addCmd.AddCommand(addModuleCmd)
	addCmd.AddCommand(addMiddlewareCmd)
	addCmd.AddCommand(addEnvCmd)
	RootCmd.AddCommand(addCmd)
}
