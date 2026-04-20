// (C) Copyright 2025-2026 Hewlett Packard Enterprise Development LP

package sweep

import (
	"testing"

	"github.com/hashicorp/terraform-plugin-testing/helper/resource"

	_ "github.com/HPE/terraform-provider-hpe/morpheus"
	"github.com/HPE/terraform-provider-hpe/morpheus/testhelpers/systemoverride"
)

func TestMain(m *testing.M) {
	systemoverride.ParseFlags()
	resource.TestMain(m)
}
