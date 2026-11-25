// (C) Copyright 2026 Hewlett Packard Enterprise Development LP

package appliancesetup

import (
	"context"
	"errors"
	"fmt"
	"net/http"
	"net/http/httptest"
	"sync/atomic"
	"testing"
	"time"

	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/hashicorp/terraform-plugin-go/tftypes"

	"github.com/HPE/terraform-provider-hpe/morpheus/model"
	"github.com/HPE/terraform-provider-hpe/morpheus/utils/clientfactory"
)

type deadlineTransport struct {
	testing *testing.T
	maximum time.Duration
}

func (transport deadlineTransport) RoundTrip(request *http.Request) (*http.Response, error) {
	deadline, bounded := request.Context().Deadline()
	if !bounded || time.Until(deadline) > transport.maximum {
		transport.testing.Error("setup request has no bounded operation deadline")
	}

	return nil, errors.New("simulated unresponsive appliance")
}

func TestUnitSetupReadAndImportHaveDeadlines(t *testing.T) {
	t.Parallel()
	for _, shorter := range []bool{false, true} {
		ctx := context.Background()
		maximum := time.Minute
		if shorter {
			maximum = time.Second
			var cancel context.CancelFunc
			ctx, cancel = context.WithTimeout(ctx, maximum)
			defer cancel()
		}
		factory := clientfactory.New(model.MorpheusProviderModel{
			URL: types.StringValue("https://appliance.example"), AccessToken: types.StringValue("test-token"),
		}, clientfactory.WithFactoryHTTPClient(&http.Client{Transport: deadlineTransport{testing: t, maximum: maximum}}))
		resourceUnderTest := &Resource{}
		var configured resource.ConfigureResponse
		resourceUnderTest.Configure(ctx, resource.ConfigureRequest{ProviderData: factory}, &configured)
		if configured.Diagnostics.HasError() {
			t.Fatal(configured.Diagnostics)
		}
		_, _, state := objectValue(t, stateValues("https://appliance.example"))
		read := resource.ReadResponse{State: state}
		resourceUnderTest.Read(ctx, resource.ReadRequest{State: state}, &read)
		if !read.Diagnostics.HasError() || !read.State.Raw.Equal(state.Raw) {
			t.Fatal("failed refresh must report an error and retain state")
		}
		_, _, empty := objectValue(t, nil)
		imported := resource.ImportStateResponse{State: empty}
		resourceUnderTest.ImportState(ctx, resource.ImportStateRequest{ID: "https://appliance.example"}, &imported)
		if !imported.Diagnostics.HasError() {
			t.Fatal("failed import must report an error")
		}
	}
}

func TestUnitSetupMissingAccountDoesNotDisableOtherGuards(t *testing.T) {
	t.Parallel()
	for _, name := range []string{"username", "email", "appliance_url"} {
		for _, removing := range []bool{false, true} {
			original := stateValues("https://appliance.example")
			delete(original, "account_name")
			_, _, state := objectValue(t, original)
			if removing {
				original[name] = tftypes.NewValue(tftypes.String, nil)
			} else {
				original[name] = str("changed")
			}
			plan, _, _ := objectValue(t, original)
			resourceUnderTest := &Resource{}
			planned := resource.ModifyPlanResponse{Plan: plan}
			resourceUnderTest.ModifyPlan(context.Background(), resource.ModifyPlanRequest{Plan: plan, State: state}, &planned)
			updated := resource.UpdateResponse{State: state}
			resourceUnderTest.Update(context.Background(), resource.UpdateRequest{Plan: plan, State: state}, &updated)
			if !planned.Diagnostics.HasError() || !updated.Diagnostics.HasError() {
				t.Fatalf("missing account_name bypassed guard for %s", name)
			}
		}
	}
}

func TestUnitSetupRemovingRecordedOptionalMetadataIsRefused(t *testing.T) {
	t.Parallel()
	for _, name := range []string{"appliance_name", "first_name", "last_name"} {
		original := stateValues("https://appliance.example")
		original[name] = str("recorded")
		_, _, state := objectValue(t, original)
		delete(original, name)
		plan, _, _ := objectValue(t, original)
		resourceUnderTest := &Resource{}
		planned := resource.ModifyPlanResponse{Plan: plan}
		resourceUnderTest.ModifyPlan(context.Background(), resource.ModifyPlanRequest{Plan: plan, State: state}, &planned)
		updated := resource.UpdateResponse{State: state}
		resourceUnderTest.Update(context.Background(), resource.UpdateRequest{Plan: plan, State: state}, &updated)
		if !planned.Diagnostics.HasError() || !updated.Diagnostics.HasError() || !updated.State.Raw.Equal(state.Raw) {
			t.Fatalf("removing %s must not erase the recorded setup identity", name)
		}
	}
}

