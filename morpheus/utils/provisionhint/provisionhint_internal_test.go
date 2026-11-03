// (C) Copyright 2026 Hewlett Packard Enterprise Development LP

package provisionhint

import (
	"context"
	"errors"
	"io"
	"net/http"
	"strings"
	"testing"

	sdk "github.com/HPE/terraform-provider-hpe/internal/sdk/oapigen"
)

// The body Morpheus returns for the case that motivated this package.
const invalidNetworkBody = `{"errors":{"networkInterfaces":"Invalid network"},"msg":"Failed to create instance","success":false}`

func TestUnitIsInvalidNetwork(t *testing.T) {
	t.Parallel()

	cases := map[string]struct {
		body string
		want bool
	}{
		"the real rejection":            {invalidNetworkBody, true},
		"other networkInterfaces error": {`{"errors":{"networkInterfaces":"You must enter an ip address"}}`, false},
		"other field":                   {`{"errors":{"plan":"Invalid network"}}`, false},
		"no errors object":              {`{"msg":"Invalid network","success":false}`, false},
		"errors not a string":           {`{"errors":{"networkInterfaces":["Invalid network"]}}`, false},
		"not json":                      {`502 Bad Gateway`, false},
		"empty":                         {``, false},
	}

	for name, tc := range cases {
		t.Run(name, func(t *testing.T) {
			t.Parallel()

			if got := isInvalidNetwork([]byte(tc.body)); got != tc.want {
				t.Errorf("isInvalidNetwork(%q) = %v, want %v", tc.body, got, tc.want)
			}
		})
	}
}

func TestUnitResourcePoolIDFromRequestConfigs(t *testing.T) {
	t.Parallel()

	str := func(s string) *string { return &s }

	cases := map[string]struct {
		config any
		want   int64
	}{
		"nil": {nil, 0},
		"instance request, typed HVM config, pool- prefix": {
			sdk.AddInstanceRequestConfig{HVMInstanceConfiguration: &sdk.HVMInstanceConfiguration{ResourcePoolId: str("pool-1702")}},
			1702,
		},
		"instance request, typed HVM config, bare": {
			sdk.AddInstanceRequestConfig{HVMInstanceConfiguration: &sdk.HVMInstanceConfiguration{ResourcePoolId: str("7")}},
			7,
		},
		"clone request, pointer to typed config": {
			&sdk.CloneInstanceRequestConfig{HVMInstanceConfiguration1: &sdk.HVMInstanceConfiguration1{ResourcePoolId: str("pool-2")}},
			2,
		},
		"clone request, nil pointer": {(*sdk.CloneInstanceRequestConfig)(nil), 0},
		"dynamic config map, string": {map[string]any{"resourcePoolId": "pool-5", "poolProviderType": "mvm"}, 5},
		"dynamic config map, number": {map[string]any{"resourcePoolId": 9}, 9},
		"dynamic config map, absent": {map[string]any{"imageId": 12}, 0},
		"pool group is not a pool":   {map[string]any{"resourcePoolId": "poolGroup-3"}, 0},
	}

	for name, tc := range cases {
		t.Run(name, func(t *testing.T) {
			t.Parallel()

			if got := resourcePoolID(tc.config); got != tc.want {
				t.Errorf("resourcePoolID = %d, want %d", got, tc.want)
			}
		})
	}
}

