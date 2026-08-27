resource "hpe_morpheus_load_balancer_profile" "client_ssl" {
  load_balancer_id = 1
  name             = "Client SSL Profile"
  service_type     = "LBClientSslProfile"

  config_client_ssl = {
    ssl_suite               = "CUSTOM"
    session_cache           = true
    session_cache_timeout   = 300
    prefer_server_cipher    = true
    supported_ssl_ciphers   = ["TLS_RSA_WITH_AES_128_GCM_SHA256"]
    supported_ssl_protocols = ["TLS_V1_2"]
  }
}
