// (C) Copyright 2026 Hewlett Packard Enterprise Development LP

package virtualswitch_test

import (
	"encoding/json"
	"fmt"
	"net/http"
	"net/http/httptest"
	"os"
	"regexp"
	"slices"
	"strconv"
	"strings"
	"sync"
	"testing"

	fwprovider "github.com/hashicorp/terraform-plugin-framework/provider"
	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
	"github.com/hashicorp/terraform-plugin-testing/plancheck"
	"github.com/hashicorp/terraform-plugin-testing/terraform"

	"github.com/HPE/terraform-provider-hpe/morpheus"
	"github.com/HPE/terraform-provider-hpe/morpheus/model"
	"github.com/HPE/terraform-provider-hpe/morpheus/testhelpers"
	"github.com/HPE/terraform-provider-hpe/morpheus/utils/clientfactory"
	"github.com/HPE/terraform-provider-hpe/provider/adapter"
)

func TestMain(m *testing.M) {
	code := testhelpers.TestMain(m)
	testhelpers.WriteMergedResults()
	os.Exit(code)
}

// fakeAppliance is a tiny in-memory Morpheus that implements just enough of the
// virtual switch API to exercise the provider's Create -> list-by-name -> poll
// -> Update-in-steps -> Delete flow, including the async process envelope, the
// sync guard, and the 404-for-unknown-cluster behaviour. It proves the
// PROVIDER's logic, not a real appliance's.
type fakeAppliance struct {
	mu        sync.Mutex
	nextSwID  int64
	nextSegID int64
	nextProc  int64
	switches  map[int64]*fakeSwitch
	// procPolls counts GET /api/processes/{id} calls per process so a process
	// reports "running" for the first two polls and "complete" thereafter.
	procPolls map[int64]int
	// requestLog records mutation paths in order for sequence assertions.
	requestLog []string
	// refuseMutations makes the next N mutations answer the sync guard (a
	// 400 with errors.sync) although no sync is in flight, as a sync that
	// starts between the provider's pre-wait and its mutation would.
	refuseMutations int
	refusals        int
	refusedPaths    []string
	// syncingPolls makes the next N GETs of any switch report status
	// "syncing" with no process behind it: a sync started by the appliance
	// itself, which the provider must wait out before mutating.
	syncingPolls int
}

type fakeSwitch struct {
	id        int64
	clusterID int64
	name      string
	baseType  string
	bondMode  *string
	mtu       int32
	desc      *string
	status    string
	segments  []*fakeSeg
	uplinks   []fakeUplink
	// mutated flips true after the first mutation; it drives a changed uplink
	// syncMessage so the detail returned before and after an update differ,
	// reproducing the "inconsistent result after apply" set-correlation bug.
	mutated bool
}

type fakeSeg struct {
	// mtu mirrors 9.1's asymmetry: a segment created with the switch keeps the
	// request's (null) mtu, a segment added by a composite edit gets vs.mtu.
	mtu         *int32
	id          int64
	trafficType string
	vlanID      *int32
	ipAddress   *string
	netmask     *string
	hostIps     []map[string]any
}

type fakeUplink struct {
	id       int64
	serverID int64
	nic1     string
	nic2     *string
}

const knownClusterID = 4242

func newFakeAppliance() *fakeAppliance {
	return &fakeAppliance{
		nextSwID:  1,
		nextSegID: 1000,
		nextProc:  500,
		switches:  map[int64]*fakeSwitch{},
		procPolls: map[int64]int{},
	}
}

func (f *fakeAppliance) log(path string) {
	f.requestLog = append(f.requestLog, path)
}

func (f *fakeAppliance) startProc() int64 {
	f.nextProc++

	return f.nextProc
}

var (
	reGet     = regexp.MustCompile(`^/api/virtual-switches/(\d+)$`)
	reRename  = regexp.MustCompile(`^/api/virtual-switches/(\d+)/rename$`)
	reEdit    = regexp.MustCompile(`^/api/virtual-switches/(\d+)/edit$`)
	reSegment = regexp.MustCompile(`^/api/virtual-switches/(\d+)/segments/(\d+)$`)
	reProcess = regexp.MustCompile(`^/api/processes/(\d+)$`)
)

