# Plugins can be imported by their numeric id, e.g.
terraform import hpe_morpheus_plugin.example 42

# After import, jar_path/jar_url and jar_sha256 are null (the API does not
# return the jar). Add one of jar_path or jar_url to the configuration before
# the next apply so Terraform can manage (and re-upload) the jar.
