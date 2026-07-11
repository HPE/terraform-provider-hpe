resource "hpe_opsramp_integration_app" "kubernetes_2_0_integration_app" {
  application  = "Kubernetes-2.0"
  version = "2.3.0"
  bypass_resource_reconciliation = true
}

resource "hpe_opsramp_integration_config" "kubernetes_2_0_integration_config" {
  integration_id = hpe_opsramp_integration_app.kubernetes_2_0_integration_app.id
  
  name = "test"
  config = jsonencode({"Etcd":true,"coreDNS":true,"KubeProxy":true,"enableLog":true,"enableEBPF":true,"kubeEvents":true,"enableTrace":true,"KubeletStats":true,"KubeAPIServer":true,"KubeScheduler":true,"ebpfFeatureFlag":false,"KubeClusterReceiver":true,"KubeControllerManager":true,"clientLevelLogPermission":false,"clientLevelTracePermission":false,"replicaCount":3,"DistributionType":"K8S","eBPFContainerEngine":"cri-o"})
  all_resources = false
}