func TestUnitNetworkIDsFromBothRequestTypes(t *testing.T) {
	t.Parallel()

	// The instance request uses one generated type, the clone request another;
	// both must be read.
	instanceIfaces := []sdk.InstancesNetworkInterfaces2{
		{Network: sdk.InstancesNetworkInterfaces1Network{Id: "1"}},
		{Network: sdk.InstancesNetworkInterfaces1Network{Id: "network-79"}},
		{Network: sdk.InstancesNetworkInterfaces1Network{Id: "networkGroup-3"}}, // skipped
		{Network: sdk.InstancesNetworkInterfaces1Network{Id: "subnet-4"}},       // skipped
		{Network: sdk.InstancesNetworkInterfaces1Network{Id: "1"}},              // duplicate
	}

	cloneIfaces := []sdk.InstancesNetworkInterfaces3{
		{Network: sdk.InstancesNetworkInterfaces3Network{Id: "network-1"}},
		{Network: sdk.InstancesNetworkInterfaces3Network{Id: "79"}},
	}

	for name, tc := range map[string]struct {
		in   any
		want []int64
	}{
		"instance request type": {instanceIfaces, []int64{1, 79}},
		"clone request type":    {cloneIfaces, []int64{1, 79}},
		"nil":                   {nil, nil},
		"not a list":            {"garbage", nil},
	} {
		t.Run(name, func(t *testing.T) {
			t.Parallel()

			got := networkIDs(tc.in)
			if len(got) != len(tc.want) {
				t.Fatalf("networkIDs = %v, want %v", got, tc.want)
			}

			for i := range got {
				if got[i] != tc.want[i] {
					t.Errorf("networkIDs = %v, want %v", got, tc.want)
				}
			}
		})
	}
}

// fakeNetworks answers get() from a table, standing in for the SDK.
type fakeNetworks map[int64]networkFacts

func (f fakeNetworks) get(_ context.Context, id int64) networkFacts {
	if facts, ok := f[id]; ok {
		return facts
	}

	return networkFacts{ID: id, LookupFailed: true, LookupStatus: http.StatusNotFound}
}

func active(b bool) *bool { return &b }

// only wraps a network table as lookups with no selectable-network or
// instance lookups — the pre-option-source behaviour.
func only(networks fakeNetworks) lookups { return lookups{networks: networks} }

// fakeSelectable answers selectable() with a fixed list, or an error, and
// records the scope it was asked about.
type fakeSelectable struct {
	list  []selectableNetwork
	err   error
	asked *scope
}

func (f *fakeSelectable) selectable(_ context.Context, s scope) ([]selectableNetwork, error) {
	f.asked = &s

	return f.list, f.err
}

// fakeInstance answers context() for one instance id.
type fakeInstance struct {
	id                         int64
	cloudID, groupID, layoutID int64
}

func (f fakeInstance) context(_ context.Context, id int64) (int64, int64, int64) {
	if id != f.id {
		return 0, 0, 0
	}

	return f.cloudID, f.groupID, f.layoutID
}

