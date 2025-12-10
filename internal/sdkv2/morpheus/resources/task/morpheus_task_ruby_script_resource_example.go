// (C) Copyright 2025 Hewlett Packard Enterprise Development LP

package task

//go:generate go run ../../../../../cmd/render -out examples/resources/morpheus_task_ruby_script/resource.tf task_ruby_script_resource.tf.tmpl Name tfexample_ruby_local Code tfexample_ruby_local Labels "\"demo\", \"terraform\"" SourceType local ScriptContent "puts \"testing\"" Retryable true RetryCount 1 RetryDelaySeconds 10 AllowCustomConfig true
//go:generate go run ../../../../../cmd/render -out examples/resources/morpheus_task_ruby_script/resource_git.tf task_ruby_script_resource_git.tf.tmpl Name tfexample_ruby_git Code tfexample_ruby_git Labels "\"demo\", \"terraform\"" SourceType repository ResultType json ScriptPath example.rb VersionRef master RepositoryId 1 Retryable true RetryCount 1 RetryDelaySeconds 10 AllowCustomConfig true
//go:generate go run ../../../../../cmd/render -out examples/resources/morpheus_task_ruby_script/resource_url.tf task_ruby_script_resource_url.tf.tmpl Name tfexample_ruby_url Code tfexample_ruby_url Labels "\"demo\", \"terraform\"" SourceType url ResultType json ScriptPath https://example.com/example.rb Retryable true RetryCount 1 RetryDelaySeconds 10 AllowCustomConfig true
