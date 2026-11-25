// (C) Copyright 2026 Hewlett Packard Enterprise Development LP

package process_test

import (
	"context"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"strings"
	"sync/atomic"
	"testing"
	"time"

	sdk "github.com/HPE/terraform-provider-hpe/internal/sdk/oapigen"
	"github.com/HPE/terraform-provider-hpe/morpheus/utils/process"
)

func newClient(t *testing.T, url string) *sdk.APIClient {
	t.Helper()

	cfg := sdk.NewConfiguration()
	cfg.Servers = sdk.ServerConfigurations{{URL: url}}

	return sdk.NewAPIClient(cfg)
}

func writeProcess(w http.ResponseWriter, status, message, errText, output string) {
	body := map[string]any{
		"process": map[string]any{
			"id":          int64(42),
			"displayName": "Apply Virtual Switch",
			"status":      status,
		},
	}
	proc := body["process"].(map[string]any)
	if message != "" {
		proc["message"] = message
	}
	if errText != "" {
		proc["error"] = errText
	}
	if output != "" {
		proc["output"] = output
	}

	w.Header().Set("Content-Type", "application/json")
	_ = json.NewEncoder(w).Encode(body)
}

func TestWaitComplete(t *testing.T) {
	t.Parallel()

	var calls int32

	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, _ *http.Request) {
		n := atomic.AddInt32(&calls, 1)
		if n < 2 {
			writeProcess(w, "running", "", "", "")

			return
		}

		writeProcess(w, process.StatusComplete, "", "", "")
	}))
	defer srv.Close()

	proc, err := process.Wait(context.Background(), newClient(t, srv.URL), 42, 30*time.Second)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	if proc == nil || proc.Status == nil || *proc.Status != process.StatusComplete {
		t.Fatalf("expected complete process, got %+v", proc)
	}
}

func TestWaitFailedWithMessage(t *testing.T) {
	t.Parallel()

	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, _ *http.Request) {
		writeProcess(w, process.StatusFailed, "host down", "", "")
	}))
	defer srv.Close()

	_, err := process.Wait(context.Background(), newClient(t, srv.URL), 42, 30*time.Second)
	if err == nil {
		t.Fatal("expected error")
	}

	if !strings.Contains(err.Error(), "host down") ||
		!strings.Contains(err.Error(), "Apply Virtual Switch") {
		t.Fatalf("error missing detail: %v", err)
	}
}

func TestWaitTimeout(t *testing.T) {
	t.Parallel()

	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, _ *http.Request) {
		writeProcess(w, "running", "", "", "")
	}))
	defer srv.Close()

	_, err := process.Wait(context.Background(), newClient(t, srv.URL), 42, 1*time.Second)
	if err == nil {
		t.Fatal("expected timeout error")
	}

	if !strings.Contains(err.Error(), "42") {
		t.Fatalf("timeout error should name process id: %v", err)
	}
}

func TestWaitContextCancel(t *testing.T) {
	t.Parallel()

	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, _ *http.Request) {
		writeProcess(w, "running", "", "", "")
	}))
	defer srv.Close()

	ctx, cancel := context.WithCancel(context.Background())
	go func() {
		time.Sleep(200 * time.Millisecond)
		cancel()
	}()

	_, err := process.Wait(ctx, newClient(t, srv.URL), 42, 30*time.Second)
	if err == nil {
		t.Fatal("expected cancellation error")
	}

	if !strings.Contains(err.Error(), "42") {
		t.Fatalf("cancel error should name process id: %v", err)
	}
}

func TestWarnings(t *testing.T) {
	t.Parallel()

	diags := process.Warnings("virtual switch", "sw1", []string{"noSpareNic", "", "another"})
	if diags.WarningsCount() != 2 {
		t.Fatalf("expected 2 warnings, got %d", diags.WarningsCount())
	}
}