func TestUnitInvalidNetworkMessage(t *testing.T) {
	t.Parallel()

	duck := poolRef{ID: 1, Name: "Duck"}

	t.Run("the customer's case: pool mismatch", func(t *testing.T) {
		t.Parallel()

		networks := fakeNetworks{1: {ID: 1, Name: "Management", CloudID: 1, Active: active(true), Pools: []poolRef{duck}}}
		req := Request{
			CloudID:    1,
			Config:     map[string]any{"resourcePoolId": "pool-1702"},
			Interfaces: []sdk.InstancesNetworkInterfaces2{{Network: sdk.InstancesNetworkInterfaces1Network{Id: "1"}}},
		}

		got := invalidNetwork(context.Background(), only(networks), req)

		want := "Morpheus rejected a network interface because its network is not selectable for this " +
			"instance's cloud, group and resource pool:\n" +
			"\n  - network 1 \"Management\" belongs to resource pool 1 \"Duck\"; this instance targets resource pool 1702." +
			"\n\nNetworks discovered on an HVM cluster are bound to that cluster's resource pool. " +
			"Provision into the cluster's pool (hpe_morpheus_cluster: permissions.resource_pool.id), " +
			"or choose a network that belongs to resource pool 1702."

		if got != want {
			t.Errorf("message\n--- want ---\n%s\n--- got ---\n%s", want, got)
		}
	})

	t.Run("pool matches: falls through to visibility", func(t *testing.T) {
		t.Parallel()

		networks := fakeNetworks{1: {ID: 1, Name: "Management", CloudID: 1, Active: active(true), Pools: []poolRef{duck}}}
		req := Request{
			CloudID:    1,
			Config:     map[string]any{"resourcePoolId": "pool-1"},
			Interfaces: []map[string]any{{"network": map[string]any{"id": "1"}}},
		}

		got := invalidNetwork(context.Background(), only(networks), req)

		wantLine := `network 1 "Management" is in this cloud and active; it is not of a type this layout's provision type ` +
			`can use, or it is not visible to the instance's group or tenant`
		if !strings.Contains(got, wantLine) {
			t.Errorf("expected the visibility line\n--- got ---\n%s", got)
		}

		if strings.Contains(got, "Networks discovered on an HVM cluster") {
			t.Errorf("closing paragraph must only appear for a pool mismatch\n--- got ---\n%s", got)
		}
	})

	t.Run("no pool requested: never blames the pool", func(t *testing.T) {
		t.Parallel()

		networks := fakeNetworks{1: {ID: 1, Name: "Management", CloudID: 1, Active: active(true), Pools: []poolRef{duck}}}
		req := Request{CloudID: 1, Interfaces: []map[string]any{{"network": map[string]any{"id": "1"}}}}

		got := invalidNetwork(context.Background(), only(networks), req)

		if strings.Contains(got, "targets resource pool") || strings.Contains(got, "HVM cluster") {
			t.Errorf("no pool in the request, so the pool must not be mentioned\n--- got ---\n%s", got)
		}
	})

	t.Run("cloud mismatch wins over pool", func(t *testing.T) {
		t.Parallel()

		networks := fakeNetworks{5: {ID: 5, Name: "other", CloudID: 3, Active: active(true), Pools: []poolRef{{ID: 9}}}}
		req := Request{
			CloudID:    1,
			Config:     map[string]any{"resourcePoolId": "pool-2"},
			Interfaces: []map[string]any{{"network": map[string]any{"id": "5"}}},
		}

		got := invalidNetwork(context.Background(), only(networks), req)

		if !strings.Contains(got, `network 5 "other" is in cloud 3; this instance targets cloud 1.`) {
			t.Errorf("expected the cloud line\n--- got ---\n%s", got)
		}
	})

	t.Run("clone: no cloud on the request, cloud check skipped", func(t *testing.T) {
		t.Parallel()

		networks := fakeNetworks{5: {ID: 5, Name: "other", CloudID: 3, Active: active(true), Pools: []poolRef{{ID: 9, Name: "nine"}}}}
		req := Request{Config: &sdk.CloneInstanceRequestConfig{}, Interfaces: []map[string]any{{"network": map[string]any{"id": "5"}}}}
		// Config marshals with no resourcePoolId, so no pool either: the only
		// remaining explanation is visibility.
		got := invalidNetwork(context.Background(), only(networks), req)

		if strings.Contains(got, "targets cloud") {
			t.Errorf("cloud must not be compared when the request has none\n--- got ---\n%s", got)
		}

		if !strings.Contains(got, "not visible to the instance's group or tenant") {
			t.Errorf("expected the visibility line\n--- got ---\n%s", got)
		}
	})

	t.Run("inactive, no pool, several pools, unreadable", func(t *testing.T) {
		t.Parallel()

		networks := fakeNetworks{
			2: {ID: 2, Name: "off", CloudID: 1, Active: active(false), Pools: []poolRef{duck}},
			3: {ID: 3, Name: "loose", CloudID: 1, Active: active(true)},
			4: {ID: 4, Name: "shared", CloudID: 1, Active: active(true), Pools: []poolRef{duck, {ID: 8, Name: "eight"}}},
			6: {ID: 6, LookupFailed: true, LookupStatus: http.StatusForbidden},
		}
		req := Request{
			CloudID: 1,
			Config:  map[string]any{"resourcePoolId": "pool-1702"},
			Interfaces: []map[string]any{
				{"network": map[string]any{"id": "2"}},
				{"network": map[string]any{"id": "3"}},
				{"network": map[string]any{"id": "4"}},
				{"network": map[string]any{"id": "6"}},
				{"network": map[string]any{"id": "7"}}, // not in the table: 404
			},
		}

		got := invalidNetwork(context.Background(), only(networks), req)

		for _, want := range []string{
			`network 2 "off" is inactive.`,
			`network 3 "loose" belongs to no resource pool, which does not exclude it from resource pool 1702; it is not of a type`,
			`network 4 "shared" belongs to resource pools 1 "Duck", 8 "eight"; this instance targets resource pool 1702.`,
			`network 6 could not be read (HTTP 403).`,
			`network 7 does not exist, or is not visible to this user.`,
			"or choose a network that belongs to resource pool 1702.",
		} {
			if !strings.Contains(got, want) {
				t.Errorf("missing %q\n--- got ---\n%s", want, got)
			}
		}
	})

	t.Run("no network interfaces: nothing to say", func(t *testing.T) {
		t.Parallel()

		req := Request{
			CloudID:    1,
			Config:     map[string]any{"resourcePoolId": "pool-1"},
			Interfaces: []map[string]any{{"network": map[string]any{"id": "networkGroup-3"}}},
		}

		if got := invalidNetwork(context.Background(), only(fakeNetworks{}), req); got != "" {
			t.Errorf("expected no explanation for a group-only request, got\n%s", got)
		}
	})
}

