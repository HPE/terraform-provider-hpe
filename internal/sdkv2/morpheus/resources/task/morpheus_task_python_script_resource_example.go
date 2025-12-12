// (C) Copyright 2025 Hewlett Packard Enterprise Development LP

package task

//go:generate go run ../../../../../cmd/render -out examples/resources/morpheus_task_python_script/resource.tf morpheus_task_python_script_resource.tf.tmpl Name tfexample_python_local Code tfexample_python_local Labels "[\"demo\", \"terraform\"]" SourceType local ScriptContent "print('morpheus')\nprint('python')" CommandArguments example AdditionalPackages pyyaml PythonBinary /usr/bin/python3 Retryable true RetryCount 1 RetryDelaySeconds 10 AllowCustomConfig true

//go:generate go run ../../../../../cmd/render -out examples/resources/morpheus_task_python_script/resource_url.tf morpheus_task_python_script_resource_url.tf.tmpl Name tfexample_python_url Code tfexample_python_url Labels "[\"demo\", \"terraform\"]" SourceType url ResultType json ScriptPath https://example.com/example.py CommandArguments example AdditionalPackages pyyaml PythonBinary /usr/bin/python3 Retryable true RetryCount 1 RetryDelaySeconds 10 AllowCustomConfig true

//go:generate go run ../../../../../cmd/render -out examples/resources/morpheus_task_python_script/resource_git.tf morpheus_task_python_script_resource_git.tf.tmpl Name tfexample_python_git Code tfexample_python_git Labels "[\"demo\", \"terraform\"]" SourceType repository ResultType json ScriptPath example.py VersionRef master RepositoryId 1 CommandArguments example AdditionalPackages pyyaml PythonBinary /usr/bin/python3 Retryable true RetryCount 1 RetryDelaySeconds 10 AllowCustomConfig true
