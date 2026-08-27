#!/bin/sh
# Helper script to generate per-type form resource examples

RENDER=../../../../bin/render

$RENDER \
  -out examples/resources/morpheus_form/resource_key_value.tf \
  form_key_value.tf.tmpl \
  Name 'demo' \
  Code 'demo' \
  Description 'demo' \
  Labels '["terraform", "demo"]' \
  OptionTypeName 'tf keyValue example' \
  OptionTypeCode 'keyValue-input' \
  OptionTypeDescription 'Terraform keyValue example' \
  OptionTypeType 'keyValue' \
  OptionTypeFieldLabel 'KeyValue' \
  OptionTypeFieldName 'keyValue' \
  OptionTypeDefaultValue 'jsonencode([{ key = "a", value = "b" }, { key = "c", value = "d" }])' \
  OptionTypeHelpBlock 'Select a key-value pair' \
  OptionTypeRequired 'true' \
  OptionTypeExportMeta 'true' \
  OptionTypeDisplayValueOnDetails 'true' \
  OptionTypeLocked 'true' \
  OptionTypeHidden 'false' \
  OptionTypeExcludeFromSearch 'true' \
  OptionTypeConvertToObject 'true' \
  OptionTypeKeyPlaceholder 'Key123' \
  OptionTypeValuePlaceholder 'Value123'

$RENDER \
  -out examples/resources/morpheus_form/resource_virtual_image.tf \
  form_virtual_image.tf.tmpl \
  Name 'demo' \
  Code 'demo' \
  Description 'demo' \
  Labels '["terraform", "demo"]' \
  OptionTypeName 'tf virtual-image example' \
  OptionTypeCode 'virtual-image' \
  OptionTypeDescription 'Terraform virtual-image example' \
  OptionTypeType 'virtual-image' \
  OptionTypeFieldLabel 'Virtual Image' \
  OptionTypeFieldName 'virtual-image' \
  OptionTypeDefaultValue '' \
  OptionTypeHelpBlock 'Select a virtual image' \
  OptionTypeCloudFieldType 'id' \
  OptionTypeCloudId '1' \
  OptionTypeRequired 'true' \
  OptionTypeExportMeta 'true' \
  OptionTypeDisplayValueOnDetails 'true' \
  OptionTypeLocked 'true' \
  OptionTypeHidden 'false' \
  OptionTypeExcludeFromSearch 'true'

$RENDER \
  -out examples/resources/morpheus_form/resource_vmw_folders.tf \
  form_vmw_folders.tf.tmpl \
  Name 'demo' \
  Code 'demo' \
  Description 'demo' \
  Labels '["terraform", "demo"]' \
  OptionTypeName 'tf vmwFolders example' \
  OptionTypeCode 'vmw-folders-input' \
  OptionTypeDescription 'Terraform vmwFolders example' \
  OptionTypeType 'vmwFolders' \
  OptionTypeFieldLabel 'VmwFolders' \
  OptionTypeFieldName 'vmwFolders' \
  OptionTypeDefaultValue '' \
  OptionTypeHelpBlock 'Select a vmwFolder' \
  OptionTypeGroupFieldType 'value' \
  OptionTypeGroupId '1' \
  OptionTypeCloudFieldType 'value' \
  OptionTypeCloudId '1' \
  OptionTypePlanFieldType 'value' \
  OptionTypePlanId '1' \
  OptionTypeRequired 'true' \
  OptionTypeExportMeta 'true' \
  OptionTypeDisplayValueOnDetails 'true' \
  OptionTypeLocked 'true' \
  OptionTypeHidden 'false' \
  OptionTypeExcludeFromSearch 'true'


$RENDER \
  -out examples/resources/morpheus_form/resource_file_content.tf \
  form_file_content.tf.tmpl \
  Name 'demo' \
  Code 'demo' \
  Description 'demo' \
  Labels '["terraform", "demo"]' \
  OptionTypeName 'tf fileContent example' \
  OptionTypeCode 'fileContent' \
  OptionTypeDescription 'Terraform fileContent example' \
  OptionTypeType 'fileContent' \
  OptionTypeFieldLabel 'FileContent' \
  OptionTypeFieldName 'fileContent' \
  OptionTypeHelpBlock 'Set fileContent' \
  OptionTypeRequired 'true' \
  OptionTypeExportMeta 'true' \
  OptionTypeDisplayValueOnDetails 'true' \
  OptionTypeLocked 'true' \
  OptionTypeHidden 'false' \
  OptionTypePlaceholder 'testing123' \
  OptionTypeExcludeFromSearch 'true'

