// (C) Copyright 2025 Hewlett Packard Enterprise Development LP

package task

//go:generate go run ../../../../../cmd/render -out examples/resources/morpheus_task_restart/resource.tf task_restart_resource.tf.tmpl Name tfexample_restart Code tfexample_restart Labels ["demo","terraform"] Retryable true RetryCount 1 RetryDelaySeconds 10 AllowCustomConfig true
