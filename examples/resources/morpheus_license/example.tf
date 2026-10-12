resource "hpe_morpheus_license" "example" {
  # key is write-only: it is sent to the appliance and never stored in state.
  key = "22324FEF3WMCDMMSWE"
}
