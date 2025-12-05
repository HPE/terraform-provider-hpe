// (C) Copyright 2025 Hewlett Packard Enterprise Development LP

package task

//go:generate go run ../../../../../cmd/render -out examples/resources/morpheus_task_groovy_script/resource.tf task_groovy_script_resource.tf.tmpl Name tfexample_groovy_local Code tfexample_groovy_local SourceType local ScriptContent "println \"hello\"" Retryable true RetryCount 1 RetryDelaySeconds 10 AllowCustomConfig true
//go:generate go run ../../../../../cmd/render -out examples/resources/morpheus_task_groovy_script/resource_url.tf task_groovy_script_resource_url.tf.tmpl Name tfexample_groovy_url Code tfexample_groovy_url SourceType url ResultType json ScriptPath https://example.com/example.groovy Retryable true RetryCount 1 RetryDelaySeconds 10 AllowCustomConfig true
//go:generate go run ../../../../../cmd/render -out examples/resources/morpheus_task_groovy_script/resource_git.tf task_groovy_script_resource_git.tf.tmpl Name tfexample_groovy_git Code tfexample_groovy_git SourceType repository ResultType json ScriptPath example.groovy VersionRef master RepositoryId 1 Retryable true RetryCount 1 RetryDelaySeconds 10 AllowCustomConfig true
