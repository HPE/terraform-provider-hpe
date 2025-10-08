# v0.2.0 Release Notes

In this release (v0.2.0) we have added the following resource functionality:
- hpe_morpheus_datastore resource has been added (Create, Delete, Read, Update)
- hpe_morpheus_service_plan Update functionality has been added
- hpe_morpheus_cloud now has a dynamic `config` block to support arbitrary cloud configuration options

We have added the following data-source functionality:
- hpe_morpheus_datastore data-source has been added

We have fixed the following issues:
- hpe_morpheus_user would force recreation if an attribute was updated, this has been fixed
- hpe_morpheus_network switchId is now supported
- hpe_morpheus_role `default persona` issue has been fixed

## New known issues

- hpe_morpheus_datastore when creating a datastore of type NFS the creation will silently fail if the NFS server is not reachable or the share is not accessible.
  The datastore will remain in a `provisioning` state indefinitely. Ensure the Morpheus appliance can reach the NFS server
  and that the share is accessible before creating.

## Known Issues from previous releases

- hpe_morpheus_instance has issues with using the same `datastore_id` with multiple volumes, please use
  a different `datastore_id` for each volume.
- hpe_morpheus_instance depending on the layout used may require one or more `volumes` to be specified,
  in these cases not specifying the correct number of `volumes` will cause instance creation to fail.
- There are intermittent issues with the provider failing to authenticate, a 500 error is returned from the Morpheus API.
  If this happens please retry the operation.  This is being investigated.


# v0.1.0 Release Notes

## New functionality

In this release (v0.1.0) the following resources have been added:
- hpe_morpheus_cloud for HPE HVM or HPE VME clouds
- hpe_morpheus_group
- hpe_morpheus_instance for HPE HVM or HPE VME instances (Create, Delete and Read - no Update)
- hpe_morpheus_network
- hpe_morpheus_role for Morpheus roles (user and tenant)
- hpe_morpheus_service_plan (Create, Delete and Read - no Update)
- hpe_morpheus_user (Create, Delete and Read - no Update)

In this release (v0.1.0) the following data sources have been added:
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
- There are intermittent issues with the provider failing to authenticate, a 500 error is returned from the Morpheus API.
  If this happens please retry the operation.  This is being investigated.