func (f *fakeAppliance) handler() http.Handler {
	mux := http.NewServeMux()
	mux.HandleFunc("/", f.route)

	return mux
}

func writeJSON(w http.ResponseWriter, code int, body any) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(code)
	_ = json.NewEncoder(w).Encode(body)
}

func asyncOK(processID int64) map[string]any {
	return map[string]any{"success": true, "processId": processID}
}

func syncGuard() (int, map[string]any) {
	return http.StatusBadRequest, map[string]any{
		"success": false,
		"errors":  map[string]any{"sync": "A sync is in flight"},
	}
}

func (f *fakeAppliance) route(w http.ResponseWriter, r *http.Request) {
	f.mu.Lock()
	defer f.mu.Unlock()

	path := r.URL.Path

	switch {
	case path == "/api/whoami":
		writeJSON(w, http.StatusOK, map[string]any{
			"appliance": map[string]any{"buildVersion": "9.1.0"},
		})

		return
	case path == "/api/setup/check":
		writeJSON(w, http.StatusOK, map[string]any{"buildVersion": "9.1.0"})

		return
	case path == "/api/virtual-switches" && r.Method == http.MethodPost:
		f.createSwitch(w, r)

		return
	case path == "/api/virtual-switches" && r.Method == http.MethodGet:
		f.listSwitches(w, r)

		return
	}

	if m := reProcess.FindStringSubmatch(path); m != nil {
		f.getProcess(w, m[1])

		return
	}

	if m := reRename.FindStringSubmatch(path); m != nil && r.Method == http.MethodPut {
		f.rename(w, r, m[1])

		return
	}

	if m := reEdit.FindStringSubmatch(path); m != nil && r.Method == http.MethodPut {
		f.compositeEdit(w, r, m[1])

		return
	}

	if m := reSegment.FindStringSubmatch(path); m != nil && r.Method == http.MethodDelete {
		f.deleteSegment(w, r, m[1], m[2])

		return
	}

	if m := reGet.FindStringSubmatch(path); m != nil {
		switch r.Method {
		case http.MethodGet:
			f.getSwitch(w, m[1])
		case http.MethodPut:
			f.updateSwitch(w, r, m[1])
		case http.MethodDelete:
			f.deleteSwitch(w, r, m[1])
		}

		return
	}

	http.NotFound(w, r)
}

func (f *fakeAppliance) getProcess(w http.ResponseWriter, idStr string) {
	id, _ := strconv.ParseInt(idStr, 10, 64)
	f.procPolls[id]++

	status := "complete"
	if f.procPolls[id] <= 2 {
		status = "running"
	} else {
		// Once the process completes, the switch becomes available.
		for _, sw := range f.switches {
			if sw.status == "syncing" {
				sw.status = "available"
			}
		}
	}

	writeJSON(w, http.StatusOK, map[string]any{
		"process": map[string]any{
			"id": id, "displayName": "Apply Virtual Switch", "status": status,
		},
	})
}

