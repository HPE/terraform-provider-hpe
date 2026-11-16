// (C) Copyright 2026 Hewlett Packard Enterprise Development LP

package network

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
	"strconv"
	"time"

	"github.com/cenkalti/backoff/v5"
	"github.com/hashicorp/terraform-plugin-framework/diag"
	"github.com/hashicorp/terraform-plugin-log/tflog"

	sdk "github.com/HPE/terraform-provider-hpe/internal/sdk/oapigen"
	"github.com/HPE/terraform-provider-hpe/morpheus/utils/errfmt"
	"github.com/HPE/terraform-provider-hpe/utils/paging"
)

const (
	// segmentKey is the config key naming the virtual switch segment an HVM
	// data network with the VSWITCH layout is built on.
	segmentKey = "virtualSwitchSegmentId"
	// segmentReadyInterval is how often the segment is re-read while waiting.
	segmentReadyInterval = 15 * time.Second
	// segmentReadyTimeout bounds the wait for a segment's host interface.
	// Morpheus fills it in on the cluster's inventory sync, which runs every
	// five minutes or so; two cycles is a comfortable margin.
	segmentReadyTimeout = 12 * time.Minute
)

var (
	// errSegmentNotOnSwitch reports that a switch, read successfully, has no
	// segment with the id asked for.
	errSegmentNotOnSwitch = errors.New("no such segment on the switch")
	// errSegmentNotFound reports that no switch in the cloud has the segment.
	errSegmentNotFound = errors.New("no virtual switch in the cloud has the segment")
)

// waitForSegmentReady holds a create whose config names a virtual switch
// segment until Morpheus knows the segment's host interface. A segment just
// added to a switch has none until the cluster's next inventory sync, and the
// HVM plugin builds a data network's interface naming from it — without it the
// create fails with an unexplained 500 — so a switch and a data network on it
// cannot otherwise be applied together.
//
// Nothing is done when the config names no segment, or no switch in the cloud
// has it: Morpheus reports a bad segment itself. Transient failures reading
// the clusters or switches are retried within the same bound; a wait that
// cannot be completed is an error, since letting the create through would
// only trade a clear message for the 500.
func waitForSegmentReady(ctx context.Context, client *sdk.APIClient, cloudID int64, cfg map[string]any) diag.Diagnostics {
	var diags diag.Diagnostics

	segmentID, ok := int64Value(cfg[segmentKey])
	if !ok {
		return diags
	}

	var switchID int64

	poll := func() (struct{}, error) {
		if switchID == 0 {
			id, err := findSegmentSwitch(ctx, client, cloudID, segmentID)
			if err != nil {
				return struct{}{}, err
			}

			if id == 0 {
				return struct{}{}, backoff.Permanent(errSegmentNotFound)
			}

			switchID = id
		}

		iface, err := segmentInterface(ctx, client, switchID, segmentID)
		if errors.Is(err, errSegmentNotOnSwitch) {
			// The segment was there a moment ago; it has been removed.
			return struct{}{}, backoff.Permanent(err)
		}

		if err != nil {
			return struct{}{}, err
		}

		if iface == "" {
			return struct{}{}, fmt.Errorf("segment %d of virtual switch %d has no host interface yet", segmentID, switchID)
		}

		return struct{}{}, nil
	}

	_, err := backoff.Retry(ctx, poll,
		backoff.WithBackOff(backoff.NewConstantBackOff(segmentReadyInterval)),
		backoff.WithMaxElapsedTime(segmentReadyTimeout))

	switch {
	case err == nil:
		return diags
	case errors.Is(err, errSegmentNotFound):
		tflog.Warn(ctx, "virtual switch segment not found in the cloud; leaving it to Morpheus", map[string]any{
			"segment": segmentID, "cloud": cloudID,
		})

		return diags
	default:
		diags.AddError("create network resource",
			fmt.Sprintf("virtual switch segment %d is not ready for a network: %s. "+
				"Morpheus records a segment's host interface on the cluster's inventory sync; "+
				"if the hosts are reachable, retry the apply once the switch's segments show one.",
				segmentID, err.Error()))

		return diags
	}
}

