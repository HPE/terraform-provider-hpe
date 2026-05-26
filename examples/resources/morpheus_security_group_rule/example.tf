resource "hpe_morpheus_security_group_rule" "example" {
  security_group_id = 1
  protocol          = "tcp"
  rule_type         = "ingress"
}
