# Create an admin user
resource "hpe_opsramp_user" "admin" {
  login_name = "testadmin"
  password   = "*******"
  first_name = "Admin"
  last_name  = "Istrator"
  email      = "email@example.com"

  time_zone = "Europe/Paris"
  country   = "Spain"

  roles = [
    hpe_opsramp_role.client_admin_role.unique_id
  ]

  # User notification preferences
  user_notifications = [
    {
      notify_type             = "Account Information"
      notify_method           = "Email"
      notify_input_type       = "Primary Email"
      notify_recurring_report = false
    },
    {
      notify_type             = "Alert Notification"
      notify_method           = "No Notify"
      notify_recurring_report = false
    },
    {
      notify_type             = "Report Notification"
      notify_method           = "No Notify"
      notify_recurring_report = false
    },
    {
      notify_type             = "Export Notification"
      notify_method           = "No Notify"
      notify_recurring_report = false
    },
    {
      notify_type             = "Login Activity Notification"
      notify_method           = "No Notify"
      notify_recurring_report = false
    }
  ]
}