$RENDER \
  -out examples/resources/morpheus_form/resource_select.tf \
  form_select.tf.tmpl \
  Name 'demo' \
  Code 'demo' \
  Description 'demo' \
  Labels '["terraform", "demo"]' \
  OptionTypeName 'tf example select' \
  OptionTypeCode 'select-input' \
  OptionTypeDescription 'Terraform select example' \
  OptionTypeType 'select' \
  OptionTypeFieldLabel 'Select Test' \
  OptionTypeFieldName 'selectTest' \
  OptionTypeDefaultValue '"level1"' \
  OptionTypePlaceholder 'Testing 123' \
  OptionTypeHelpBlock 'Select an option' \
  OptionTypeOptionListId '1' \
  OptionTypeRequired 'true' \
  OptionTypeExportMeta 'true' \
  OptionTypeDisplayValueOnDetails 'true' \
  OptionTypeLocked 'true' \
  OptionTypeHidden 'true' \
  OptionTypeExcludeFromSearch 'true'

$RENDER \
  -out examples/resources/morpheus_form/resource_radio.tf \
  form_radio.tf.tmpl \
  Name 'demo' \
  Code 'demo' \
  Description 'demo' \
  Labels '["terraform", "demo"]' \
  OptionTypeName 'tf radio example' \
  OptionTypeCode 'radio-input' \
  OptionTypeDescription 'Terraform radio example' \
  OptionTypeType 'radio' \
  OptionTypeFieldLabel 'Radio Test' \
  OptionTypeFieldName 'radioTest' \
  OptionTypeDefaultValue 'Demo123' \
  OptionTypePlaceholder 'Testing 123' \
  OptionTypeHelpBlock 'Select an option' \
  OptionTypeOptionListId '1' \
  OptionTypeRequired 'true' \
  OptionTypeExportMeta 'true' \
  OptionTypeDisplayValueOnDetails 'true' \
  OptionTypeLocked 'true' \
  OptionTypeHidden 'true' \
  OptionTypeExcludeFromSearch 'true'

$RENDER \
  -out examples/resources/morpheus_form/resource_text.tf \
  form_text.tf.tmpl \
  Name 'demo' \
  Code 'demo' \
  Description 'demo' \
  Labels '["terraform", "demo"]' \
  OptionTypeName 'tf text example' \
  OptionTypeCode 'test-input' \
  OptionTypeDescription 'Terraform text example' \
  OptionTypeType 'text' \
  OptionTypeFieldLabel 'Testin' \
  OptionTypeFieldName 'test' \
  OptionTypeDefaultValue 'Demo123' \
  OptionTypePlaceholder 'Testing 123' \
  OptionTypeHelpBlock 'Help block example' \
  OptionTypeRequired 'true' \
  OptionTypeExportMeta 'true' \
  OptionTypeDisplayValueOnDetails 'true' \
  OptionTypeLocked 'true' \
  OptionTypeHidden 'true' \
  OptionTypeExcludeFromSearch 'true'

$RENDER \
  -out examples/resources/morpheus_form/resource_checkbox.tf \
  form_checkbox.tf.tmpl \
  Name 'demo' \
  Code 'demo' \
  Description 'demo' \
  Labels '["terraform", "demo"]' \
  OptionTypeName 'tf checkbox example' \
  OptionTypeCode 'checkbox-input' \
  OptionTypeDescription 'Terraform checkbox example' \
  OptionTypeType 'checkbox' \
  OptionTypeFieldLabel 'checkbox input' \
  OptionTypeFieldName 'checkboxInput' \
  OptionTypeDefaultChecked 'true' \
  OptionTypePlaceholder 'Testing 123' \
  OptionTypeHelpBlock 'Help block example' \
  OptionTypeRequired 'true' \
  OptionTypeExportMeta 'true' \
  OptionTypeDisplayValueOnDetails 'true' \
  OptionTypeLocked 'true' \
  OptionTypeHidden 'true' \
  OptionTypeExcludeFromSearch 'true'

$RENDER \
  -out examples/resources/morpheus_form/resource_hidden.tf \
  form_hidden.tf.tmpl \
  Name 'demo' \
  Code 'demo' \
  Description 'demo' \
  Labels '["terraform", "demo"]' \
  OptionTypeName 'tf hidden input example' \
  OptionTypeCode 'hidden-input' \
  OptionTypeDescription 'Terraform hidden input example' \
  OptionTypeType 'hidden' \
  OptionTypeFieldLabel 'hidden input' \
  OptionTypeFieldName 'hiddenInput' \
  OptionTypeDefaultValue 'test' \
  OptionTypePlaceholder 'Testing 123' \
  OptionTypeHelpBlock 'Help block example' \
  OptionTypeRequired 'true' \
  OptionTypeExportMeta 'true' \
  OptionTypeDisplayValueOnDetails 'true' \
  OptionTypeLocked 'true' \
  OptionTypeHidden 'true' \
  OptionTypeExcludeFromSearch 'true'

