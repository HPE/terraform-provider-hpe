// (C) Copyright 2025 Hewlett Packard Enterprise Development LP

// Package sweep allows deletion of dangling test resources
package sweep

import (
	"context"
	"fmt"
	"log"
	"net/http"
	"strings"

	"github.com/hashicorp/terraform-plugin-testing/helper/resource"

	"github.com/HPE/terraform-provider-hpe/internal/subproviders/morpheus/errors"
)

// Datastores whose name begins with this string will be eligible for deletion
const (
	testDatastorePrefix   = "TestAccMorpheusDatastore"
	enableDatastoreDelete = false
)

func Datastores() {
	resource.AddTestSweepers(
		"hpe_morpheus_datastore",
		&resource.Sweeper{
			Name: "hpe_morpheus_datastore",
			F:    testSweepMorpheusDatastores,
		})
}

func testSweepMorpheusDatastores(_ string) error {
	ctx := context.Background()

	client, err := NewSweepClient(ctx)
	if err != nil {
		return err
	}

	datastores, hresp, err := client.DatastoresAPI.ListDatastores(ctx).
		Phrase(testDatastorePrefix).Execute()
	if err != nil {
		// Handle 404 and 403 as "no matches found" rather than an error
		if hresp != nil && (hresp.StatusCode == http.StatusNotFound || hresp.StatusCode == http.StatusForbidden) {
			log.Printf("[INFO] No datastores found matching prefix (status %d): %s", hresp.StatusCode, testDatastorePrefix)

			return nil
		}

		return fmt.Errorf("failed to list datastores: %s", errors.ErrMsg(err, hresp))
	}
	if hresp.StatusCode != http.StatusOK {
		return fmt.Errorf("failed to list datastores: %s", errors.ErrMsg(err, hresp))
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
			_, hresp, err := client.DatastoresAPI.DeleteDatastores(ctx, *id).Execute()
			if err != nil || hresp.StatusCode != http.StatusOK {
				errMsg := fmt.Sprintf(
					"failed to delete datastore %s (id: %d): %s",
					*name, *id, errors.ErrMsg(err, hresp),
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
