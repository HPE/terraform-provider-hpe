# The HVM cluster layout 2.0 of a Morpheus 9.1.0 appliance: the layout
# version selects the virtual-switch generation, and the code prefix
# distinguishes the HVM cluster layouts from the Kubernetes-on-HVM ones.
# Values are regular expressions, so the version is anchored: "2.0" alone
# would also match "12.0" or "2.0.1".
data "hpe_morpheus_cluster_layouts" "hvm_2" {
  filter {
    name   = "cluster_version"
    values = ["^2\\.0$"]
  }
  filter {
    name   = "code"
    values = ["^mvm-"]
  }
}