$RENDER \
  -out examples/resources/morpheus_form/resource_number.tf \
  form_number.tf.tmpl \
  Name 'demo' \
  Code 'demo' \
  Description 'demo' \
  Labels '["terraform", "demo"]' \
  OptionTypeName 'tf number input example' \
  OptionTypeCode 'number-input' \
  OptionTypeDescription 'Terraform number example' \
  OptionTypeType 'number' \
  OptionTypeFieldLabel 'number input' \
  OptionTypeFieldName 'numberInput' \
  OptionTypeDefaultValue '4' \
  OptionTypePlaceholder 'Testing 123' \
  OptionTypeHelpBlock 'Help block example' \
  OptionTypeRequired 'true' \
  OptionTypeExportMeta 'true' \
  OptionTypeDisplayValueOnDetails 'true' \
  OptionTypeLocked 'true' \
  OptionTypeHidden 'true' \
  OptionTypeExcludeFromSearch 'true' \
  OptionTypeMinValue '3' \
  OptionTypeMaxValue '44' \
  OptionTypeStep '2'

$RENDER \
  -out examples/resources/morpheus_form/resource_network_manager.tf \
  form_network_manager.tf.tmpl \
  Name 'demo' \
  Code 'demo' \
  Description 'demo' \
  Labels '["terraform", "demo"]' \
  OptionTypeName 'tf network manager example' \
  OptionTypeCode 'network-manager-input' \
  OptionTypeDescription 'Terraform network manager example' \
  OptionTypeType 'networkManager' \
  OptionTypeFieldLabel 'network input' \
  OptionTypeFieldName 'networkInput' \
  OptionTypeDefaultValue 'jsonencode([{ primaryInterface = true, displayOrder = 0, ipMode = "", ipAddress = "", networkInterfaceTypeId = "4", network = { id = "network-216", pool = "{id: \"\"}" } }, { primaryInterface = false, displayOrder = 1, ipMode = "dhcp", ipAddress = "", networkInterfaceTypeId = 4, network = { id = "network-216", pool = "{id: \"\"}" } }])' \
  OptionTypePlaceholder 'Select network' \
  OptionTypeHelpBlock 'Select a network' \
  OptionTypeRequired 'true' \
  OptionTypeExportMeta 'true' \
  OptionTypeDisplayValueOnDetails 'true' \
  OptionTypeLocked 'true' \
  OptionTypeHidden 'false' \
  OptionTypeExcludeFromSearch 'true' \
  OptionTypeShowNetworkTypeSelection 'true' \
  OptionTypeEnableIPModeSelection 'true' \
  OptionTypeGroupFieldType 'value' \
  OptionTypeGroupId '1' \
  OptionTypeCloudFieldType 'value' \
  OptionTypeCloudId '1' \
  OptionTypePoolFieldType 'value' \
  OptionTypePoolId '1' \
  OptionTypeLayoutFieldType 'value' \
  OptionTypeLayoutId '1'

$RENDER \
  -out examples/resources/morpheus_form/resource_cloud.tf \
  form_cloud.tf.tmpl \
  Name 'demo' \
  Code 'demo' \
  Description 'demo' \
  Labels '["terraform", "demo"]' \
  OptionTypeName 'tf cloud example' \
  OptionTypeCode 'cloud-input' \
  OptionTypeDescription 'Terraform cloud example' \
  OptionTypeType 'cloud' \
  OptionTypeFieldLabel 'cloud input' \
  OptionTypeFieldName 'cloudInput' \
  OptionTypeDefaultValue 'test123' \
  OptionTypePlaceholder 'Select cloud' \
  OptionTypeHelpBlock 'Select a cloud' \
  OptionTypeRequired 'true' \
  OptionTypeExportMeta 'true' \
  OptionTypeDisplayValueOnDetails 'true' \
  OptionTypeLocked 'true' \
  OptionTypeHidden 'false' \
  OptionTypeExcludeFromSearch 'true' \
  OptionTypeFilterFromResource 'true' \
  OptionTypeGroupFieldType 'value' \
  OptionTypeGroupId '1' \
  OptionTypeInstanceTypeFieldType 'value' \
  OptionTypeInstanceTypeCode 'apache' \
  OptionTypeCloudType '4'

$RENDER \
  -out examples/resources/morpheus_form/resource_layout.tf \
  form_layout.tf.tmpl \
  Name 'demo' \
  Code 'demo' \
  Description 'demo' \
  Labels '["terraform", "demo"]' \
  OptionTypeName 'tf layout example' \
  OptionTypeCode 'layout-input' \
  OptionTypeDescription 'Terraform layout example' \
  OptionTypeType 'layout' \
  OptionTypeFieldLabel 'layout input' \
  OptionTypeFieldName 'layoutInput' \
  OptionTypeDefaultValue '' \
  OptionTypePlaceholder 'Select layout' \
  OptionTypeHelpBlock 'Select a layout' \
  OptionTypeRequired 'true' \
  OptionTypeExportMeta 'true' \
  OptionTypeDisplayValueOnDetails 'true' \
  OptionTypeLocked 'true' \
  OptionTypeHidden 'false' \
  OptionTypeExcludeFromSearch 'true' \
  OptionTypeGroupFieldType 'value' \
  OptionTypeGroupId '1' \
  OptionTypeCloudFieldType 'value' \
  OptionTypeCloudId '1' \
  OptionTypeInstanceTypeFieldType 'value' \
  OptionTypeInstanceTypeCode 'apache'

