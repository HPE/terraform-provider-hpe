resource "hpe_morpheus_monitoring_contact" "example" {
  name          = "Ops Team"
  email_address = "ops-team@example.com"
  sms_address   = "+15551234567"
}
