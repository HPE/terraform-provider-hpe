resource "hpe_opsramp_first_response_policy" "seasonality_first_response_policy" {
  name = "Seasonal-based Suppression"

  enabled_mode = "OBSERVED"
  filter_query = ""

  pattern_actions = {
    seasonality_time_frame = "7D"
    suppress = {
      seasonal_alerts = true
    }
  }
}