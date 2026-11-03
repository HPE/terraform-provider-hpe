# Install a plugin from a local .jar file. Setting jar_sha256 with filesha256()
# verifies the jar before upload and re-uploads it when the file changes.
resource "hpe_morpheus_plugin" "from_path" {
  jar_path   = "${path.module}/plugins/my-plugin-1.0.0.jar"
  jar_sha256 = filesha256("${path.module}/plugins/my-plugin-1.0.0.jar")
  enabled    = true
}

# Install a plugin from the Morpheus Exchange. jar_sha256 is the Exchange
# "shaHash" for the version, which verifies the download before upload.
resource "hpe_morpheus_plugin" "from_exchange" {
  jar_url    = "https://share.morpheusdata.com/api/packages/pull?code=vault-plugin&version=1.4.2"
  jar_sha256 = "PUT_THE_EXCHANGE_shaHash_HERE"
  enabled    = true
}