func (f *fakeAppliance) createSwitch(w http.ResponseWriter, r *http.Request) {
	clusterID, _ := strconv.ParseInt(r.URL.Query().Get("clusterId"), 10, 64)
	if clusterID != knownClusterID {
		writeJSON(w, http.StatusNotFound, map[string]any{"success": false})

		return
	}

	var body struct {
		VirtualSwitch struct {
			Name     string  `json:"name"`
			BaseType *string `json:"baseType"`
			BondMode *string `json:"bondMode"`
			Mtu      *int32  `json:"mtu"`
			Segments []struct {
				TrafficType string           `json:"trafficType"`
				VlanId      *int32           `json:"vlanId"`
				IpAddress   *string          `json:"ipAddress"`
				Netmask     *string          `json:"netmask"`
				HostIps     []map[string]any `json:"hostIps"`
			} `json:"segments"`
			Uplinks []struct {
				ServerId int64   `json:"serverId"`
				Nic1     string  `json:"nic1"`
				Nic2     *string `json:"nic2"`
			} `json:"uplinks"`
		} `json:"virtualSwitch"`
	}
	_ = json.NewDecoder(r.Body).Decode(&body)

	sw := &fakeSwitch{
		id:        f.nextSwID,
		clusterID: clusterID,
		name:      body.VirtualSwitch.Name,
		baseType:  "general",
		bondMode:  body.VirtualSwitch.BondMode,
		mtu:       1500,
		status:    "syncing",
	}
	f.nextSwID++

	if body.VirtualSwitch.BaseType != nil {
		sw.baseType = *body.VirtualSwitch.BaseType
	}

	if body.VirtualSwitch.Mtu != nil {
		sw.mtu = *body.VirtualSwitch.Mtu
	}

	for _, s := range body.VirtualSwitch.Segments {
		seg := &fakeSeg{
			id:          f.nextSegID,
			trafficType: s.TrafficType,
			vlanID:      s.VlanId,
			ipAddress:   s.IpAddress,
			netmask:     s.Netmask,
			hostIps:     s.HostIps,
		}
		f.nextSegID++
		sw.segments = append(sw.segments, seg)
	}

	for _, u := range body.VirtualSwitch.Uplinks {
		sw.uplinks = append(sw.uplinks, fakeUplink{
			id: f.nextSegID, serverID: u.ServerId, nic1: u.Nic1, nic2: u.Nic2,
		})
		f.nextSegID++
	}

	f.switches[sw.id] = sw
	f.log("POST /virtual-switches")

	writeJSON(w, http.StatusOK, asyncOK(f.startProc()))
}

func (f *fakeAppliance) listSwitches(w http.ResponseWriter, r *http.Request) {
	clusterID, _ := strconv.ParseInt(r.URL.Query().Get("clusterId"), 10, 64)
	if clusterID != knownClusterID {
		writeJSON(w, http.StatusNotFound, map[string]any{"success": false})

		return
	}

	var out []map[string]any
	for _, sw := range f.switches {
		if sw.clusterID == clusterID {
			out = append(out, map[string]any{"id": sw.id, "name": sw.name})
		}
	}

	writeJSON(w, http.StatusOK, map[string]any{"virtualSwitches": out})
}

func (f *fakeAppliance) getSwitch(w http.ResponseWriter, idStr string) {
	id, _ := strconv.ParseInt(idStr, 10, 64)

	sw, ok := f.switches[id]
	if !ok {
		writeJSON(w, http.StatusNotFound, map[string]any{"success": false})

		return
	}

	detail := f.switchDetail(sw)
	if f.syncingPolls > 0 {
		f.syncingPolls--
		detail["status"] = "syncing"
	}

	writeJSON(w, http.StatusOK, map[string]any{"virtualSwitch": detail})
}

