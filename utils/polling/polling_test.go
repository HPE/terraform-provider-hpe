package polling

import (
	"context"
	"fmt"
	"net/http"
	"testing"
	"time"
)

func TestWaitForStatus_Success(t *testing.T) {
	callCount := 0
	fn := func(ctx context.Context) (string, *http.Response, error) {
		callCount++
		if callCount < 3 {
			return "provisioning", nil, nil
		}

		return "running", nil, nil
	}

	cfg := Config{
		Interval:       10 * time.Millisecond,
		MaxTimeout:     5 * time.Second,
		TargetStatuses: []string{"running"},
		ErrorStatuses:  []string{"failed"},
	}

	status, err := WaitForStatus(context.Background(), cfg, fn)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if status != "running" {
		t.Fatalf("expected status 'running', got %q", status)
	}
	if callCount != 3 {
		t.Fatalf("expected 3 calls, got %d", callCount)
	}
}

func TestWaitForStatus_ErrorStatus(t *testing.T) {
	fn := func(ctx context.Context) (string, *http.Response, error) {
		return "failed", nil, nil
	}

	cfg := Config{
		Interval:       10 * time.Millisecond,
		MaxTimeout:     1 * time.Second,
		TargetStatuses: []string{"running"},
		ErrorStatuses:  []string{"failed"},
	}

	_, err := WaitForStatus(context.Background(), cfg, fn)
	if err == nil {
		t.Fatal("expected error, got nil")
	}
	if err.Error() != "resource reached error status: failed" {
		t.Fatalf("unexpected error message: %v", err)
	}
}

func TestWaitForStatus_Timeout(t *testing.T) {
	fn := func(ctx context.Context) (string, *http.Response, error) {
		return "provisioning", nil, nil
	}

	cfg := Config{
		Interval:       10 * time.Millisecond,
		MaxTimeout:     50 * time.Millisecond,
		TargetStatuses: []string{"running"},
		ErrorStatuses:  []string{"failed"},
	}

	_, err := WaitForStatus(context.Background(), cfg, fn)
	if err == nil {
		t.Fatal("expected timeout error, got nil")
	}
}

func TestWaitForStatus_NotFound(t *testing.T) {
	fn := func(ctx context.Context) (string, *http.Response, error) {
		return "", &http.Response{StatusCode: http.StatusNotFound}, fmt.Errorf("not found")
	}

	cfg := Config{
		Interval:       10 * time.Millisecond,
		MaxTimeout:     1 * time.Second,
		TargetStatuses: []string{"running"},
		ErrorStatuses:  []string{"failed"},
	}

	_, err := WaitForStatus(context.Background(), cfg, fn)
	if err == nil {
		t.Fatal("expected error, got nil")
	}
}

func TestWaitForDeletion_Success(t *testing.T) {
	callCount := 0
	fn := func(ctx context.Context) (bool, error) {
		callCount++
		if callCount < 3 {
			return true, nil
		}

		return false, nil
	}

	cfg := Config{
		Interval:   10 * time.Millisecond,
		MaxTimeout: 5 * time.Second,
	}

	err := WaitForDeletion(context.Background(), cfg, fn)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if callCount != 3 {
		t.Fatalf("expected 3 calls, got %d", callCount)
	}
}

func TestWaitForDeletion_Timeout(t *testing.T) {
	fn := func(ctx context.Context) (bool, error) {
		return true, nil // always exists
	}

	cfg := Config{
		Interval:   10 * time.Millisecond,
		MaxTimeout: 50 * time.Millisecond,
	}

	err := WaitForDeletion(context.Background(), cfg, fn)
	if err == nil {
		t.Fatal("expected timeout error, got nil")
	}
}

func TestWaitForStatusOrDeletion_Deleted(t *testing.T) {
	fn := func(ctx context.Context) (string, *http.Response, error) {
		return "", &http.Response{StatusCode: http.StatusNotFound}, fmt.Errorf("not found")
	}

	cfg := Config{
		Interval:       10 * time.Millisecond,
		MaxTimeout:     1 * time.Second,
		TargetStatuses: []string{"stopped"},
		ErrorStatuses:  []string{"failed"},
	}

	status, err := WaitForStatusOrDeletion(context.Background(), cfg, fn)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if status != "deleted" {
		t.Fatalf("expected 'deleted', got %q", status)
	}
}

func TestWaitForStatusOrDeletion_Status(t *testing.T) {
	fn := func(ctx context.Context) (string, *http.Response, error) {
		return "stopped", nil, nil
	}

	cfg := Config{
		Interval:       10 * time.Millisecond,
		MaxTimeout:     1 * time.Second,
		TargetStatuses: []string{"stopped"},
		ErrorStatuses:  []string{"failed"},
	}

	status, err := WaitForStatusOrDeletion(context.Background(), cfg, fn)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if status != "stopped" {
		t.Fatalf("expected 'stopped', got %q", status)
	}
}
