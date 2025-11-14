resource "hpe_morpheus_script_boot" "tf_example_boot_script" {
  name    = "TF Example Boot Script"
  content = "ls"
}