func (f *fakeAppliance) switchDetail(sw *fakeSwitch) map[string]any {
	segs := make([]map[string]any, 0, len(sw.segments))
	for _, s := range sw.segments {
		seg := map[string]any{
			"id": s.id, "trafficType": s.trafficType, "status": "available",
		}
		if s.vlanID != nil {
			seg["vlanId"] = *s.vlanID
		}
		if s.mtu != nil {
			seg["mtu"] = *s.mtu
		}
		// Morpheus' cluster refresh backfills a segment's ipAddress from its
		// first host entry when none was given; mirror that so the test proves
		// ip_address is tolerated as computed rather than reported inconsistent.
		if s.ipAddress != nil {
			seg["ipAddress"] = *s.ipAddress
		} else if len(s.hostIps) > 0 {
			if ip, ok := s.hostIps[0]["ip"].(string); ok && ip != "" {
				seg["ipAddress"] = ip
			}
		}
		if s.netmask != nil {
			seg["netmask"] = *s.netmask
		}
		if len(s.hostIps) > 0 {
			seg["hostIps"] = s.hostIps
		}
		segs = append(segs, seg)
	}

	// A real appliance changes an uplink's syncMessage as work runs; emulate
	// that by returning a non-empty message before the first mutation and an
	// empty one afterwards. The provider must mark computed uplink fields unknown
	// on update so this change does not break set correlation at apply time.
	syncMsg := "Interface configuration applied successfully."
	if sw.mutated {
		syncMsg = ""
	}

	ups := make([]map[string]any, 0, len(sw.uplinks))
	for _, u := range sw.uplinks {
		up := map[string]any{
			"id": u.id, "serverId": u.serverID, "nic1": u.nic1,
			"syncStatus": "synced", "syncMessage": syncMsg,
		}
		if u.nic2 != nil {
			up["nic2"] = *u.nic2
		}
		ups = append(ups, up)
	}

	detail := map[string]any{
		"id": sw.id, "name": sw.name, "baseType": sw.baseType, "mtu": sw.mtu,
		"status": sw.status, "active": true, "clusterId": sw.clusterID,
		"segments": segs, "uplinks": ups,
	}
	if sw.bondMode != nil {
		detail["bondMode"] = *sw.bondMode
	}
	if sw.desc != nil {
		detail["description"] = *sw.desc
	}

	return detail
}

func (f *fakeAppliance) guardOrStart(w http.ResponseWriter, sw *fakeSwitch, what string) (int64, bool) {
	refuse := f.refuseMutations > 0
	if refuse {
		f.refuseMutations--
	}

	// Every guard answer is counted, whatever caused it: a forced refusal,
	// a mutation of the fake's own in flight, or an appliance-started sync
	// (syncingPolls) that the provider should have waited out first.
	if sw.status == "syncing" || f.syncingPolls > 0 || refuse {
		f.refusals++
		f.refusedPaths = append(f.refusedPaths, what)

		code, body := syncGuard()
		writeJSON(w, code, body)

		return 0, false
	}

	sw.status = "syncing"
	sw.mutated = true

	return f.startProc(), true
}

func (f *fakeAppliance) rename(w http.ResponseWriter, r *http.Request, idStr string) {
	id, _ := strconv.ParseInt(idStr, 10, 64)
	sw := f.switches[id]

	var body struct {
		Name string `json:"name"`
	}
	_ = json.NewDecoder(r.Body).Decode(&body)

	proc, ok := f.guardOrStart(w, sw, "PUT /rename")
	if !ok {
		return
	}

	sw.name = body.Name
	f.log("PUT /rename")
	writeJSON(w, http.StatusOK, asyncOK(proc))
}

func (f *fakeAppliance) compositeEdit(w http.ResponseWriter, r *http.Request, idStr string) {
	id, _ := strconv.ParseInt(idStr, 10, 64)
	sw := f.switches[id]

	var body struct {
		SwitchUpdate *struct {
			BondMode *string `json:"bondMode"`
			Mtu      *int32  `json:"mtu"`
		} `json:"switchUpdate"`
		NewSegments []struct {
			TrafficType string           `json:"trafficType"`
			VlanId      *int32           `json:"vlanId"`
			IpAddress   *string          `json:"ipAddress"`
			Netmask     *string          `json:"netmask"`
			HostIps     []map[string]any `json:"hostIps"`
		} `json:"newSegments"`
		EditedSegments []struct {
			SegId       int64   `json:"segId"`
			TrafficType *string `json:"trafficType"`
			VlanId      *int32  `json:"vlanId"`
		} `json:"editedSegments"`
	}
	_ = json.NewDecoder(r.Body).Decode(&body)

	proc, ok := f.guardOrStart(w, sw, "PUT /edit")
	if !ok {
		return
	}

	if body.SwitchUpdate != nil {
		if body.SwitchUpdate.BondMode != nil {
			if *body.SwitchUpdate.BondMode == "none" {
				sw.bondMode = nil
			} else {
				bm := *body.SwitchUpdate.BondMode
				sw.bondMode = &bm
			}
		}
		if body.SwitchUpdate.Mtu != nil {
			sw.mtu = *body.SwitchUpdate.Mtu
		}
	}

	for _, ns := range body.NewSegments {
		mtu := sw.mtu
		seg := &fakeSeg{
			id: f.nextSegID, trafficType: ns.TrafficType, vlanID: ns.VlanId,
			ipAddress: ns.IpAddress, netmask: ns.Netmask, hostIps: ns.HostIps,
			mtu: &mtu,
		}
		f.nextSegID++
		sw.segments = append(sw.segments, seg)
	}

	for _, es := range body.EditedSegments {
		// Morpheus skips an edited segment with no trafficType before it even
		// looks at segId; without it the edit is silently lost.
		if es.TrafficType == nil || *es.TrafficType == "" {
			continue
		}
		for _, seg := range sw.segments {
			if seg.id == es.SegId {
				seg.vlanID = es.VlanId
			}
		}
	}

	f.log("PUT /edit")
	writeJSON(w, http.StatusOK, asyncOK(proc))
}

