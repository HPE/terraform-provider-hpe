// (C) Copyright 2026 Hewlett Packard Enterprise Development LP

package optionsource

import (
	"context"
	"encoding/json"
	"errors"
	"net/http"
	"net/http/httptest"
	"net/url"
	"strings"
	"testing"

	sdk "github.com/HPE/terraform-provider-hpe/internal/sdk/oapigen"
)

// Response bodies below are as returned by a 9.0.1 appliance, abbreviated.

const zonePoolsBody = `{"success":true,"data":[
  {"name":"Resource Pools","isGroup":true,"group":"pool"},
  {"enabled":true,"externalId":null,"group":"pool","id":1,"isDefault":false,
   "name":"Duck","providerType":"mvm","type":"mvm-cluster","value":"pool-1"},
  {"enabled":true,"externalId":null,"group":"pool","id":1702,"isDefault":false,
   "name":"jasontest","providerType":"mvm","type":"mvm-cluster","value":"pool-1702"}
]}`

const zoneNetworkOptionsBody = `{"success":true,"data":{
  "networks":[
    {"id":"network-1","name":"Management","dhcpServer":true,"allowStaticOverride":false,"pool":null},
    {"id":"network-91","name":"eoghan-demo-network","dhcpServer":true,"allowStaticOverride":false,
     "pool":{"id":3,"name":"p"}}],
  "networkGroups":[],
  "networkSubnets":[],
  "networkTypes":[{"id":4,"code":"kvmBridge","name":"Bridge"}],
  "hasNetworks":true,"maxNetworks":0,"enableNetworkTypeSelection":"on","supportsNetworkSelection":true
}}`

const gasketBody = `{"success":false,"msg":"Looks like the server threw a gasket"}`

func TestUnitParseDataArrayDropsHeaderRowsAndStringifies(t *testing.T) {
	t.Parallel()

	var env envelope
	if err := json.Unmarshal([]byte(zonePoolsBody), &env); err != nil {
		t.Fatal(err)
	}

	values, err := parseData("zonePools", env.Data, "")
	if err != nil {
		t.Fatal(err)
	}

	if len(values) != 2 {
		t.Fatalf("want 2 values (header row dropped), got %d: %+v", len(values), values)
	}

	first := values[0]
	if first.Name != "Duck" || first.Value != "pool-1" {
		t.Errorf("first = %q/%q, want Duck/pool-1", first.Name, first.Value)
	}

	for k, want := range map[string]string{
		"id": "1702", "providerType": "mvm", "type": "mvm-cluster", "enabled": "true", "isDefault": "false", "group": "pool",
	} {
		if got := values[1].Attributes[k]; got != want {
			t.Errorf("attributes[%q] = %q, want %q", k, got, want)
		}
	}

	if _, present := values[1].Attributes["externalId"]; present {
		t.Error("null attribute must be omitted, not rendered as a string")
	}

	if _, present := values[1].Attributes["name"]; present {
		t.Error("name must not be duplicated into attributes")
	}
}

