# A Kubernetes worker on a VMware cloud, added via the free-form config map
# (the type code is the cloud's worker type from /api/server-types?creatable=true).
resource "hpe_morpheus_host" "k8s_worker" {
  name                = "k8s-worker-1"
  cloud_id            = 2
  group_id            = 1
  compute_server_type = "vmwareKubeWorker"
  plan_id             = 104

  config = {
    resourcePoolId = "pool-123"
  }
}
