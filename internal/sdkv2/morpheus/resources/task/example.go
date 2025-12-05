// (C) Copyright 2025 Hewlett Packard Enterprise Development LP

package task

//go:generate go run ../../../../../cmd/render -out examples/resources/morpheus_task_write_attributes/resource.tf task_write_attributes_resource.tf.tmpl Name tfexample_write_attributes Code tfexample_write_attributes Label1 demo Label2 terraform Attributes {\"demo\":\"test\"} Retryable true RetryCount 1 RetryDelaySeconds 10 AllowCustomConfig true
