// (C) Copyright 2025 Hewlett Packard Enterprise Development LP

package catalogitem

//go:generate go run ../../../../../cmd/render -out examples/resources/morpheus_catalog_item_instance/resource.tf morpheus_catalog_item_instance_resource.tf.tmpl Name tfexample_instance_catalog Description "terraform example instance catalog item" ImagePath tfexample.png ImageName tfexample.png Enabled true Featured true Content "{\"name\":\"test\"}" Config "{\"name\":\"test\"}" Visibility private
