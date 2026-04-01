// (C) Copyright 2025 Hewlett Packard Enterprise Development LP

package ostype

import (
	"context"
	"fmt"
	"log"
	"net/http"
	"strings"

	"github.com/hashicorp/terraform-plugin-testing/helper/resource"

	"github.com/HewlettPackard/hpe-morpheus-go-sdk/oapigen/sdk"

	"github.com/HPE/terraform-provider-hpe/morpheus/utils/errfmt"
)

const testOsTypePrefix = "TestAccMorpheusOsType"

type osTypeSweeper struct {
	client *sdk.APIClient
}

func NewOsTypeSweeper(client *sdk.APIClient) {
	s := &osTypeSweeper{
		client: client,
	}

	resource.AddTestSweepers(
		"hpe_morpheus_os_type",
		&resource.Sweeper{
			Name: "hpe_morpheus_os_type",
			F:    s.Sweep,
		})
}

func (s *osTypeSweeper) Sweep(_ string) error {
	ctx := context.Background()

	if s.client == nil {
		log.Printf("[INFO] No client provided, skipping os type sweep")

		return nil
	}

	osTypes, hresp, err := s.client.LibraryAPI.ListOsTypes(ctx).
		Phrase(testOsTypePrefix).Execute()
	if err != nil || hresp.StatusCode != http.StatusOK {
		return fmt.Errorf("failed to list os types: %s", errfmt.ErrMsg(err, hresp))
	}

	osTypeList := osTypes.GetOsTypes()
	var sweptCount int
	var sweepErrors []string

	for _, osType := range osTypeList {
		name, ok := osType.GetNameOk()
		if !ok || name == nil {
			continue
		}

		if !strings.HasPrefix(*name, testOsTypePrefix) {
			log.Printf("[INFO] Skipping os type (name): %s", *name)

			continue
		}

		id, ok := osType.GetIdOk()
		if !ok || id == nil {
			log.Printf("[INFO] Skipping os type (id): %s", *name)

			continue
		}

		log.Printf("[INFO] Sweeping os type: %s (id: %d)", *name, *id)

		_, hresp, err := s.client.LibraryAPI.DeleteOsType(ctx, *id).Execute()
		if err != nil || hresp.StatusCode != http.StatusOK {
			errMsg := fmt.Sprintf(
				"failed to delete os type %s (id: %d): %s",
				*name, *id, errfmt.ErrMsg(err, hresp),
			)
			log.Printf("[ERROR] %s", errMsg)
			sweepErrors = append(sweepErrors, errMsg)

			continue
		}

		sweptCount++
	}

	log.Printf(
		"[INFO] OS type sweep completed. OS types swept: %d, errors: %d",
		sweptCount, len(sweepErrors),
	)

	if len(sweepErrors) > 0 {
		return fmt.Errorf("%s", strings.Join(sweepErrors, "\n"))
	}

	return nil
}
