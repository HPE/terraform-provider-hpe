// (C) Copyright 2025 Hewlett Packard Enterprise Development LP

package datastore

import (
	"context"
	"fmt"
	"log"
	"net/http"
	"strings"

	"github.com/HewlettPackard/hpe-morpheus-go-sdk/oapigen/sdk"
	"github.com/hashicorp/terraform-plugin-testing/helper/resource"

	morpheuserrors "github.com/HPE/terraform-provider-hpe/internal/framework/subproviders/morpheus/errors"
)

// Datastores whose name begins with this string will be eligible for deletion
const (
	testDatastorePrefix   = "TestAccMorpheusDatastore"
	enableDatastoreDelete = false
)

// datastoreSweeper handles cleanup of test datastores
type datastoreSweeper struct {
	client *sdk.APIClient
}

// NewDatastoreSweeper creates and registers a datastore sweeper
func NewDatastoreSweeper(client *sdk.APIClient) {
	s := &datastoreSweeper{
		client: client,
	}

	resource.AddTestSweepers(
		"hpe_morpheus_datastore",
		&resource.Sweeper{
			Name: "hpe_morpheus_datastore",
			F:    s.Sweep,
		})
}

// Sweep cleans up test datastores
func (s *datastoreSweeper) Sweep(_ string) error {
	ctx := context.Background()

	if s.client == nil {
		log.Printf("[INFO] No client provided, skipping datastore sweep")

		return nil
	}

	datastores, hresp, err := s.client.DatastoresAPI.ListDatastores(ctx).
		Phrase(testDatastorePrefix).Execute()
	if err != nil {
		// Handle 404 and 403 as "no matches found" rather than an error
		if hresp != nil && (hresp.StatusCode == http.StatusNotFound || hresp.StatusCode == http.StatusForbidden) {
			log.Printf("[INFO] No datastores found matching prefix (status %d): %s", hresp.StatusCode, testDatastorePrefix)

			return nil
		}

		return fmt.Errorf("failed to list datastores: %s", morpheuserrors.ErrMsg(err, hresp))
	}
	if hresp.StatusCode != http.StatusOK {
		return fmt.Errorf("failed to list datastores: %s", morpheuserrors.ErrMsg(err, hresp))
	}

	datastoreList := datastores.GetDatastores()
	var sweptCount int
	var sweepErrors []string

	for _, datastore := range datastoreList {
		name, ok := datastore.GetNameOk()
		if !ok || name == nil {
			continue
		}

		if !strings.HasPrefix(*name, testDatastorePrefix) {
			log.Printf("[INFO] Skipping datastore (name): %s", *name)

			continue
		}

		id, ok := datastore.GetIdOk()
		if !ok || id == nil {
			log.Printf("[INFO] Skipping datastore (id): %s", *name)

			continue
		}

		// Delete the datastore
		if enableDatastoreDelete {
			log.Printf("[INFO] Sweeping datastore: %s (id: %d)", *name, *id)
			_, hresp, err := s.client.DatastoresAPI.DeleteDatastores(ctx, *id).Execute()
			if err != nil || hresp.StatusCode != http.StatusOK {
				errMsg := fmt.Sprintf(
					"failed to delete datastore %s (id: %d): %s",
					*name, *id, morpheuserrors.ErrMsg(err, hresp),
				)
				log.Printf("[ERROR] %s", errMsg)
				sweepErrors = append(sweepErrors, errMsg)

				continue
			}

			sweptCount++
		} else {
			log.Printf("[INFO] datastore: %s (id: %d): delete disabled", *name, *id)
		}
	}

	log.Printf(
		"[INFO] Datastore sweep completed. Datastores swept: %d, errors: %d",
		sweptCount, len(sweepErrors),
	)

	if len(sweepErrors) > 0 {
		return fmt.Errorf("%s", strings.Join(sweepErrors, "\n"))
	}

	return nil
}
