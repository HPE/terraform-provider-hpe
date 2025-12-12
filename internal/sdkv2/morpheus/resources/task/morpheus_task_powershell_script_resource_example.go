// (C) Copyright 2025 Hewlett Packard Enterprise Development LP

package task

//go:generate go run ../../../../../cmd/render -out examples/resources/morpheus_task_powershell_script/resource.tf morpheus_task_powershell_script_resource.tf.tmpl Name tfexample_powershell_local Code tfexample_powershell_local Labels "\"demo\", \"terraform\"" SourceType local ScriptContent "Write-Output \"testing\"" ElevatedShell true Retryable true RetryCount 1 RetryDelaySeconds 10 AllowCustomConfig true
//go:generate go run ../../../../../cmd/render -out examples/resources/morpheus_task_powershell_script/resource_git.tf morpheus_task_powershell_script_resource_git.tf.tmpl Name tfexample_powershell_git Code tfexample_powershell_git Labels "\"demo\", \"terraform\"" SourceType repository ResultType json ScriptPath example.ps VersionRef master RepositoryId 1 ElevatedShell true Retryable true RetryCount 1 RetryDelaySeconds 10 AllowCustomConfig true
//go:generate go run ../../../../../cmd/render -out examples/resources/morpheus_task_powershell_script/resource_url.tf morpheus_task_powershell_script_resource_url.tf.tmpl Name tfexample_powershell_url Code tfexample_powershell_url Labels "\"demo\", \"terraform\"" SourceType url ResultType json ScriptPath https://example.com/example.ps ElevatedShell true Retryable true RetryCount 1 RetryDelaySeconds 10 AllowCustomConfig true
