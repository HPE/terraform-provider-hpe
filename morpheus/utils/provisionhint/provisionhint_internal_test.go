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

		got := invalidNetwork(context.Background(), networks, req)

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

		got := invalidNetwork(context.Background(), networks, req)

		wantLine := `network 1 "Management" is in this cloud and active; it is not visible to the instance's group or tenant`
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

		got := invalidNetwork(context.Background(), networks, req)

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

		got := invalidNetwork(context.Background(), networks, req)

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
		got := invalidNetwork(context.Background(), networks, req)

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

		got := invalidNetwork(context.Background(), networks, req)

		for _, want := range []string{
			`network 2 "off" is inactive.`,
			`network 3 "loose" belongs to no resource pool; this instance targets resource pool 1702.`,
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

		if got := invalidNetwork(context.Background(), fakeNetworks{}, req); got != "" {
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
