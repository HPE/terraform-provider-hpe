resource "hpe_morpheus_security_group_rule" "example" {
  security_group_id = 1
  name              = "Allow HTTPS"
  protocol          = "tcp"
  rule_type         = "customRule"
  direction         = "ingress"
  port_range        = "443"
  source            = "0.0.0.0/0"
  destination       = "0.0.0.0/0"
  policy            = "accept"
}
