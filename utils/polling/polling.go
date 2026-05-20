package polling

import (
	"context"
	"fmt"
	"net/http"
	"time"

	"github.com/cenkalti/backoff/v5"
)

// StatusFunc fetches the current status of a resource.
// Returns the status string, the HTTP response, and any error.
type StatusFunc func(ctx context.Context) (status string, httpResp *http.Response, err error)

// ExistenceFunc checks if a resource still exists.
// Returns true if it exists, false if gone (404).
type ExistenceFunc func(ctx context.Context) (exists bool, err error)

// Config holds polling configuration.
type Config struct {
	// Interval between polls. Default: 5s
	Interval time.Duration
	// MaxTimeout is the maximum time to wait. Default: 45m
	MaxTimeout time.Duration
	// TargetStatuses - polling succeeds when status matches one of these
	TargetStatuses []string
	// ErrorStatuses - polling fails immediately when status matches one of these
	ErrorStatuses []string
}

// Common status sets for reuse across resources.
var (
	// InstanceRunningTargets - instance is successfully provisioned
	InstanceRunningTargets = []string{"running"}

	// InstanceErrorStatuses - instance hit a terminal error
	InstanceErrorStatuses = []string{
		"denied", "cancelled", "failed", "stopped",
		"suspended", "removing", "pendingRemoval",
	}

	// ProvisionedTargets - resource is provisioned but not necessarily running
	ProvisionedTargets = []string{"provisioned", "stopped", "running"}

	// GenericErrorStatuses - common error statuses for most resources
	GenericErrorStatuses = []string{"denied", "cancelled", "failed"}
)

// DefaultConfig returns sensible defaults.
func DefaultConfig() Config {
	return Config{
		Interval:   5 * time.Second,
		MaxTimeout: 45 * time.Minute,
	}
}

// ForCreate returns a Config suitable for resource creation polling.
func ForCreate(timeout time.Duration) Config {
	return Config{
		Interval:       5 * time.Second,
		MaxTimeout:     timeout,
		TargetStatuses: InstanceRunningTargets,
		ErrorStatuses:  InstanceErrorStatuses,
	}
}

// ForDelete returns a Config suitable for resource deletion polling.
func ForDelete(timeout time.Duration) Config {
	return Config{
		Interval:       5 * time.Second,
		MaxTimeout:     timeout,
		TargetStatuses: []string{"deleted"},
		ErrorStatuses:  GenericErrorStatuses,
	}
}

// ForStatus returns a Config with custom target and error statuses.
func ForStatus(timeout time.Duration, targets []string, errors []string) Config {
	return Config{
		Interval:       5 * time.Second,
		MaxTimeout:     timeout,
		TargetStatuses: targets,
		ErrorStatuses:  errors,
	}
}

// WaitForStatus polls until the resource reaches a target status or errors.
// Returns the final status string or an error.
func WaitForStatus(ctx context.Context, cfg Config, fn StatusFunc) (string, error) {
	poll := func() (string, error) {
		status, httpResp, err := fn(ctx)
		if err != nil {
			if httpResp != nil && httpResp.StatusCode == http.StatusNotFound {
				return "", backoff.Permanent(fmt.Errorf("resource not found during status poll"))
			}

			return "", backoff.Permanent(err)
		}

		return status, checkStatus(status, cfg.TargetStatuses, cfg.ErrorStatuses)
	}

	return backoff.Retry(
		ctx, poll,
		backoff.WithBackOff(backoff.NewConstantBackOff(cfg.Interval)),
		backoff.WithMaxElapsedTime(cfg.MaxTimeout),
	)
}

// WaitForDeletion polls until the resource returns exists=false.
// Used for delete operations.
func WaitForDeletion(ctx context.Context, cfg Config, fn ExistenceFunc) error {
	poll := func() (struct{}, error) {
		exists, err := fn(ctx)
		if err != nil {
			return struct{}{}, backoff.Permanent(err)
		}
		if !exists {
			return struct{}{}, nil
		}

		return struct{}{}, fmt.Errorf("resource still exists")
	}

	_, err := backoff.Retry(
		ctx, poll,
		backoff.WithBackOff(backoff.NewConstantBackOff(cfg.Interval)),
		backoff.WithMaxElapsedTime(cfg.MaxTimeout),
	)

	return err
}

// WaitForStatusOrDeletion polls until the resource reaches a target status OR
// is deleted (404). Useful for delete flows where you wait for "stopped" then delete.
func WaitForStatusOrDeletion(ctx context.Context, cfg Config, fn StatusFunc) (string, error) {
	poll := func() (string, error) {
		status, httpResp, err := fn(ctx)
		if err != nil {
			if httpResp != nil && httpResp.StatusCode == http.StatusNotFound {
				return "deleted", nil
			}

			return "", backoff.Permanent(err)
		}

		return status, checkStatus(status, cfg.TargetStatuses, cfg.ErrorStatuses)
	}

	return backoff.Retry(
		ctx, poll,
		backoff.WithBackOff(backoff.NewConstantBackOff(cfg.Interval)),
		backoff.WithMaxElapsedTime(cfg.MaxTimeout),
	)
}

// checkStatus evaluates the current status against target and error sets.
func checkStatus(status string, targets, errors []string) error {
	for _, s := range errors {
		if status == s {
			return backoff.Permanent(fmt.Errorf("resource reached error status: %s", status))
		}
	}
	for _, s := range targets {
		if status == s {
			return nil
		}
	}

	return fmt.Errorf("status %q not yet in target set", status)
}
