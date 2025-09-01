variable "prefix" {
  default = "documentation"
}

data "hpe_morpheus_network" "example" {
  name = "${var.prefix}-example"
}
