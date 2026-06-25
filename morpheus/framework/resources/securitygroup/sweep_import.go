// (C) Copyright 2026 Hewlett Packard Enterprise Development LP

//go:build sweep

package securitygroup

import (
	_ "github.com/HPE/terraform-provider-hpe/morpheus/framework/resources/securitygroup/sweep"
	_ "github.com/HPE/terraform-provider-hpe/morpheus/framework/resources/securitygrouprule/sweep"
)
