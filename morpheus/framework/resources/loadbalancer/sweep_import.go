// (C) Copyright 2026 Hewlett Packard Enterprise Development LP

//go:build sweep

package loadbalancer

import (
	_ "github.com/HPE/terraform-provider-hpe/morpheus/framework/resources/loadbalancer/sweep"
	_ "github.com/HPE/terraform-provider-hpe/morpheus/framework/resources/loadbalancermonitor/sweep"
	_ "github.com/HPE/terraform-provider-hpe/morpheus/framework/resources/loadbalancervirtualserver/sweep"
)
