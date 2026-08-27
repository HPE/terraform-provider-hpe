resource "hpe_morpheus_form" "example" {
  name        = "demo"
  code        = "demo"
  description = "demo"
  labels      = ["terraform", "demo"]

  option_type {
    name                          = "tf disk manager example"
    code                          = "disk-manager-input"
    description                   = "Terraform disk manager example"
    type                          = "diskManager"
    field_label                   = "disk manager input"
    field_name                    = "diskManagerInput"
    default_value                 = jsonencode([{ rootVolume = true, name = "root", size = 10, sizeBytes = 10737418240, minStorage = 0, displayOrder = 0, storageType = 1, datastoreId = "52" }, { rootVolume = false, name = "data-1", size = 20, sizeBytes = 21474836480, minStorage = 0, displayOrder = 1, datastoreId = "autoCluster", storageType = 1 }])
    help_block                    = "Configure disks"
    required                      = true
    export_meta                   = true
    display_value_on_details      = true
    locked                        = true
    hidden                        = false
    exclude_from_search           = true
    group_field_type              = "value"
    group_id                      = "1"
    cloud_field_type              = "value"
    cloud_id                      = "1"
    plan_field_type               = "value"
    plan_id                       = "1"
    layout_field_type             = "value"
    layout_id                     = "1"
    pool_field_type               = "value"
    pool_id                       = "1"
    virtual_image_field_type      = "value"
    image_id                      = "1"
    enable_disk_type_selection    = true
    enable_storage_type_selection = true
    enable_datastore_selection    = true
  }
}
