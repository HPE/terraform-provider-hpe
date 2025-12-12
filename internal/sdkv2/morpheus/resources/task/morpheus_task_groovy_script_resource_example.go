// (C) Copyright 2025 Hewlett Packard Enterprise Development LP

package task

//go:generate sh -c "go run ../../../../../cmd/render -out examples/resources/morpheus_task_groovy_script/resource.tf morpheus_task_groovy_script_resource.tf.tmpl Name tfexample_groovy_local Code tfexample_groovy_local Labels '[\"demo\",\"terraform\"]' SourceType local ScriptContent 'println \"hello\"' Retryable true RetryCount 1 RetryDelaySeconds 10 AllowCustomConfig true"
//go:generate sh -c "go run ../../../../../cmd/render -out examples/resources/morpheus_task_groovy_script/resource_git.tf morpheus_task_groovy_script_resource_git.tf.tmpl Name tfexample_groovy_git Code tfexample_groovy_git Labels '[\"demo\",\"terraform\"]' SourceType repository ResultType json ScriptPath example.groovy VersionRef master RepositoryId 1 Retryable true RetryCount 1 RetryDelaySeconds 10 AllowCustomConfig true"
//go:generate sh -c "go run ../../../../../cmd/render -out examples/resources/morpheus_task_groovy_script/resource_url.tf morpheus_task_groovy_script_resource_url.tf.tmpl Name tfexample_groovy_url Code tfexample_groovy_url Labels '[\"demo\",\"terraform\"]' SourceType url ResultType json ScriptPath https://example.com/example.groovy Retryable true RetryCount 1 RetryDelaySeconds 10 AllowCustomConfig true"
