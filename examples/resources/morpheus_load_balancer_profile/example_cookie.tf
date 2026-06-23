resource "hpe_morpheus_load_balancer_profile" "cookie" {
  load_balancer_id = 1
  name             = "Cookie Profile"
  service_type     = "LBCookiePersistenceProfile"

  config_cookie_persistence = {
    cookie_mode       = "INSERT"
    cookie_type       = "session"
    cookie_name       = "SERVERID"
    cookie_fallback   = true
    cookie_garbling   = true
    share_persistence = false
  }
}