func TestUnitParseDataObjectNeedsCollection(t *testing.T) {
	t.Parallel()

	var env envelope
	if err := json.Unmarshal([]byte(zoneNetworkOptionsBody), &env); err != nil {
		t.Fatal(err)
	}

	t.Run("no collection: error lists the lists", func(t *testing.T) {
		t.Parallel()

		_, err := parseData("zoneNetworkOptions", env.Data, "")

		var ce *CollectionError
		if !errors.As(err, &ce) {
			t.Fatalf("want CollectionError, got %v", err)
		}

		want := []string{"networkGroups", "networkSubnets", "networkTypes", "networks"}
		if strings.Join(ce.Available, ",") != strings.Join(want, ",") {
			t.Errorf("available = %v, want %v", ce.Available, want)
		}

		if !strings.Contains(err.Error(), "set collection to one of: networkGroups, networkSubnets, networkTypes, networks") {
			t.Errorf("message: %s", err)
		}
	})

	t.Run("wrong collection", func(t *testing.T) {
		t.Parallel()

		_, err := parseData("zoneNetworkOptions", env.Data, "hasNetworks")

		var ce *CollectionError
		if !errors.As(err, &ce) || ce.Collection != "hasNetworks" {
			t.Fatalf("want CollectionError for hasNetworks (a scalar), got %v", err)
		}
	})

	t.Run("networks: id becomes value and is kept in attributes", func(t *testing.T) {
		t.Parallel()

		values, err := parseData("zoneNetworkOptions", env.Data, "networks")
		if err != nil {
			t.Fatal(err)
		}

		if len(values) != 2 {
			t.Fatalf("want 2, got %d", len(values))
		}

		if values[0].Value != "network-1" || values[0].Name != "Management" {
			t.Errorf("first = %q/%q", values[0].Name, values[0].Value)
		}

		if values[0].Attributes["id"] != "network-1" {
			t.Errorf("attributes[id] = %q, want network-1", values[0].Attributes["id"])
		}

		if _, present := values[0].Attributes["pool"]; present {
			t.Error("null pool must be omitted")
		}

		if got := values[1].Attributes["pool"]; got != `{"id":3,"name":"p"}` {
			t.Errorf("nested object must be compact JSON, got %q", got)
		}
	})

	t.Run("empty collection is an empty list, not an error", func(t *testing.T) {
		t.Parallel()

		values, err := parseData("zoneNetworkOptions", env.Data, "networkGroups")
		if err != nil || len(values) != 0 {
			t.Errorf("want empty list, got %v / %v", values, err)
		}
	})
}

func TestUnitParseDataShapes(t *testing.T) {
	t.Parallel()

	t.Run("collection on an array is a ShapeError", func(t *testing.T) {
		t.Parallel()

		_, err := parseData("groups", json.RawMessage(`[{"name":"a","value":"1"}]`), "networks")

		var se *ShapeError
		if !errors.As(err, &se) {
			t.Fatalf("want ShapeError, got %v", err)
		}
	})

	t.Run("null data is empty", func(t *testing.T) {
		t.Parallel()

		values, err := parseData("x", json.RawMessage(`null`), "")
		if err != nil || len(values) != 0 {
			t.Errorf("got %v / %v", values, err)
		}
	})

	t.Run("bare strings are options", func(t *testing.T) {
		t.Parallel()

		values, err := parseData("x", json.RawMessage(`["a","b"]`), "")
		if err != nil || len(values) != 2 || values[1].Name != "b" || values[1].Value != "b" {
			t.Errorf("got %v / %v", values, err)
		}
	})

	t.Run("numeric value keeps its text", func(t *testing.T) {
		t.Parallel()

		values, err := parseData("x", json.RawMessage(`[{"name":"n","value":1702}]`), "")
		if err != nil || values[0].Value != "1702" {
			t.Errorf("got %v / %v", values, err)
		}
	})

	t.Run("isGroup as a string is a header row", func(t *testing.T) {
		t.Parallel()

		values, err := parseData("x", json.RawMessage(`[{"name":"h","isGroup":"true"},{"name":"n","value":"v"}]`), "")
		if err != nil || len(values) != 1 {
			t.Errorf("got %v / %v", values, err)
		}
	})
}

func TestUnitRequestPathAndQuery(t *testing.T) {
	t.Parallel()

	req := Request{Source: "networks", Parameters: map[string]string{"zoneId": "1", "phrase": "a b"}}

	if got := req.Path(); got != "/api/options/networks" {
		t.Errorf("path = %q", got)
	}

	twoSegment := Request{Source: "scvmmLibraryShares", ServiceType: "scvmm"}
	if got := twoSegment.Path(); got != "/api/options/scvmm/scvmmLibraryShares" {
		t.Errorf("two-segment path = %q", got)
	}

	// auto: both positions, sorted keys, brackets percent-encoded.
	if got := req.Query(); got != "config%5Bphrase%5D=a+b&config%5BzoneId%5D=1&phrase=a+b&zoneId=1" {
		t.Errorf("auto query = %q", got)
	}

	req.Style = StyleTopLevel
	if got := req.Query(); got != "phrase=a+b&zoneId=1" {
		t.Errorf("top_level query = %q", got)
	}

	req.Style = StyleConfig
	if got := req.Query(); got != "config%5Bphrase%5D=a+b&config%5BzoneId%5D=1" {
		t.Errorf("config query = %q", got)
	}

	if _, err := ParseStyle("both"); err == nil {
		t.Error("ParseStyle must reject an unknown style")
	}

	if s, err := ParseStyle(""); err != nil || s != StyleAuto {
		t.Errorf("empty style must be auto, got %v / %v", s, err)
	}
}

