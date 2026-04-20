// (C) Copyright 2026 Hewlett Packard Enterprise Development LP

package testhelpers

import (
	"context"
	"errors"
	"fmt"
	"log"
	"net/http"
	"os"
	"reflect"
	"slices"
	"strings"

	"github.com/HewlettPackard/hpe-morpheus-go-sdk/oapigen/sdk"
	"github.com/hashicorp/terraform-plugin-testing/helper/resource"

	"github.com/HPE/terraform-provider-hpe/morpheus/utils/clientfactory"
	"github.com/HPE/terraform-provider-hpe/morpheus/utils/errfmt"
)

type SweepFilter func(ctx context.Context, client *sdk.APIClient, item any) (bool, string, error)

type SweepDelete func(ctx context.Context, client *sdk.APIClient, id int64, item any) (*http.Response, error)

type SweepOption func(*sweepConfig)

type sweepConfig struct {
	nameMethod         string
	idMethod           string
	filter             SweepFilter
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

func RegisterAPISweeper(
	resourceName string,
	namePrefix string,
	list func(ctx context.Context, client *sdk.APIClient, prefix string) (any, *http.Response, error),
	listItemsMethod string,
	delete SweepDelete,
	options ...SweepOption,
) {
	config := sweepConfig{
		nameMethod: "GetNameOk",
		idMethod:   "GetIdOk",
	}

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

		itemsResponse, hresp, err := list(ctx, client, namePrefix)
		if err != nil {
			if hresp != nil && slices.Contains(config.ignoreListStatuses, hresp.StatusCode) {
				log.Printf("[INFO] No %s found matching prefix (status %d): %s", resourceName, hresp.StatusCode, namePrefix)
				return nil
			}

			return fmt.Errorf("failed to list %s: %s", resourceName, errfmt.ErrMsg(err, hresp))
		}

		if hresp == nil || hresp.StatusCode != http.StatusOK {
			return fmt.Errorf("failed to list %s: %s", resourceName, errfmt.ErrMsg(err, hresp))
		}

		items, err := getItems(itemsResponse, listItemsMethod)
		if err != nil {
			return fmt.Errorf("failed to read %s list response: %w", resourceName, err)
		}

		var sweptCount int
		var sweepErrors []string

		for _, item := range items {
			name, ok := getStringPtrMethod(item, config.nameMethod)
			if !ok {
				continue
			}

			if !strings.HasPrefix(name, namePrefix) {
				log.Printf("[INFO] Skipping %s (name): %s", resourceName, name)
				continue
			}

			if config.filter != nil {
				allowed, reason, err := config.filter(ctx, client, item)
				if err != nil {
					errMsg := fmt.Sprintf("failed to evaluate %s %s: %s", resourceName, name, err)
					log.Printf("[ERROR] %s", errMsg)
					sweepErrors = append(sweepErrors, errMsg)
					continue
				}

				if !allowed {
					log.Printf("[INFO] Skipping %s (%s): %s", resourceName, reason, name)
					continue
				}
			}

			id, ok := getInt64PtrMethod(item, config.idMethod)
			if !ok {
				log.Printf("[INFO] Skipping %s (id): %s", resourceName, name)
				continue
			}

			log.Printf("[INFO] Sweeping %s: %s (id: %d)", resourceName, name, id)

			hresp, err := delete(ctx, client, id, item)
			if err != nil || hresp == nil || hresp.StatusCode != http.StatusOK {
				errMsg := fmt.Sprintf(
					"failed to delete %s %s (id: %d): %s",
					resourceName, name, id, errfmt.ErrMsg(err, hresp),
				)
				log.Printf("[ERROR] %s", errMsg)
				sweepErrors = append(sweepErrors, errMsg)
				continue
			}

			sweptCount++
		}

		log.Printf("[INFO] %s sweep completed. Resources swept: %d, errors: %d", resourceName, sweptCount, len(sweepErrors))

		return SweepErrorsToError(sweepErrors)
	})
}

