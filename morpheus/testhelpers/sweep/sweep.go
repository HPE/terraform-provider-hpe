// (C) Copyright 2026 Hewlett Packard Enterprise Development LP

package sweep

import (
	"context"
	"errors"
	"fmt"
	"log"
	"net/http"
	"slices"

	"github.com/HewlettPackard/hpe-morpheus-go-sdk/oapigen/sdk"
	"github.com/hashicorp/terraform-plugin-testing/helper/resource"

	"github.com/HPE/terraform-provider-hpe/morpheus/testhelpers"
	"github.com/HPE/terraform-provider-hpe/morpheus/utils/errfmt"
)

// TestResourcePrefix identifies acceptance test resources eligible for sweeping.
const TestResourcePrefix = "TestAccMorpheus"

// TypedSweepList returns all candidate resources that could be swept.
type TypedSweepList[T any] func(ctx context.Context, client *sdk.APIClient) ([]T, *http.Response, error)

// TypedResourceCheck decides whether a listed item is a test resource.
type TypedResourceCheck[T any] func(item T) bool

// TypedSweepDelete deletes a resource item selected for sweeping.
type TypedSweepDelete[T any] func(ctx context.Context, client *sdk.APIClient, item T) (*http.Response, error)

// TypedSweepOption configures optional sweep behavior.
type TypedSweepOption[T any] func(*typedSweepConfig[T])

// TypedSweepFilter applies additional checks before delete is attempted.
type TypedSweepFilter[T any] func(ctx context.Context, client *sdk.APIClient, item T) (bool, string, error)

type typedSweepConfig[T any] struct {
	filter             TypedSweepFilter[T]
	ignoreListStatuses []int
}

func registerSweeper(resourceName string, sweep func(string) error) {
	resource.AddTestSweepers(
		resourceName,
		&resource.Sweeper{
			Name: resourceName,
			F: func(system string) (retErr error) {
				defer recoverSweepPanic(resourceName, &retErr)

				return sweep(system)
			},
		},
	)
}

// RegisterTypedAPISweeper registers a typed sweeper with an explicit contract:
//   - listResource: list all resources that could be swept.
//   - isTestResource: decide if a listed item is a test resource.
//   - deleteResource: delete a matched test resource.
//
// Optional behavior can be added via options (for example, additional filters or
// ignored list status codes).
func RegisterTypedAPISweeper[T any](
	resourceName string,
	listResource TypedSweepList[T],
	isTestResource TypedResourceCheck[T],
	deleteResource TypedSweepDelete[T],
	options ...TypedSweepOption[T],
) {
	config := typedSweepConfig[T]{}

	for _, option := range options {
		option(&config)
	}

	registerSweeper(resourceName, func(system string) error {
		ctx := context.Background()

		client, err := testhelpers.NewClientForServer(ctx, system)
		if err != nil {
			log.Printf("[WARN] Cannot create sweep client for %q: %v", system, err)

			return nil
		}

		items, hresp, err := listResource(ctx, client)
		if err != nil {
			if hresp != nil && slices.Contains(config.ignoreListStatuses, hresp.StatusCode) {
				log.Printf("[INFO] No %s found (status %d)", resourceName, hresp.StatusCode)

				return nil
			}

			return fmt.Errorf("failed to list %s: %s", resourceName, errfmt.ErrMsg(err, hresp))
		}

		if hresp == nil || hresp.StatusCode != http.StatusOK {
			return fmt.Errorf("failed to list %s: %s", resourceName, errfmt.ErrMsg(err, hresp))
		}

		var sweptCount int
		var sweepErr error
		errCount := 0

		for _, item := range items {
			if !isTestResource(item) {
				continue
			}

			if config.filter != nil {
				allowed, reason, err := config.filter(ctx, client, item)
				if err != nil {
					errMsg := fmt.Sprintf("failed to evaluate filter for %s: %s", resourceName, err)
					log.Printf("[ERROR] %s", errMsg)
					sweepErr = errors.Join(sweepErr, errors.New(errMsg))
					errCount++

					continue
				}

				if !allowed {
					log.Printf("[INFO] Skipping %s (%s)", resourceName, reason)

					continue
				}
			}

			log.Printf("[INFO] Sweeping %s", resourceName)

			hresp, err := deleteResource(ctx, client, item)
			if err != nil || hresp == nil || hresp.StatusCode != http.StatusOK {
				errMsg := fmt.Sprintf("failed to delete %s: %s", resourceName, errfmt.ErrMsg(err, hresp))
				log.Printf("[ERROR] %s", errMsg)
				sweepErr = errors.Join(sweepErr, errors.New(errMsg))
				errCount++

				continue
			}

			sweptCount++
		}

		log.Printf("[INFO] %s sweep completed. Resources swept: %d, errors: %d", resourceName, sweptCount, errCount)

		return sweepErr
	})
}

// WithFilter adds an optional post-check before deleteResource is called.
func WithFilter[T any](filter TypedSweepFilter[T]) TypedSweepOption[T] {
	return func(config *typedSweepConfig[T]) {
		config.filter = filter
	}
}

// WithIgnoreListStatuses treats listed HTTP status codes from listResource as
// non-fatal and returns success for the sweep.
func WithIgnoreListStatuses[T any](statuses ...int) TypedSweepOption[T] {
	return func(config *typedSweepConfig[T]) {
		config.ignoreListStatuses = append(config.ignoreListStatuses, statuses...)
	}
}

// RecoverSweepPanic recovers from panics during sweep execution and converts them to errors.
func recoverSweepPanic(resourceName string, retErr *error) {
	if r := recover(); r != nil {
		*retErr = fmt.Errorf("panic during %s sweep: %v", resourceName, r)
		log.Printf("[ERROR] %v", *retErr)
	}
}
