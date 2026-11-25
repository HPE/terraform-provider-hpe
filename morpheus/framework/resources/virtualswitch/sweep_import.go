// (C) Copyright 2026 Hewlett Packard Enterprise Development LP

//go:build sweep

package virtualswitch

import (
	_ "github.com/HPE/terraform-provider-hpe/morpheus/framework/resources/network/sweep"
	_ "github.com/HPE/terraform-provider-hpe/morpheus/framework/resources/virtualswitch/sweep"
)
