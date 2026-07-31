
terraform {
  required_providers {
    hpe = {
      source  = "HPE/hpe"
      version = ">= 1.6.0"
    }
  }
}

provider "hpe" {
  opsramp {
    client_id     = "abcdefghijklmnopqrstuvwxyz123456"
    client_secret = "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ1234567890ab"
    endpoint      = "tenant.api.pov.opsramp.com"
    tenant        = "abcdefgh-1234-5678-90ab-cdefghijklmn"
  }
}

# Permission Sets
resource "hpe_opsramp_permission_set" "client_admin_perms" {
  name        = "Admin Permission Set"
  description = "Full administrative access, yes!"

  permissions = [
    {
      "name" : "Alerts",
      "type" : "Manage Alerts "
    },
    {
      "name" : "Administration",
      "type" : "Administration"
    },
    {
      "name" : "Reports",
      "type" : "Manage Reports"
    },
    {
      "name" : "Users",
      "type" : "View Users "
    },
    {
      "name" : "Roles",
      "type" : "View Roles "
    },
    {
      "name" : "Credentials",
      "type" : "Credentials View"
    },
    {
      "name" : "Dashboards",
      "type" : "Manage Dashboard"
    },
    {
      "name" : "Scheduled Maintenance",
      "type" : "Manage Scheduled Maintenance"
    },
    {
      "name" : "Metrics",
      "type" : "Metrics Manage"
    },
    {
      "name" : "Devices",
      "type" : "Manage Device "
    },
    {
      "name" : "Custom Attributes",
      "type" : "Custom Attributes Manage"
    },
    {
      "name" : "My Profile",
      "type" : "My Profile Edit"
    },
    {
      "name" : "Copilot",
      "type" : "Copilot View"
    },
    {
      "name" : "Gateway Firmware",
      "type" : "Allow Gateway Firmware Update"
    },
    {
      "name" : "Manage Management Profile ",
      "type" : "Manage Management Profile "
    },
    {
      "name" : "Commands",
      "type" : "Allow To Run Commands"
    },
    {
      "name" : "Monitors",
      "type" : "Manage Monitors "
    },
    {
      "name" : "Integration",
      "type" : "Manage Integration"
    },
    {
      "name" : "Device Monitor Template Configuration",
      "type" : "Customize Monitor Templates"
    },
    {
      "name" : "OpsQ",
      "type" : "OpsQ Manage"
    },
    {
      "name" : "Incident",
      "type" : "Manage Incident"
    },
    {
      "name" : "Change Request",
      "type" : "Manage Change Request"
    },
    {
      "name" : "Problem",
      "type" : "Manage Problem"
    },
    {
      "name" : "Projects",
      "type" : "Projects Manage"
    },
    {
      "name" : "Service Desk",
      "type" : "Manage Service Desk"
    },
    {
      "name" : "Service Request",
      "type" : "Manage Service Request"
    },
    {
      "name" : "Task Request",
      "type" : "Manage Task Request"
    },
    {
      "name" : "Time-Bound Request",
      "type" : "Manage Time-Bound Request"
    },
    {
      "name" : "Knowledge Base",
      "type" : "Manage Knowledge Base "
    },
    {
      "name" : "Jobs",
      "type" : "Manage Jobs "
    },
    {
      "name" : "Patch approvals",
      "type" : "Manage Patch Approval "
    },
    {
      "name" : "Recordings audit",
      "type" : "All Recordings: Play, Search, Edit"
    },
    {
      "name" : "Scripts",
      "type" : "Manage Scripts "
    },
    {
      "name" : "Network Configuration Management",
      "type" : "Approve NCM"
    },
    {
      "name" : "Network Performance Management",
      "type" : "NPM Manage"
    }
  ]
}

