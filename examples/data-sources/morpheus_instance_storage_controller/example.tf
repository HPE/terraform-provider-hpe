data "hpe_morpheus_instance_storage_controller" "example" {
  controller_name  = "SCSI VMware Paravirtual"
  bus_number       = 1
  interface_number = 0
}
