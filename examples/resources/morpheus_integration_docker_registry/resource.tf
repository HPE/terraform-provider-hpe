resource "hpe_morpheus_integration_docker_registry" "tf_example_docker_registry_integration" {
  name     = "tfexampledockerregistry"
  enabled  = true
  url      = "https://index.docker.io/v1/"
  username = "admin"
  password = "password123"
}
