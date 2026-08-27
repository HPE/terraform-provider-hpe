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