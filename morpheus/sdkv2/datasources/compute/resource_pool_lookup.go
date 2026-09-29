// (C) Copyright 2026 Hewlett Packard Enterprise Development LP

package compute

import (
	"fmt"
	"log"
	"strconv"
	"strings"

	morpheus "github.com/HPE/terraform-provider-hpe/internal/sdk/legacy"
)

// The cloud's resource-pool listing (GET /api/zones/{id}/resource-pools) returns
// cloud-level pools only. The pool Morpheus creates for an HVM cluster is
// attached to the cluster, not the cloud, so a lookup by name cannot find it
// even though it exists and is the pool the user almost certainly wants. The
// messages below say so, and where possible give the id outright.

// resourcePoolNotFoundByNameMessage is the error for a name lookup that matched
// nothing. hint, when non-empty, is appended as its own paragraph.
func resourcePoolNotFoundByNameMessage(name string, cloudID int, hint string) string {
	var b strings.Builder

	fmt.Fprintf(&b, "No resource pool named %q was found in cloud %d.\n\n", name, cloudID)
	b.WriteString("This lookup uses the cloud's resource-pool listing, which returns cloud-level pools only. ")
	b.WriteString("The resource pool Morpheus creates for an HVM cluster is attached to the cluster, not the cloud, and is not listed. ")
	b.WriteString("To provision into a cluster, read the pool from the hpe_morpheus_cluster data source ")
	b.WriteString("(permissions.resource_pool.id), or look it up here by id.")

	if hint != "" {
		b.WriteString("\n\n")
		b.WriteString(hint)
	}

	return b.String()
}

// clusterResourcePoolHint is the paragraph added when a cluster with the looked-up
// name exists in the same cloud: it gives the id of that cluster's pool.
func clusterResourcePoolHint(name string, cloudID int, poolID int64) string {
	return fmt.Sprintf(
		"A cluster named %q exists in cloud %d; its provisioning pool has id %d. "+
			"Set id = %d, or use the hpe_morpheus_cluster data source.",
		name, cloudID, poolID, poolID,
	)
}

// resourcePoolMultipleMatchesMessage is the error for a name that matched more
// than one pool.
func resourcePoolMultipleMatchesMessage(name string, cloudID int, ids []int64) string {
	idStrs := make([]string, 0, len(ids))
	for _, id := range ids {
		idStrs = append(idStrs, strconv.FormatInt(id, 10))
	}

	return fmt.Sprintf(
		"%d resource pools named %q were found in cloud %d (ids %s). Use id to select one.",
		len(ids), name, cloudID, strings.Join(idStrs, ", "),
	)
}

// resourcePoolNotFoundByIDMessage is the error for an id lookup that returned 404.
func resourcePoolNotFoundByIDMessage(id, cloudID int) string {
	return fmt.Sprintf("Resource pool id %d was not found in cloud %d.", id, cloudID)
}

// findClusterResourcePoolID looks for exactly one cluster named name in cloud
// cloudID and returns the id of its resource pool. It is called only after a
// name lookup has already failed, to enrich the error, so it never fails
// loudly: any error, ambiguity or absence returns (0, false) and the caller
// falls back to the plain message. Surfacing a failure here as a diagnostic
// would risk burying the primary error, so each reason for giving up is
// logged instead and can be read with TF_LOG=DEBUG.
func findClusterResourcePoolID(client *morpheus.Client, name string, cloudID int) (int64, bool) {
	skip := func(reason string, args ...any) (int64, bool) {
		log.Printf("[DEBUG] resource pool %q in cloud %d: cluster hint skipped: %s",
			name, cloudID, fmt.Sprintf(reason, args...))

		return 0, false
	}

	resp, err := client.ListClusters(&morpheus.Request{
		QueryParams: map[string]string{"name": name},
	})
	if err != nil || resp == nil {
		return skip("listing clusters failed: %v", err)
	}

	list, ok := resp.Result.(*morpheus.ListClustersResult)
	if !ok || list.Clusters == nil {
		return skip("cluster list response had an unexpected shape")
	}

	var matched []morpheus.Cluster

	for _, c := range *list.Clusters {
		if strings.EqualFold(c.Name, name) && c.Zone.Id == int64(cloudID) {
			matched = append(matched, c)
		}
	}

	if len(matched) != 1 {
		return skip("%d cluster(s) of that name in the cloud", len(matched))
	}

	// permissions.resourcePool is only present on the single-cluster GET.
	getResp, err := client.GetCluster(matched[0].ID, &morpheus.Request{})
	if err != nil || getResp == nil {
		return skip("reading cluster %d failed: %v", matched[0].ID, err)
	}

	got, ok := getResp.Result.(*morpheus.GetClusterResult)
	if !ok || got.Cluster == nil || got.Cluster.Permissions.ResourcePool == nil {
		return skip("cluster %d reported no resource pool", matched[0].ID)
	}

	return got.Cluster.Permissions.ResourcePool.ID, true
}