// newTestClient points an SDK client at a test server, with a bearer token in
// the default headers as the provider's client factory does.
func newTestClient(t *testing.T, srv *httptest.Server) *sdk.APIClient {
	t.Helper()

	cfg := sdk.NewConfiguration()
	cfg.Servers = sdk.ServerConfigurations{{URL: srv.URL}}
	cfg.HTTPClient = srv.Client()
	cfg.AddDefaultHeader("Authorization", "Bearer test-token")

	return sdk.NewAPIClient(cfg)
}

func TestUnitFetch(t *testing.T) {
	t.Parallel()

	var seen struct {
		path  string
		query url.Values
		auth  string
	}

	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		seen.path = r.URL.Path
		seen.query = r.URL.Query()
		seen.auth = r.Header.Get("Authorization")

		switch r.URL.Path {
		case "/api/options/zonePools":
			_, _ = w.Write([]byte(zonePoolsBody))
		case "/api/options/zoneNetworkOptions":
			_, _ = w.Write([]byte(zoneNetworkOptionsBody))
		case "/api/options/noSuchSource":
			w.WriteHeader(http.StatusInternalServerError)
			_, _ = w.Write([]byte(gasketBody))
		case "/api/options/quietFailure":
			_, _ = w.Write([]byte(`{"success":false,"msg":"no permission"}`))
		case "/api/options/html":
			w.WriteHeader(http.StatusBadGateway)
			_, _ = w.Write([]byte("<html>502</html>"))
		default:
			w.WriteHeader(http.StatusNotFound)
		}
	}))
	defer srv.Close()

	client := newTestClient(t, srv)
	ctx := context.Background()

	t.Run("sends auth, both parameter positions, and parses", func(t *testing.T) {
		res, err := Fetch(ctx, client, Request{
			Source: "zonePools", Parameters: map[string]string{"zoneId": "1", "layoutId": "1193"},
		})
		if err != nil {
			t.Fatal(err)
		}

		if seen.auth != "Bearer test-token" {
			t.Errorf("auth header = %q", seen.auth)
		}

		if seen.query.Get("zoneId") != "1" || seen.query.Get("config[zoneId]") != "1" ||
			seen.query.Get("layoutId") != "1193" || seen.query.Get("config[layoutId]") != "1193" {
			t.Errorf("query = %v", seen.query)
		}

		if len(res.Values) != 2 || res.Values[0].Value != "pool-1" {
			t.Errorf("values = %+v", res.Values)
		}

		if res.Path != "/api/options/zonePools" || !strings.Contains(res.Query, "zoneId=1") {
			t.Errorf("result path/query = %q / %q", res.Path, res.Query)
		}
	})

	t.Run("object response honours collection", func(t *testing.T) {
		res, err := Fetch(ctx, client, Request{Source: "zoneNetworkOptions", Collection: "networks"})
		if err != nil {
			t.Fatal(err)
		}

		if len(res.Values) != 2 || res.Values[1].Value != "network-91" {
			t.Errorf("values = %+v", res.Values)
		}
	})

	t.Run("gasket 500 is a ServerError that names the request", func(t *testing.T) {
		_, err := Fetch(ctx, client, Request{Source: "noSuchSource", Parameters: map[string]string{"zoneId": "1"}})

		var se *ServerError
		if !errors.As(err, &se) {
			t.Fatalf("want ServerError, got %v", err)
		}

		if se.Status != 500 || !se.Gasket() {
			t.Errorf("status/gasket = %d/%v", se.Status, se.Gasket())
		}

		msg := err.Error()
		for _, want := range []string{
			`option source "noSuchSource" returned HTTP 500`,
			"threw a gasket",
			"unknown source name and for a source whose parameters are missing or malformed alike",
			"GET /api/options/noSuchSource?",
			"zoneId=1",
		} {
			if !strings.Contains(msg, want) {
				t.Errorf("message missing %q:\n%s", want, msg)
			}
		}
	})

	t.Run("success:false with 200 is a ServerError", func(t *testing.T) {
		_, err := Fetch(ctx, client, Request{Source: "quietFailure"})

		var se *ServerError
		if !errors.As(err, &se) || se.Message != "no permission" || se.Gasket() {
			t.Fatalf("got %v", err)
		}
	})

	t.Run("non-JSON error body is excerpted", func(t *testing.T) {
		_, err := Fetch(ctx, client, Request{Source: "html"})

		var se *ServerError
		if !errors.As(err, &se) || se.Status != 502 || !strings.Contains(se.Message, "<html>") {
			t.Fatalf("got %v", err)
		}
	})

	t.Run("nil client and empty source are rejected", func(t *testing.T) {
		if _, err := Fetch(ctx, nil, Request{Source: "x"}); err == nil {
			t.Error("nil client accepted")
		}

		if _, err := Fetch(ctx, client, Request{Source: "  "}); err == nil {
			t.Error("empty source accepted")
		}
	})
}