$RENDER \
  -out examples/resources/morpheus_form/resource_group.tf \
  form_group.tf.tmpl \
  Name 'demo' \
  Code 'demo' \
  Description 'demo' \
  Labels '["terraform", "demo"]' \
  OptionTypeName 'tf group example' \
  OptionTypeCode 'group-input' \
  OptionTypeDescription 'Terraform group example' \
  OptionTypeType 'group' \
  OptionTypeFieldLabel 'group input' \
  OptionTypeFieldName 'groupInput' \
  OptionTypeDefaultValue 'test123' \
  OptionTypePlaceholder 'Select group' \
  OptionTypeHelpBlock 'Select a group' \
  OptionTypeRequired 'true' \
  OptionTypeExportMeta 'true' \
  OptionTypeDisplayValueOnDetails 'true' \
  OptionTypeLocked 'true' \
  OptionTypeHidden 'false' \
  OptionTypeExcludeFromSearch 'true' \
  OptionTypeAllowReadOnly 'true'

$RENDER \
  -out examples/resources/morpheus_form/resource_disk_manager.tf \
  form_disk_manager.tf.tmpl \
  Name 'demo' \
  Code 'demo' \
  Description 'demo' \
  Labels '["terraform", "demo"]' \
  OptionTypeName 'tf disk manager example' \
  OptionTypeCode 'disk-manager-input' \
  OptionTypeDescription 'Terraform disk manager example' \
  OptionTypeType 'diskManager' \
  OptionTypeFieldLabel 'disk manager input' \
  OptionTypeFieldName 'diskManagerInput' \
  OptionTypeDefaultValue 'jsonencode([{ rootVolume = true, name = "root", size = 10, sizeBytes = 10737418240, minStorage = 0, displayOrder = 0, storageType = 1, datastoreId = "52" }, { rootVolume = false, name = "data-1", size = 20, sizeBytes = 21474836480, minStorage = 0, displayOrder = 1, datastoreId = "autoCluster", storageType = 1 }])' \
  OptionTypeHelpBlock 'Configure disks' \
  OptionTypeRequired 'true' \
  OptionTypeExportMeta 'true' \
  OptionTypeDisplayValueOnDetails 'true' \
  OptionTypeLocked 'true' \
  OptionTypeHidden 'false' \
  OptionTypeExcludeFromSearch 'true' \
  OptionTypeGroupFieldType 'value' \
  OptionTypeGroupId '1' \
  OptionTypeCloudFieldType 'value' \
  OptionTypeCloudId '1' \
  OptionTypePlanFieldType 'value' \
  OptionTypePlanId '1' \
  OptionTypeLayoutFieldType 'value' \
  OptionTypeLayoutId '1' \
  OptionTypePoolFieldType 'value' \
  OptionTypePoolId '1' \
  OptionTypeVirtualImageFieldType 'value' \
  OptionTypeImageId '1' \
  OptionTypeEnableDiskTypeSelection 'true' \
  OptionTypeEnableStorageTypeSelection 'true' \
  OptionTypeEnableDatastoreSelection 'true'

$RENDER \
  -out examples/resources/morpheus_form/resource_plan.tf \
  form_plan.tf.tmpl \
  Name 'demo' \
  Code 'demo' \
  Description 'demo' \
  Labels '["terraform", "demo"]' \
  OptionTypeName 'tf plan example' \
  OptionTypeCode 'plan-input' \
  OptionTypeDescription 'Terraform plan example' \
  OptionTypeType 'plan' \
  OptionTypeFieldLabel 'plan input' \
  OptionTypeFieldName 'planInput' \
  OptionTypeDefaultValue 'jsonencode({ id = 1088, maxMemory = 8589934592, maxCores = "4", coresPerSocket = "2" })' \
  OptionTypePlaceholder 'Select plan' \
  OptionTypeHelpBlock 'Select a plan' \
  OptionTypeRequired 'true' \
  OptionTypeExportMeta 'true' \
  OptionTypeDisplayValueOnDetails 'true' \
  OptionTypeLocked 'true' \
  OptionTypeHidden 'false' \
  OptionTypeExcludeFromSearch 'true' \
  OptionTypeShowPricing 'false' \
  OptionTypeGroupFieldType 'value' \
  OptionTypeGroupId '1' \
  OptionTypeCloudFieldType 'value' \
  OptionTypeCloudId '1' \
  OptionTypeLayoutFieldType 'value' \
  OptionTypeLayoutId '1' \
  OptionTypePoolFieldType 'value' \
  OptionTypePoolId '1'

