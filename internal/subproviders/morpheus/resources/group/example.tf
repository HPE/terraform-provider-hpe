resource "hpe_morpheus_group" "example" {
  name = "TestGroup"
  location = "here"
  code = "aCode"
  labels = ["aLabel1", "aLabel2"]
}
