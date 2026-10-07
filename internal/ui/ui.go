package ui

import (
	"fmt"
	"os"
	"strings"
	"time"

	"github.com/briandowns/spinner"
	"github.com/fatih/color"
)

var (
	Cyan    = color.New(color.FgCyan, color.Bold).SprintFunc()
	Green   = color.New(color.FgGreen, color.Bold).SprintFunc()
	Yellow  = color.New(color.FgYellow, color.Bold).SprintFunc()
	Red     = color.New(color.FgRed, color.Bold).SprintFunc()
	Magenta = color.New(color.FgMagenta, color.Bold).SprintFunc()
	Gray    = color.New(color.FgHiBlack).SprintFunc()
	Bold    = color.New(color.Bold).SprintFunc()
	White   = color.New(color.FgWhite).SprintFunc()
)

const Divider = "────────────────────────────────────────────────────"
const SubDivider = "┄┄┄┄┄┄┄┄┄┄┄┄┄┄┄┄┄┄┄┄┄┄┄┄┄┄┄┄┄┄┄┄┄┄┄┄┄┄┄┄┄┄┄┄┄┄"

func PrintBanner(version string) {
	fmt.Println()
	fmt.Println(Gray(Divider))
	fmt.Println()
	fmt.Printf("   %s   %s  %s\n", Cyan("CEM"), White("create-express-modular (Go)"), Gray("v"+version))
	fmt.Printf("  %s\n", Gray("Modular Express + TypeScript — experimental native engine"))
	fmt.Println()
	fmt.Println(Gray(Divider))
	fmt.Println()
}

func SectionHeader(title string) {
	fmt.Println()
	fmt.Println(Gray("  " + SubDivider))
	fmt.Printf("  %s\n", Bold(title))
	fmt.Println(Gray("  " + SubDivider))
	fmt.Println()
}

func Bullet(key, val string) {
	fmt.Printf("  %s  %-10s %s\n", Cyan("◆"), key, White(val))
}

func Substep(text string) {
	fmt.Printf("     %s  %s\n", Gray("·"), Gray(text))
}

func Warn(msg string) {
	fmt.Printf("  %s  %s\n", Yellow("⚠"), msg)
}

func Err(msg string) {
	fmt.Printf("  %s  %s\n", Red("✖"), Red(msg))
}

func Success(msg string) {
	fmt.Printf("  %s  %s\n", Green("✔"), msg)
}

func Abort(msg string) {
	fmt.Println()
	Err(msg)
	fmt.Println()
	os.Exit(1)
}

func NewSpinner(msg string) *spinner.Spinner {
	s := spinner.New(spinner.CharSets[14], 80*time.Millisecond)
	s.Prefix = "  "
	s.Suffix = "  " + msg
	_ = s.Color("cyan")
	return s
}

func PrintSummary(name, db, validator string, auth, docker, swagger bool) {
	W := 48
	border := func(l, m, r string) string {
		return "  " + Cyan(l+strings.Repeat(m, W-2)+r)
	}
	blank := func() string {
		return "  " + Cyan("│") + strings.Repeat(" ", W-2) + Cyan("│")
	}
	row := func(label, val string) string {
		content := fmt.Sprintf("  %-12s %s", label, val)
		pad := W - 2 - len(label) - len(val) - 4
		if pad < 0 {
			pad = 0
		}
		return "  " + Cyan("│") + content + strings.Repeat(" ", pad) + Cyan("│")
	}

	authStr := Gray("no")
	if auth {
		authStr = Green("yes")
	}
	dockerStr := Gray("no")
	if docker {
		dockerStr = Green("yes")
	}
	swaggerStr := Gray("no")
	if swagger {
		swaggerStr = Green("yes")
	}

	fmt.Println()
	fmt.Println(border("╭", "─", "╮"))
	fmt.Println(blank())
	title := Bold(Cyan("Project ready!"))
	fmt.Printf("  %s  %-44s%s\n", Cyan("│"), title, Cyan("│"))
	fmt.Println(blank())
	fmt.Println(border("├", "─", "┤"))
	fmt.Println(blank())
	fmt.Println(row("Name", name))
	fmt.Println(row("Database", db))
	fmt.Println(row("Validator", validator))
	fmt.Println(row("Auth", authStr))
	fmt.Println(row("Docker", dockerStr))
	fmt.Println(row("Swagger", swaggerStr))
	fmt.Println(blank())
	fmt.Println(border("╰", "─", "╯"))
	fmt.Println()
}

func PrintNextSteps(projectName string) {
	fmt.Println(Bold("  Next steps"))
	fmt.Println()
	fmt.Printf("  %s  %s %s\n", Gray("1."), Cyan("cd"), White(projectName))
	fmt.Printf("  %s  %s %s\n", Gray("2."), Cyan("cem"), White("dev"))
	fmt.Printf("  %s  %s %s\n", Gray("3."), Cyan("cem add module"), White("Product"))
	fmt.Printf("  %s  %s\n", Gray("4."), Cyan("cem build"))
	fmt.Printf("  %s  %s\n", Gray("5."), Cyan("cem start"))
	fmt.Println()
	fmt.Printf("  %s  %s\n", Gray("Docs →"), Gray("https://create-express-modular.lovable.app/docs"))
	fmt.Println()
	fmt.Printf("  %s\n", Gray("Anonymous telemetry helps improve CEM (no personal data). Disable: CEM_TELEMETRY=off"))
	fmt.Println()
}
