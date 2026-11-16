// (C) Copyright 2026 Hewlett Packard Enterprise Development LP

package appliancesetup

import (
	"context"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"strings"
	"sync"
	"testing"
	"time"

	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-framework/tfsdk"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/hashicorp/terraform-plugin-go/tftypes"

	"github.com/HPE/terraform-provider-hpe/morpheus/model"
	"github.com/HPE/terraform-provider-hpe/morpheus/utils/clientfactory"
)

// fakeAppliance is the setup API of one appliance. Before `ready` polls it
// answers like a front-end proxy whose application is still starting; once
// set up it refuses init and blanks the version for unauthenticated callers.
type fakeAppliance struct {
	mu          sync.Mutex
	notReady    int // number of check calls answered with a loading page
	setUp       bool
	refuseMsg   string            // when set, init is refused with this message
	refuseErrs  map[string]string // per-field errors sent with the refusal
	checks      int
	inits       []initPayload
	whoamiCalls int
}

func (f *fakeAppliance) handler(t *testing.T) http.Handler {
	t.Helper()

	mux := http.NewServeMux()

	mux.HandleFunc("GET "+checkPath, func(w http.ResponseWriter, r *http.Request) {
		f.mu.Lock()
		defer f.mu.Unlock()

		f.checks++

		if f.checks <= f.notReady {
			w.Header().Set("Content-Type", "text/html")
			w.WriteHeader(http.StatusBadGateway)
			_, _ = w.Write([]byte("<html><body><h1>Morpheus Loading</h1></body></html>"))

			return
		}

		if r.Header.Get("Authorization") != "" {
			t.Errorf("check must be unauthenticated; got Authorization header")
		}

		answer := map[string]any{"success": true, "buildVersion": ""}
		if !f.setUp {
			answer["setupNeeded"] = true
			answer["applianceUrl"] = "https://appliance.example"
		}

		_ = json.NewEncoder(w).Encode(answer)
	})

	mux.HandleFunc("POST "+initPath, func(w http.ResponseWriter, r *http.Request) {
		f.mu.Lock()
		defer f.mu.Unlock()

		var payload initPayload
		if err := json.NewDecoder(r.Body).Decode(&payload); err != nil {
			t.Errorf("init body: %v", err)
		}

		f.inits = append(f.inits, payload)

		if f.setUp {
			w.WriteHeader(http.StatusBadRequest)
			_ = json.NewEncoder(w).Encode(map[string]any{"success": false, "msg": "Appliance has already been setup"})

			return
		}

		if f.refuseMsg != "" {
			w.WriteHeader(http.StatusBadRequest)
			_ = json.NewEncoder(w).Encode(map[string]any{"success": false, "msg": f.refuseMsg, "errors": f.refuseErrs})

			return
		}

		f.setUp = true

		_ = json.NewEncoder(w).Encode(map[string]any{"success": true, "buildVersion": "9.1.0", "applianceUrl": payload.ApplianceURL})
	})

	mux.HandleFunc("GET /api/whoami", func(w http.ResponseWriter, r *http.Request) {
		f.mu.Lock()
		defer f.mu.Unlock()

		f.whoamiCalls++

		if r.Header.Get("Authorization") == "" {
			w.WriteHeader(http.StatusUnauthorized)

			return
		}

		w.Header().Set("Content-Type", "application/json")
		_ = json.NewEncoder(w).Encode(map[string]any{"appliance": map[string]any{"buildVersion": "9.1.0.2"}})
	})

	return mux
}

// newResource returns a configured resource pointed at the fake appliance.
func newResource(t *testing.T, f *fakeAppliance) *Resource {
	t.Helper()

	srv := httptest.NewServer(f.handler(t))
	t.Cleanup(srv.Close)

	cf := clientfactory.New(model.MorpheusProviderModel{
		URL:         types.StringValue(srv.URL + "/"),
		AccessToken: types.StringValue("unit-test"),
	})

	r := &Resource{}

	var resp resource.ConfigureResponse

	r.Configure(context.Background(), resource.ConfigureRequest{ProviderData: cf}, &resp)
	if resp.Diagnostics.HasError() {
		t.Fatal(resp.Diagnostics)
	}

	return r
}