$RENDER \
  -out examples/resources/morpheus_form/resource_field_groups.tf \
  form_field_groups.tf.tmpl \
  Name 'demo' \
  Code 'demo' \
  Description 'demo' \
  Labels '["terraform", "demo"]' \
  FieldGroup1Name 'fg1' \
  FieldGroup1Description 'testin' \
  FieldGroup1Collapsible 'true' \
  FieldGroup1CollapsedByDefault 'true' \
  FieldGroup1OptionTypeName 'tf field group 1 text input example' \
  FieldGroup1OptionTypeCode 'test-input-1' \
  FieldGroup1OptionTypeDescription 'Terraform text input example' \
  FieldGroup1OptionTypeType 'text' \
  FieldGroup1OptionTypeFieldLabel 'Testing 1' \
  FieldGroup1OptionTypeFieldName 'test1' \
  FieldGroup1OptionTypeDefaultValue 'Demo123' \
  FieldGroup1OptionTypePlaceholder 'Testing 123' \
  FieldGroup1OptionTypeHelpBlock 'Help block example' \
  FieldGroup1OptionTypeRequired 'true' \
  FieldGroup1OptionTypeExportMeta 'true' \
  FieldGroup1OptionTypeDisplayValueOnDetails 'true' \
  FieldGroup1OptionTypeLocked 'true' \
  FieldGroup1OptionTypeHidden 'false' \
  FieldGroup1OptionTypeExcludeFromSearch 'true' \
  FieldGroup2Name 'fg2' \
  FieldGroup2Description 'testin' \
  FieldGroup2Collapsible 'true' \
  FieldGroup2CollapsedByDefault 'true' \
  FieldGroup2OptionTypeName 'tf field group 2 text input example' \
  FieldGroup2OptionTypeCode 'test-input-2' \
  FieldGroup2OptionTypeDescription 'Terraform text input example' \
  FieldGroup2OptionTypeType 'text' \
  FieldGroup2OptionTypeFieldLabel 'Testing 2' \
  FieldGroup2OptionTypeFieldName 'test2' \
  FieldGroup2OptionTypeDefaultValue 'Demo123' \
  FieldGroup2OptionTypePlaceholder 'Testing 123' \
  FieldGroup2OptionTypeHelpBlock 'Help block example' \
  FieldGroup2OptionTypeRequired 'true' \
  FieldGroup2OptionTypeExportMeta 'true' \
  FieldGroup2OptionTypeDisplayValueOnDetails 'true' \
  FieldGroup2OptionTypeLocked 'true' \
  FieldGroup2OptionTypeHidden 'false' \
  FieldGroup2OptionTypeExcludeFromSearch 'true'

$RENDER \
  -out examples/resources/morpheus_form/resource_environment.tf \
  form_environment.tf.tmpl \
  Name 'demo' \
  Code 'demo' \
  Description 'demo' \
  Labels '["terraform", "demo"]' \
  OptionTypeName 'tf environment example' \
  OptionTypeCode 'environment-input' \
  OptionTypeDescription 'Terraform environment example' \
  OptionTypeType 'environment' \
  OptionTypeFieldLabel 'Environment' \
  OptionTypeFieldName 'environment' \
  OptionTypeDefaultValue 'staging' \
  OptionTypePlaceholder '' \
  OptionTypeHelpBlock 'Select an environment' \
  OptionTypeRequired 'true' \
  OptionTypeExportMeta 'true' \
  OptionTypeDisplayValueOnDetails 'true' \
  OptionTypeLocked 'true' \
  OptionTypeHidden 'false' \
  OptionTypeExcludeFromSearch 'true'

$RENDER \
  -out examples/resources/morpheus_form/resource_servers_input.tf \
  form_servers_input.tf.tmpl \
  Name 'demo' \
  Code 'demo' \
  Description 'demo' \
  Labels '["terraform", "demo"]' \
  OptionTypeName 'tf servers-input example' \
  OptionTypeCode 'servers-input' \
  OptionTypeDescription 'Terraform servers-input example' \
  OptionTypeType 'servers-input' \
  OptionTypeFieldLabel 'Server' \
  OptionTypeFieldName 'server' \
  OptionTypeDefaultValue '' \
  OptionTypeHelpBlock 'Select a server' \
  OptionTypeCloudFieldType 'value' \
  OptionTypeCloudId '1' \
  OptionTypeRequired 'true' \
  OptionTypeExportMeta 'true' \
  OptionTypeDisplayValueOnDetails 'true' \
  OptionTypeLocked 'true' \
  OptionTypeHidden 'false' \
  OptionTypeExcludeFromSearch 'true'

