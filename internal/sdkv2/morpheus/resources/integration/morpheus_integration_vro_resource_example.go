// (C) Copyright 2025 Hewlett Packard Enterprise Development LP

package integration

//go:generate go run ../../../../../cmd/render -out examples/resources/morpheus_integration_vro/resource.tf morpheus_integration_vro_resource.tf.tmpl Name "tfexample vro" Enabled true Url https://myvro/vco/api Username my-vro-username Password my-vro-password AuthType basic Tenant vsphere.local
