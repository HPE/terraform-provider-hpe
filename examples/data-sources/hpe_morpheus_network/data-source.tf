variable "prefix" {
  default = "documentation"
}

resource "hpe_morpheus_network" "example" {
  name = "${var.prefix}-example"
}
