// (C) Copyright 2025 Hewlett Packard Enterprise Development LP

package integration

//go:generate go run ../../../../../cmd/render -out examples/resources/morpheus_integration_servicenow/resource.tf morpheus_integration_servicenow_resource.tf.tmpl Name "terraform servicenow integration" Enabled true Url "https://servicenowprod.service-now.com" Username "my-snow-username" Password "my-snow-password" DefaultCmdbBusinessClass "demo"
