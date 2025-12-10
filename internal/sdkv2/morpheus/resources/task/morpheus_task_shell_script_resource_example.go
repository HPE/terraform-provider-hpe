// (C) Copyright 2025 Hewlett Packard Enterprise Development LP

package task

//go:generate go run ../../../../../cmd/render -out examples/resources/morpheus_task_shell_script/resource.tf task_shell_script_resource.tf.tmpl Name tfexample_shell_local Code tfexample_shell_local Labels "\"demo\", \"terraform\"" SourceType local ScriptContent "  echo \"testing\"" Sudo true Retryable true RetryCount 1 RetryDelaySeconds 10 AllowCustomConfig true
//go:generate go run ../../../../../cmd/render -out examples/resources/morpheus_task_shell_script/resource_git.tf task_shell_script_resource_git.tf.tmpl Name tfexample_shell_git Code tfexample_shell_git Labels "\"demo\", \"terraform\"" SourceType repository ResultType json ScriptPath example.sh VersionRef master RepositoryId 1 Sudo true Retryable true RetryCount 1 RetryDelaySeconds 10 AllowCustomConfig true
//go:generate go run ../../../../../cmd/render -out examples/resources/morpheus_task_shell_script/resource_url.tf task_shell_script_resource_url.tf.tmpl Name tfexample_shell_url Code tfexample_shell_url Labels "\"demo\", \"terraform\"" SourceType url ResultType json ScriptPath https://example.com/example.sh Sudo true Retryable true RetryCount 1 RetryDelaySeconds 10 AllowCustomConfig true
