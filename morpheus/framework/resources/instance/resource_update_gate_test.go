// (C) Copyright 2026 Hewlett Packard Enterprise Development LP

package instance_test

import (
	"sort"
	"strings"
	"sync"
	"testing"

	"github.com/hashicorp/terraform-plugin-testing/helper/acctest"
	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
	"github.com/hashicorp/terraform-plugin-testing/plancheck"
	"github.com/hashicorp/terraform-plugin-testing/terraform"

	"github.com/HPE/terraform-provider-hpe/morpheus/framework/resources/instance"
	"github.com/HPE/terraform-provider-hpe/morpheus/testhelpers"
	"github.com/HPE/terraform-provider-hpe/morpheus/testhelpers/capabilities"
	"github.com/HPE/terraform-provider-hpe/provider/adapter"
)

// callRecorder captures the appliance calls the instance update path issues,
// via instance.UpdateCallObserver. The provider runs in-process under the
// acceptance harness, so the observer is visible to the test.
//
// The observer is a single package-level hook, so the tests that use it must
// not run in parallel with each other; they install it for their own duration
// and serialise on recorderMu.
type callRecorder struct {
	mu    sync.Mutex
	calls []string
}

var recorderMu sync.Mutex

func (r *callRecorder) record(call string) {
	r.mu.Lock()
	defer r.mu.Unlock()
	r.calls = append(r.calls, call)
}

func (r *callRecorder) reset() {
	r.mu.Lock()
	defer r.mu.Unlock()
	r.calls = nil
}

func (r *callRecorder) sorted() []string {
	r.mu.Lock()
	defer r.mu.Unlock()
	out := append([]string(nil), r.calls...)
	sort.Strings(out)

	return out
}

// installRecorder installs a fresh recorder for the duration of the test.
func installRecorder(t *testing.T) *callRecorder {
	t.Helper()

	recorderMu.Lock()
	rec := &callRecorder{}
	instance.UpdateCallObserver = rec.record
	t.Cleanup(func() {
		instance.UpdateCallObserver = nil
		recorderMu.Unlock()
	})

	return rec
}

// resetCalls discards whatever the create step recorded, so the following
// update step is counted from zero.
func resetCalls(rec *callRecorder) resource.TestCheckFunc {
	return func(_ *terraform.State) error {
		rec.reset()

		return nil
	}
}

// expectCalls asserts the exact set of appliance calls made since the last
// reset, then resets. An empty want means the step must have made no call.
func expectCalls(t *testing.T, rec *callRecorder, want ...string) resource.TestCheckFunc {
	t.Helper()

	return func(_ *terraform.State) error {
		got := rec.sorted()
		sort.Strings(want)
		if strings.Join(got, ",") != strings.Join(want, ",") {
			t.Errorf("update issued calls %v, want %v", got, want)
		}
		rec.reset()

		return nil
	}
}

// TestAccMorpheusInstanceResourceTimeoutsOnlyUpdateMakesNoCall is the
// MORPH-17751 acceptance criterion: editing only the timeouts block on an
// existing instance is an in-place update that makes no instance update or
// resize call. Before the rework it made both -- the gate compared at the model
// layer, where optional-and-computed volume children the practitioner never set
// are unknown in the plan and known in state, so every update looked like a
// change; and the resize path attached an (empty) service-plan-options object
// unconditionally, so the resize fired on every update too.
func TestAccMorpheusInstanceResourceTimeoutsOnlyUpdateMakesNoCall(t *testing.T) {
	defer testhelpers.RecordResult(t)

	capabilities.MustHaveOrSkip(t, capabilities.All)

	if testing.Short() {
		t.Skip("Skipping slow test in short mode")
	}

	// Not t.Parallel(): installs the package-level call observer.
	rec := installRecorder(t)

	providerConfig := testhelpers.ProviderBlock()
	name := acctest.RandomWithPrefix(t.Name())

	initial, err := instance.RenderInstanceConfig(t, map[string]string{
		"Name":          name,
		"CreateTimeout": "45m",
	})
	if err != nil {
		t.Fatal(err)
	}

	timeoutsOnly, err := instance.RenderInstanceConfig(t, map[string]string{
		"Name":          name,
		"CreateTimeout": "90m",
	})
	if err != nil {
		t.Fatal(err)
	}

	resource.Test(t, resource.TestCase{
		ProtoV6ProviderFactories: testhelpers.GetAccTestFactories(t, adapter.NewMorpheus(), nil),
		Steps: []resource.TestStep{
			{
				Config: providerConfig + initial,
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr(instanceResourceName, "timeouts.create", "45m"),
					// Create is not an update; whatever it did, start the count clean.
					resetCalls(rec),
				),
			},
			{
				Config: providerConfig + timeoutsOnly,
				ConfigPlanChecks: resource.ConfigPlanChecks{
					PreApply: []plancheck.PlanCheck{
						plancheck.ExpectResourceAction(instanceResourceName, plancheck.ResourceActionUpdate),
					},
				},
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr(instanceResourceName, "timeouts.create", "90m"),
					// The whole point: no appliance call for a timeouts-only change.
					expectCalls(t, rec),
				),
			},
			{
				Config:             providerConfig + timeoutsOnly,
				PlanOnly:           true,
				ExpectNonEmptyPlan: false,
			},
		},
	})
}