func TestUnitKnownAndAdvice(t *testing.T) {
	t.Parallel()

	if !Known("zonePools") || !Known("scvmmLibraryShares") || Known("noSuchSource") {
		t.Error("catalogue membership wrong for known probes")
	}

	if got := ServiceTypesFor("scvmmLibraryShares"); len(got) != 1 || got[0] != "scvmm" {
		t.Errorf("ServiceTypesFor(scvmmLibraryShares) = %v", got)
	}

	if got := ServiceTypesFor("zonePools"); len(got) != 1 || got[0] != "" {
		t.Errorf("ServiceTypesFor(zonePools) = %v", got)
	}

	cases := map[string]struct {
		source, serviceType string
		wantContains        string
	}{
		"known, no service type needed": {"zonePools", "", ""},
		"unknown": {
			"noSuchSource", "",
			`is not in the provider's catalogue of sources known for Morpheus ` + KnownSourcesRelease,
		},
		"needs service type, none given": {
			"scvmmLibraryShares", "",
			`is defined for service_type "scvmm" and is not reachable without it`,
		},
		"needs service type, right given": {"scvmmLibraryShares", "scvmm", ""},
		"needs service type, wrong given": {
			"scvmmLibraryShares", "vmware",
			`is not defined for service_type "vmware"`,
		},
	}

	for name, tc := range cases {
		t.Run(name, func(t *testing.T) {
			t.Parallel()

			got := Advice(tc.source, tc.serviceType)
			if tc.wantContains == "" && got != "" {
				t.Errorf("want no advice, got %q", got)
			}

			if tc.wantContains != "" && !strings.Contains(got, tc.wantContains) {
				t.Errorf("advice %q does not contain %q", got, tc.wantContains)
			}
		})
	}
}

func TestUnitMissingContextHint(t *testing.T) {
	t.Parallel()

	got := MissingContextHint("zonePools", map[string]string{"zoneId": "1"})
	if !strings.Contains(got, "set layout_id or provision_type_code") {
		t.Errorf("zonePools without layout: %q", got)
	}

	if got := MissingContextHint("zonePools", map[string]string{"zoneId": "1", "layoutId": "5"}); got != "" {
		t.Errorf("zonePools with layout must have no hint, got %q", got)
	}

	withProvisionType := map[string]string{"zoneId": "1", "provisionType": "kvm"}
	if got := MissingContextHint("zoneNetworkOptions", withProvisionType); got != "" {
		t.Errorf("zoneNetworkOptions with provisionType must have no hint, got %q", got)
	}

	if got := MissingContextHint("groups", nil); got != "" {
		t.Errorf("unlisted source must have no hint, got %q", got)
	}
}

