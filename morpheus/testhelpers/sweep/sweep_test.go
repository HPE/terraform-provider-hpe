// (C) Copyright 2025-2026 Hewlett Packard Enterprise Development LP

package sweep_test

import (
	"testing"

	"github.com/hashicorp/terraform-plugin-testing/helper/resource"

	_ "github.com/HPE/terraform-provider-hpe/morpheus"
)

func TestMain(m *testing.M) {
	resource.TestMain(m)
}
