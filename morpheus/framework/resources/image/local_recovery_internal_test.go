package image

import (
	"context"
	"fmt"
	"net/http"
	"net/http/httptest"
	"strings"
	"sync/atomic"
	"testing"

	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/hashicorp/terraform-plugin-go/tftypes"

	"github.com/HPE/terraform-provider-hpe/morpheus/model"
	"github.com/HPE/terraform-provider-hpe/morpheus/utils/clientfactory"
)

func newPrivate[Private any](_ *Private) *Private {
	return new(Private)
}

func TestUnitImageUploadFailureRetainsIdentity(t *testing.T) {
	t.Parallel()
	for _, operation := range []string{"create", "update"} {
		for _, failure := range []string{"upload", "listing", "cleanup"} {
			t.Run(operation+"/"+failure, func(t *testing.T) {
				t.Parallel()
				ctx := context.Background()
				store := newFakeImageStore(map[string][]byte{"older.iso": []byte("old")})
				var uploaded atomic.Bool
				handler := store.handler(t)
				server := httptest.NewServer(http.HandlerFunc(func(writer http.ResponseWriter, request *http.Request) {
					writer.Header().Set("Content-Type", "application/json")
					if request.Method == http.MethodPost && request.URL.Path == "/api/virtual-images" {
						fmt.Fprint(writer, `{"success":true,"virtualImage":{"id":42}}`)

						return
					}
					if (failure == "upload" && request.Method == http.MethodPost) ||
						(failure == "listing" && uploaded.Load() && request.Method == http.MethodGet) ||
						(failure == "cleanup" && request.Method == http.MethodDelete) {
						writer.WriteHeader(http.StatusInternalServerError)
						fmt.Fprint(writer, `{"success":false,"msg":"injected failure"}`)

						return
					}
					handler.ServeHTTP(writer, request)
					if request.Method == http.MethodPost {
						uploaded.Store(true)
					}
				}))
				defer server.Close()
				underTest := &Resource{}
				var configured resource.ConfigureResponse
				underTest.Configure(ctx, resource.ConfigureRequest{ProviderData: clientfactory.New(model.MorpheusProviderModel{
					URL: types.StringValue(server.URL), AccessToken: types.StringValue("test-token"),
				})}, &configured)
				if configured.Diagnostics.HasError() {
					t.Fatal(configured.Diagnostics)
				}
				values := map[string]tftypes.Value{
					"name": tftypes.NewValue(tftypes.String, "image"), "image_type": tftypes.NewValue(tftypes.String, "iso"),
					"file_path": tftypes.NewValue(tftypes.String, writeTempFile(t, "disk.iso", []byte("bytes"))),
				}
				config, plan, state := localObject(t, values)
				if operation == "create" {
					state.Raw = tftypes.NewValue(state.Raw.Type(), nil)
					response := resource.CreateResponse{State: state}
					response.Private = newPrivate(response.Private)
					underTest.Create(ctx, resource.CreateRequest{Config: config, Plan: plan}, &response)
					if !response.Diagnostics.HasError() || response.State.Raw.IsNull() || !response.State.Raw.IsFullyKnown() {
						t.Fatalf("failed create did not retain recoverable state: %v", response.Diagnostics)
					}
					var actual ImageModel
					if diagnostics := response.State.Get(ctx, &actual); diagnostics.HasError() || actual.Id.ValueInt64() != 42 {
						t.Fatalf("failed create lost image ID: %v", diagnostics)
					}
				} else {
					values["id"] = tftypes.NewValue(tftypes.Number, 42)
					_, _, state = localObject(t, values)
					response := resource.UpdateResponse{State: state}
					underTest.Update(ctx, resource.UpdateRequest{Config: config, Plan: plan, State: state}, &response)
					if !response.Diagnostics.HasError() || !response.State.Raw.Equal(state.Raw) {
						t.Fatalf("failed update must preserve prior image state: %v", response.Diagnostics)
					}
				}
				store.mu.Lock()
				defer store.mu.Unlock()
				if store.files["older.iso"] == nil {
					t.Fatal("failure removed unverified older file")
				}
				for _, call := range store.calls {
					if strings.HasPrefix(call, "DELETE /api/virtual-images/42?") {
						t.Fatal("upload failure deleted the image record")
					}
				}
			})
		}
	}
}