// objectValue builds a value of the resource's object type with the given
// attributes set and every other one null.
func objectValue(t *testing.T, set map[string]tftypes.Value) (tfsdk.Plan, tfsdk.Config, tfsdk.State) {
	t.Helper()

	ctx := context.Background()

	var sr resource.SchemaResponse

	(&Resource{}).Schema(ctx, resource.SchemaRequest{}, &sr)

	objectType, ok := sr.Schema.Type().TerraformType(ctx).(tftypes.Object)
	if !ok {
		t.Fatal("schema type is not an object")
	}

	values := make(map[string]tftypes.Value, len(objectType.AttributeTypes))
	for name, typ := range objectType.AttributeTypes {
		values[name] = tftypes.NewValue(typ, nil)
	}

	for name, v := range set {
		values[name] = v
	}

	raw := tftypes.NewValue(objectType, values)

	return tfsdk.Plan{Schema: sr.Schema, Raw: raw},
		tfsdk.Config{Schema: sr.Schema, Raw: raw},
		tfsdk.State{Schema: sr.Schema, Raw: raw}
}

func str(v string) tftypes.Value { return tftypes.NewValue(tftypes.String, v) }

var unknownString = tftypes.NewValue(tftypes.String, tftypes.UnknownValue)

func inputs() map[string]tftypes.Value {
	return map[string]tftypes.Value{
		"account_name":   str("HPE"),
		"username":       str("admin"),
		"password_wo":    str("Sup3r-secret"),
		"email":          str("admin@example.com"),
		"first_name":     str("Lab"),
		"last_name":      str("Admin"),
		"hub_mode":       str("skip"),
		"license_key_wo": str("LICENSE-KEY"),
		"appliance_url":  unknownString,
		"id":             unknownString,
	}
}

func create(t *testing.T, r *Resource, set map[string]tftypes.Value) (ApplianceSetupModel, resource.CreateResponse) {
	t.Helper()

	_, config, state := objectValue(t, set)

	// The write-only values are null in the plan, as the framework sends
	// them; only the configuration carries them.
	planValues := map[string]tftypes.Value{}
	for k, v := range set {
		planValues[k] = v
	}

	planValues["password_wo"] = tftypes.NewValue(tftypes.String, nil)
	planValues["license_key_wo"] = tftypes.NewValue(tftypes.String, nil)
	plan, _, _ := objectValue(t, planValues)

	resp := resource.CreateResponse{State: state}

	r.Create(context.Background(), resource.CreateRequest{Plan: plan, Config: config}, &resp)

	var out ApplianceSetupModel
	if !resp.Diagnostics.HasError() {
		if d := resp.State.Get(context.Background(), &out); d.HasError() {
			t.Fatal(d)
		}
	}

	return out, resp
}

func TestUnitCreateInitialisesAFreshAppliance(t *testing.T) {
	t.Parallel()

	f := &fakeAppliance{}
	r := newResource(t, f)

	out, resp := create(t, r, inputs())
	if resp.Diagnostics.HasError() {
		t.Fatal(resp.Diagnostics)
	}

	f.mu.Lock()
	defer f.mu.Unlock()

	if len(f.inits) != 1 {
		t.Fatalf("init calls = %d, want 1", len(f.inits))
	}

	sent := f.inits[0]
	if sent.Password != "Sup3r-secret" || sent.LicenseKey != "LICENSE-KEY" || sent.HubMode != "skip" ||
		sent.AccountName != "HPE" || sent.Username != "admin" || sent.FirstName != "Lab" {
		t.Errorf("init payload = %+v", sent)
	}

	if sent.ApplianceURL == "" || strings.HasSuffix(sent.ApplianceURL, "/") {
		t.Errorf("applianceUrl should default to the provider url without a trailing slash; got %q", sent.ApplianceURL)
	}

	if out.AlreadyInitialised.ValueBool() {
		t.Error("already_initialised should be false after an init")
	}

	if out.Id.ValueString() != sent.ApplianceURL || out.ApplianceURL.ValueString() != sent.ApplianceURL {
		t.Errorf("id/appliance_url = %q/%q, want %q", out.Id.ValueString(), out.ApplianceURL.ValueString(), sent.ApplianceURL)
	}

	if out.BuildVersion.ValueString() != "9.1.0" {
		t.Errorf("build_version = %q, want the init answer's", out.BuildVersion.ValueString())
	}

	if !out.PasswordWo.IsNull() || !out.LicenseKeyWo.IsNull() {
		t.Error("write-only values must not reach state")
	}

	if len(resp.Diagnostics) != 0 {
		t.Errorf("unexpected diagnostics: %v", resp.Diagnostics)
	}
}

