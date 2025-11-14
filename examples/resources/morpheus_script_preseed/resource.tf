resource "hpe_morpheus_script_preseed" "tf_example_preseed_script" {
  name    = "TF Example Preseed Script"
  content = "ls"
}