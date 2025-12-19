// (C) Copyright 2025 Hewlett Packard Enterprise Development LP

package instancetypelayout

//go:generate go run ../../../../../../cmd/render -out examples/data-sources/morpheus_instance_type_layout/example-id.tf example-id.tf.tmpl Id 99
//go:generate go run ../../../../../../cmd/render -out examples/data-sources/morpheus_instance_type_layout/example-name.tf example-name.tf.tmpl Name "Example name"
//go:generate go run ../../../../../../cmd/render -out examples/data-sources/morpheus_instance_type_layout/example-name-version.tf example-name-version.tf.tmpl Name "Example name" Version "1.2.3"
