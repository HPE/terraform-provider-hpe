data "hpe_morpheus_storage_controller_type" "example" {
  controller_name  = "SCSI VMware Paravirtual"
  bus_number       = 1
  interface_number = 0
}
