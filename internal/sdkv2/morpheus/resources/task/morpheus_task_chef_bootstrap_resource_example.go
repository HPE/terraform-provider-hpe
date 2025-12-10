// (C) Copyright 2025 Hewlett Packard Enterprise Development LP

package task

//go:generate go run ../../../../../cmd/render -out examples/resources/morpheus_task_chef_bootstrap/resource.tf task_chef_bootstrap_resource.tf.tmpl Name terraform_example_chef Code terraform_example_chef Labels "\"demo\", \"terraform\"" ChefServerId 9 Environment dev RunList role[web] DataBagKey test123 DataBagKeyPath /etc/chef/databag_secret NodeName demonode NodeAttributes "\"test\":\"demo\"" Retryable true RetryCount 1 RetryDelaySeconds 10 AllowCustomConfig true Visibility public
