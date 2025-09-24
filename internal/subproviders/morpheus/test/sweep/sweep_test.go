// (C) Copyright 2025 Hewlett Packard Enterprise Development LP

package sweep

import (
	"testing"

	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
)

func init() {
	Networks()
	Users()
}

func TestMain(m *testing.M) {
	resource.TestMain(m)
}
