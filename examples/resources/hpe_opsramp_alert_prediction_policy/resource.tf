resource "hpe_opsramp_alert_prediction_policy" "default_alert_prediction_policy" {
  name = "Default Prediction Policy"

  enabled_mode = "OFF"
  filter_query = ""

  seasonality_time_frame    = "7D"
  generate_prediction_alert = true
}