func (f *fakeAppliance) updateSwitch(w http.ResponseWriter, r *http.Request, idStr string) {
	id, _ := strconv.ParseInt(idStr, 10, 64)
	sw := f.switches[id]

	var body struct {
		VirtualSwitch struct {
			Description *string `json:"description"`
		} `json:"virtualSwitch"`
	}
	_ = json.NewDecoder(r.Body).Decode(&body)

	proc, ok := f.guardOrStart(w, sw, "PUT /virtual-switches")
	if !ok {
		return
	}

	sw.desc = body.VirtualSwitch.Description
	f.log("PUT /update")
	writeJSON(w, http.StatusOK, asyncOK(proc))
}

func (f *fakeAppliance) deleteSegment(w http.ResponseWriter, r *http.Request, idStr, segStr string) {
	id, _ := strconv.ParseInt(idStr, 10, 64)
	segID, _ := strconv.ParseInt(segStr, 10, 64)
	sw := f.switches[id]

	proc, ok := f.guardOrStart(w, sw, "DELETE /segments")
	if !ok {
		return
	}

	kept := sw.segments[:0]
	for _, s := range sw.segments {
		if s.id != segID {
			kept = append(kept, s)
		}
	}
	sw.segments = kept

	f.log(fmt.Sprintf("DELETE /segments/%d", segID))
	writeJSON(w, http.StatusOK, asyncOK(proc))
}

// deleteSwitch is guarded like every mutation (VirtualSwitchService refuses
// a delete while a sync is in flight) unless force is set, as the API does.
func (f *fakeAppliance) deleteSwitch(w http.ResponseWriter, r *http.Request, idStr string) {
	id, _ := strconv.ParseInt(idStr, 10, 64)

	sw, ok := f.switches[id]
	if !ok {
		writeJSON(w, http.StatusNotFound, map[string]any{"success": false})

		return
	}

	var proc int64
	if r.URL.Query().Get("force") == "true" {
		proc = f.startProc()
	} else {
		var started bool
		if proc, started = f.guardOrStart(w, sw, "DELETE /virtual-switches"); !started {
			return
		}
	}

	delete(f.switches, id)
	f.log("DELETE /virtual-switches")
	writeJSON(w, http.StatusOK, asyncOK(proc))
}

func fakeProviderConfig(url string) string {
	return fmt.Sprintf(`
provider "hpe" {
  morpheus {
    url          = %q
    access_token = "fake-token"
    insecure     = true
  }
}
`, url)
}