func TestUnitCreateWaitsForAStartingAppliance(t *testing.T) {
	t.Parallel()

	f := &fakeAppliance{notReady: 3}
	r := newResource(t, f)

	unauth, err := r.NewUnauthenticated()
	if err != nil {
		t.Fatal(err)
	}

	c := newSetupClient(unauth)
	c.pollInterval = 5 * time.Millisecond

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	check, err := c.awaitAPI(ctx)
	if err != nil {
		t.Fatal(err)
	}

	if !check.SetupNeeded {
		t.Error("a fresh appliance should report setupNeeded")
	}

	f.mu.Lock()
	defer f.mu.Unlock()

	if f.checks != 4 {
		t.Errorf("check calls = %d, want 4 (three loading pages, one answer)", f.checks)
	}
}

func TestUnitAwaitAPITimesOut(t *testing.T) {
	t.Parallel()

	f := &fakeAppliance{notReady: 1 << 30}
	r := newResource(t, f)

	unauth, err := r.NewUnauthenticated()
	if err != nil {
		t.Fatal(err)
	}

	c := newSetupClient(unauth)
	c.pollInterval = 5 * time.Millisecond

	ctx, cancel := context.WithTimeout(context.Background(), 50*time.Millisecond)
	defer cancel()

	_, err = c.awaitAPI(ctx)
	if err == nil || !strings.Contains(err.Error(), "timed out") || !strings.Contains(err.Error(), "Morpheus Loading") {
		t.Fatalf("err = %v, want a timeout quoting the last answer", err)
	}
}

func TestUnitCreateAdoptsASetUpAppliance(t *testing.T) {
	t.Parallel()

	f := &fakeAppliance{setUp: true}
	r := newResource(t, f)

	out, resp := create(t, r, inputs())
	if resp.Diagnostics.HasError() {
		t.Fatal(resp.Diagnostics)
	}

	f.mu.Lock()
	defer f.mu.Unlock()

	if len(f.inits) != 0 {
		t.Fatalf("init must not be called on a set-up appliance; calls = %d", len(f.inits))
	}

	if !out.AlreadyInitialised.ValueBool() {
		t.Error("already_initialised should be true")
	}

	if out.BuildVersion.ValueString() != "9.1.0.2" {
		t.Errorf("build_version = %q, want the whoami answer's", out.BuildVersion.ValueString())
	}

	if resp.Diagnostics.WarningsCount() != 1 {
		t.Errorf("want one warning about adopting; got %v", resp.Diagnostics)
	}
}

func TestUnitCreateReportsARefusal(t *testing.T) {
	t.Parallel()

	f := &fakeAppliance{refuseMsg: "Failed to create master user", refuseErrs: map[string]string{
		"password": "Password must contain a symbol", "email": "invalid",
	}}
	r := newResource(t, f)

	_, resp := create(t, r, inputs())
	if !resp.Diagnostics.HasError() {
		t.Fatal("want an error")
	}

	text := resp.Diagnostics.Errors()[0].Detail()
	if !strings.Contains(text, "Failed to create master user (email: invalid; password: Password must contain a symbol)") {
		t.Errorf("detail = %q", text)
	}

	if strings.Contains(text, "Sup3r-secret") || strings.Contains(text, "LICENSE-KEY") {
		t.Error("the diagnostic must not carry the secrets")
	}
}

func TestUnitCreateRaceWithAnotherSetupAdopts(t *testing.T) {
	t.Parallel()

	f := &fakeAppliance{refuseMsg: "Appliance has already been setup"}
	r := newResource(t, f)

	out, resp := create(t, r, inputs())
	if resp.Diagnostics.HasError() || !out.AlreadyInitialised.ValueBool() {
		t.Fatalf("concurrent setup was not adopted: %v", resp.Diagnostics)
	}
	if out.Id.IsNull() || out.BuildVersion.ValueString() != "9.1.0.2" || resp.Diagnostics.WarningsCount() != 1 {
		t.Fatal("adoption must retain state, refresh the version and warn")
	}
	if !out.PasswordWo.IsNull() || !out.LicenseKeyWo.IsNull() {
		t.Fatal("write-only credentials reached state")
	}
	f.mu.Lock()
	defer f.mu.Unlock()
	if len(f.inits) != 1 {
		t.Fatalf("init calls = %d, want 1", len(f.inits))
	}
}

func read(t *testing.T, r *Resource, set map[string]tftypes.Value) resource.ReadResponse {
	t.Helper()

	_, _, state := objectValue(t, set)
	resp := resource.ReadResponse{State: state}

	r.Read(context.Background(), resource.ReadRequest{State: state}, &resp)

	return resp
}