// TestAccMorpheusInstanceResourceNameOnlyUpdateMakesOnlyUpdateCall pins the
// other half of the split: a change carried by the instance update call issues
// exactly that call and no resize. Before the rework every update also issued a
// resize (and waited for the instance to settle), roughly doubling update time.
func TestAccMorpheusInstanceResourceNameOnlyUpdateMakesOnlyUpdateCall(t *testing.T) {
	defer testhelpers.RecordResult(t)

	capabilities.MustHaveOrSkip(t, capabilities.All)

	if testing.Short() {
		t.Skip("Skipping slow test in short mode")
	}

	rec := installRecorder(t)

	providerConfig := testhelpers.ProviderBlock()
	name := acctest.RandomWithPrefix(t.Name())

	initial, err := instance.RenderInstanceConfig(t, map[string]string{"Name": name})
	if err != nil {
		t.Fatal(err)
	}
	renamed, err := instance.RenderInstanceConfig(t, map[string]string{"Name": name + "-renamed"})
	if err != nil {
		t.Fatal(err)
	}

	resource.Test(t, resource.TestCase{
		ProtoV6ProviderFactories: testhelpers.GetAccTestFactories(t, adapter.NewMorpheus(), nil),
		Steps: []resource.TestStep{
			{
				Config: providerConfig + initial,
				Check:  resetCalls(rec),
			},
			{
				Config: providerConfig + renamed,
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr(instanceResourceName, "name", name+"-renamed"),
					expectCalls(t, rec, "update"),
				),
			},
		},
	})
}

// TestAccMorpheusInstanceResourcePlanIdChangeResizes pins that a plan_id
// change is applied. Before the rework it was silently dropped twice over: the
// gate did not consider plan_id, and the resize path's trigger did not include
// a plan change even when the gate fired for another reason.
func TestAccMorpheusInstanceResourcePlanIdChangeResizes(t *testing.T) {
	defer testhelpers.RecordResult(t)

	capabilities.MustHaveOrSkip(t, capabilities.All)

	if testing.Short() {
		t.Skip("Skipping slow test in short mode")
	}

	rec := installRecorder(t)

	providerConfig := testhelpers.ProviderBlock()
	name := acctest.RandomWithPrefix(t.Name())

	smallPlan, err := instance.RenderInstanceConfig(t, map[string]string{"Name": name})
	if err != nil {
		t.Fatal(err)
	}
	// A stock plan every appliance seeds, so the second step is a genuine
	// plan change; TF_VAR_testacc_morpheus_hvm_plan_name must name another.
	biggerPlan, err := instance.RenderInstanceConfig(t, map[string]string{
		"Name":     name,
		"PlanName": "1 CPU, 2GB Memory",
	})
	if err != nil {
		t.Fatal(err)
	}

	resource.Test(t, resource.TestCase{
		ProtoV6ProviderFactories: testhelpers.GetAccTestFactories(t, adapter.NewMorpheus(), nil),
		Steps: []resource.TestStep{
			{
				Config: providerConfig + smallPlan,
				Check:  resetCalls(rec),
			},
			{
				Config: providerConfig + biggerPlan,
				ConfigPlanChecks: resource.ConfigPlanChecks{
					PreApply: []plancheck.PlanCheck{
						plancheck.ExpectResourceAction(instanceResourceName, plancheck.ResourceActionUpdate),
					},
				},
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttrPair(instanceResourceName, "plan_id",
						"data.hpe_morpheus_service_plan.vme_plan", "id"),
					// A plan change is a resize, and only a resize.
					expectCalls(t, rec, "resize"),
				),
			},
			{
				Config:             providerConfig + biggerPlan,
				PlanOnly:           true,
				ExpectNonEmptyPlan: false,
			},
		},
	})
}
