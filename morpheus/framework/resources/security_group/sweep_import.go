// (C) Copyright 2026 Hewlett Packard Enterprise Development LP

//go:build sweep

package security_group

import (
	_ "github.com/HPE/terraform-provider-hpe/morpheus/framework/resources/security_group/sweep"
	_ "github.com/HPE/terraform-provider-hpe/morpheus/framework/resources/security_group_rule/sweep"
)
