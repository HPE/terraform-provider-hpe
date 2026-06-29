<img src="https://morpheusdata.com/wp-content/uploads/2020/04/morpheus-logo-v2.svg" width="200px">

# go-morpheus-sdk
[![GoReportCard][report-badge]][report]
[![GitHub release](https://img.shields.io/github/release/gomorpheus/morpheus-go-sdk.svg)](https://github.com/gomorpheus/morpheus-go-sdk/releases/)
[![GoDoc](https://pkg.go.dev/badge/badge/github.com/gomorpheus/morpheus-go-sdk?utm_source=godoc)](https://godoc.org/github.com/gomorpheus/morpheus-go-sdk)

[report-badge]: https://goreportcard.com/badge/github.com/gomorpheus/morpheus-go-sdk
[report]: https://goreportcard.com/report/github.com/gomorpheus/morpheus-go-sdk

- Website: https://www.morpheusdata.com/
- Docs: [Morpheus Documentation](https://docs.morpheusdata.com)
- Support: [Morpheus Support](https://support.morpheusdata.com)

This package provides the official [Go](https://golang.org/) library for the [Morpheus API](https://apidocs.morpheusdata.com/).

This is being developed in conjunction with the [Morpheus Terraform Provider](https://github.com/gomorpheus/terraform-provider-morpheus).

## Setup

Install Go, export environment variables, go get the morpheus package and begin executing requests.

## Requirements

* [Go](https://golang.org/dl/) | 1.17

### Environment Variables

Be sure to setup your Go environment variables.

```bash
export GOPATH=$HOME/gocode
export PATH=$PATH:$GOPATH/bin
```

### Installation

Use go get to retrieve the SDK to add it to your GOPATH workspace, or project's Go module dependencies.

```sh
go get github.com/gomorpheus/morpheus-go-sdk
```

To update the SDK use go get -u to retrieve the latest version of the SDK.

```sh
go get -u github.com/gomorpheus/morpheus-go-sdk
```

## Usage

Here are some examples of how to use `morpheus.Client`.

### New Client

Instantiate a new client and authenticate.

```go
import "github.com/gomorpheus/morpheus-go-sdk"
client := morpheus.NewClient("https://yourmorpheus.com")
client.SetUsernameAndPassword("username", "password")
resp, err := client.Login()
if err != nil {
    fmt.Println("LOGIN ERROR: ", err)
}
fmt.Println("LOGIN RESPONSE:", resp)
```

You can also create a client with a valid access token, instead of authenticating with a username and password.

```go
import "github.com/gomorpheus/morpheus-go-sdk"
client := morpheus.NewClient("https://yourmorpheus.com")
client.SetAccessToken("a3a4c6ea-fb54-42af-109b-63bdd19e5ae1", "", 0, "write")
resp, err := client.Whoami()
if err != nil {
    fmt.Println("WHOAMI ERROR: ", err)
}
fmt.Println("WHOAMI RESPONSE:", resp)
```

**NOTE** It is not necessary to call `client.Login()` explicitely. The client will attempt to authenticate, if needed, whenever `Execute()` is called.

### Execute Any Request

You can also use the `Execute` method to execute an arbitrary api request, using any http method, path parameters, and body.

```go
resp, err := client.Execute(&morpheus.Request{
    Method: "GET",
    Path: "/api/instances",
    QueryParams:map[string]string{
        "name": "tftest",
    },
})
if err != nil {
    fmt.Println("API ERROR: ", err)
}
fmt.Println("API RESPONSE:", resp)
```

### List Instances

Fetch a list of instances.

```go
resp, err := client.ListInstances(&morpheus.Request{})
// parse JSON and fetch the first one by ID
listInstancesResult := resp.Result.(*morpheus.ListInstancesResult)
instancesCount := listInstancesResult.Meta.Total
fmt.Sprintf("Found %d Instances.", instancesCount)
```

**NOTE:** This may be simplified so that typecasting the result is not always needed.

## Testing

You can execute the latest tests using:

```sh
go test
```

The above command will (ideally) print results like this:

```
Initializing test client for tfplugin @ https://yourmorpheus.com
PASS
ok      github.com/gomorpheus/morpheus-go-sdk   1.098s
```

Running `go test` will fail with a panic right away if you have not yet setup your test environment variables.  

```bash
export MORPHEUS_TEST_URL=https://yourmorpheus.com
export MORPHEUS_TEST_USERNAME=gotest
export MORPHEUS_TEST_PASSWORD=19830B3f489
export MORPHEUS_TEST_TOKEN=8c6380df-4cwf-40qd-9fm6-hj16a0357094
```
**Be Careful running this test suite**. It creates and destroys data. Never point at any URL other than a test environment. Although, in reality, tests will not modify or destroy any pre-existing data. It could still orphan some test some data, or cause otherwise undesired effects.

You can run an individual test like this:

```sh
go test -run TestGroupsCRUD
```


```bash
go test -v
```

## Contribution

This library is currently under development.  Eventually every API endpoint will have a corresponding method defined by [Client](client.go) with the request and response types defined.

Feel free to contribute by implementing the list of missing endpoints. See [Coverage](#coverage).

### Code Structure

The main type this package exposes is [Client](../blob/master/client.go), implemented in client.go.  

Each resource is defined in its own file eg. [instances.go](../blob/master/instances.go)  which extends the `Client` type by defining a function for each endpoint the resource has, such as GetInstance(), ListInstances(), CreateInstance(), UpdateInstance, DeleteInstance(), etc. The request and response payload types used by those methods are also defined here.

#### Test Files

Be sure to add a `_test.go` file with unit tests for each new resource that is implemented.

### External Resources

Link | Description
--------- | -----------
[Morpheus API](https://apidocs.morpheusdata.com/) | The Morpheus API documentation.


## Coverage

API | Available?
--------- | -----------
account_groups | n/a
accounts | [Accounts](accounts.go)
activity | [Activity](activity.go)
appliance_settings | [Appliance Settings](appliance_settings.go)
approvals | [Approvals](approvals.go)
apps | [Apps](apps.go)
archive_buckets | n/a
archive_files | n/a
auth | n/a
blueprints | [Blueprints](blueprints.go)
budgets | [Budgets](budgets.go)
cloud_datastores | n/a
cloud_folders | n/a
cloud_policies | n/a
cloud_resource_pools | n/a
clouds | [Clouds](clouds.go)
clusters | [Clusters](clusters.go)
containers | n/a
custom_instance_types | n/a
cypher | [Cypher](cypher.go)
dashboard | n/a
deploy | n/a
deployments | [Deployments](deployments.go)
credentials | [Credentials](credentials.go)
environments | [Environments](environments.go)
execute_schedules | [Execute Schedules](execute_schedules.go)
execution_request | n/a
file_copy_request | n/a
group_policies | n/a
groups | [Groups](groups.go)
guidance_settings | [Guidance Settings](guidance_settings.go)
image_builder | n/a
instances | [Instances](instances.go)
integrations | [Integrations](integrations.go)
key_pairs | [Key Pairs](key_pairs.go)
library_cluster_layouts | [Cluster Layouts](cluster_layouts.go)
library_compute_type_layouts | n/a
library_container_scripts | [Script Templates](script_templates.go)
library_container_templates | [File Templates](file_templates.go)
library_container_types |[Node Types](node_types.go)
library_container_upgrades | n/a
library_instance_types | [Instance Types](instance_types.go)
library_layouts | n/a
library_spec_templates | [Spec Templates](spec_templates.go)
license | [License](license.go)
load_balancer_monitors | [Load Balancer Monitors](load_balancer_monitors.go)
load_balancer_pools | [Load Balancer Pools](load_balancer_pools.go)
load_balancer_monitors | [Load Balancer Profiles](load_balancer_profiles.go)
load_balancer_types | [Load Balancer Types](load_balancer_monitors.go)
load_balancer_virtual_servers | [Load Balancer Virtual Servers](load_balancer_virtual_servers.go)
load_balancers | [Load Balancers](load_balancers.go)
logs | n/a
log_settings | [Log Settings](log_settings.go)
monitoring | n/a
monitoring.checks | [Checks](checks.go)
monitoring.groups | [Check Groups](check_groups.go)
monitoring.apps | [Monitoring Apps](monitoring_apps.go)
monitoring_settings | [Monitoring Settings](monitoring_settings.go)
monitoring.incidents | [Incidents](incidents.go)
monitoring.alerts | [Alerts](alerts.go)
monitoring.contacts | [Contacts](contacts.go)
network_domain_records | n/a
network_domains | [Network Domains](network_domains.go)
network_groups | [Network Groups](network_groups.go)
network_pool_ips | n/a
network_pool_servers | [Network Pool Servers](network_pool_servers.go)
network_pools | [Network Pools](network_pools.go)
network_proxies | [Network Proxies](network_proxies.go)
network_services | n/a
network_subnet_types | n/a
network_subnets | n/a
network_types | n/a
networks | [Networks](networks.go)
option_type_lists | [Option Type Lists](option_type_lists.go)
option_types | [Option Types](option_types.go)
plans | [Plans](plans.go)
plugins | [Plugins](plugins.go)
policies | [Policies](policies.go)
power_schedules | [Power Schedules](power_schedules.go)
prices | [Prices](prices.go)
price_sets | [Price Sets](price_sets.go)
processes | n/a
provision_types | [Provision Types](provision_types.go)
refresh_token | n/a
reports | n/a
resource_pool_groups | [Resource Pool Groups](resource_pool_groups.go)
roles | [Roles](roles.go)
scale_thresholds | [Scale Thresholds](scale_thresholds.go)
security_group_rules | n/a
security_groups | n/a
security_packages | [Security Packages](security_packages.go)
security_scans | [Security Scans](security_scans.go)
server_types | n/a
servers | n/a
service_plans | [Service Plans](service_plans.go)
setup | [Setup](setup.go)
software_licenses | [Software Licenses](software_licenses.go)
storage_buckets | [Storage Buckets](storage_buckets.go)
storage_providers | n/a
storage_servers | [Storage Servers](storage_servers.go)
subnets | n/a
task_sets | [Task Sets](task_sets.go)
tasks | [Tasks](tasks.go)
user_groups | [User Groups](user_groups.go)
user_settings | n/a
user_sources | [Identity Sources](identity_sources.go)
users | [Users](users.go)
vdi_allocations | [VDI Allocations](vdi_allocations.go)
vdi_apps | [VDI Apps](vdi_apps.go)
vdi_gateways | [VDI Gateways](vdi_gateway.go)
vdi_pools | [VDI Pools](vdi_pools.go)
virtual_images | [Virtual Images](virtual_images.go)
whoami | [Whoami](whoami.go)
whitelabel_settings | [Whitelabel Settings](whitelabel_settings.go)
wikis | [Wikis](wikis.go)
