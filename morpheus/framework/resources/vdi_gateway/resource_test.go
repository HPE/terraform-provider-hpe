// (C) Copyright 2026 Hewlett Packard Enterprise Development LP

package vdi_gateway_test

import (
	"os"
	"testing"

	"github.com/HPE/terraform-provider-hpe/morpheus/testhelpers"
)

func TestMain(m *testing.M) {
	code := testhelpers.TestMain(m)
	testhelpers.WriteMergedResults()
	os.Exit(code)
}
