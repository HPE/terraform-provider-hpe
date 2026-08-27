resource "hpe_morpheus_app_blueprint_arm" "example" {
  name               = "example_app_arm_blueprint_json"
  description        = "example arm app blueprint"
  category           = "armtemplates"
  source_type        = "json"
  install_agent      = true
  cloud_init_enabled = true
  os_type            = "linux"
  blueprint_content  = <<EOF
{
 "type": "Microsoft.Storage/storageAccounts",
 "apiVersion": "2019-04-01",
 "name": "string",
 "location": "string",
 "tags": {
 "tagName1": "tagValue1",
 "tagName2": "tagValue2"
 },
 "sku": {
 "name": "string",
 "restrictions": [
 {
 "reasonCode": "string"
 }
 ]
 },
 "kind": "string",
 "identity": {
 "type": "SystemAssigned"
 },
 "properties": {
 "accessTier": "string",
 "allowBlobPublicAccess": "bool",
 "allowSharedKeyAccess": "bool",
 "azureFilesIdentityBasedAuthentication": {
 "activeDirectoryProperties": {
 "azureStorageSid": "string",
 "domainGuid": "string",
 "domainName": "string",
 "domainSid": "string",
 "forestName": "string",
 "netBiosDomainName": "string"
 },
 "directoryServiceOptions": "string"
 },
 "customDomain": {
 "name": "string",
 "useSubDomainName": "bool"
 },
 "encryption": {
 "keySource": "string",
 "keyvaultproperties": {
 "keyname": "string",
 "keyvaulturi": "string",
 "keyversion": "string"
 },
 "services": {
 "blob": {
 "enabled": "bool"
 },
 "file": {
 "enabled": "bool"
 }
 }
 },
 "isHnsEnabled": "bool",
 "largeFileSharesState": "string",
 "minimumTlsVersion": "string",
 "networkAcls": {
 "bypass": "string",
 "defaultAction": "string",
 "ipRules": [
 {
 "action": "Allow",
 "value": "string"
 }
 ],
 "virtualNetworkRules": [
 {
 "action": "Allow",
 "id": "string",
 "state": "string"
 }
 ]
 },
 "supportsHttpsTrafficOnly": "bool"
 }
}
EOF
}
