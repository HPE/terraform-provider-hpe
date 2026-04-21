// (C) Copyright 2026 Hewlett Packard Enterprise Development LP

package sweep

import (
	"context"
	"errors"
	"fmt"
	"log"
	"net/http"
	"os"
	"slices"
	"strings"

	"github.com/HewlettPackard/hpe-morpheus-go-sdk/oapigen/sdk"
	"github.com/hashicorp/terraform-plugin-testing/helper/resource"

	"github.com/HPE/terraform-provider-hpe/morpheus/utils/clientfactory"
	"github.com/HPE/terraform-provider-hpe/morpheus/utils/errfmt"
)

type TypedSweepFilter[T any] func(ctx context.Context, client *sdk.APIClient, item T) (bool, string, error)

type TypedSweepDelete[T any] func(ctx context.Context, client *sdk.APIClient, item T) (*http.Response, error)

type TypedSweepList[T any] func(ctx context.Context, client *sdk.APIClient) ([]T, *http.Response, error)

type TypedResourceCheck[T any] func(item T) bool

type TypedSweepOption[T any] func(*typedSweepConfig[T])

type typedSweepConfig[T any] struct {
	filter             TypedSweepFilter[T]
	ignoreListStatuses []int
}

func RegisterSweeper(resourceName string, sweep func() error) {
	resource.AddTestSweepers(
		resourceName,
		&resource.Sweeper{
			Name: resourceName,
			F: func(_ string) (retErr error) {
				defer RecoverSweepPanic(resourceName, &retErr)

				return sweep()
			},
		},
	)
}

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

	RegisterSweeper(resourceName, func() error {
		ctx := context.Background()

		client, err := NewSweepClient(ctx)
		if err != nil {
			log.Printf("[WARN] Cannot create sweep client: %v", err)

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

				continue
			}

			sweptCount++
		}

		errCount := 0
		if sweepErr != nil {
			errCount = len(strings.Split(sweepErr.Error(), "\n"))
		}

		log.Printf("[INFO] %s sweep completed. Resources swept: %d, errors: %d", resourceName, sweptCount, errCount)

		return sweepErr
	})
}

func WithFilter[T any](filter TypedSweepFilter[T]) TypedSweepOption[T] {
	return func(config *typedSweepConfig[T]) {
		config.filter = filter
	}
}

func WithIgnoreListStatuses[T any](statuses ...int) TypedSweepOption[T] {
	return func(config *typedSweepConfig[T]) {
		config.ignoreListStatuses = append(config.ignoreListStatuses, statuses...)
	}
}

func NewSweepClient(ctx context.Context) (*sdk.APIClient, error) {
	var username, password string

	url, ok := os.LookupEnv("TF_VAR_testacc_morpheus_url")
	if !ok {
		return nil, errors.New("TF_VAR_testacc_morpheus_url not set")
	}

	token, ok := os.LookupEnv("TF_VAR_testacc_morpheus_access_token")
	if !ok {
		username, ok = os.LookupEnv("TF_VAR_testacc_morpheus_username")
		if !ok {
			return nil, errors.New(
				"one of TF_VAR_testacc_morpheus_access_token or " +
					"TF_VAR_testacc_morpheus_username must be set",
			)
		}

		password, ok = os.LookupEnv("TF_VAR_testacc_morpheus_password")
		if !ok {
			return nil, errors.New(
				"one of TF_VAR_testacc_morpheus_access_token or " +
					"TF_VAR_testacc_morpheus_password must be set",
			)
		}
	}

	_, insecure := os.LookupEnv("TF_VAR_testacc_morpheus_insecure")
	var opts []clientfactory.ClientOption
	if insecure {
		opts = append(opts, clientfactory.WithInsecureTLS())
	}

	client := clientfactory.NewAPIClient(
		ctx,
		url,
		username,
		password,
		"",
		token,
		opts...,
	)

	return client, nil
}

// RecoverSweepPanic recovers from panics during sweep execution and converts them to errors.
func RecoverSweepPanic(resourceName string, retErr *error) {
	if r := recover(); r != nil {
		*retErr = fmt.Errorf("panic during %s sweep: %v", resourceName, r)
		log.Printf("[ERROR] %v", *retErr)
	}
}
