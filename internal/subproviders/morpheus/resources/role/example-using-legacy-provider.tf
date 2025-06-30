data "morpheus_task" "example_legacy_task" {
  name = "example_task"
}

resource "hpe_morpheus_role" "example_with_legacy_provider" {
  name = "ExampleRoleWithLegacyProvider"
  description = "An example role using legacy provider"
  role_type = "user"
  permissions = jsonencode({
    "taskPermissions" : [
      {
        "id" = data.morpheus_task.example_legacy_task.id
        "access" = "full"
      }
    ]
    }
  )
}
