// (C) Copyright 2025 Hewlett Packard Enterprise Development LP

package task

//go:generate go run ../../../../../cmd/render -out examples/resources/morpheus_task_javascript/resource.tf task_javascript_resource.tf.tmpl Name tfexample_javascript Code tfexample_javascript Labels ["demo","terraform"] ScriptContent console.log("testing") Retryable true RetryCount 1 RetryDelaySeconds 10 AllowCustomConfig true