$RENDER \
  -out examples/resources/morpheus_form/resource_resource_pool.tf \
  form_resource_pool.tf.tmpl \
  Name 'demo' \
  Code 'demo' \
  Description 'demo' \
  Labels '["terraform", "demo"]' \
  OptionTypeName 'tf resourcePool example' \
  OptionTypeCode 'resource-pool-input' \
  OptionTypeDescription 'Terraform resourcePool example' \
  OptionTypeType 'resourcePool' \
  OptionTypeFieldLabel 'Resource Pool' \
  OptionTypeFieldName 'resourcePool' \
  OptionTypeDefaultValue '' \
  OptionTypeHelpBlock 'Select a resource pool' \
  OptionTypeGroupFieldType 'value' \
  OptionTypeGroupId '1' \
  OptionTypeCloudFieldType 'value' \
  OptionTypeCloudId '1' \
  OptionTypePlanFieldType 'value' \
  OptionTypePlanId '1' \
  OptionTypeLayoutFieldType 'value' \
  OptionTypeLayoutId '1' \
  OptionTypeRequired 'true' \
  OptionTypeExportMeta 'true' \
  OptionTypeDisplayValueOnDetails 'true' \
  OptionTypeLocked 'true' \
  OptionTypeHidden 'false' \
  OptionTypeExcludeFromSearch 'true'


$RENDER \
  -out examples/resources/morpheus_form/resource_sec_group.tf \
  form_sec_group.tf.tmpl \
  Name 'demo' \
  Code 'demo' \
  Description 'demo' \
  Labels '["terraform", "demo"]' \
  OptionTypeName 'tf secGroup example' \
  OptionTypeCode 'sec-group-input' \
  OptionTypeDescription 'Terraform secGroup example' \
  OptionTypeType 'secGroup' \
  OptionTypeFieldLabel 'Security Groups' \
  OptionTypeFieldName 'securityGroups' \
  OptionTypeDefaultValue 'jsonencode([{ id = "sec-group-default" }])' \
  OptionTypeHelpBlock 'Select security groups' \
  OptionTypeRequired 'true' \
  OptionTypeExportMeta 'true' \
  OptionTypeDisplayValueOnDetails 'true' \
  OptionTypeLocked 'true' \
  OptionTypeHidden 'false' \
  OptionTypeExcludeFromSearch 'true' \
  OptionTypeCloudFieldType 'value' \
  OptionTypeCloudId '1' \
  OptionTypePoolField 'resourcePool'

$RENDER \
  -out examples/resources/morpheus_form/resource_tag.tf \
  form_tag.tf.tmpl \
  Name 'demo' \
  Code 'demo' \
  Description 'demo' \
  Labels '["terraform", "demo"]' \
  OptionTypeName 'tf tag example' \
  OptionTypeCode 'tag-input' \
  OptionTypeDescription 'Terraform tag example' \
  OptionTypeType 'tag' \
  OptionTypeFieldLabel 'Tags' \
  OptionTypeFieldName 'tags' \
  OptionTypeDefaultValue 'jsonencode([{ name = "Sample Name", value = "Sample Value" }])' \
  OptionTypeHelpBlock 'Configure tags' \
  OptionTypeRequired 'true' \
  OptionTypeExportMeta 'true' \
  OptionTypeDisplayValueOnDetails 'true' \
  OptionTypeLocked 'true' \
  OptionTypeHidden 'false' \
  OptionTypeExcludeFromSearch 'true' \
  OptionTypeGroupFieldType 'value' \
  OptionTypeGroupId '1' \
  OptionTypeCloudFieldType 'value' \
  OptionTypeCloudId '1'

$RENDER \
  -out examples/resources/morpheus_form/resource_instances_input.tf \
  form_instances_input.tf.tmpl \
  Name 'demo' \
  Code 'demo' \
  Description 'demo' \
  Labels '["terraform", "demo"]' \
  OptionTypeName 'tf instances-input example' \
  OptionTypeCode 'instances-input' \
  OptionTypeDescription 'Terraform instances-input example' \
  OptionTypeType 'instances-input' \
  OptionTypeFieldLabel 'Instance' \
  OptionTypeFieldName 'instance' \
  OptionTypeDefaultValue '' \
  OptionTypeHelpBlock 'Select an instance' \
  OptionTypeCloudFieldType 'value' \
  OptionTypeCloudId '1' \
  OptionTypeRequired 'true' \
  OptionTypeExportMeta 'true' \
  OptionTypeDisplayValueOnDetails 'true' \
  OptionTypeLocked 'true' \
  OptionTypeHidden 'false' \
  OptionTypeExcludeFromSearch 'true'