// findSegmentSwitch returns the id of the virtual switch in the cloud that
// carries the segment, or 0 when none does. A failure to read the clusters or
// a switch is returned, so the caller can retry it rather than conclude the
// segment is nowhere.
func findSegmentSwitch(ctx context.Context, client *sdk.APIClient, cloudID, segmentID int64) (int64, error) {
	// Every page of the cloud's clusters: a segment on a cluster past the
	// first page would otherwise be taken as absent.
	clusters, err := paging.Collect(ctx, func(
		ctx context.Context,
		offset, max int64,
	) ([]sdk.ListClusters200ResponseAllOfClustersInner, int64, error) {
		rs, hresp, err := client.ClustersAPI.ListClusters(ctx).ZoneId(cloudID).Max(max).Offset(offset).Execute()
		if err != nil || rs == nil || hresp.StatusCode != http.StatusOK {
			return nil, 0, fmt.Errorf("listing the clusters of cloud %d: %s", cloudID, errfmt.ErrMsg(err, hresp))
		}

		total := int64(0)
		if rs.Meta != nil && rs.Meta.Total != nil {
			total = *rs.Meta.Total
		}

		return rs.Clusters, total, nil
	})
	if err != nil {
		return 0, err
	}

	for _, c := range clusters {
		if c.Id == nil {
			continue
		}

		switches, hresp, err := client.VirtualSwitchesAPI.ListVirtualSwitches(ctx).ClusterId(*c.Id).Execute()
		if hresp != nil && hresp.StatusCode == http.StatusNotFound {
			// A cluster that is not HVM 2.0 has no switches; the listing says
			// the cluster is not found. It is not where the segment is.
			continue
		}

		if err != nil || switches == nil || hresp.StatusCode != http.StatusOK {
			return 0, fmt.Errorf("listing the virtual switches of cluster %d: %s", *c.Id, errfmt.ErrMsg(err, hresp))
		}

		// The listing names only each segment's traffic type; the ids are on
		// the switch itself.
		for _, vs := range switches.VirtualSwitches {
			if vs.Id == nil {
				continue
			}

			_, err := segmentInterface(ctx, client, *vs.Id, segmentID)
			if errors.Is(err, errSegmentNotOnSwitch) {
				continue
			}

			if err != nil {
				return 0, err
			}

			return *vs.Id, nil
		}
	}

	return 0, nil
}

// segmentInterface reads the host interface Morpheus records for a segment:
// empty when it has none yet, errSegmentNotOnSwitch when the switch has no
// such segment, and any other error when the switch could not be read.
func segmentInterface(ctx context.Context, client *sdk.APIClient, switchID, segmentID int64) (string, error) {
	vs, hresp, err := client.VirtualSwitchesAPI.GetVirtualSwitch(ctx, switchID).Execute()
	if err != nil || vs == nil || vs.VirtualSwitch == nil || hresp.StatusCode != http.StatusOK {
		return "", fmt.Errorf("reading virtual switch %d: %s", switchID, errfmt.ErrMsg(err, hresp))
	}

	for _, seg := range vs.VirtualSwitch.Segments {
		if seg.Id != nil && *seg.Id == segmentID {
			if seg.InterfaceName.IsSet() && seg.InterfaceName.Get() != nil {
				return *seg.InterfaceName.Get(), nil
			}

			return "", nil
		}
	}

	return "", fmt.Errorf("virtual switch %d: %w", switchID, errSegmentNotOnSwitch)
}

// int64Value reads an integer however the config value was decoded.
func int64Value(v any) (int64, bool) {
	switch n := v.(type) {
	case int64:
		return n, true
	case int:
		return int64(n), true
	case float64:
		return int64(n), true
	case json.Number:
		i, err := n.Int64()

		return i, err == nil
	case string:
		i, err := strconv.ParseInt(n, 10, 64)

		return i, err == nil
	default:
		return 0, false
	}
}