func TestUnitVirtualSwitchFakeAppliance(t *testing.T) {
	t.Parallel()
	defer testhelpers.RecordResult(t)

	fake := newFakeAppliance()
	srv := httptest.NewServer(fake.handler())
	defer srv.Close()

	factories := testhelpers.GetAccTestFactories(t, newTestProvider(), nil)

	baseCfg := fakeProviderConfig(srv.URL)

	createCfg := baseCfg + fmt.Sprintf(`
resource "hpe_morpheus_virtual_switch" "sw" {
  cluster_id = %d
  name       = "vsw1"

  uplinks = [
    { server_id = 7, nic1 = "eth2" },
    { server_id = 8, nic1 = "eth2" },
  ]

  segments = {
    mgmt = { traffic_type = "management", ip_address = "10.0.0.2", netmask = "255.255.255.0" }
    vm   = { traffic_type = "vm" }
  }
}
`, knownClusterID)

	// renameOnlyCfg changes only the name. It reproduces the set-correlation bug:
	// after the rename the appliance returns a different uplink syncMessage than
	// the create read recorded, so without marking computed uplink fields unknown
	// the apply fails with "Provider produced inconsistent result after apply".
	renameOnlyCfg := baseCfg + fmt.Sprintf(`
resource "hpe_morpheus_virtual_switch" "sw" {
  cluster_id = %d
  name       = "vsw1b"

  uplinks = [
    { server_id = 7, nic1 = "eth2" },
    { server_id = 8, nic1 = "eth2" },
  ]

  segments = {
    mgmt = { traffic_type = "management", ip_address = "10.0.0.2", netmask = "255.255.255.0" }
    vm   = { traffic_type = "vm" }
  }
}
`, knownClusterID)

	updateCfg := baseCfg + fmt.Sprintf(`
resource "hpe_morpheus_virtual_switch" "sw" {
  cluster_id = %d
  name       = "vsw2"
  bond_mode  = "active-backup"

  uplinks = [
    { server_id = 7, nic1 = "eth2" },
    { server_id = 8, nic1 = "eth2" },
  ]

  segments = {
    mgmt = { traffic_type = "management", ip_address = "10.0.0.2", netmask = "255.255.255.0" }
    vm   = { traffic_type = "vm" }
    nfs  = {
      traffic_type = "data-nfs"
      vlan_id      = 100
      netmask      = "255.255.255.0"
      host_ips     = [{ server_id = 7, ip = "10.0.0.7" }]
    }
  }
}
`, knownClusterID)

	vlanEditCfg := strings.Replace(updateCfg, "vlan_id      = 100", "vlan_id      = 200", 1)

	removeSegCfg := baseCfg + fmt.Sprintf(`
resource "hpe_morpheus_virtual_switch" "sw" {
  cluster_id = %d
  name       = "vsw2"
  bond_mode  = "active-backup"

  uplinks = [
    { server_id = 7, nic1 = "eth2" },
    { server_id = 8, nic1 = "eth2" },
  ]

  segments = {
    mgmt = { traffic_type = "management", ip_address = "10.0.0.2", netmask = "255.255.255.0" }
    vm   = { traffic_type = "vm" }
  }
}
`, knownClusterID)

	// The final configuration: applied by the replacement step, destroyed last.
	lastCfg := baseCfg + fmt.Sprintf(`
resource "hpe_morpheus_virtual_switch" "sw" {
  cluster_id = %d
  name       = "vsw2"
  bond_mode  = "active-backup"

  uplinks = [
    { server_id = 7, nic1 = "eth2" },
    { server_id = 8, nic1 = "eth2" },
  ]

  segments = {
    mgmt = { traffic_type = "live-migration" }
    vm   = { traffic_type = "vm" }
  }
}
`, knownClusterID)

	resource.UnitTest(t, resource.TestCase{
		ProtoV6ProviderFactories: factories,
		Steps: []resource.TestStep{
			{
				Config: createCfg,
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr("hpe_morpheus_virtual_switch.sw", "name", "vsw1"),
					resource.TestCheckResourceAttr("hpe_morpheus_virtual_switch.sw", "status", "available"),
					resource.TestCheckResourceAttr("hpe_morpheus_virtual_switch.sw", "segments.%", "2"),
				),
			},
			{
				// Name-only update: the appliance now returns a changed uplink
				// syncMessage. This step passes only because ModifyPlan marks the
				// volatile computed uplink fields unknown in the plan.
				Config: renameOnlyCfg,
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr("hpe_morpheus_virtual_switch.sw", "name", "vsw1b"),
					resource.TestCheckResourceAttr("hpe_morpheus_virtual_switch.sw", "uplinks.#", "2"),
				),
			},
			{
				// rename + add segment + bond change in one apply; assert the
				// mutation sequence was rename, then compositeEdit. The first
				// mutation (the rename) is refused once by the guard and must be
				// retried, with the composite edit still after it.
				PreConfig: func() { fake.beginStep(); fake.setRefuseMutations(1) },
				Config:    updateCfg,
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr("hpe_morpheus_virtual_switch.sw", "name", "vsw2"),
					resource.TestCheckResourceAttr("hpe_morpheus_virtual_switch.sw", "bond_mode", "active-backup"),
					resource.TestCheckResourceAttr("hpe_morpheus_virtual_switch.sw", "segments.%", "3"),
					func(_ *terraform.State) error { return assertSequence(fake, "PUT /rename", "PUT /edit") },
					func(_ *terraform.State) error { return assertRefusals(fake, 1, "PUT /rename") },
				),
			},
			{
				// vlan change: editedSegments with segId, no DELETE. The switch
				// reports an appliance-started sync for the first two GETs, so
				// the update's pre-wait must poll it out before the first
				// mutation: a mutation sent during it is answered by the guard
				// and counted, so assertRefusals(0) fails without the pre-wait.
				PreConfig: func() { fake.beginStep(); fake.setSyncingPolls(2) },
				Config:    vlanEditCfg,
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr("hpe_morpheus_virtual_switch.sw", "segments.nfs.vlan_id", "200"),
					func(_ *terraform.State) error { return assertNoDelete(fake) },
					func(_ *terraform.State) error { return assertRefusals(fake, 0) },
					func(_ *terraform.State) error { return assertSyncWaitedOut(fake) },
				),
			},
			{
				// remove a segment: DELETE segments/{segId}. The first attempt is
				// answered with the sync guard (as the real appliance does when a
				// sync starts after the pre-wait); the provider must retry it, so
				// the apply succeeds and exactly one refusal was served.
				PreConfig: func() { fake.beginStep(); fake.setRefuseMutations(1) },
				Config:    removeSegCfg,
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr("hpe_morpheus_virtual_switch.sw", "segments.%", "2"),
					func(_ *terraform.State) error { return assertHasDelete(fake) },
					func(_ *terraform.State) error { return assertRefusals(fake, 1, "DELETE /segments") },
				),
			},
			{
				// Removing a host uplink is a plan-time error.
				Config: baseCfg + fmt.Sprintf(`
resource "hpe_morpheus_virtual_switch" "sw" {
  cluster_id = %d
  name       = "vsw2"
  bond_mode  = "active-backup"

  uplinks = [
    { server_id = 7, nic1 = "eth2" },
  ]

  segments = {
    mgmt = { traffic_type = "management", ip_address = "10.0.0.2", netmask = "255.255.255.0" }
    vm   = { traffic_type = "vm" }
  }
}
`, knownClusterID),
				ExpectError: regexp.MustCompile("Cannot remove a host uplink"),
			},
			{
				// Changing a segment traffic_type forces replacement.
				Config: lastCfg,
				ConfigPlanChecks: resource.ConfigPlanChecks{
					PreApply: []plancheck.PlanCheck{
						plancheck.ExpectResourceAction(
							"hpe_morpheus_virtual_switch.sw",
							plancheck.ResourceActionReplace,
						),
					},
				},
			},
			{
				ResourceName:      "hpe_morpheus_virtual_switch.sw",
				ImportState:       true,
				ImportStateVerify: true,
				ImportStateVerifyIgnore: []string{
					"timeouts", "segments", // synthesised keys differ from user keys
				},
			},
			{
				// Destroy while a sync is in flight: the DELETE is refused once
				// by the guard (the appliance's delete is guarded too, unless
				// forced) and must be retried once the switch settles.
				PreConfig: func() { fake.beginStep(); fake.setRefuseMutations(1) },
				Config:    lastCfg,
				Destroy:   true,
				Check: resource.ComposeAggregateTestCheckFunc(
					func(_ *terraform.State) error { return assertRefusals(fake, 1, "DELETE /virtual-switches") },
					func(_ *terraform.State) error { return assertHasPath(fake, "DELETE /virtual-switches") },
				),
			},
		},
	})
}

