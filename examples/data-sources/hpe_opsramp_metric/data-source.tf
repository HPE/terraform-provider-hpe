data "hpe_opsramp_metric" "example" {
  metric_name = "haproxy.denied_resp_rate"
  scope       = "GLOBAL"
}