data "hpe_opsramp_template" "example" {
  name  = "Active Directory Performance Counters DotNet v4"
  scope = "GLOBAL"
}

output "template_id" {
  value = data.hpe_opsramp_template.example.id
}