func assertHasPath(f *fakeAppliance, want string) error {
	f.mu.Lock()
	defer f.mu.Unlock()

	if slices.Contains(f.requestLog, want) {
		return nil
	}

	return fmt.Errorf("expected %q in %v", want, f.requestLog)
}

func assertSequence(f *fakeAppliance, first, second string) error {
	f.mu.Lock()
	defer f.mu.Unlock()

	firstIdx, secondIdx := -1, -1
	for i, p := range f.requestLog {
		if p == first && firstIdx == -1 {
			firstIdx = i
		}
		if p == second && firstIdx != -1 && secondIdx == -1 {
			secondIdx = i
		}
	}

	if firstIdx == -1 || secondIdx == -1 || firstIdx >= secondIdx {
		return fmt.Errorf("expected %q before %q in %v", first, second, f.requestLog)
	}

	return nil
}

func assertNoDelete(f *fakeAppliance) error {
	f.mu.Lock()
	defer f.mu.Unlock()

	for _, p := range f.requestLog {
		if strings.HasPrefix(p, "DELETE /segments/") {
			return fmt.Errorf("unexpected segment DELETE in %v", f.requestLog)
		}
	}

	return nil
}

func assertHasDelete(f *fakeAppliance) error {
	f.mu.Lock()
	defer f.mu.Unlock()

	for _, p := range f.requestLog {
		if strings.HasPrefix(p, "DELETE /segments/") {
			return nil
		}
	}

	return fmt.Errorf("expected a segment DELETE in %v", f.requestLog)
}

