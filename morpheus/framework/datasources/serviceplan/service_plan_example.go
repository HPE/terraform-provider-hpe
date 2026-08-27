// (C) Copyright 2025 Hewlett Packard Enterprise Development LP

package serviceplan

//go:generate ../../../../bin/render -out examples/data-sources/morpheus_service_plan/example-id.tf example-id.tf.tmpl Id 99
//go:generate ../../../../bin/render -out examples/data-sources/morpheus_service_plan/example-name-provision.tf example-name-provision.tf.tmpl Name "Example name" ProvisionTypeCode "arm"
//go:generate ../../../../bin/render -out examples/data-sources/morpheus_service_plan/example-name-provision-cloud.tf example-name-provision-cloud.tf.tmpl Name "Example name" ProvisionTypeCode "arm" CloudId 5
