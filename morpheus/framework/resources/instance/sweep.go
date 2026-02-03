// (C) Copyright 2025 Hewlett Packard Enterprise Development LP

package instance

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

// Instances whose name begins with this string will be eligible for deletion
const testInstancePrefix = "TestAccMorpheusInstance"

// All of these tags must be present with value "true" for the instance to be deleted
var requiredSweepTags = []string{
	"terraform",
	"acctest",
	"hpe_morpheus_instance",
	"sweepable",
}

func hasRequiredTags(tags []sdk.AddInstance200ResponseAllOfOneOfInstanceTagsInner) bool {
	if tags == nil {
		return false
	}

	tagMap := make(map[string]string)
	for _, tag := range tags {
		if name, ok := tag.GetNameOk(); ok && name != nil {
			if value, ok := tag.GetValueOk(); ok && value != nil {
				tagMap[*name] = *value
			}
		}
	}

	for _, requiredTag := range requiredSweepTags {
		if value, exists := tagMap[requiredTag]; !exists || value != "true" {
			return false
		}
	}

	return true
}

// instanceSweeper handles cleanup of test instances
type instanceSweeper struct {
	client *sdk.APIClient
}

// NewInstanceSweeper creates and registers an instance sweeper
func NewInstanceSweeper(client *sdk.APIClient) {
	s := &instanceSweeper{
		client: client,
	}

	resource.AddTestSweepers(
		"hpe_morpheus_instance",
		&resource.Sweeper{
			Name: "hpe_morpheus_instance",
			F:    s.Sweep,
		})
}

// Sweep cleans up test instances
func (s *instanceSweeper) Sweep(_ string) error {
	ctx := context.Background()

	if s.client == nil {
		log.Printf("[INFO] No client provided, skipping instance sweep")

		return nil
	}

	instances, hresp, err := s.client.InstancesAPI.ListInstances(ctx).
		Phrase(testInstancePrefix).Execute()
	if err != nil || hresp.StatusCode != http.StatusOK {
		return fmt.Errorf("failed to list instances: %s", errfmt.ErrMsg(err, hresp))
	}

	instanceList := instances.GetInstances()
	var sweptCount int
	var sweepErrors []string

	for _, instance := range instanceList {
		name, ok := instance.GetNameOk()
		if !ok || name == nil {
			continue
		}

		if !strings.HasPrefix(*name, testInstancePrefix) {
			log.Printf("[INFO] Skipping instance (name): %s", *name)

			continue
		}

		id, ok := instance.GetIdOk()
		if !ok || id == nil {
			log.Printf("[INFO] Skipping instance (id): %s", *name)

			continue
		}

		// Get instance details to check tags
		instanceDetail, hresp, err := s.client.InstancesAPI.GetInstance(ctx, *id).Execute()
		if err != nil || hresp.StatusCode != http.StatusOK {
			log.Printf("[INFO] Skipping instance (failed to get details): %s", *name)

			continue
		}

		inst := instanceDetail.GetInstance()
		tags, ok := inst.GetTagsOk()
		if !ok || !hasRequiredTags(tags) {
			log.Printf("[INFO] Skipping instance (tags): %s", *name)

			continue
		}

		log.Printf("[INFO] Sweeping instance: %s (id: %d)", *name, *id)

		// Delete the instance
		_, hresp, err = s.client.InstancesAPI.DeleteInstance(ctx, *id).Execute()
		if err != nil || hresp.StatusCode != http.StatusOK {
			errMsg := fmt.Sprintf(
				"failed to delete instance %s (id: %d): %s",
				*name, *id, errfmt.ErrMsg(err, hresp),
			)
			log.Printf("[ERROR] %s", errMsg)
			sweepErrors = append(sweepErrors, errMsg)

			continue
		}

		sweptCount++
	}

	log.Printf(
		"[INFO] Instance sweep completed. Instances swept: %d, errors: %d",
		sweptCount, len(sweepErrors),
	)

	if len(sweepErrors) > 0 {
		return fmt.Errorf("%s", strings.Join(sweepErrors, "\n"))
	}

	return nil
}
