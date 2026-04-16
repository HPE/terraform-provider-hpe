// (C) Copyright 2026 Hewlett Packard Enterprise Development LP

package cluster

import (
	"context"
	"fmt"
	"log"
	"net/http"
	"strings"

	"github.com/HewlettPackard/hpe-morpheus-go-sdk/oapigen/sdk"
	"github.com/hashicorp/terraform-plugin-testing/helper/resource"

	"github.com/HPE/terraform-provider-hpe/morpheus/utils/errfmt"
)

// Clusters whose name begins with this string will be eligible for deletion.
const testClusterPrefix = "TestAccMorpheusCluster"

// clusterSweeper handles cleanup of test clusters.
type clusterSweeper struct {
	client *sdk.APIClient
}

// NewClusterSweeper creates and registers a cluster sweeper.
func NewClusterSweeper(client *sdk.APIClient) {
	s := &clusterSweeper{
		client: client,
	}

	resource.AddTestSweepers(
		"hpe_morpheus_cluster",
		&resource.Sweeper{
			Name: "hpe_morpheus_cluster",
			F:    s.Sweep,
		},
	)
}

// Sweep cleans up test clusters.
func (s *clusterSweeper) Sweep(_ string) error {
	ctx := context.Background()

	if s.client == nil {
		log.Printf("[INFO] No client provided, skipping cluster sweep")

		return nil
	}

	clusters, hresp, err := s.client.ClustersAPI.ListClusters(ctx).
		Phrase(testClusterPrefix).Execute()
	if err != nil {
		// Handle 404 and 403 as "no matches found" rather than an error.
		if hresp != nil && (hresp.StatusCode == http.StatusNotFound || hresp.StatusCode == http.StatusForbidden) {
			log.Printf("[INFO] No clusters found matching prefix (status %d): %s", hresp.StatusCode, testClusterPrefix)

			return nil
		}

		return fmt.Errorf("failed to list clusters: %s", errfmt.ErrMsg(err, hresp))
	}
	if hresp.StatusCode != http.StatusOK {
		return fmt.Errorf("failed to list clusters: %s", errfmt.ErrMsg(err, hresp))
	}

	clusterList := clusters.GetClusters()
	var sweptCount int
	var sweepErrors []string

	for _, cluster := range clusterList {
		name, ok := cluster.GetNameOk()
		if !ok || name == nil {
			continue
		}

		if !strings.HasPrefix(*name, testClusterPrefix) {
			log.Printf("[INFO] Skipping cluster (name): %s", *name)

			continue
		}

		id, ok := cluster.GetIdOk()
		if !ok || id == nil {
			log.Printf("[INFO] Skipping cluster (id): %s", *name)

			continue
		}

		log.Printf("[INFO] Sweeping cluster: %s (id: %d)", *name, *id)

		_, hresp, err := s.client.ClustersAPI.DeleteCluster(ctx, *id).Force("on").Execute()
		if err != nil || hresp.StatusCode != http.StatusOK {
			errMsg := fmt.Sprintf(
				"failed to delete cluster %s (id: %d): %s",
				*name, *id, errfmt.ErrMsg(err, hresp),
			)
			log.Printf("[ERROR] %s", errMsg)
			sweepErrors = append(sweepErrors, errMsg)

			continue
		}

		sweptCount++
	}

	log.Printf(
		"[INFO] Cluster sweep completed. Clusters swept: %d, errors: %d",
		sweptCount, len(sweepErrors),
	)

	if len(sweepErrors) > 0 {
		return fmt.Errorf("%s", strings.Join(sweepErrors, "\n"))
	}

	return nil
}