$RENDER \
  -out examples/resources/morpheus_form/resource_ports.tf \
  form_ports.tf.tmpl \
  Name 'demo' \
  Code 'demo' \
  Description 'demo' \
  Labels '["terraform", "demo"]' \
  OptionTypeName 'tf ports example' \
  OptionTypeCode 'ports-input' \
  OptionTypeDescription 'Terraform ports example' \
  OptionTypeType 'ports' \
  OptionTypeFieldLabel 'Exposed Ports' \
  OptionTypeFieldName 'ports' \
  OptionTypeDefaultValue 'jsonencode([{ name = "standard", externalPort = "80", loadBalanceProtocol = "HTTP" }, { name = "ssl-title", externalPort = "443", loadBalanceProtocol = "HTTPS" }, { name = "tcp", externalPort = "40", loadBalanceProtocol = "TCP" }])' \
  OptionTypeHelpBlock 'Configure exposed ports' \
  OptionTypeGroupField 'myGroup' \
  OptionTypeCloudField 'myCloud' \
  OptionTypeLayoutField 'myLayout' \
  OptionTypeRequired 'true' \
  OptionTypeExportMeta 'true' \
  OptionTypeDisplayValueOnDetails 'true' \
  OptionTypeLocked 'true' \
  OptionTypeHidden 'false' \
  OptionTypeExcludeFromSearch 'true'

$RENDER \
  -out examples/resources/morpheus_form/resource_httpheader.tf \
  form_httpheader.tf.tmpl \
  Name 'demo' \
  Code 'demo' \
  Description 'demo' \
  Labels '["terraform", "demo"]' \
  OptionTypeName 'tf httpheader example' \
  OptionTypeCode 'httpheader-input' \
  OptionTypeDescription 'Terraform HTTP header input example' \
  OptionTypeType 'httpHeader' \
  OptionTypeFieldLabel 'HTTP Headers' \
  OptionTypeFieldName 'httpHeaders' \
  OptionTypeDefaultValue '[{ name = "header1", value = "value1", masked = false }]' \
  OptionTypeHelpBlock 'Configure HTTP headers' \
  OptionTypeRequired 'true' \
  OptionTypeExportMeta 'true' \
  OptionTypeDisplayValueOnDetails 'true' \
  OptionTypeLocked 'true' \
  OptionTypeHidden 'false' \
  OptionTypeExcludeFromSearch 'true'

$RENDER \
  -out examples/resources/morpheus_form/resource_logo_selector.tf \
  form_logo_selector.tf.tmpl \
  Name 'demo' \
  Code 'demo' \
  Description 'demo' \
  Labels '["terraform", "demo"]' \
  OptionTypeName 'tf logo selector example' \
  OptionTypeCode 'logo-selector-input' \
  OptionTypeDescription 'Terraform logo selector example' \
  OptionTypeType 'logoSelector' \
  OptionTypeFieldLabel 'Select Logo' \
  OptionTypeFieldName 'logoSelector' \
  OptionTypeDefaultValue 'jsonencode({ value = "identicon", settings = { type = "identicon", iconLabel = "example" } })' \
  OptionTypePlaceholder '' \
  OptionTypeHelpBlock 'Select or upload a logo' \
  OptionTypeRequired 'true' \
  OptionTypeExportMeta 'true' \
  OptionTypeDisplayValueOnDetails 'true' \
  OptionTypeLocked 'true' \
  OptionTypeHidden 'false' \
  OptionTypeExcludeFromSearch 'true'

$RENDER \
  -out examples/resources/morpheus_form/resource_bytesize.tf \
  form_bytesize.tf.tmpl \
  Name 'demo' \
  Code 'demo' \
  Description 'demo' \
  Labels '["terraform", "demo"]' \
  OptionTypeName 'tf byteSize example' \
  OptionTypeCode 'bytesize-input' \
  OptionTypeDescription 'Terraform byteSize example' \
  OptionTypeType 'byteSize' \
  OptionTypeFieldLabel 'Byte Size' \
  OptionTypeFieldName 'byteSize' \
  OptionTypeDefaultValue '48318382080' \
  OptionTypePlaceholder '' \
  OptionTypeHelpBlock 'Select byte size display' \
  OptionTypeRequired 'true' \
  OptionTypeExportMeta 'true' \
  OptionTypeDisplayValueOnDetails 'true' \
  OptionTypeLocked 'true' \
  OptionTypeHidden 'false' \
  OptionTypeExcludeFromSearch 'true' \
  OptionTypeDisplay 'GB' \
  OptionTypeLockDisplay 'false'

$RENDER \
  -out examples/resources/morpheus_form/resource_code_editor.tf \
  form_code_editor.tf.tmpl \
  Name 'demo' \
  Code 'demo' \
  Description 'demo' \
  Labels '["terraform", "demo"]' \
  OptionTypeName 'tf code-editor example' \
  OptionTypeCode 'code-editor-input' \
  OptionTypeDescription 'Terraform code-editor example' \
  OptionTypeType 'code-editor' \
  OptionTypeFieldLabel 'Code Editor' \
  OptionTypeFieldName 'codeEditor' \
  OptionTypeDefaultValue 'echo hello world' \
  OptionTypePlaceholder '' \
  OptionTypeHelpBlock 'Enter code' \
  OptionTypeRequired 'true' \
  OptionTypeExportMeta 'true' \
  OptionTypeDisplayValueOnDetails 'true' \
  OptionTypeLocked 'true' \
  OptionTypeHidden 'false' \
  OptionTypeExcludeFromSearch 'true' \
  OptionTypeShowLineNumbers 'true' \
  OptionTypeCodeLanguage 'bash'

