# v0.0.1 Release Notes

## New functionality

In this release (v0.0.1) the following resources have been added:
- hpe_morpheus_cloud for HPE HVM or HPE VME clouds
- hpe_morpheus_group
- hpe_morpheus_instance for HPE HVM or HPE VME instances (Create, Delete and Read - no Update)
- hpe_morpheus_network
- hpe_morpheus_role for Morpheus roles (user and tenant)
- hpe_morpheus_service_plan (Create, Delete and Read - no Update)
- hpe_morpheus_user (Create, Delete and Read - no Update)

In this release (v0.0.1) the following data sources have been added:
- hpe_morpheus_cloud
- hpe_morpheus_environment
- hpe_morpheus_group
- hpe_morpheus_instance_type_layout
- hpe_morpheus_network
- hpe_morpheus_role
- hpe_morpheus_service_plan

## Known Issues

- hpe_morpheus_instance has issues with using the same `datastore_id` with multiple volumes, please use
  a different `datastore_id` for each volume.
- hpe_morpheus_instance depending on the layout used may require one or more `volumes` to be specified,
  in these cases not specifying the correct number of `volumes` will cause instance creation to fail.
- hpe_morpheus_network switchId is not supported yet, prevents creating some network types, eg OVS Port Group
- hpe_morpheus_user will force recreation if an attribute is updated
