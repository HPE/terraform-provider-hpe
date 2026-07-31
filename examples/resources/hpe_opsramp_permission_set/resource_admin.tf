
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