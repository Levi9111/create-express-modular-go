package cmd

import (
	"bytes"
	"strings"
	"testing"
)

func TestRootCmdHelp(t *testing.T) {
	buf := new(bytes.Buffer)
	RootCmd.SetOut(buf)
	RootCmd.SetErr(buf)
	RootCmd.SetArgs([]string{"--help"})

	err := RootCmd.Execute()
	if err != nil {
		t.Fatalf("expected nil error on --help, got: %v", err)
	}

	out := buf.String()
	if !strings.Contains(out, "CEM CLI") && !strings.Contains(out, "Usage:") {
		t.Fatalf("expected help message in output, got: %s", out)
	}
}

func TestSubcommandsRegistered(t *testing.T) {
	expectedCmds := []string{
		"add",
		"remove",
		"list",
		"dev",
		"build",
		"start",
		"check",
		"fix",
		"eject",
	}

	cmdMap := make(map[string]bool)
	for _, c := range RootCmd.Commands() {
		cmdMap[c.Name()] = true
	}

	for _, name := range expectedCmds {
		if !cmdMap[name] {
			t.Errorf("expected subcommand %q to be registered on RootCmd", name)
		}
	}
}

func TestAddCommandHelp(t *testing.T) {
	buf := new(bytes.Buffer)
	RootCmd.SetOut(buf)
	RootCmd.SetErr(buf)
	RootCmd.SetArgs([]string{"add", "--help"})

	err := RootCmd.Execute()
	if err != nil {
		t.Fatalf("expected nil error on add --help, got: %v", err)
	}

	out := buf.String()
	if !strings.Contains(out, "add") {
		t.Fatalf("expected 'add' in help output, got: %s", out)
	}
}
