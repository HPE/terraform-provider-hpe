resource "hpe_morpheus_role" "example" {
	name = "ExampleRole"
	multitenant = false
        description = "An example role"
        role_type = "user"
}