$RENDER \
  -out examples/resources/morpheus_form/resource_password.tf \
  form_password.tf.tmpl \
  Name 'demo' \
  Code 'demo' \
  Description 'demo' \
  Labels '["terraform", "demo"]' \
  OptionTypeName 'tf password example' \
  OptionTypeCode 'password-input' \
  OptionTypeDescription 'Terraform password example' \
  OptionTypeType 'password' \
  OptionTypeFieldLabel 'Password' \
  OptionTypeFieldName 'password' \
  OptionTypeDefaultValue '' \
  OptionTypePlaceholder 'Enter password' \
  OptionTypeHelpBlock 'Enter a secure password' \
  OptionTypeRequired 'true' \
  OptionTypeExportMeta 'true' \
  OptionTypeDisplayValueOnDetails 'true' \
  OptionTypeLocked 'true' \
  OptionTypeHidden 'false' \
  OptionTypeExcludeFromSearch 'true' \
  OptionTypeAllowPasswordPeek 'true'

$RENDER \
  -out examples/resources/morpheus_form/resource_textarea.tf \
  form_textarea.tf.tmpl \
  Name 'demo' \
  Code 'demo' \
  Description 'demo' \
  Labels '["terraform", "demo"]' \
  OptionTypeName 'tf textarea example' \
  OptionTypeCode 'textarea-input' \
  OptionTypeDescription 'Terraform textarea example' \
  OptionTypeType 'textarea' \
  OptionTypeFieldLabel 'Text Area' \
  OptionTypeFieldName 'textArea' \
  OptionTypeDefaultValue 'Sample text' \
  OptionTypePlaceholder 'Enter text' \
  OptionTypeHelpBlock 'Enter multiple lines of text' \
  OptionTypeRequired 'true' \
  OptionTypeExportMeta 'true' \
  OptionTypeDisplayValueOnDetails 'true' \
  OptionTypeLocked 'true' \
  OptionTypeHidden 'false' \
  OptionTypeExcludeFromSearch 'true' \
  OptionTypeTextRows '5'

$RENDER \
  -out examples/resources/morpheus_form/resource_text_array.tf \
  form_text_array.tf.tmpl \
  Name 'demo' \
  Code 'demo' \
  Description 'demo' \
  Labels '["terraform", "demo"]' \
  OptionTypeName 'tf textArray example' \
  OptionTypeCode 'text-array-input' \
  OptionTypeDescription 'Terraform textArray example' \
  OptionTypeType 'textArray' \
  OptionTypeFieldLabel 'Text Array' \
  OptionTypeFieldName 'textArray' \
  OptionTypeDefaultValue 'jsonencode(["item1", "item2", "item3"])' \
  OptionTypeHelpBlock 'Enter comma-separated values' \
  OptionTypeRequired 'true' \
  OptionTypeExportMeta 'true' \
  OptionTypeDisplayValueOnDetails 'true' \
  OptionTypeLocked 'true' \
  OptionTypeHidden 'false' \
  OptionTypeExcludeFromSearch 'true' \
  OptionTypeDelimiter ','

$RENDER \
  -out examples/resources/morpheus_form/resource_typeahead.tf \
  form_typeahead.tf.tmpl \
  Name 'demo' \
  Code 'demo' \
  Description 'demo' \
  Labels '["terraform", "demo"]' \
  OptionTypeName 'tf typeahead example' \
  OptionTypeCode 'typeahead-input' \
  OptionTypeDescription 'Terraform typeahead example' \
  OptionTypeType 'typeahead' \
  OptionTypeFieldLabel 'Typeahead' \
  OptionTypeFieldName 'typeahead' \
  OptionTypeDefaultValue 'test' \
  OptionTypePlaceholder 'Search...' \
  OptionTypeHelpBlock 'Select an option from the list' \
  OptionTypeOptionListId '1' \
  OptionTypeRequired 'true' \
  OptionTypeExportMeta 'true' \
  OptionTypeDisplayValueOnDetails 'true' \
  OptionTypeLocked 'true' \
  OptionTypeHidden 'false' \
  OptionTypeExcludeFromSearch 'true' \
  OptionTypeSortable 'true' \
  OptionTypeAllowDuplicates 'false' \
  OptionTypeCustomData '{}' \
  OptionTypeAllowMultipleSelections 'false'
