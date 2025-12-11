resource "hpe_morpheus_wiki_page" "tfexample_wiki_page" {
  name     = "tfexample_wiki_page"
  category = "morpheus-terraform"
  content  = <<EOF
# Terraform Example

This is an example of using the Morpheus terraform provider.
EOF
}
