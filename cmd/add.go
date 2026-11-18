package cmd

import (
	"os"

	"github.com/Levi9111/create-express-modular-go/internal/modules"
	"github.com/Levi9111/create-express-modular-go/internal/ui"
	"github.com/spf13/cobra"
)

var (
	flagFields string
	flagAddYes bool
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
		items, err := modules.AddModuleWithFields(cwd, args, flagFields)
		if err != nil {
			ui.Abort(err.Error())
		}
		if ui.GetOutputMode() == ui.ModeJSON {
			ui.PrintJSON(map[string]any{"ok": true, "action": "add", "type": "module", "items": items})
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
		items, err := modules.AddMiddleware(cwd, args)
		if err != nil {
			ui.Abort(err.Error())
		}
		if ui.GetOutputMode() == ui.ModeJSON {
			ui.PrintJSON(map[string]any{"ok": true, "action": "add", "type": "middleware", "items": items})
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
		items, err := modules.AddEnv(cwd, args)
		if err != nil {
			ui.Abort(err.Error())
		}
		if ui.GetOutputMode() == ui.ModeJSON {
			ui.PrintJSON(map[string]any{"ok": true, "action": "add", "type": "env", "items": items})
		}
	},
}

func init() {
	addModuleCmd.Flags().StringVar(&flagFields, "fields", "", "Field definitions (e.g. \"name:string,price:number,inStock:boolean?\")")
	addModuleCmd.Flags().BoolVarP(&flagAddYes, "yes", "y", false, "Skip interactive prompts and use defaults")
	addCmd.AddCommand(addModuleCmd)
	addCmd.AddCommand(addMiddlewareCmd)
	addCmd.AddCommand(addEnvCmd)
	RootCmd.AddCommand(addCmd)
}