func stateValues(id string) map[string]tftypes.Value {
	v := inputs()
	v["id"] = str(id)
	v["appliance_url"] = str(id)
	v["already_initialised"] = tftypes.NewValue(tftypes.Bool, false)
	v["build_version"] = str("9.1.0")
	delete(v, "password_wo")
	delete(v, "license_key_wo")

	return v
}

func TestUnitReadKeepsASetUpApplianceAndRefreshesTheVersion(t *testing.T) {
	t.Parallel()

	f := &fakeAppliance{setUp: true}
	r := newResource(t, f)

	resp := read(t, r, stateValues("https://appliance.example"))
	if resp.Diagnostics.HasError() {
		t.Fatal(resp.Diagnostics)
	}

	if resp.State.Raw.IsNull() {
		t.Fatal("resource removed from state")
	}

	var out ApplianceSetupModel
	if d := resp.State.Get(context.Background(), &out); d.HasError() {
		t.Fatal(d)
	}

	if out.BuildVersion.ValueString() != "9.1.0.2" {
		t.Errorf("build_version = %q, want the whoami answer's", out.BuildVersion.ValueString())
	}
}

func TestUnitReadRemovesAReinstalledAppliance(t *testing.T) {
	t.Parallel()

	f := &fakeAppliance{}
	r := newResource(t, f)

	resp := read(t, r, stateValues("https://appliance.example"))
	if resp.Diagnostics.HasError() {
		t.Fatal(resp.Diagnostics)
	}

	if !resp.State.Raw.IsNull() {
		t.Fatal("an appliance that needs setup again should leave state")
	}
}

func TestUnitReadFailsWhenTheAPIDoesNotAnswer(t *testing.T) {
	t.Parallel()

	f := &fakeAppliance{notReady: 1 << 30}
	r := newResource(t, f)

	resp := read(t, r, stateValues("https://appliance.example"))
	if !resp.Diagnostics.HasError() {
		t.Fatal("want an error, not a removal, while the appliance does not answer")
	}

	if resp.State.Raw.IsNull() {
		t.Fatal("resource must stay in state")
	}
}

func TestUnitImportRequiresTheProviderURLAndASetUpAppliance(t *testing.T) {
	t.Parallel()

	f := &fakeAppliance{setUp: true}
	r := newResource(t, f)

	unauth, err := r.NewUnauthenticated()
	if err != nil {
		t.Fatal(err)
	}

	_, _, state := objectValue(t, nil)

	resp := resource.ImportStateResponse{State: state}
	r.ImportState(context.Background(), resource.ImportStateRequest{ID: "https://other.example"}, &resp)
	if !resp.Diagnostics.HasError() {
		t.Error("an import ID other than the provider url should fail")
	}

	resp = resource.ImportStateResponse{State: state}
	r.ImportState(context.Background(), resource.ImportStateRequest{ID: unauth.BaseURL + "/"}, &resp)
	if resp.Diagnostics.HasError() {
		t.Fatal(resp.Diagnostics)
	}

	var out ApplianceSetupModel
	if d := resp.State.Get(context.Background(), &out); d.HasError() {
		t.Fatal(d)
	}

	if out.Id.ValueString() != unauth.BaseURL || !out.AlreadyInitialised.ValueBool() || out.BuildVersion.ValueString() != "9.1.0.2" {
		t.Errorf("imported state = %+v", out)
	}

	f.mu.Lock()
	f.setUp = false
	f.mu.Unlock()

	resp = resource.ImportStateResponse{State: state}
	r.ImportState(context.Background(), resource.ImportStateRequest{ID: unauth.BaseURL}, &resp)
	if !resp.Diagnostics.HasError() {
		t.Error("importing an appliance that still needs setup should fail")
	}
}

func TestUnitInitResultMessage(t *testing.T) {
	t.Parallel()

	if got := (initResult{}).message(); got != "no message" {
		t.Errorf("empty = %q", got)
	}

	got := initResult{Msg: "refused", Errors: map[string]string{"b": "2", "a": "1"}}.message()
	if got != "refused (a: 1; b: 2)" {
		t.Errorf("got %q", got)
	}
}

func TestUnitExcerpt(t *testing.T) {
	t.Parallel()

	if got := excerpt(nil); got != "empty body" {
		t.Errorf("nil = %q", got)
	}

	if got := excerpt([]byte("<html>\n  <body>  Loading </body>\n</html>")); got != "<html> <body> Loading </body> </html>" {
		t.Errorf("html = %q", got)
	}

	long := strings.Repeat("x", 200)
	if got := excerpt([]byte(long)); len(got) != 123 || !strings.HasSuffix(got, "...") {
		t.Errorf("long = %q", got)
	}
}