func TestUnitInvalidNetworkIgnoresOtherFailures(t *testing.T) {
	t.Parallel()

	// A client is required by the signature but must never be called when the
	// body is not the rejection we explain; a nil-configured client proves it.
	client := sdk.NewAPIClient(sdk.NewConfiguration())
	body := []byte(`{"errors":{"plan":"Plan not found"},"success":false}`)

	req := Request{CloudID: 1, Interfaces: []map[string]any{{"network": map[string]any{"id": "1"}}}}

	if got := InvalidNetwork(context.Background(), client, body, req); got != "" {
		t.Errorf("expected \"\" for a non-network failure, got %q", got)
	}
}

func TestUnitErrorBody(t *testing.T) {
	t.Parallel()

	t.Run("from the response, leaving it readable", func(t *testing.T) {
		t.Parallel()

		resp := &http.Response{StatusCode: 400, Body: io.NopCloser(strings.NewReader(invalidNetworkBody))}

		got := ErrorBody(errors.New("400 Bad Request"), resp)
		if string(got) != invalidNetworkBody {
			t.Fatalf("ErrorBody = %q", got)
		}

		again, err := io.ReadAll(resp.Body)
		if err != nil || string(again) != invalidNetworkBody {
			t.Errorf("response body must still be readable afterwards; got %q, %v", again, err)
		}
	})

	t.Run("nil response", func(t *testing.T) {
		t.Parallel()

		if got := ErrorBody(errors.New("dial tcp: connection refused"), nil); got != nil {
			t.Errorf("expected nil, got %q", got)
		}
	})
}

func TestUnitFactsFromNetwork(t *testing.T) {
	t.Parallel()

	id := int64(1)
	name := "Management"
	zoneID := int64(1)
	poolID := int64(1)
	poolName := "Duck"

	n := &sdk.GetNetwork200ResponseNetwork{
		Id:       &id,
		Name:     &name,
		Active:   active(true),
		Zone:     &sdk.GetNetwork200ResponseNetworkZone{Id: &zoneID},
		ZonePool: &sdk.GetNetwork200ResponseNetworkZonePool{Id: &poolID, Name: &poolName},
		AdditionalProperties: map[string]any{
			// As decoded from JSON: numbers are float64. Pool 1 repeats zonePool
			// and must not be listed twice; 3 is new.
			"assignedZonePools": []any{
				map[string]any{"id": float64(3), "name": "three"},
				map[string]any{"id": float64(1), "name": "Duck"},
			},
		},
	}

	got := factsFromNetwork(1, n)

	if got.Name != "Management" || got.CloudID != 1 || got.Active == nil || !*got.Active {
		t.Errorf("basic facts wrong: %+v", got)
	}

	if len(got.Pools) != 2 || got.Pools[0] != (poolRef{ID: 1, Name: "Duck"}) || got.Pools[1] != (poolRef{ID: 3, Name: "three"}) {
		t.Errorf("pools = %+v, want [1 Duck, 3 three]", got.Pools)
	}
}

