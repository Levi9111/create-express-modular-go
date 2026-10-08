package telemetry

import (
	"encoding/json"
	"regexp"
	"testing"
)

func TestRandomUUID(t *testing.T) {
	uuidRegex := regexp.MustCompile(`^[0-9a-f]{8}-[0-9a-f]{4}-4[0-9a-f]{3}-[89ab][0-9a-f]{3}-[0-9a-f]{12}$`)

	for i := 0; i < 10; i++ {
		id := randomUUID()
		if !uuidRegex.MatchString(id) {
			t.Fatalf("expected valid v4 UUID, got %q", id)
		}
	}
}

func TestPayloadSerialization(t *testing.T) {
	payload := Payload{
		Command:        "create",
		CliVersion:     "3.3.10-go.exp",
		PackageManager: "pnpm",
		NodeVersion:    "go1.22.0",
		OS:             "linux-amd64",
		AnonID:         randomUUID(),
	}

	data, err := json.Marshal(payload)
	if err != nil {
		t.Fatalf("failed to marshal payload: %v", err)
	}

	var parsed Payload
	if err := json.Unmarshal(data, &parsed); err != nil {
		t.Fatalf("failed to unmarshal payload: %v", err)
	}

	if parsed.Command != payload.Command || parsed.PackageManager != payload.PackageManager {
		t.Fatalf("mismatched payload fields: %+v vs %+v", parsed, payload)
	}
}

func TestReportInstallOptOut(t *testing.T) {
	t.Setenv("CEM_TELEMETRY", "off")
	// Should return immediately without network access
	ReportInstall("test", "3.3.10", "npm")
}
