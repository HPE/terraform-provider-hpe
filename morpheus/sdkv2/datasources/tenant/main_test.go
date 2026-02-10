// (C) Copyright 2026 Hewlett Packard Enterprise Development LP

package tenant_test

import (
	"os"
	"testing"

	"github.com/HPE/terraform-provider-hpe/morpheus/testhelpers"
	"github.com/HPE/terraform-provider-hpe/morpheus/testhelpers/systemoverride"
)

func TestMain(m *testing.M) {
	systemoverride.ParseFlags()
	code := m.Run()

	testhelpers.WriteMergedResults()

	os.Exit(code)
}
