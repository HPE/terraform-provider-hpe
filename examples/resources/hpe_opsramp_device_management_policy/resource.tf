resource "hpe_opsramp_device_management_policy" "test_policy" {
  name         = "mysql-servers"
  search_query = "name CONTAINS \"mysql\""

  actions = [
    {
      action = "ASSIGN MONITORINGTEMPLATE"
      items = [
        {
          item = "4053"
          parameters = {
            "Packet Count"    = "3"
            "Packet interval" = "1000"
          }
        }
      ]
    },
    {
      action = "Assign Availability Rule"
      items = [
        {
          item = "AND"
        }
      ]
    }
  ]
}
