resource "hpe_morpheus_app_blueprint_terraform" "tfapp_blueprint" {
  name              = "tfappbluedemo"
  description       = "testing terraform"
  category          = "terraformdemo"
  source_type       = "hcl"
  blueprint_content = <<EOF
...
EOF
  terraform_version = "1.1.1"
  terraform_options = "-var 'foo=bar'"
  tfvar_secret      = "tfvars/rdsdemo-secrets"
}