func TestUnitInvalidNetworkReasonMatchesText(t *testing.T) {
	t.Parallel()

	client := sdk.NewAPIClient(sdk.NewConfiguration())
	req := Request{Interfaces: []map[string]any{{"network": map[string]any{"id": "1"}}}}

	// Text that is not the rejection: nothing is looked up, nothing is said.
	if got := InvalidNetworkReason(context.Background(), client, "Failed to allocate storage", req); got != "" {
		t.Errorf("expected \"\" for an unrelated failure, got %q", got)
	}

	// Nil client: never explain, never panic.
	if got := InvalidNetworkReason(context.Background(), nil, "Invalid network", req); got != "" {
		t.Errorf("expected \"\" with a nil client, got %q", got)
	}
}

// The selectable-network paragraph: what Morpheus would have accepted, from
// the same query the validation ran.
func TestUnitInvalidNetworkSelectableParagraph(t *testing.T) {
	t.Parallel()

	duck := poolRef{ID: 1, Name: "Duck"}
	management := fakeNetworks{1: {ID: 1, Name: "Management", CloudID: 1, Active: active(true), Pools: []poolRef{duck}}}
	wrongPool := Request{
		CloudID:    1,
		GroupID:    265,
		LayoutID:   1193,
		Config:     map[string]any{"resourcePoolId": "pool-1702"},
		Interfaces: []map[string]any{{"network": map[string]any{"id": "1"}}},
	}

	t.Run("the incident: a pool with nothing selectable is said outright", func(t *testing.T) {
		t.Parallel()

		lister := &fakeSelectable{list: []selectableNetwork{}}
		got := invalidNetwork(context.Background(), lookups{networks: management, selectable: lister}, wrongPool)

		want := "Morpheus rejected a network interface because its network is not selectable for this " +
			"instance's cloud, group and resource pool:\n" +
			"\n  - network 1 \"Management\" belongs to resource pool 1 \"Duck\"; this instance targets resource pool 1702." +
			"\n\nNetworks discovered on an HVM cluster are bound to that cluster's resource pool. " +
			"Provision into the cluster's pool (hpe_morpheus_cluster: permissions.resource_pool.id)." +
			"\n\nNo network at all is selectable for this instance (cloud 1, group 265, layout 1193, resource pool 1702), " +
			"so no network interface can be provisioned with this combination; change the resource pool, group or " +
			"layout. Data source hpe_morpheus_option_source_values, source \"zoneNetworkOptions\", collection " +
			"\"networks\", lists what is selectable for a given cloud_id, group_id, layout_id and pool_id."

		if got != want {
			t.Errorf("message\n--- want ---\n%s\n--- got ---\n%s", want, got)
		}

		if lister.asked == nil || *lister.asked != (scope{CloudID: 1, GroupID: 265, LayoutID: 1193, PoolID: 1702}) {
			t.Errorf("the lister must be asked about the request's exact scope; asked %+v", lister.asked)
		}
	})

	t.Run("networks that would be accepted are listed", func(t *testing.T) {
		t.Parallel()

		lister := &fakeSelectable{list: []selectableNetwork{{ID: 5, Name: "five"}, {ID: 9}}}
		got := invalidNetwork(context.Background(), lookups{networks: management, selectable: lister}, wrongPool)

		for _, want := range []string{
			"Provision into the cluster's pool (hpe_morpheus_cluster: permissions.resource_pool.id), " +
				"or choose one of the networks listed below.",
			"Networks Morpheus would accept for this instance (cloud 1, group 265, layout 1193, resource pool 1702): " +
				"network 5 \"five\", network 9. Data source hpe_morpheus_option_source_values",
		} {
			if !strings.Contains(got, want) {
				t.Errorf("missing %q\n--- got ---\n%s", want, got)
			}
		}

		if strings.Contains(got, "or choose a network that belongs to resource pool 1702") {
			t.Errorf("the unknown-list tail must not appear when the list is known\n--- got ---\n%s", got)
		}
	})

	t.Run("a long list is capped", func(t *testing.T) {
		t.Parallel()

		list := make([]selectableNetwork, 0, 40)
		for i := int64(1); i <= 40; i++ {
			list = append(list, selectableNetwork{ID: i})
		}

		got := invalidNetwork(context.Background(), lookups{networks: management, selectable: &fakeSelectable{list: list}}, wrongPool)

		if !strings.Contains(got, "network 15, and 25 more.") || strings.Contains(got, "network 16,") {
			t.Errorf("expected 15 listed then \"and 25 more\"\n--- got ---\n%s", got)
		}
	})

	t.Run("lookup failure: paragraph omitted, the rest stands", func(t *testing.T) {
		t.Parallel()

		lister := &fakeSelectable{err: errors.New("option source: HTTP 500")}
		got := invalidNetwork(context.Background(), lookups{networks: management, selectable: lister}, wrongPool)

		if strings.Contains(got, "Networks Morpheus would accept") || strings.Contains(got, "No network at all") {
			t.Errorf("a failed lookup must not produce a selectable paragraph\n--- got ---\n%s", got)
		}

		if !strings.Contains(got, "or choose a network that belongs to resource pool 1702.") {
			t.Errorf("with the list unknown the original tail stands\n--- got ---\n%s", got)
		}
	})

	t.Run("no layout: not asked, since the source would answer nothing", func(t *testing.T) {
		t.Parallel()

		lister := &fakeSelectable{list: []selectableNetwork{}}
		req := wrongPool
		req.LayoutID = 0

		got := invalidNetwork(context.Background(), lookups{networks: management, selectable: lister}, req)

		if lister.asked != nil {
			t.Errorf("must not ask without a layout; asked %+v", *lister.asked)
		}

		if strings.Contains(got, "No network at all") {
			t.Errorf("an unasked question must not be reported as an empty answer\n--- got ---\n%s", got)
		}
	})

	t.Run("no cloud: not asked, since the source would answer for every cloud", func(t *testing.T) {
		t.Parallel()

		lister := &fakeSelectable{list: []selectableNetwork{{ID: 5}}}
		req := wrongPool
		req.CloudID = 0

		invalidNetwork(context.Background(), lookups{networks: management, selectable: lister}, req)

		if lister.asked != nil {
			t.Errorf("must not ask without a cloud; asked %+v", *lister.asked)
		}
	})

	t.Run("pool matches, list known: no HVM paragraph, list still given", func(t *testing.T) {
		t.Parallel()

		lister := &fakeSelectable{list: []selectableNetwork{{ID: 91, Name: "other"}}}
		req := wrongPool
		req.Config = map[string]any{"resourcePoolId": "pool-1"}

		got := invalidNetwork(context.Background(), lookups{networks: management, selectable: lister}, req)

		if strings.Contains(got, "HVM cluster") {
			t.Errorf("no pool mismatch, so no HVM paragraph\n--- got ---\n%s", got)
		}

		if !strings.Contains(got, "(cloud 1, group 265, layout 1193, resource pool 1): network 91 \"other\".") {
			t.Errorf("expected the list for the matching pool\n--- got ---\n%s", got)
		}
	})
}

