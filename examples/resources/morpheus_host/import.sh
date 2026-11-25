# Hosts are imported by numeric ID. Create-only inputs the API does not return
# (the config_* blocks, network_interfaces and volumes) are not populated and
# would plan a replacement if configured, so leave them out of an imported
# host's configuration. delete_options may be set, and ssh_password_wo is
# updated in place (see the Import section for its digest form).
terraform import hpe_morpheus_host.example 123