// TestUnitTypedKeysCoverEveryTypedField guards the three places a typed
// attribute must be declared: the Typed struct, typedKeys, and the entries list
// in MergeParameters. A field set on Typed must reach the request.
func TestUnitTypedKeysCoverEveryTypedField(t *testing.T) {
	t.Parallel()

	i := func(v int64) *int64 { return &v }
	s := func(v string) *string { return &v }

	all := Typed{
		CloudID: i(1), Phrase: s("p"), ParentID: i(3), CredentialID: i(4), GroupID: i(5),
		LayoutID: i(6), ProvisionTypeCode: s("c"), ProvisionTypeID: i(7), PoolID: s("pool-8"), PlanID: i(9),
		ServerGroupID: i(10), ProviderID: i(11),
	}

	got, err := MergeParameters(nil, all)
	if err != nil {
		t.Fatal(err)
	}

	for attribute, keys := range typedKeys {
		for _, k := range keys {
			if _, ok := got[k]; !ok {
				t.Errorf("typed attribute %s: key %q was not set although the field is populated", attribute, k)
			}
		}
	}

	// Every key that was set must be explained by typedKeys.
	known := map[string]bool{}
	for _, keys := range typedKeys {
		for _, k := range keys {
			known[k] = true
		}
	}

	for k := range got {
		if !known[k] {
			t.Errorf("key %q set by MergeParameters is not declared in typedKeys", k)
		}
	}
}

func TestUnitMergeParameters(t *testing.T) {
	t.Parallel()

	i := func(v int64) *int64 { return &v }
	s := func(v string) *string { return &v }

	t.Run("typed attributes fan out to their keys", func(t *testing.T) {
		t.Parallel()

		got, err := MergeParameters(map[string]string{"phrase": "web"}, Typed{
			CloudID: i(1), GroupID: i(265), LayoutID: i(1193), ProvisionTypeCode: s("kvm"), ProvisionTypeID: i(3),
			PoolID: s("pool-1"), CredentialID: i(12), ParentID: i(7), PlanID: i(19), ServerGroupID: i(1), ProviderID: i(4),
		})
		if err != nil {
			t.Fatal(err)
		}

		want := map[string]string{
			"phrase": "web", "zoneId": "1", "cloudId": "1", "groupId": "265", "siteId": "265", "layoutId": "1193",
			"provisionTypeCode": "kvm", "provisionType": "kvm", "provisionTypeId": "3", "poolId": "pool-1",
			"credential.type": "12", "parentId": "7", "parentid": "7", "planId": "19", "serverGroupId": "1",
			"providerId": "4",
		}

		if len(got) != len(want) {
			t.Fatalf("got %v, want %v", got, want)
		}

		for k, v := range want {
			if got[k] != v {
				t.Errorf("%s = %q, want %q", k, got[k], v)
			}
		}
	})

	t.Run("a key set both ways is an error", func(t *testing.T) {
		t.Parallel()

		_, err := MergeParameters(map[string]string{"zoneId": "2"}, Typed{CloudID: i(1)})
		if err == nil || !strings.Contains(err.Error(), `parameter "zoneId" is set both by cloud_id and in parameters`) {
			t.Errorf("got %v", err)
		}
	})

	t.Run("describe sorts keys", func(t *testing.T) {
		t.Parallel()

		if got := DescribeParameters(map[string]string{"zoneId": "1", "layoutId": "5"}); got != "layoutId=5, zoneId=1" {
			t.Errorf("got %q", got)
		}

		if got := DescribeParameters(nil); got != "(none)" {
			t.Errorf("got %q", got)
		}
	})
}

// The server replaces accountId (and a few internal keys) with its own values
// on every option-source call; the data source must be able to warn.
func TestUnitServerSetParameters(t *testing.T) {
	t.Parallel()

	got := ServerSetParameters(map[string]string{"zoneId": "1", "accountId": "451", "namespace": "x"})
	if len(got) != 2 || got[0] != "accountId" || got[1] != "namespace" {
		t.Fatalf("ServerSetParameters = %v, want [accountId namespace]", got)
	}

	if got := ServerSetParameters(map[string]string{"zoneId": "1"}); len(got) != 0 {
		t.Errorf("ServerSetParameters = %v, want none", got)
	}

	if !strings.Contains(ServerSetAdvice("accountId"), "authenticate as a user of that tenant") {
		t.Errorf("accountId advice must say how to read as another tenant: %q", ServerSetAdvice("accountId"))
	}

	if !strings.Contains(ServerSetAdvice("namespace"), "namespace") {
		t.Errorf("generic advice must name the key: %q", ServerSetAdvice("namespace"))
	}
}
