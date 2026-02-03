// (C) Copyright 2025 Hewlett Packard Enterprise Development LP

package template

//go:generate ../../../../../bin/render -out examples/resources/morpheus_script_template/resource.tf morpheus_script_template_resource.tf.tmpl Name tf-terraform-script-template Labels "[\"demo\", \"template\", \"terraform\"]" ScriptType bash ScriptPhase provision ScriptContent "echo \"testing\"" RunAsUser root Sudo true