func TestUnitBootstrapRefusesRedirects(t *testing.T) {
	t.Parallel()
	for _, status := range []int{http.StatusTemporaryRedirect, http.StatusPermanentRedirect} {
		for _, injected := range []bool{false, true} {
			t.Run(fmt.Sprintf("%d/injected=%t", status, injected), func(t *testing.T) {
				t.Parallel()
				var received atomic.Int64
				target := httptest.NewServer(http.HandlerFunc(func(writer http.ResponseWriter, _ *http.Request) {
					received.Add(1)
					writer.WriteHeader(http.StatusOK)
				}))
				defer target.Close()
				origin := httptest.NewServer(http.HandlerFunc(func(writer http.ResponseWriter, request *http.Request) {
					http.Redirect(writer, request, target.URL, status)
				}))
				defer origin.Close()
				custom := &http.Client{Timeout: 0}
				var options []clientfactory.FactoryOption
				if injected {
					options = append(options, clientfactory.WithFactoryHTTPClient(custom))
				}
				factory := clientfactory.New(model.MorpheusProviderModel{URL: types.StringValue(origin.URL)}, options...)
				unauthenticated, err := factory.NewUnauthenticated()
				if err != nil {
					t.Fatal(err)
				}
				_, err = newSetupClient(unauthenticated).init(context.Background(), initPayload{
					Password: "test-password", LicenseKey: "test-license",
				})
				if err == nil || received.Load() != 0 {
					t.Fatal("bootstrap redirect was not refused before forwarding credentials")
				}
				if custom.CheckRedirect != nil {
					t.Fatal("factory mutated the supplied HTTP client")
				}
			})
		}
	}
}

func TestUnitSetupIdentityChangesAreRefused(t *testing.T) {
	t.Parallel()
	ctx := context.Background()
	for _, name := range []string{
		"account_name", "username", "email", "appliance_url", "appliance_name",
		"first_name", "last_name", "hub_mode",
	} {
		t.Run(name, func(t *testing.T) {
			t.Parallel()
			original := stateValues("https://appliance.example")
			original[name] = str("recorded")
			_, _, state := objectValue(t, original)
			original[name] = str("changed")
			plan, _, _ := objectValue(t, original)
			resourceUnderTest := &Resource{}
			planned := resource.ModifyPlanResponse{Plan: plan}
			resourceUnderTest.ModifyPlan(ctx, resource.ModifyPlanRequest{Plan: plan, State: state}, &planned)
			if !planned.Diagnostics.HasError() {
				t.Fatal("plan accepted an immutable setup change")
			}
			updated := resource.UpdateResponse{State: state}
			resourceUnderTest.Update(ctx, resource.UpdateRequest{Plan: plan, State: state}, &updated)
			if !updated.Diagnostics.HasError() || !updated.State.Raw.Equal(state.Raw) {
				t.Fatal("update accepted an immutable setup change")
			}
		})
	}
}

func TestUnitSetupImportMetadataCanBeAdopted(t *testing.T) {
	t.Parallel()
	ctx := context.Background()
	_, _, imported := objectValue(t, map[string]tftypes.Value{
		"id": str("https://appliance.example"), "appliance_url": str("https://appliance.example"),
		"already_initialised": tftypes.NewValue(tftypes.Bool, true),
	})
	plan, _, _ := objectValue(t, stateValues("https://appliance.example"))
	resourceUnderTest := &Resource{}
	planned := resource.ModifyPlanResponse{Plan: plan}
	resourceUnderTest.ModifyPlan(ctx, resource.ModifyPlanRequest{Plan: plan, State: imported}, &planned)
	updated := resource.UpdateResponse{State: imported}
	resourceUnderTest.Update(ctx, resource.UpdateRequest{Plan: plan, State: imported}, &updated)
	if planned.Diagnostics.HasError() || updated.Diagnostics.HasError() {
		t.Fatalf("import adoption failed: %v %v", planned.Diagnostics, updated.Diagnostics)
	}
	var actual ApplianceSetupModel
	if diagnostics := updated.State.Get(ctx, &actual); diagnostics.HasError() {
		t.Fatal(diagnostics)
	}
	if actual.Id.ValueString() != "https://appliance.example" ||
		actual.ApplianceURL.ValueString() != "https://appliance.example" || !actual.AlreadyInitialised.ValueBool() {
		t.Fatal("import adoption lost identity or initialization status")
	}
}
