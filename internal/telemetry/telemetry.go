package telemetry

import (
	"bytes"
	"crypto/rand"
	"encoding/json"
	"fmt"
	"net/http"
	"os"
	"runtime"
	"time"
)

const Endpoint = "https://create-express-modular.lovable.app/api/public/install-ping"

type Payload struct {
	Command        string `json:"command"`
	CliVersion     string `json:"cliVersion"`
	PackageManager string `json:"packageManager"`
	NodeVersion    string `json:"nodeVersion"`
	OS             string `json:"os"`
	AnonID         string `json:"anonId"`
}

func randomUUID() string {
	b := make([]byte, 16)
	_, _ = rand.Read(b)
	b[6] = (b[6] & 0x0f) | 0x40 // Version 4
	b[8] = (b[8] & 0x3f) | 0x80 // Variant 10
	return fmt.Sprintf("%x-%x-%x-%x-%x", b[0:4], b[4:6], b[6:8], b[8:10], b[10:])
}

// ReportInstall sends an anonymous ping upon successful scaffolding.
// Wrapped so it never blocks or fails project generation.
func ReportInstall(command, cliVersion, pm string) {
	if os.Getenv("CEM_TELEMETRY") == "off" {
		return
	}

	payload := Payload{
		Command:        command,
		CliVersion:     cliVersion,
		PackageManager: pm,
		NodeVersion:    runtime.Version(),
		OS:             fmt.Sprintf("%s-%s", runtime.GOOS, runtime.GOARCH),
		AnonID:         randomUUID(),
	}

	data, err := json.Marshal(payload)
	if err != nil {
		return
	}

	client := &http.Client{Timeout: 1500 * time.Millisecond}
	req, err := http.NewRequest("POST", Endpoint, bytes.NewBuffer(data))
	if err != nil {
		return
	}
	req.Header.Set("Content-Type", "application/json")

	resp, err := client.Do(req)
	if err == nil && resp != nil {
		_ = resp.Body.Close()
	}
}
