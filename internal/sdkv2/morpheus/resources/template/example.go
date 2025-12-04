// (C) Copyright 2025 Hewlett Packard Enterprise Development LP

package template

//go:generate go run ../../../../../cmd/render -out examples/resources/morpheus_file_template/resource.tf hpe_morpheus_file_template_resource.tf.tmpl Name tf-terraform-file-template Labels ["demo","template","terraform"] FileName tfcustom.cnf FilePath /etc/my.cnf.d Phase preProvision FileOwner root SettingName myCnf SettingCategory master
