package cmd

import (
	"os"

	"github.com/Levi9111/create-express-modular-go/internal/modules"
	"github.com/Levi9111/create-express-modular-go/internal/ui"
	"github.com/spf13/cobra"
)

var removeCmd = &cobra.Command{
	Use:     "remove",
	Aliases: []string{"rm"},
	Short:   "Remove modules, middlewares, or environment variables",
}

var removeModuleCmd = &cobra.Command{
	Use:   "module [names...]",
	Short: "Delete module directory(ies) and unwire routes",
	Args:  cobra.MinimumNArgs(1),
	Run: func(cmd *cobra.Command, args []string) {
		cwd, err := os.Getwd()
		if err != nil {
			ui.Abort(err.Error())
		}
		items, err := modules.RemoveModule(cwd, args)
		if err != nil {
			ui.Abort(err.Error())
		}
		if ui.GetOutputMode() == ui.ModeJSON {
			ui.PrintJSON(map[string]any{"ok": true, "action": "remove", "type": "module", "items": items})
		}
	},
}

var removeMiddlewareCmd = &cobra.Command{
	Use:   "middleware [names...]",
	Short: "Delete custom middleware file(s)",
	Args:  cobra.MinimumNArgs(1),
	Run: func(cmd *cobra.Command, args []string) {
		cwd, err := os.Getwd()
		if err != nil {
			ui.Abort(err.Error())
		}
		items, err := modules.RemoveMiddleware(cwd, args)
		if err != nil {
			ui.Abort(err.Error())
		}
		if ui.GetOutputMode() == ui.ModeJSON {
			ui.PrintJSON(map[string]any{"ok": true, "action": "remove", "type": "middleware", "items": items})
		}
	},
}

var removeEnvCmd = &cobra.Command{
	Use:   "env [keys...]",
	Short: "Remove environment variable(s) from .env and .env.example",
	Args:  cobra.MinimumNArgs(1),
	Run: func(cmd *cobra.Command, args []string) {
		cwd, err := os.Getwd()
		if err != nil {
			ui.Abort(err.Error())
		}
		items, err := modules.RemoveEnv(cwd, args)
		if err != nil {
			ui.Abort(err.Error())
		}
		if ui.GetOutputMode() == ui.ModeJSON {
			ui.PrintJSON(map[string]any{"ok": true, "action": "remove", "type": "env", "items": items})
		}
	},
}

func init() {
	removeCmd.AddCommand(removeModuleCmd)
	removeCmd.AddCommand(removeMiddlewareCmd)
	removeCmd.AddCommand(removeEnvCmd)
	RootCmd.AddCommand(removeCmd)
}
