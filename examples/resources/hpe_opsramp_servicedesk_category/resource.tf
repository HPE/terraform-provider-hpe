# Create multiple categories
resource "hpe_opsramp_servicedesk_category" "category1" {
  name        = "Category1"
  description = "Category1 Description"
  ticket_type = "serviceRequests"
}