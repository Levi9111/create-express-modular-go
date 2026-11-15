package ui

import (
	"encoding/json"
	"fmt"
	"io"
	"os"
	"strings"
	"time"

	"github.com/briandowns/spinner"
	"github.com/fatih/color"
)

type OutputMode int

const (
	ModeNormal OutputMode = iota
	ModeQuiet
	ModeJSON
)

var currentOutputMode = ModeNormal

func SetOutputMode(mode OutputMode) {
	currentOutputMode = mode
}

func GetOutputMode() OutputMode {
	return currentOutputMode
}

func IsSilent() bool {
	return currentOutputMode != ModeNormal
}

func PrintJSON(v any) {
	data, err := json.Marshal(v)
	if err != nil {
		fmt.Printf(`{"ok":false,"error":%q}`+"\n", err.Error())
		return
	}
	fmt.Println(string(data))
}

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
	if IsSilent() {
		return
	}
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
	if IsSilent() {
		return
	}
	fmt.Println()
	fmt.Println(Gray("  " + SubDivider))
	fmt.Printf("  %s\n", Bold(title))
	fmt.Println(Gray("  " + SubDivider))
	fmt.Println()
}

func Bullet(key, val string) {
	if IsSilent() {
		return
	}
	fmt.Printf("  %s  %-10s %s\n", Cyan("◆"), key, White(val))
}

func Substep(text string) {
	if IsSilent() {
		return
	}
	fmt.Printf("     %s  %s\n", Gray("·"), Gray(text))
}

func Warn(msg string) {
	if IsSilent() {
		return
	}
	fmt.Printf("  %s  %s\n", Yellow("⚠"), msg)
}

func Err(msg string) {
	if currentOutputMode == ModeJSON {
		return
	}
	fmt.Printf("  %s  %s\n", Red("✖"), Red(msg))
}

func Success(msg string) {
	if IsSilent() {
		return
	}
	fmt.Printf("  %s  %s\n", Green("✔"), msg)
}

func Abort(msg string) {
	if currentOutputMode == ModeJSON {
		PrintJSON(map[string]any{"ok": false, "error": msg})
		os.Exit(1)
	}
	if currentOutputMode == ModeQuiet {
		fmt.Fprintln(os.Stderr, msg)
		os.Exit(1)
	}
	fmt.Println()
	Err(msg)
	fmt.Println()
	os.Exit(1)
}

func NewSpinner(msg string) *spinner.Spinner {
	s := spinner.New(spinner.CharSets[14], 80*time.Millisecond)
	if IsSilent() {
		s.Writer = io.Discard
	}
	s.Prefix = "  "
	s.Suffix = "  " + msg
	_ = s.Color("cyan")
	return s
}

func PrintSummary(name, db, validator string, auth, docker, swagger bool) {
	if IsSilent() {
		return
	}
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
	if IsSilent() {
		return
	}
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

// PrintGoNotice displays a Golang-engine notice and, when an update is available,
// points to the Go GitHub releases rather than npm.
func PrintGoNotice(current, latest string, updateAvailable bool) {
	if IsSilent() {
		return
	}
	if !updateAvailable {
		fmt.Printf("  %s  %s\n\n", Cyan("◆"), Gray(fmt.Sprintf("Running the Golang engine (create-express-modular-go v%s)", strings.TrimPrefix(current, "v"))))
		return
	}

	const W = 62
	border := func(l, m, r string) string {
		return "  " + Yellow(l+strings.Repeat(m, W)+r)
	}
	blank := func() string {
		return "  " + Yellow("│") + strings.Repeat(" ", W) + Yellow("│")
	}
	printRow := func(colored string, plainLen int) {
		pad := W - plainLen - 4
		if pad < 0 {
			pad = 0
		}
		fmt.Printf("  %s    %s%s%s\n", Yellow("│"), colored, strings.Repeat(" ", pad), Yellow("│"))
	}

	titlePlain := fmt.Sprintf("Update available (Go engine)! %s → %s", current, latest)
	titleColored := fmt.Sprintf("Update available (Go engine)! %s → %s", Gray(current), Green(latest))

	enginePlain := "You are using the Golang version (create-express-modular-go)."
	engineColored := Cyan(enginePlain)

	cmdPlain := "Download: github.com/Levi9111/create-express-modular-go/releases"
	cmdColored := fmt.Sprintf("Download: %s", White("github.com/Levi9111/create-express-modular-go/releases"))

	fmt.Println(border("╭", "─", "╮"))
	fmt.Println(blank())
	printRow(titleColored, len([]rune(titlePlain)))
	printRow(engineColored, len([]rune(enginePlain)))
	fmt.Println(blank())
	printRow(cmdColored, len([]rune(cmdPlain)))
	fmt.Println(blank())
	fmt.Println(border("╰", "─", "╯"))
	fmt.Println()
}
