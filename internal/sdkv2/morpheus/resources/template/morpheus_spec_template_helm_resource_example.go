// (C) Copyright 2025 Hewlett Packard Enterprise Development LP

package template

//go:generate go run ../../../../../cmd/render -out examples/resources/morpheus_spec_template_helm/resource_url.tf spec_template_helm_resource_url.tf.tmpl Name tf-helm-spec-example-url SourceType url SpecPath http://example.com/chart.yaml
//go:generate go run ../../../../../cmd/render -out examples/resources/morpheus_spec_template_helm/resource_local.tf spec_template_helm_resource_local.tf.tmpl Name tf-helm-spec-example-local SourceType local SpecContent "apiVersion: v1\nkind: Service\nmetadata:\nname: {{ template \"fullname\" . }}\nlabels:\n    chart: \"{{ .Chart.Name }}-{{ .Chart.Version | replace \"+\" \"_\" }}\"\nspec:\ntype: {{ .Values.service.type }}\nports:\n- port: {{ .Values.service.externalPort }}\n    targetPort: {{ .Values.service.internalPort }}\n    protocol: TCP\n    name: {{ .Values.service.name }}\nselector:\n    app: {{ template \"fullname\" . }}"
//go:generate go run ../../../../../cmd/render -out examples/resources/morpheus_spec_template_helm/resource_git.tf spec_template_helm_resource_git.tf.tmpl Name tf-helm-spec-example-git SourceType repository RepositoryId 2 VersionRef main SpecPath ./spec.yaml
