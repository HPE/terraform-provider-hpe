#!/bin/sh
# Helper script to generate per-type form resource examples

RENDER=../../../../bin/render

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
  OptionTypeDefaultValue 'test123' \
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
  OptionTypeHelpBlock 'Is this working now' \
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
  OptionTypeHelpBlock 'Is this working now' \
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
  OptionTypeHelpBlock 'Is this working now' \
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
  OptionTypeHelpBlock 'Is this working now' \
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
  OptionTypeDefaultValue 'test123' \
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
  OptionTypeDefaultValue '' \
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
  FieldGroup1OptionTypeCode 'test-input' \
  FieldGroup1OptionTypeDescription 'Terraform text input example' \
  FieldGroup1OptionTypeType 'text' \
  FieldGroup1OptionTypeFieldLabel 'Testin' \
  FieldGroup1OptionTypeFieldName 'test' \
  FieldGroup1OptionTypeDefaultValue 'Demo123' \
  FieldGroup1OptionTypePlaceholder 'Testing 123' \
  FieldGroup1OptionTypeHelpBlock 'Is this working now' \
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
  FieldGroup2OptionTypeCode 'test-input' \
  FieldGroup2OptionTypeDescription 'Terraform text input example' \
  FieldGroup2OptionTypeType 'text' \
  FieldGroup2OptionTypeFieldLabel 'Testin' \
  FieldGroup2OptionTypeFieldName 'test' \
  FieldGroup2OptionTypeDefaultValue 'Demo123' \
  FieldGroup2OptionTypePlaceholder 'Testing 123' \
  FieldGroup2OptionTypeHelpBlock 'Is this working now' \
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
  OptionTypeDefaultValue '' \
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
  OptionTypeDefaultValue '' \
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
