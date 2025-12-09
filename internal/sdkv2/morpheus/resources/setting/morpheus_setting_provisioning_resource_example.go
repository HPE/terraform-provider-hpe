// (C) Copyright 2025 Hewlett Packard Enterprise Development LP

package setting

//go:generate go run ../../../../../cmd/render -out examples/resources/morpheus_setting_provisioning/resource.tf morpheus_setting_provisioning_resource.tf.tmpl AllowZoneSelection false AllowHostSelection false RequireEnvironments false ShowPricing true HideDatastoreStats true CrossTenantNamingPolicies false CloudinitUsername cloudinit CloudinitPassword Pa55w0rd! WindowsPassword Pa55w0rd! PxeRootPassword Pa55w0rd!