// newTestProvider builds the real Morpheus framework provider. The virtual
// switch resource is registered in morpheus/resources.go, so no additional
// wiring is needed here. It is an adapter-wrapped provider so the
// "hpe { morpheus { ... } }" nesting works exactly as in production.
func newTestProvider() fwprovider.Provider {
	inner := &morpheus.MorpheusProvider{
		NewClientFactory: func(m model.MorpheusProviderModel) *clientfactory.ClientFactory {
			return clientfactory.New(m)
		},
	}

	return adapter.NewAdaptedProvider(inner)
}

// beginStep resets what the assertions scan -- the request log and the
// refusal counters -- so a step is judged on its own requests only.
func (f *fakeAppliance) beginStep() {
	f.mu.Lock()
	defer f.mu.Unlock()

	f.requestLog = nil
	f.refusals, f.refusedPaths = 0, nil
}

func (f *fakeAppliance) setRefuseMutations(n int) {
	f.mu.Lock()
	defer f.mu.Unlock()

	f.refuseMutations = n
}

func (f *fakeAppliance) setSyncingPolls(n int) {
	f.mu.Lock()
	defer f.mu.Unlock()

	f.syncingPolls = n
}

// assertRefusals checks how many sync-guard refusals were served since the
// counters were last reset, and which mutations they hit (in order).
func assertRefusals(f *fakeAppliance, want int, paths ...string) error {
	f.mu.Lock()
	defer f.mu.Unlock()

	if f.refusals != want || f.refuseMutations != 0 {
		return fmt.Errorf("sync-guard refusals served = %d (pending %d), want %d served and none pending",
			f.refusals, f.refuseMutations, want)
	}

	if len(paths) > 0 && !slices.Equal(f.refusedPaths, paths) {
		return fmt.Errorf("refused mutations = %v, want %v", f.refusedPaths, paths)
	}

	return nil
}

// assertSyncWaitedOut checks that the appliance-started sync was polled to
// its end before any mutation.
func assertSyncWaitedOut(f *fakeAppliance) error {
	f.mu.Lock()
	defer f.mu.Unlock()

	if f.syncingPolls != 0 {
		return fmt.Errorf("%d syncing polls left unconsumed: the pre-wait did not poll the sync out", f.syncingPolls)
	}

	return nil
}
