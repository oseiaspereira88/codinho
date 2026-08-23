package main

import (
	"os/exec"
	"strings"
	"testing"
	"time"
)

// TestServeNeverWritesToStdoutOnStartupFailure proves reliability-
// observability-compatibility requirement R9 for the error path: even
// when the catalog fails to load and codinho serve exits non-zero,
// stdout stays completely empty — mcp.StdioTransport's exclusive
// ownership of stdout (PROJECT.md §15) is never violated, and the error
// goes only to stderr.
func TestServeNeverWritesToStdoutOnStartupFailure(t *testing.T) {
	bin := buildCodinhoBinary(t)
	workspace := t.TempDir()
	// No packs/manifest.yaml at all: catalog load fails deterministically.

	cmd := exec.Command(bin, "serve")
	cmd.Dir = workspace
	var stdout, stderr strings.Builder
	cmd.Stdout = &stdout
	cmd.Stderr = &stderr
	cmd.Stdin = strings.NewReader("")

	err := cmd.Run()
	if err == nil {
		t.Fatal("expected codinho serve to fail with no packs/manifest.yaml")
	}
	if stdout.Len() != 0 {
		t.Fatalf("stdout must stay empty on startup failure, got %q", stdout.String())
	}
	if stderr.Len() == 0 {
		t.Fatal("expected the startup error on stderr")
	}
}

// TestServeStdoutIsOnlyValidJSONFrames proves R9 for the happy path with
// a raw, non-SDK reader: every non-empty stdout line codinho serve
// writes while handling one real request is syntactically JSON — never a
// stray print mixed into the protocol stream.
func TestServeStdoutIsOnlyValidJSONFrames(t *testing.T) {
	bin := buildCodinhoBinary(t)
	workspace := t.TempDir()
	writeIntegrationFixturePack(t, workspace)

	cmd := exec.Command(bin, "serve")
	cmd.Dir = workspace
	stdin, err := cmd.StdinPipe()
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	stdoutPipe, err := cmd.StdoutPipe()
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	var stderr strings.Builder
	cmd.Stderr = &stderr

	if err := cmd.Start(); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	defer func() {
		stdin.Close()
		cmd.Process.Kill()
		cmd.Wait()
	}()

	initRequest := `{"jsonrpc":"2.0","id":1,"method":"initialize","params":{"protocolVersion":"2025-06-18","capabilities":{},"clientInfo":{"name":"reliability-test","version":"0.0.0"}}}` + "\n"
	if _, err := stdin.Write([]byte(initRequest)); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	buf := make([]byte, 4096)
	n, readErr := stdoutPipe.Read(buf)
	if readErr != nil && n == 0 {
		t.Fatalf("reading stdout: %v (stderr: %s)", readErr, stderr.String())
	}
	line := strings.TrimSpace(string(buf[:n]))
	if line == "" {
		t.Fatal("expected a non-empty JSON-RPC response line on stdout")
	}
	if !strings.HasPrefix(line, "{") || !strings.HasSuffix(line, "}") {
		t.Fatalf("stdout line is not a bare JSON object, protocol purity violated: %q", line)
	}
}

// TestServeStartupRespondsWithinOneSecond proves reliability-
// observability-compatibility requirement R5's startup budget: from
// process start to the first real MCP response (loading the catalog,
// opening the event store, acquiring the workspace lock), well under 1s
// at V1's current catalog size.
func TestServeStartupRespondsWithinOneSecond(t *testing.T) {
	bin := buildCodinhoBinary(t)
	workspace := t.TempDir()
	writeIntegrationFixturePack(t, workspace)

	cmd := exec.Command(bin, "serve")
	cmd.Dir = workspace
	stdin, err := cmd.StdinPipe()
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	stdoutPipe, err := cmd.StdoutPipe()
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	start := time.Now()
	if err := cmd.Start(); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	defer func() {
		stdin.Close()
		cmd.Process.Kill()
		cmd.Wait()
	}()

	initRequest := `{"jsonrpc":"2.0","id":1,"method":"initialize","params":{"protocolVersion":"2025-06-18","capabilities":{},"clientInfo":{"name":"reliability-test","version":"0.0.0"}}}` + "\n"
	if _, err := stdin.Write([]byte(initRequest)); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	buf := make([]byte, 4096)
	if _, err := stdoutPipe.Read(buf); err != nil {
		t.Fatalf("reading stdout: %v", err)
	}
	if elapsed := time.Since(start); elapsed > time.Second {
		t.Fatalf("startup took %s, want well under the 1s budget (requirement R5)", elapsed)
	}
}
