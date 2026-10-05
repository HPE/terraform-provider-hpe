# The import ID is the licence's key id: the first eight characters of the key,
# as GET /api/license reports it under keyId. The key itself cannot be read back,
# so the first apply after the import sets it (re-installing an installed key is a no-op).
terraform import hpe_morpheus_license.tf_example_license 22324FEF
