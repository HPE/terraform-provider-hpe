// (C) Copyright 2025 Hewlett Packard Enterprise Development LP

package template

//go:generate go run ../../../../../cmd/render -out examples/resources/morpheus_spec_template_kubernetes/resource_git.tf morpheus_spec_template_kubernetes_resource_git.tf.tmpl 'Name' 'tf-kubernetes-spec-example-git' 'SourceType' 'repository' 'RepositoryId' '2' 'VersionRef' 'main' 'SpecPath' './spec.yaml'

//go:generate go run ../../../../../cmd/render -out examples/resources/morpheus_spec_template_kubernetes/resource_local.tf morpheus_spec_template_kubernetes_resource_local.tf.tmpl 'Name' 'tf-terraform-spec-example-local' 'SourceType' 'local' 'SpecContent' "---\napiVersion: apps/v1\nkind: Deployment\nmetadata:\n  name: nginx-deployment\n  labels:\n    app: nginx\nspec:\n  replicas: 3\n  selector:\n    matchLabels:\n      app: nginx\n  template:\n    metadata:\n      labels:\n        app: nginx\n    spec:\n      containers:\n      - name: nginx\n        image: nginx:1.14.2\n        ports:\n        - containerPort: 80"

//go:generate go run ../../../../../cmd/render -out examples/resources/morpheus_spec_template_kubernetes/resource_url.tf morpheus_spec_template_kubernetes_resource_url.tf.tmpl 'Name' 'tf-kubernetes-spec-example-url' 'SourceType' 'url' 'SpecPath' 'http://example.com/spec.yaml'