// A pool-less network passes the server's pool filter, so the line must not
// claim the pool excluded it — and must not trigger the HVM paragraph.
func TestUnitPoollessNetworkIsNotAPoolMismatch(t *testing.T) {
	t.Parallel()

	networks := fakeNetworks{139: {ID: 139, Name: "docker bridge", CloudID: 1, Active: active(true)}}
	req := Request{
		CloudID:    1,
		Config:     map[string]any{"resourcePoolId": "pool-1"},
		Interfaces: []map[string]any{{"network": map[string]any{"id": "139"}}},
	}

	got := invalidNetwork(context.Background(), only(networks), req)

	want := `network 139 "docker bridge" belongs to no resource pool, which does not exclude it from resource pool 1; ` +
		`it is not of a type this layout's provision type can use, or it is not visible to the instance's group or ` +
		`tenant (check the network's type, group access and tenant permissions).`
	if !strings.Contains(got, want) {
		t.Errorf("missing the corrected line\n--- got ---\n%s", got)
	}

	if strings.Contains(got, "HVM cluster") || strings.Contains(got, "targets resource pool") {
		t.Errorf("a pool-less network is not a pool mismatch\n--- got ---\n%s", got)
	}
}

// A clone inherits the cloud, group and layout it did not set from its source.
func TestUnitCloneInheritsScopeFromSource(t *testing.T) {
	t.Parallel()

	duck := poolRef{ID: 1, Name: "Duck"}
	networks := fakeNetworks{1: {ID: 1, Name: "Management", CloudID: 1, Active: active(true), Pools: []poolRef{duck}}}
	source := fakeInstance{id: 42, cloudID: 1, groupID: 265, layoutID: 1193}

	t.Run("nothing set on the clone: all three inherited", func(t *testing.T) {
		t.Parallel()

		lister := &fakeSelectable{list: []selectableNetwork{}}
		req := Request{
			SourceInstanceID: 42,
			Config:           map[string]any{"resourcePoolId": "pool-1702"},
			Interfaces:       []map[string]any{{"network": map[string]any{"id": "1"}}},
		}

		got := invalidNetwork(context.Background(), lookups{networks: networks, selectable: lister, instances: source}, req)

		if lister.asked == nil || *lister.asked != (scope{CloudID: 1, GroupID: 265, LayoutID: 1193, PoolID: 1702}) {
			t.Fatalf("expected the source's scope; asked %+v", lister.asked)
		}

		if !strings.Contains(got, "(cloud 1, group 265, layout 1193, resource pool 1702)") {
			t.Errorf("expected the inherited scope in the message\n--- got ---\n%s", got)
		}
	})

	t.Run("clone names its own cloud and group: those win, layout inherited", func(t *testing.T) {
		t.Parallel()

		lister := &fakeSelectable{list: []selectableNetwork{}}
		req := Request{
			SourceInstanceID: 42,
			CloudID:          7,
			GroupID:          8,
			Config:           map[string]any{"resourcePoolId": "pool-1702"},
			Interfaces:       []map[string]any{{"network": map[string]any{"id": "1"}}},
		}

		invalidNetwork(context.Background(), lookups{networks: networks, selectable: lister, instances: source}, req)

		if lister.asked == nil || *lister.asked != (scope{CloudID: 7, GroupID: 8, LayoutID: 1193, PoolID: 1702}) {
			t.Errorf("expected cloud 7, group 8 from the clone and layout 1193 from the source; asked %+v", lister.asked)
		}
	})

	t.Run("source unreadable: no scope, so not asked", func(t *testing.T) {
		t.Parallel()

		lister := &fakeSelectable{list: []selectableNetwork{}}
		req := Request{
			SourceInstanceID: 99, // not the fake's instance
			Config:           map[string]any{"resourcePoolId": "pool-1702"},
			Interfaces:       []map[string]any{{"network": map[string]any{"id": "1"}}},
		}

		got := invalidNetwork(context.Background(), lookups{networks: networks, selectable: lister, instances: source}, req)

		if lister.asked != nil {
			t.Errorf("must not ask with no cloud or layout; asked %+v", *lister.asked)
		}

		if !strings.Contains(got, "or choose a network that belongs to resource pool 1702.") {
			t.Errorf("the unknown-list tail must stand\n--- got ---\n%s", got)
		}
	})
}

func TestUnitScopeDescribe(t *testing.T) {
	t.Parallel()

	for name, tc := range map[string]struct {
		s    scope
		want string
	}{
		"all": {
			scope{CloudID: 1, GroupID: 265, LayoutID: 1193, PoolID: 1702},
			"cloud 1, group 265, layout 1193, resource pool 1702",
		},
		"no group":    {scope{CloudID: 1, LayoutID: 1193, PoolID: 1}, "cloud 1, layout 1193, resource pool 1"},
		"cloud alone": {scope{CloudID: 3}, "cloud 3"},
		"nothing":     {scope{}, ""},
	} {
		t.Run(name, func(t *testing.T) {
			t.Parallel()

			if got := tc.s.describe(); got != tc.want {
				t.Errorf("describe() = %q, want %q", got, tc.want)
			}
		})
	}
}