func WithNameMethod(nameMethod string) SweepOption {
	return func(config *sweepConfig) {
		config.nameMethod = nameMethod
	}
}

func WithIDMethod(idMethod string) SweepOption {
	return func(config *sweepConfig) {
		config.idMethod = idMethod
	}
}

func WithFilter(filter SweepFilter) SweepOption {
	return func(config *sweepConfig) {
		config.filter = filter
	}
}

func WithIgnoreListStatuses(statuses ...int) SweepOption {
	return func(config *sweepConfig) {
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

func getItems(response any, accessor string) ([]any, error) {
	sliceValue, err := getSliceValue(response, accessor)
	if err != nil {
		return nil, err
	}

	items := make([]any, 0, sliceValue.Len())
	for i := 0; i < sliceValue.Len(); i++ {
		items = append(items, sliceValue.Index(i).Interface())
	}

	return items, nil
}

func getSliceValue(target any, accessor string) (reflect.Value, error) {
	method, err := getMethodValue(target, accessor)
	if err == nil {
		results := method.Call(nil)
		if len(results) != 1 {
			return reflect.Value{}, fmt.Errorf("method %s returned %d values", accessor, len(results))
		}

		sliceValue := results[0]
		if sliceValue.Kind() != reflect.Slice {
			return reflect.Value{}, fmt.Errorf("method %s did not return a slice", accessor)
		}

		return sliceValue, nil
	}

	value := reflect.ValueOf(target)
	if value.Kind() == reflect.Pointer {
		value = value.Elem()
	}

	if !value.IsValid() {
		return reflect.Value{}, fmt.Errorf("invalid target for accessor %s", accessor)
	}

	field := value.FieldByName(accessor)
	if !field.IsValid() {
		return reflect.Value{}, fmt.Errorf("accessor %s not found as method or field", accessor)
	}

	if field.Kind() != reflect.Slice {
		return reflect.Value{}, fmt.Errorf("field %s is not a slice", accessor)
	}

	return field, nil
}

func getStringPtrMethod(target any, methodName string) (string, bool) {
	method, err := getMethodValue(target, methodName)
	if err != nil {
		return "", false
	}

	results := method.Call(nil)
	if len(results) != 2 || results[0].Kind() != reflect.Pointer || results[0].IsNil() || !results[1].Bool() {
		return "", false
	}

	stringValue, ok := results[0].Interface().(*string)
	if !ok || stringValue == nil {
		return "", false
	}

	return *stringValue, true
}

func getInt64PtrMethod(target any, methodName string) (int64, bool) {
	method, err := getMethodValue(target, methodName)
	if err != nil {
		return 0, false
	}

	results := method.Call(nil)
	if len(results) != 2 || results[0].Kind() != reflect.Pointer || results[0].IsNil() || !results[1].Bool() {
		return 0, false
	}

	intValue, ok := results[0].Interface().(*int64)
	if !ok || intValue == nil {
		return 0, false
	}

	return *intValue, true
}

func getMethodValue(target any, methodName string) (reflect.Value, error) {
	value := reflect.ValueOf(target)
	method := value.MethodByName(methodName)
	if method.IsValid() {
		return method, nil
	}

	if value.Kind() != reflect.Pointer && value.CanAddr() {
		method = value.Addr().MethodByName(methodName)
		if method.IsValid() {
			return method, nil
		}
	}

	return reflect.Value{}, fmt.Errorf("method %s not found", methodName)
}

// RecoverSweepPanic recovers from panics during sweep execution and converts them to errors.
func RecoverSweepPanic(resourceName string, retErr *error) {
	if r := recover(); r != nil {
		*retErr = fmt.Errorf("panic during %s sweep: %v", resourceName, r)
		log.Printf("[ERROR] %v", *retErr)
	}
}

// SweepErrorsToError converts a slice of sweep error strings into a single error, or nil if empty.
func SweepErrorsToError(sweepErrors []string) error {
	if len(sweepErrors) > 0 {
		return fmt.Errorf("%s", strings.Join(sweepErrors, "\n"))
	}

	return nil
}