# Permission Sets
resource "hpe_opsramp_permission_set" "client_view_perms" {
  name        = "View-only Permission Set"
  description = "View-only access for all modules, no changes allowed"

  permissions = [
    {
      "name" : "Custom Attributes",
      "type" : "Custom Attributes View"
    },
    {
      "name" : "Patch approvals",
      "type" : "View Patch Approval "
    },
    {
      "name" : "Recordings audit",
      "type" : "Play, Search All Recordings "
    },
    {
      "name" : "Projects",
      "type" : "Projects View"
    },
    {
      "name" : "Scheduled Maintenance",
      "type" : "View Scheduled Maintenance"
    },
    {
      "name" : "Network Performance Management",
      "type" : "NPM View"
    },
    {
      "name" : "Dashboards",
      "type" : "View Dashboard"
    },
    {
      "name" : "Roles",
      "type" : "View Roles "
    },
    {
      "name" : "Jobs",
      "type" : "View Jobs "
    },
    {
      "name" : "Users",
      "type" : "View Users "
    },
    {
      "name" : "Device Monitor Template Configuration",
      "type" : "Apply Monitor Templates"
    },
    {
      "name" : "Incident",
      "type" : "View Incident"
    },
    {
      "name" : "Alerts",
      "type" : "View Alerts "
    },
    {
      "name" : "Credentials",
      "type" : "Credentials View"
    },
    {
      "name" : "Change Request",
      "type" : "View Change Request"
    },
    {
      "name" : "Integration",
      "type" : "View Integration"
    },
    {
      "name" : "OpsQ",
      "type" : "OpsQ View"
    },
    {
      "name" : "Knowledge Base",
      "type" : "View Knowledge Base "
    },
    {
      "name" : "Copilot",
      "type" : "Copilot View"
    },
    {
      "name" : "Devices",
      "type" : "View Device "
    },
    {
      "name" : "Time-Bound Request",
      "type" : "View Time-Bound Request"
    },
    {
      "name" : "Service Request",
      "type" : "View Service Request"
    },
    {
      "name" : "Problem",
      "type" : "View Problem"
    },
    {
      "name" : "Network Configuration Management",
      "type" : "View NCM"
    },
    {
      "name" : "My Profile",
      "type" : "My Profile Edit"
    },
    {
      "name" : "Manage Management Profile ",
      "type" : "View Management Profile "
    },
    {
      "name" : "Commands",
      "type" : "Allow To Run Commands"
    },
    {
      "name" : "Task Request",
      "type" : "View Task Request"
    },
    {
      "name" : "Monitors",
      "type" : "Monitors View"
    },
    {
      "name" : "Scripts",
      "type" : "View Scripts "
    },
    {
      "name" : "Metrics",
      "type" : "Metrics Manage"
    },
    {
      "name" : "Gateway Firmware",
      "type" : "Allow Gateway Firmware Update"
    },
    {
      "name" : "Reports",
      "type" : "View Reports "
    },
    {
      "name" : "Service Desk",
      "type" : "View Service Desk"
    }
  ]
}


# Roles
resource "hpe_opsramp_role" "client_admin_role" {
  name        = "Admin Role"
  description = "Administrative role with full permissions"

  permissions = [
    hpe_opsramp_permission_set.client_admin_perms.unique_id
  ]
}

resource "hpe_opsramp_role" "client_view_role" {
  name        = "View Only Role"
  description = "View-only role with permissions to view all modules but not make changes"

  permissions = [
    hpe_opsramp_permission_set.client_view_perms.unique_id
  ]
}

# Create an admin user
resource "hpe_opsramp_user" "admin" {
  login_name = "testadmin"
  password   = var.password
  first_name = "Admin"
  last_name  = "Istrator"
  email      = var.email

  alt_email     = ""
  phone_number  = ""
  mobile_number = ""
  time_zone     = "Europe/Paris"
  designation   = "Observability Advance Operations, Administrator"
  address       = ""
  city          = "Valencia"
  state         = ""
  zip           = ""
  country       = "Spain"

  roles = [
  ]

  # User notification preferences
  user_notifications = [{
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
  }]

  change_password = false
}

# User Group for Client Admins
resource "hpe_opsramp_user_group" "client_admin_group" {
  name        = "Admin User Group"
  description = "User group for client administrators"

  roles = [
    hpe_opsramp_role.client_admin_role.id
  ]

  users = [
    hpe_opsramp_user.admin.id
  ]
}