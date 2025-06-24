package compare_test

import (
	"encoding/json"
	"errors"
	"fmt"
	"testing"

	"github.com/HewlettPackard/hpe-morpheus-go-sdk/sdk"

	"github.com/HPE/terraform-provider-hpe/internal/subproviders/morpheus/compare"

	"github.com/stretchr/testify/assert"
)

func TestContainsSubsetPermissionsJSON(t *testing.T) {
	//nolint:lll
	testCases := []struct {
		name            string
		permissionsPlan string
		permissionsAPI  string
		expectEq        bool
		expectErr       error
	}{
		{
			name:            "is subset: user featurePermissions vs user featurePermissions",
			permissionsPlan: permissionsTestUserFeaturePermissions,
			permissionsAPI:  permissionsTestUserFeaturePermissions,
			expectEq:        true,
			expectErr:       nil,
		},
		{
			name:            "is subset: computed API permissions vs computed API permissions",
			permissionsPlan: permissionsTestAPIComputedFull,
			permissionsAPI:  permissionsTestAPIComputedFull,
			expectEq:        true,
			expectErr:       nil,
		},
		{
			name:            "is subset: user featurePermissions vs minimal out of order API featurePermissions",
			permissionsPlan: permissionsTestUserFeaturePermissions,
			permissionsAPI:  permissionsTestAPIFeaturePermissionsOutOfOrder,
			expectEq:        true,
			expectErr:       nil,
		},
		{
			name:            "is subset: user featurePermissions vs partial computed API featurePermissions",
			permissionsPlan: permissionsTestUserFeaturePermissions,
			permissionsAPI:  permissionsTestAPIComputedPartial,
			expectEq:        true,
			expectErr:       nil,
		},
		{
			name:            "is subset: user mixed permissions vs mixed partial computed API permissions",
			permissionsPlan: permissionsTestUserMixedPermissions,
			permissionsAPI:  permissionsTestAPIMixedPermissionsComputedPartial,
			expectEq:        true,
			expectErr:       nil,
		},
		{
			name:            "is subset: statefile computed permissions vs computed API permissions",
			permissionsPlan: permissionsTestAPIComputedFullStatefile,
			permissionsAPI:  permissionsTestAPIComputedFull,
			expectEq:        true,
			expectErr:       nil,
		},
		{
			name:            "not subset: string value mismatch",
			permissionsPlan: permissionsTestUserDefaultPermissions,
			permissionsAPI:  permissionsTestAPIDefaultPermissionsMismatch,
			expectEq:        false,
			expectErr:       errors.New(compare.ErrorNotSubset),
		},
		{
			name:            "not subset: no matching string key found",
			permissionsPlan: permissionsTestUserMinimum,
			permissionsAPI:  permissionsTestAPINoMatchingStringKey,
			expectEq:        false,
			expectErr:       errors.New(compare.ErrorNotSubset),
		},
		{
			name:            "not subset: no matching string value found",
			permissionsPlan: permissionsTestUserMinimum,
			permissionsAPI:  permissionsTestAPINoMatchingStringValue,
			expectEq:        false,
			expectErr:       errors.New(compare.ErrorNotSubset),
		},
		{
			name:            "not subset: no matching array key found",
			permissionsPlan: permissionsTestUserMinimum,
			permissionsAPI:  permissionsTestAPINoMatchingArrayKey,
			expectEq:        false,
			expectErr:       errors.New(compare.ErrorNotSubset),
		},
		{
			name:            "not subset: no matching array value found",
			permissionsPlan: permissionsTestUserMinimum,
			permissionsAPI:  permissionsTestAPINoMatchingArrayValue,
			expectEq:        false,
			expectErr:       errors.New(compare.ErrorNotSubset),
		},
		{
			name:            "not subset: array element key value mismatch",
			permissionsPlan: permissionsTestUserMinimum,
			permissionsAPI:  permissionsTestAPIMinimumAccessMismatch,
			expectEq:        false,
			expectErr:       errors.New(compare.ErrorNotSubset),
		},
	}

	// for map
	for _, tc := range testCases {
		t.Run(fmt.Sprintf("test as map: %s", tc.name), func(t *testing.T) {
			var planStruct, apiStruct map[string]any

			err := json.Unmarshal([]byte(tc.permissionsPlan), &planStruct)
			assert.NoError(t, err)

			err = json.Unmarshal([]byte(tc.permissionsAPI), &apiStruct)
			assert.NoError(t, err)

			eq, err := compare.ContainsSubset(apiStruct, planStruct)
			assert.Equal(t, tc.expectEq, eq)
			assert.Equal(t, tc.expectErr, err)
		})
	}

	// for struct
	for _, tc := range testCases {
		t.Run(fmt.Sprintf("test as struct: %s", tc.name), func(t *testing.T) {
			var planStruct, apiStruct sdk.GetRole200Response

			err := json.Unmarshal([]byte(tc.permissionsPlan), &planStruct)
			assert.NoError(t, err)

			err = json.Unmarshal([]byte(tc.permissionsAPI), &apiStruct)
			assert.NoError(t, err)

			eq, err := compare.ContainsSubset(apiStruct, planStruct)
			assert.Equal(t, tc.expectEq, eq)
			assert.Equal(t, tc.expectErr, err)
		})
	}
}

func TestContainsSubsetStruct(t *testing.T) {
	testCases := []struct {
		name      string
		super     any
		sub       any
		expectEq  bool
		expectErr error
	}{
		{
			name: "is subset: equal structs with private fields",
			super: struct {
				private string
				Public  string
			}{Public: "foo"},
			sub: struct {
				private string
				Public  string
			}{Public: "foo"},
			expectEq:  true,
			expectErr: nil,
		},
		{
			name: "not subset: unequal structs with private fields",
			super: struct {
				private string
				Public  string
			}{Public: "foo"},
			sub: struct {
				private string
				Public  string
			}{Public: "bar"},
			expectEq:  false,
			expectErr: errors.New(compare.ErrorNotSubset),
		},
	}
	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			eq, err := compare.ContainsSubset(tc.super, tc.sub)
			assert.Equal(t, tc.expectEq, eq)
			assert.Equal(t, tc.expectErr, err)
		})
	}
}

const (
	permissionsTestUserDefaultPermissions = `
{
  "globalInstanceTypeAccess": "full"
}
`
	permissionsTestAPIDefaultPermissionsMismatch = `
{
  "globalInstanceTypeAccess": "none"
}
`
	permissionsTestUserFeaturePermissions = `
{
  "featurePermissions": [
    {
      "code": "integrations-ansible",
      "access": "full"
    },
    {
      "code": "admin-appliance",
      "access": "none"
    },
    {
      "code": "app-templates",
      "access": "none"
    }
  ]
}
`
	// admin-appliance and app-template not in same order as the user set
	permissionsTestAPIFeaturePermissionsOutOfOrder = `
{
  "featurePermissions": [
    {
      "id": 159,
      "code": "integrations-ansible",
      "name": "Ansible",
      "access": "full",
      "subCategory": "admin"
    },
    {
      "id": 34,
      "code": "app-templates",
      "name": "App Blueprints",
      "access": "none",
      "subCategory": "library"
    },
    {
      "id": 14,
      "code": "admin-appliance",
      "name": "Appliance Settings",
      "access": "none",
      "subCategory": "admin"
    }
  ]
}
`

	// test that a mix of the string and array values works
	permissionsTestUserMixedPermissions = `
{
  "featurePermissions": [
    {
      "code": "integrations-ansible",
      "access": "full"
    },
    {
      "code": "admin-appliance",
      "access": "none"
    }
  ],
  "globalZoneAccess": "full",
  "globalPersonaAccess": "full",
  "personaPermissions": [
    {
      "code": "serviceCatalog",
      "access": "full"
    }
  ]
}
`
	permissionsTestUserMinimum = `
{
  "featurePermissions": [
    {
      "code": "integrations-ansible",
      "access": "full"
    }
  ]
}
`
	permissionsTestAPIMinimumAccessMismatch = `
{
  "featurePermissions": [
    {
      "code": "integrations-ansible",
      "access": "none"
    }
  ]
}
`
	// this should fail because it's missing the "code" key
	permissionsTestAPINoMatchingStringKey = `
{
  "featurePermissions": [
    {
      "id": 159,
      "foo": "integrations-foo",
      "name": "Ansible",
      "access": "full",
      "subCategory": "admin"
    }
  ]
}
`
	// this should fail because even though it has the "code" key, its value is not equal
	permissionsTestAPINoMatchingStringValue = `
{
  "featurePermissions": [
    {
      "id": 159,
      "code": "integrations-foo",
      "name": "Ansible",
      "access": "full",
      "subCategory": "admin"
    }
  ]
}
`
	// this should fail because the key of the JSON array, "fooPermissions"
	// does not match "featurePermissions" when the algorithm searches for it
	permissionsTestAPINoMatchingArrayKey = `
{
  "fooPermissions": [
    {
      "id": 159,
      "code": "integrations-ansible",
      "name": "Ansible",
      "access": "full",
      "subCategory": "admin"
    }
  ]
}
`
	// this should fail because the algorithm cannot find the array values sub in this.
	// this is supposed to contain those values at minimum.
	permissionsTestAPINoMatchingArrayValue = `
{
  "featurePermissions": []
}
`
	// permissions obtained from GET to API after a create using permissionsTestUserFeaturePermissions
	permissionsTestAPIComputedPartial = `
{
  "featurePermissions": [
    {
      "id": 45,
      "code": "activity",
      "name": "Activity",
      "access": "none",
      "subCategory": "operations"
    },
    {
      "id": 72,
      "code": "provisioning-admin",
      "name": "Administrator",
      "access": "none",
      "subCategory": "provisioning"
    },
    {
      "id": 83,
      "code": "library-advanced-node-type-options",
      "name": "Advanced Node Type Options",
      "access": "none",
      "subCategory": "provisioning"
    },
    {
      "id": 48,
      "code": "operations-alarms",
      "name": "Alarms",
      "access": "none",
      "subCategory": "operations"
    },
    {
      "id": 123,
      "code": "reports-analytics",
      "name": "Analytics",
      "access": "none",
      "subCategory": "operations"
    },
    {
      "id": 159,
      "code": "integrations-ansible",
      "name": "Ansible",
      "access": "full",
      "subCategory": "admin"
    },
    {
      "id": 34,
      "code": "app-templates",
      "name": "App Blueprints",
      "access": "none",
      "subCategory": "library"
    },
    {
      "id": 14,
      "code": "admin-appliance",
      "name": "Appliance Settings",
      "access": "none",
      "subCategory": "admin"
    },
    {
      "id": 50,
      "code": "operations-approvals",
      "name": "Approvals",
      "access": "none",
      "subCategory": "operations"
    },
    {
      "id": 33,
      "code": "apps",
      "name": "Apps",
      "access": "none",
      "subCategory": "provisioning"
    },
    {
      "id": 116,
      "code": "services-archives",
      "name": "Archives",
      "access": "none",
      "subCategory": "tools"
    },
    {
      "id": 17,
      "code": "admin-backupSettings",
      "name": "Backup Settings",
      "access": "none",
      "subCategory": "admin"
    },
    {
      "id": 40,
      "code": "backups",
      "name": "Backups",
      "access": "none",
      "subCategory": "backups"
    },
    {
      "id": 55,
      "code": "billing",
      "name": "Billing",
      "access": "none",
      "subCategory": "api"
    },
    {
      "id": 35,
      "code": "arm-template",
      "name": "Blueprints - ARM",
      "access": "none",
      "subCategory": "library"
    },
    {
      "id": 37,
      "code": "cloudFormation-template",
      "name": "Blueprints - CloudFormation",
      "access": "none",
      "subCategory": "library"
    },
    {
      "id": 39,
      "code": "helm-template",
      "name": "Blueprints - Helm",
      "access": "none",
      "subCategory": "library"
    },
    {
      "id": 38,
      "code": "kubernetes-template",
      "name": "Blueprints - Kubernetes",
      "access": "none",
      "subCategory": "library"
    },
    {
      "id": 36,
      "code": "terraform-template",
      "name": "Blueprints - Terraform",
      "access": "none",
      "subCategory": "library"
    },
    {
      "id": 145,
      "code": "infrastructure-boot",
      "name": "Boot",
      "access": "none",
      "subCategory": "infrastructure"
    },
    {
      "id": 51,
      "code": "operations-budgets",
      "name": "Budgets",
      "access": "none",
      "subCategory": "operations"
    },
    {
      "id": 97,
      "code": "service-catalog",
      "name": "Catalog",
      "access": "none",
      "subCategory": "catalog"
    },
    {
      "id": 96,
      "code": "catalog",
      "name": "Catalog Items",
      "access": "none",
      "subCategory": "library"
    },
    {
      "id": 111,
      "code": "admin-certificates",
      "name": "Certificates",
      "access": "none",
      "subCategory": "infrastructure"
    },
    {
      "id": 31,
      "code": "admin-clients",
      "name": "Clients",
      "access": "none",
      "subCategory": "admin"
    },
    {
      "id": 3,
      "code": "admin-zones",
      "name": "Clouds",
      "access": "none",
      "subCategory": "infrastructure"
    },
    {
      "id": 153,
      "code": "infrastructure-cluster",
      "name": "Clusters",
      "access": "none",
      "subCategory": "infrastructure"
    },
    {
      "id": 100,
      "code": "deployments",
      "name": "Code Deployments",
      "access": "none",
      "subCategory": "provisioning"
    },
    {
      "id": 101,
      "code": "deployment-services",
      "name": "Code Integrations",
      "access": "none",
      "subCategory": "provisioning"
    },
    {
      "id": 102,
      "code": "code-repositories",
      "name": "Code Repositories",
      "access": "none",
      "subCategory": "provisioning"
    },
    {
      "id": 1,
      "code": "admin-servers",
      "name": "Compute",
      "access": "none",
      "subCategory": "infrastructure"
    },
    {
      "id": 120,
      "code": "services-vdi-copy",
      "name": "Copy/Paste",
      "access": "none",
      "subCategory": "virtual-desktop"
    },
    {
      "id": 109,
      "code": "credentials",
      "name": "Credentials",
      "access": "none",
      "subCategory": "infrastructure"
    },
    {
      "id": 115,
      "code": "services-cypher",
      "name": "Cypher",
      "access": "none",
      "subCategory": "tools"
    },
    {
      "id": 44,
      "code": "dashboard",
      "name": "Dashboard",
      "access": "none",
      "subCategory": "operations"
    },
    {
      "id": 98,
      "code": "service-catalog-dashboard",
      "name": "Dashboard",
      "access": "none",
      "subCategory": "catalog"
    },
    {
      "id": 130,
      "code": "infrastructure-network-dhcp-relay",
      "name": "DHCP Relays",
      "access": "none",
      "subCategory": "networks"
    },
    {
      "id": 129,
      "code": "infrastructure-network-dhcp-server",
      "name": "DHCP Servers",
      "access": "none",
      "subCategory": "networks"
    },
    {
      "id": 25,
      "code": "admin-distributed-workers",
      "name": "Distributed Workers",
      "access": "none",
      "subCategory": "admin"
    },
    {
      "id": 133,
      "code": "infrastructure-domains",
      "name": "Domains",
      "access": "none",
      "subCategory": "networks"
    },
    {
      "id": 20,
      "code": "admin-environments",
      "name": "Environment Settings",
      "access": "none",
      "subCategory": "admin"
    },
    {
      "id": 65,
      "code": "provisioning-environment",
      "name": "Environment Variables",
      "access": "none",
      "subCategory": "lifecycle"
    },
    {
      "id": 68,
      "code": "provisioning-execute-script",
      "name": "Execute Script",
      "access": "none",
      "subCategory": "provisioning"
    },
    {
      "id": 69,
      "code": "provisioning-execute-task",
      "name": "Execute Task",
      "access": "none",
      "subCategory": "provisioning"
    },
    {
      "id": 70,
      "code": "provisioning-execute-workflow",
      "name": "Execute Workflow",
      "access": "none",
      "subCategory": "provisioning"
    },
    {
      "id": 157,
      "code": "execution-request",
      "name": "Execution Request",
      "access": "none",
      "subCategory": "api"
    },
    {
      "id": 103,
      "code": "executions",
      "name": "Executions",
      "access": "none",
      "subCategory": "provisioning"
    },
    {
      "id": 6,
      "code": "admin-export-import",
      "name": "Export/Import",
      "access": "none",
      "subCategory": "admin"
    },
    {
      "id": 95,
      "code": "lifecycle-extend",
      "name": "Extend Expirations",
      "access": "none",
      "subCategory": "lifecycle"
    },
    {
      "id": 137,
      "code": "infrastructure-network-firewalls",
      "name": "Firewalls",
      "access": "none",
      "subCategory": "networks"
    },
    {
      "id": 136,
      "code": "infrastructure-floating-ips",
      "name": "Floating IPs",
      "access": "none",
      "subCategory": "networks"
    },
    {
      "id": 2,
      "code": "admin-groups",
      "name": "Groups",
      "access": "none",
      "subCategory": "infrastructure"
    },
    {
      "id": 54,
      "code": "guidance",
      "name": "Guidance",
      "access": "none",
      "subCategory": "operations"
    },
    {
      "id": 16,
      "code": "admin-guidanceSettings",
      "name": "Guidance Settings",
      "access": "none",
      "subCategory": "admin"
    },
    {
      "id": 47,
      "code": "admin-health",
      "name": "Health",
      "access": "none",
      "subCategory": "admin"
    },
    {
      "id": 12,
      "code": "admin-identity-sources",
      "name": "Identity Source",
      "access": "none",
      "subCategory": "admin"
    },
    {
      "id": 118,
      "code": "services-image-builder",
      "name": "Image Builder",
      "access": "none",
      "subCategory": "tools"
    },
    {
      "id": 67,
      "code": "provisioning-import-image",
      "name": "Import Image",
      "access": "none",
      "subCategory": "provisioning"
    },
    {
      "id": 9,
      "code": "admin-containers",
      "name": "Instance Types",
      "access": "none",
      "subCategory": "library"
    },
    {
      "id": 59,
      "code": "provisioning-add",
      "name": "Instances: Add",
      "access": "none",
      "subCategory": "provisioning"
    },
    {
      "id": 66,
      "code": "provisioning-clone",
      "name": "Instances: Clone",
      "access": "none",
      "subCategory": "provisioning"
    },
    {
      "id": 61,
      "code": "provisioning-delete",
      "name": "Instances: Delete",
      "access": "none",
      "subCategory": "provisioning"
    },
    {
      "id": 60,
      "code": "provisioning-edit",
      "name": "Instances: Edit",
      "access": "none",
      "subCategory": "provisioning"
    },
    {
      "id": 73,
      "code": "provisioning-force-delete",
      "name": "Instances: Force Delete",
      "access": "none",
      "subCategory": "provisioning"
    },
    {
      "id": 58,
      "code": "provisioning",
      "name": "Instances: List",
      "access": "none",
      "subCategory": "provisioning"
    },
    {
      "id": 62,
      "code": "provisioning-lock",
      "name": "Instances: Lock/Unlock",
      "access": "none",
      "subCategory": "provisioning"
    },
    {
      "id": 74,
      "code": "provisioning-remove-control",
      "name": "Instances: Remove From Control",
      "access": "none",
      "subCategory": "provisioning"
    },
    {
      "id": 63,
      "code": "provisioning-scale",
      "name": "Instances: Scale",
      "access": "none",
      "subCategory": "provisioning"
    },
    {
      "id": 64,
      "code": "provisioning-settings",
      "name": "Instances: Settings",
      "access": "none",
      "subCategory": "provisioning"
    },
    {
      "id": 128,
      "code": "infrastructure-network-integrations",
      "name": "Integration",
      "access": "none",
      "subCategory": "networks"
    },
    {
      "id": 23,
      "code": "admin-cm",
      "name": "Integrations",
      "access": "none",
      "subCategory": "admin"
    },
    {
      "id": 43,
      "code": "backup-services",
      "name": "Integrations",
      "access": "none",
      "subCategory": "backups"
    },
    {
      "id": 105,
      "code": "automation-services",
      "name": "Integrations",
      "access": "none",
      "subCategory": "library"
    },
    {
      "id": 158,
      "code": "operations-invoices",
      "name": "Invoices",
      "access": "none",
      "subCategory": "operations"
    },
    {
      "id": 135,
      "code": "infrastructure-ippools",
      "name": "IP Pools",
      "access": "none",
      "subCategory": "networks"
    },
    {
      "id": 106,
      "code": "job-executions",
      "name": "Job Executions",
      "access": "none",
      "subCategory": "provisioning"
    },
    {
      "id": 107,
      "code": "job-templates",
      "name": "Jobs",
      "access": "none",
      "subCategory": "provisioning"
    },
    {
      "id": 112,
      "code": "admin-keypairs",
      "name": "Keypairs",
      "access": "none",
      "subCategory": "infrastructure"
    },
    {
      "id": 117,
      "code": "services-kubernetes",
      "name": "Kubernetes",
      "access": "none",
      "subCategory": "tools"
    },
    {
      "id": 154,
      "code": "infrastructure-kube-cntl",
      "name": "Kubernetes Control",
      "access": "none",
      "subCategory": "infrastructure"
    },
    {
      "id": 161,
      "code": "library-packages",
      "name": "Library Packages",
      "access": "none",
      "subCategory": "library"
    },
    {
      "id": 21,
      "code": "admin-licenses",
      "name": "License Settings",
      "access": "none",
      "subCategory": "admin"
    },
    {
      "id": 126,
      "code": "infrastructure-loadbalancer",
      "name": "Load Balancers",
      "access": "none",
      "subCategory": "infrastructure"
    },
    {
      "id": 121,
      "code": "services-vdi-printer",
      "name": "Local Printer",
      "access": "none",
      "subCategory": "virtual-desktop"
    },
    {
      "id": 22,
      "code": "admin-logSettings",
      "name": "Log Settings",
      "access": "none",
      "subCategory": "admin"
    },
    {
      "id": 49,
      "code": "logs",
      "name": "Logs",
      "access": "none",
      "subCategory": "monitoring"
    },
    {
      "id": 29,
      "code": "admin-motd",
      "name": "Message of the day",
      "access": "none",
      "subCategory": "admin"
    },
    {
      "id": 57,
      "code": "monitoring",
      "name": "Monitoring",
      "access": "none",
      "subCategory": "monitoring"
    },
    {
      "id": 18,
      "code": "admin-monitorSettings",
      "name": "Monitoring Settings",
      "access": "none",
      "subCategory": "admin"
    },
    {
      "id": 155,
      "code": "infrastructure-move-server",
      "name": "Move Servers",
      "access": "none",
      "subCategory": "infrastructure"
    },
    {
      "id": 127,
      "code": "infrastructure-networks",
      "name": "Networks",
      "access": "none",
      "subCategory": "infrastructure"
    },
    {
      "id": 10,
      "code": "library-options",
      "name": "Options",
      "access": "none",
      "subCategory": "library"
    },
    {
      "id": 99,
      "code": "service-catalog-inventory",
      "name": "Order History",
      "access": "none",
      "subCategory": "catalog"
    },
    {
      "id": 26,
      "code": "admin-packages",
      "name": "Packages",
      "access": "none",
      "subCategory": "admin"
    },
    {
      "id": 24,
      "code": "admin-plugins",
      "name": "Plugins",
      "access": "none",
      "subCategory": "admin"
    },
    {
      "id": 27,
      "code": "admin-policies",
      "name": "Policies",
      "access": "none",
      "subCategory": "infrastructure"
    },
    {
      "id": 28,
      "code": "admin-global-policies",
      "name": "Policies",
      "access": "none",
      "subCategory": "admin"
    },
    {
      "id": 93,
      "code": "provisioning-power",
      "name": "Power Control",
      "access": "none",
      "subCategory": "lifecycle"
    },
    {
      "id": 30,
      "code": "admin-profiles",
      "name": "Profiles",
      "access": "none",
      "subCategory": "admin"
    },
    {
      "id": 156,
      "code": "projects",
      "name": "Projects",
      "access": "none",
      "subCategory": "projects"
    },
    {
      "id": 19,
      "code": "admin-provisioningSettings",
      "name": "Provisioning Settings",
      "access": "none",
      "subCategory": "admin"
    },
    {
      "id": 134,
      "code": "infrastructure-proxies",
      "name": "Proxies",
      "access": "none",
      "subCategory": "networks"
    },
    {
      "id": 84,
      "code": "provisioning-reconfigure",
      "name": "Reconfigure",
      "access": "none",
      "subCategory": "lifecycle"
    },
    {
      "id": 85,
      "code": "provisioning-reconfigure-change-plan",
      "name": "Reconfigure: Change Plan",
      "access": "none",
      "subCategory": "lifecycle"
    },
    {
      "id": 90,
      "code": "provisioning-reconfigure-add-disk",
      "name": "Reconfigure: Disk Add",
      "access": "none",
      "subCategory": "lifecycle"
    },
    {
      "id": 92,
      "code": "provisioning-reconfigure-disk-type",
      "name": "Reconfigure: Disk Change Type",
      "access": "none",
      "subCategory": "lifecycle"
    },
    {
      "id": 89,
      "code": "provisioning-reconfigure-modify-disk",
      "name": "Reconfigure: Disk Modify",
      "access": "none",
      "subCategory": "lifecycle"
    },
    {
      "id": 91,
      "code": "provisioning-reconfigure-remove-disk",
      "name": "Reconfigure: Disk Remove",
      "access": "none",
      "subCategory": "lifecycle"
    },
    {
      "id": 87,
      "code": "provisioning-reconfigure-add-network",
      "name": "Reconfigure: Network Add",
      "access": "none",
      "subCategory": "lifecycle"
    },
    {
      "id": 86,
      "code": "provisioning-reconfigure-modify-network",
      "name": "Reconfigure: Network Modify",
      "access": "none",
      "subCategory": "lifecycle"
    },
    {
      "id": 88,
      "code": "provisioning-reconfigure-remove-network",
      "name": "Reconfigure: Network Remove",
      "access": "none",
      "subCategory": "lifecycle"
    },
    {
      "id": 76,
      "code": "terminal",
      "name": "Remote Console",
      "access": "none",
      "subCategory": "provisioning"
    },
    {
      "id": 77,
      "code": "terminal-access",
      "name": "Remote Console Auto Login",
      "access": "none",
      "subCategory": "provisioning"
    },
    {
      "id": 122,
      "code": "reports",
      "name": "Reports",
      "access": "none",
      "subCategory": "operations"
    },
    {
      "id": 94,
      "code": "lifecycle-retry-cancel",
      "name": "Retry/Cancel",
      "access": "none",
      "subCategory": "lifecycle"
    },
    {
      "id": 8,
      "code": "admin-roles",
      "name": "Roles",
      "access": "none",
      "subCategory": "admin"
    },
    {
      "id": 141,
      "code": "infrastructure-network-router-firewalls",
      "name": "Router Firewalls",
      "access": "none",
      "subCategory": "networks"
    },
    {
      "id": 140,
      "code": "infrastructure-network-router-interfaces",
      "name": "Router Interfaces",
      "access": "none",
      "subCategory": "networks"
    },
    {
      "id": 138,
      "code": "infrastructure-nat",
      "name": "Router NAT",
      "access": "none",
      "subCategory": "networks"
    },
    {
      "id": 142,
      "code": "infrastructure-network-router-redistribution",
      "name": "Router Redistribution",
      "access": "none",
      "subCategory": "networks"
    },
    {
      "id": 139,
      "code": "infrastructure-network-router-routes",
      "name": "Router Routes",
      "access": "none",
      "subCategory": "networks"
    },
    {
      "id": 132,
      "code": "infrastructure-routers",
      "name": "Routers",
      "access": "none",
      "subCategory": "networks"
    },
    {
      "id": 160,
      "code": "security-scan",
      "name": "Scanning",
      "access": "none",
      "subCategory": "security"
    },
    {
      "id": 53,
      "code": "scheduling-execute",
      "name": "Scheduling - Execute",
      "access": "none",
      "subCategory": "library"
    },
    {
      "id": 52,
      "code": "scheduling-power",
      "name": "Scheduling - Power",
      "access": "none",
      "subCategory": "library"
    },
    {
      "id": 144,
      "code": "infrastructure-securityGroups",
      "name": "Security Groups",
      "access": "none",
      "subCategory": "infrastructure"
    },
    {
      "id": 143,
      "code": "infrastructure-network-server-groups",
      "name": "Server Groups",
      "access": "none",
      "subCategory": "networks"
    },
    {
      "id": 114,
      "code": "services-network-registry",
      "name": "Service Mesh",
      "access": "none",
      "subCategory": "provisioning"
    },
    {
      "id": 7,
      "code": "admin-servicePlans",
      "name": "Service Plans",
      "access": "none",
      "subCategory": "admin"
    },
    {
      "id": 41,
      "code": "snapshots",
      "name": "Snapshots",
      "access": "none",
      "subCategory": "snapshots"
    },
    {
      "id": 42,
      "code": "snapshots-linked-clone",
      "name": "Snapshots: Linked Clone",
      "access": "none",
      "subCategory": "snapshots"
    },
    {
      "id": 71,
      "code": "provisioning-state",
      "name": "State",
      "access": "none",
      "subCategory": "provisioning"
    },
    {
      "id": 146,
      "code": "infrastructure-state",
      "name": "State",
      "access": "none",
      "subCategory": "infrastructure"
    },
    {
      "id": 131,
      "code": "infrastructure-network-dhcp-routes",
      "name": "Static Routes",
      "access": "none",
      "subCategory": "networks"
    },
    {
      "id": 124,
      "code": "infrastructure-storage",
      "name": "Storage",
      "access": "none",
      "subCategory": "infrastructure"
    },
    {
      "id": 125,
      "code": "infrastructure-storage-browser",
      "name": "Storage Browser",
      "access": "none",
      "subCategory": "infrastructure"
    },
    {
      "id": 104,
      "code": "tasks",
      "name": "Tasks",
      "access": "none",
      "subCategory": "library"
    },
    {
      "id": 108,
      "code": "task-scripts",
      "name": "Tasks - Script Engines",
      "access": "none",
      "subCategory": "library"
    },
    {
      "id": 11,
      "code": "library-templates",
      "name": "Templates",
      "access": "none",
      "subCategory": "library"
    },
    {
      "id": 4,
      "code": "admin-accounts",
      "name": "Tenant",
      "access": "none",
      "subCategory": "admin"
    },
    {
      "id": 5,
      "code": "admin-accounts-users",
      "name": "Tenant - Impersonate Users",
      "access": "none",
      "subCategory": "admin"
    },
    {
      "id": 82,
      "code": "thresholds",
      "name": "Thresholds",
      "access": "none",
      "subCategory": "library"
    },
    {
      "id": 110,
      "code": "trust-services",
      "name": "Trust Integrations",
      "access": "none",
      "subCategory": "infrastructure"
    },
    {
      "id": 56,
      "code": "account-usage",
      "name": "Usage",
      "access": "none",
      "subCategory": "operations"
    },
    {
      "id": 13,
      "code": "admin-users",
      "name": "Users",
      "access": "none",
      "subCategory": "admin"
    },
    {
      "id": 119,
      "code": "services-vdi-pools",
      "name": "VDI Pools",
      "access": "none",
      "subCategory": "virtual-desktop"
    },
    {
      "id": 113,
      "code": "virtual-images",
      "name": "Virtual Images",
      "access": "none",
      "subCategory": "library"
    },
    {
      "id": 15,
      "code": "admin-whitelabel",
      "name": "Whitelabel Settings",
      "access": "none",
      "subCategory": "admin"
    },
    {
      "id": 46,
      "code": "operations-wiki",
      "name": "Wiki",
      "access": "none",
      "subCategory": "operations"
    }
  ],
  "globalSiteAccess": "none",
  "sites": [],
  "globalZoneAccess": "none",
  "zones": [],
  "globalInstanceTypeAccess": "none",
  "instanceTypePermissions": [],
  "globalAppTemplateAccess": "none",
  "appTemplatePermissions": [],
  "globalCatalogItemTypeAccess": "none",
  "catalogItemTypePermissions": [],
  "globalPersonaAccess": "none",
  "personaPermissions": [],
  "globalVdiPoolAccess": "none",
  "vdiPoolPermissions": [],
  "globalReportTypeAccess": "none",
  "reportTypePermissions": [],
  "globalTaskAccess": "none",
  "taskPermissions": [],
  "globalTaskSetAccess": "none",
  "taskSetPermissions": [],
  "globalClusterTypeAccess": "none",
  "clusterTypePermissions": []
}
`
	// the permissions received from a GET to the API after
	// creating a role using permissionsTestUserMixedPermissions
	permissionsTestAPIMixedPermissionsComputedPartial = `
{
  "featurePermissions": [
    {
      "id": 45,
      "code": "activity",
      "name": "Activity",
      "access": "none",
      "subCategory": "operations"
    },
    {
      "id": 72,
      "code": "provisioning-admin",
      "name": "Administrator",
      "access": "none",
      "subCategory": "provisioning"
    },
    {
      "id": 83,
      "code": "library-advanced-node-type-options",
      "name": "Advanced Node Type Options",
      "access": "none",
      "subCategory": "provisioning"
    },
    {
      "id": 48,
      "code": "operations-alarms",
      "name": "Alarms",
      "access": "none",
      "subCategory": "operations"
    },
    {
      "id": 123,
      "code": "reports-analytics",
      "name": "Analytics",
      "access": "none",
      "subCategory": "operations"
    },
    {
      "id": 159,
      "code": "integrations-ansible",
      "name": "Ansible",
      "access": "full",
      "subCategory": "admin"
    },
    {
      "id": 34,
      "code": "app-templates",
      "name": "App Blueprints",
      "access": "none",
      "subCategory": "library"
    },
    {
      "id": 14,
      "code": "admin-appliance",
      "name": "Appliance Settings",
      "access": "none",
      "subCategory": "admin"
    },
    {
      "id": 50,
      "code": "operations-approvals",
      "name": "Approvals",
      "access": "none",
      "subCategory": "operations"
    },
    {
      "id": 33,
      "code": "apps",
      "name": "Apps",
      "access": "none",
      "subCategory": "provisioning"
    },
    {
      "id": 116,
      "code": "services-archives",
      "name": "Archives",
      "access": "none",
      "subCategory": "tools"
    },
    {
      "id": 17,
      "code": "admin-backupSettings",
      "name": "Backup Settings",
      "access": "none",
      "subCategory": "admin"
    },
    {
      "id": 40,
      "code": "backups",
      "name": "Backups",
      "access": "none",
      "subCategory": "backups"
    },
    {
      "id": 55,
      "code": "billing",
      "name": "Billing",
      "access": "none",
      "subCategory": "api"
    },
    {
      "id": 35,
      "code": "arm-template",
      "name": "Blueprints - ARM",
      "access": "none",
      "subCategory": "library"
    },
    {
      "id": 37,
      "code": "cloudFormation-template",
      "name": "Blueprints - CloudFormation",
      "access": "none",
      "subCategory": "library"
    },
    {
      "id": 39,
      "code": "helm-template",
      "name": "Blueprints - Helm",
      "access": "none",
      "subCategory": "library"
    },
    {
      "id": 38,
      "code": "kubernetes-template",
      "name": "Blueprints - Kubernetes",
      "access": "none",
      "subCategory": "library"
    },
    {
      "id": 36,
      "code": "terraform-template",
      "name": "Blueprints - Terraform",
      "access": "none",
      "subCategory": "library"
    },
    {
      "id": 145,
      "code": "infrastructure-boot",
      "name": "Boot",
      "access": "none",
      "subCategory": "infrastructure"
    },
    {
      "id": 51,
      "code": "operations-budgets",
      "name": "Budgets",
      "access": "none",
      "subCategory": "operations"
    },
    {
      "id": 97,
      "code": "service-catalog",
      "name": "Catalog",
      "access": "none",
      "subCategory": "catalog"
    },
    {
      "id": 96,
      "code": "catalog",
      "name": "Catalog Items",
      "access": "none",
      "subCategory": "library"
    },
    {
      "id": 111,
      "code": "admin-certificates",
      "name": "Certificates",
      "access": "none",
      "subCategory": "infrastructure"
    },
    {
      "id": 31,
      "code": "admin-clients",
      "name": "Clients",
      "access": "none",
      "subCategory": "admin"
    },
    {
      "id": 3,
      "code": "admin-zones",
      "name": "Clouds",
      "access": "none",
      "subCategory": "infrastructure"
    },
    {
      "id": 153,
      "code": "infrastructure-cluster",
      "name": "Clusters",
      "access": "none",
      "subCategory": "infrastructure"
    },
    {
      "id": 100,
      "code": "deployments",
      "name": "Code Deployments",
      "access": "none",
      "subCategory": "provisioning"
    },
    {
      "id": 101,
      "code": "deployment-services",
      "name": "Code Integrations",
      "access": "none",
      "subCategory": "provisioning"
    },
    {
      "id": 102,
      "code": "code-repositories",
      "name": "Code Repositories",
      "access": "none",
      "subCategory": "provisioning"
    },
    {
      "id": 1,
      "code": "admin-servers",
      "name": "Compute",
      "access": "none",
      "subCategory": "infrastructure"
    },
    {
      "id": 120,
      "code": "services-vdi-copy",
      "name": "Copy/Paste",
      "access": "none",
      "subCategory": "virtual-desktop"
    },
    {
      "id": 109,
      "code": "credentials",
      "name": "Credentials",
      "access": "none",
      "subCategory": "infrastructure"
    },
    {
      "id": 115,
      "code": "services-cypher",
      "name": "Cypher",
      "access": "none",
      "subCategory": "tools"
    },
    {
      "id": 44,
      "code": "dashboard",
      "name": "Dashboard",
      "access": "none",
      "subCategory": "operations"
    },
    {
      "id": 98,
      "code": "service-catalog-dashboard",
      "name": "Dashboard",
      "access": "none",
      "subCategory": "catalog"
    },
    {
      "id": 130,
      "code": "infrastructure-network-dhcp-relay",
      "name": "DHCP Relays",
      "access": "none",
      "subCategory": "networks"
    },
    {
      "id": 129,
      "code": "infrastructure-network-dhcp-server",
      "name": "DHCP Servers",
      "access": "none",
      "subCategory": "networks"
    },
    {
      "id": 25,
      "code": "admin-distributed-workers",
      "name": "Distributed Workers",
      "access": "none",
      "subCategory": "admin"
    },
    {
      "id": 133,
      "code": "infrastructure-domains",
      "name": "Domains",
      "access": "none",
      "subCategory": "networks"
    },
    {
      "id": 20,
      "code": "admin-environments",
      "name": "Environment Settings",
      "access": "none",
      "subCategory": "admin"
    },
    {
      "id": 65,
      "code": "provisioning-environment",
      "name": "Environment Variables",
      "access": "none",
      "subCategory": "lifecycle"
    },
    {
      "id": 68,
      "code": "provisioning-execute-script",
      "name": "Execute Script",
      "access": "none",
      "subCategory": "provisioning"
    },
    {
      "id": 69,
      "code": "provisioning-execute-task",
      "name": "Execute Task",
      "access": "none",
      "subCategory": "provisioning"
    },
    {
      "id": 70,
      "code": "provisioning-execute-workflow",
      "name": "Execute Workflow",
      "access": "none",
      "subCategory": "provisioning"
    },
    {
      "id": 157,
      "code": "execution-request",
      "name": "Execution Request",
      "access": "none",
      "subCategory": "api"
    },
    {
      "id": 103,
      "code": "executions",
      "name": "Executions",
      "access": "none",
      "subCategory": "provisioning"
    },
    {
      "id": 6,
      "code": "admin-export-import",
      "name": "Export/Import",
      "access": "none",
      "subCategory": "admin"
    },
    {
      "id": 95,
      "code": "lifecycle-extend",
      "name": "Extend Expirations",
      "access": "none",
      "subCategory": "lifecycle"
    },
    {
      "id": 137,
      "code": "infrastructure-network-firewalls",
      "name": "Firewalls",
      "access": "none",
      "subCategory": "networks"
    },
    {
      "id": 136,
      "code": "infrastructure-floating-ips",
      "name": "Floating IPs",
      "access": "none",
      "subCategory": "networks"
    },
    {
      "id": 2,
      "code": "admin-groups",
      "name": "Groups",
      "access": "none",
      "subCategory": "infrastructure"
    },
    {
      "id": 54,
      "code": "guidance",
      "name": "Guidance",
      "access": "none",
      "subCategory": "operations"
    },
    {
      "id": 16,
      "code": "admin-guidanceSettings",
      "name": "Guidance Settings",
      "access": "none",
      "subCategory": "admin"
    },
    {
      "id": 47,
      "code": "admin-health",
      "name": "Health",
      "access": "none",
      "subCategory": "admin"
    },
    {
      "id": 12,
      "code": "admin-identity-sources",
      "name": "Identity Source",
      "access": "none",
      "subCategory": "admin"
    },
    {
      "id": 118,
      "code": "services-image-builder",
      "name": "Image Builder",
      "access": "none",
      "subCategory": "tools"
    },
    {
      "id": 67,
      "code": "provisioning-import-image",
      "name": "Import Image",
      "access": "none",
      "subCategory": "provisioning"
    },
    {
      "id": 9,
      "code": "admin-containers",
      "name": "Instance Types",
      "access": "none",
      "subCategory": "library"
    },
    {
      "id": 59,
      "code": "provisioning-add",
      "name": "Instances: Add",
      "access": "none",
      "subCategory": "provisioning"
    },
    {
      "id": 66,
      "code": "provisioning-clone",
      "name": "Instances: Clone",
      "access": "none",
      "subCategory": "provisioning"
    },
    {
      "id": 61,
      "code": "provisioning-delete",
      "name": "Instances: Delete",
      "access": "none",
      "subCategory": "provisioning"
    },
    {
      "id": 60,
      "code": "provisioning-edit",
      "name": "Instances: Edit",
      "access": "none",
      "subCategory": "provisioning"
    },
    {
      "id": 73,
      "code": "provisioning-force-delete",
      "name": "Instances: Force Delete",
      "access": "none",
      "subCategory": "provisioning"
    },
    {
      "id": 58,
      "code": "provisioning",
      "name": "Instances: List",
      "access": "none",
      "subCategory": "provisioning"
    },
    {
      "id": 62,
      "code": "provisioning-lock",
      "name": "Instances: Lock/Unlock",
      "access": "none",
      "subCategory": "provisioning"
    },
    {
      "id": 74,
      "code": "provisioning-remove-control",
      "name": "Instances: Remove From Control",
      "access": "none",
      "subCategory": "provisioning"
    },
    {
      "id": 63,
      "code": "provisioning-scale",
      "name": "Instances: Scale",
      "access": "none",
      "subCategory": "provisioning"
    },
    {
      "id": 64,
      "code": "provisioning-settings",
      "name": "Instances: Settings",
      "access": "none",
      "subCategory": "provisioning"
    },
    {
      "id": 128,
      "code": "infrastructure-network-integrations",
      "name": "Integration",
      "access": "none",
      "subCategory": "networks"
    },
    {
      "id": 23,
      "code": "admin-cm",
      "name": "Integrations",
      "access": "none",
      "subCategory": "admin"
    },
    {
      "id": 43,
      "code": "backup-services",
      "name": "Integrations",
      "access": "none",
      "subCategory": "backups"
    },
    {
      "id": 105,
      "code": "automation-services",
      "name": "Integrations",
      "access": "none",
      "subCategory": "library"
    },
    {
      "id": 158,
      "code": "operations-invoices",
      "name": "Invoices",
      "access": "none",
      "subCategory": "operations"
    },
    {
      "id": 135,
      "code": "infrastructure-ippools",
      "name": "IP Pools",
      "access": "none",
      "subCategory": "networks"
    },
    {
      "id": 106,
      "code": "job-executions",
      "name": "Job Executions",
      "access": "none",
      "subCategory": "provisioning"
    },
    {
      "id": 107,
      "code": "job-templates",
      "name": "Jobs",
      "access": "none",
      "subCategory": "provisioning"
    },
    {
      "id": 112,
      "code": "admin-keypairs",
      "name": "Keypairs",
      "access": "none",
      "subCategory": "infrastructure"
    },
    {
      "id": 117,
      "code": "services-kubernetes",
      "name": "Kubernetes",
      "access": "none",
      "subCategory": "tools"
    },
    {
      "id": 154,
      "code": "infrastructure-kube-cntl",
      "name": "Kubernetes Control",
      "access": "none",
      "subCategory": "infrastructure"
    },
    {
      "id": 161,
      "code": "library-packages",
      "name": "Library Packages",
      "access": "none",
      "subCategory": "library"
    },
    {
      "id": 21,
      "code": "admin-licenses",
      "name": "License Settings",
      "access": "none",
      "subCategory": "admin"
    },
    {
      "id": 126,
      "code": "infrastructure-loadbalancer",
      "name": "Load Balancers",
      "access": "none",
      "subCategory": "infrastructure"
    },
    {
      "id": 121,
      "code": "services-vdi-printer",
      "name": "Local Printer",
      "access": "none",
      "subCategory": "virtual-desktop"
    },
    {
      "id": 22,
      "code": "admin-logSettings",
      "name": "Log Settings",
      "access": "none",
      "subCategory": "admin"
    },
    {
      "id": 49,
      "code": "logs",
      "name": "Logs",
      "access": "none",
      "subCategory": "monitoring"
    },
    {
      "id": 29,
      "code": "admin-motd",
      "name": "Message of the day",
      "access": "none",
      "subCategory": "admin"
    },
    {
      "id": 57,
      "code": "monitoring",
      "name": "Monitoring",
      "access": "none",
      "subCategory": "monitoring"
    },
    {
      "id": 18,
      "code": "admin-monitorSettings",
      "name": "Monitoring Settings",
      "access": "none",
      "subCategory": "admin"
    },
    {
      "id": 155,
      "code": "infrastructure-move-server",
      "name": "Move Servers",
      "access": "none",
      "subCategory": "infrastructure"
    },
    {
      "id": 127,
      "code": "infrastructure-networks",
      "name": "Networks",
      "access": "none",
      "subCategory": "infrastructure"
    },
    {
      "id": 10,
      "code": "library-options",
      "name": "Options",
      "access": "none",
      "subCategory": "library"
    },
    {
      "id": 99,
      "code": "service-catalog-inventory",
      "name": "Order History",
      "access": "none",
      "subCategory": "catalog"
    },
    {
      "id": 26,
      "code": "admin-packages",
      "name": "Packages",
      "access": "none",
      "subCategory": "admin"
    },
    {
      "id": 24,
      "code": "admin-plugins",
      "name": "Plugins",
      "access": "none",
      "subCategory": "admin"
    },
    {
      "id": 27,
      "code": "admin-policies",
      "name": "Policies",
      "access": "none",
      "subCategory": "infrastructure"
    },
    {
      "id": 28,
      "code": "admin-global-policies",
      "name": "Policies",
      "access": "none",
      "subCategory": "admin"
    },
    {
      "id": 93,
      "code": "provisioning-power",
      "name": "Power Control",
      "access": "none",
      "subCategory": "lifecycle"
    },
    {
      "id": 30,
      "code": "admin-profiles",
      "name": "Profiles",
      "access": "none",
      "subCategory": "admin"
    },
    {
      "id": 156,
      "code": "projects",
      "name": "Projects",
      "access": "none",
      "subCategory": "projects"
    },
    {
      "id": 19,
      "code": "admin-provisioningSettings",
      "name": "Provisioning Settings",
      "access": "none",
      "subCategory": "admin"
    },
    {
      "id": 134,
      "code": "infrastructure-proxies",
      "name": "Proxies",
      "access": "none",
      "subCategory": "networks"
    },
    {
      "id": 84,
      "code": "provisioning-reconfigure",
      "name": "Reconfigure",
      "access": "none",
      "subCategory": "lifecycle"
    },
    {
      "id": 85,
      "code": "provisioning-reconfigure-change-plan",
      "name": "Reconfigure: Change Plan",
      "access": "none",
      "subCategory": "lifecycle"
    },
    {
      "id": 90,
      "code": "provisioning-reconfigure-add-disk",
      "name": "Reconfigure: Disk Add",
      "access": "none",
      "subCategory": "lifecycle"
    },
    {
      "id": 92,
      "code": "provisioning-reconfigure-disk-type",
      "name": "Reconfigure: Disk Change Type",
      "access": "none",
      "subCategory": "lifecycle"
    },
    {
      "id": 89,
      "code": "provisioning-reconfigure-modify-disk",
      "name": "Reconfigure: Disk Modify",
      "access": "none",
      "subCategory": "lifecycle"
    },
    {
      "id": 91,
      "code": "provisioning-reconfigure-remove-disk",
      "name": "Reconfigure: Disk Remove",
      "access": "none",
      "subCategory": "lifecycle"
    },
    {
      "id": 87,
      "code": "provisioning-reconfigure-add-network",
      "name": "Reconfigure: Network Add",
      "access": "none",
      "subCategory": "lifecycle"
    },
    {
      "id": 86,
      "code": "provisioning-reconfigure-modify-network",
      "name": "Reconfigure: Network Modify",
      "access": "none",
      "subCategory": "lifecycle"
    },
    {
      "id": 88,
      "code": "provisioning-reconfigure-remove-network",
      "name": "Reconfigure: Network Remove",
      "access": "none",
      "subCategory": "lifecycle"
    },
    {
      "id": 76,
      "code": "terminal",
      "name": "Remote Console",
      "access": "none",
      "subCategory": "provisioning"
    },
    {
      "id": 77,
      "code": "terminal-access",
      "name": "Remote Console Auto Login",
      "access": "none",
      "subCategory": "provisioning"
    },
    {
      "id": 122,
      "code": "reports",
      "name": "Reports",
      "access": "none",
      "subCategory": "operations"
    },
    {
      "id": 94,
      "code": "lifecycle-retry-cancel",
      "name": "Retry/Cancel",
      "access": "none",
      "subCategory": "lifecycle"
    },
    {
      "id": 8,
      "code": "admin-roles",
      "name": "Roles",
      "access": "none",
      "subCategory": "admin"
    },
    {
      "id": 141,
      "code": "infrastructure-network-router-firewalls",
      "name": "Router Firewalls",
      "access": "none",
      "subCategory": "networks"
    },
    {
      "id": 140,
      "code": "infrastructure-network-router-interfaces",
      "name": "Router Interfaces",
      "access": "none",
      "subCategory": "networks"
    },
    {
      "id": 138,
      "code": "infrastructure-nat",
      "name": "Router NAT",
      "access": "none",
      "subCategory": "networks"
    },
    {
      "id": 142,
      "code": "infrastructure-network-router-redistribution",
      "name": "Router Redistribution",
      "access": "none",
      "subCategory": "networks"
    },
    {
      "id": 139,
      "code": "infrastructure-network-router-routes",
      "name": "Router Routes",
      "access": "none",
      "subCategory": "networks"
    },
    {
      "id": 132,
      "code": "infrastructure-routers",
      "name": "Routers",
      "access": "none",
      "subCategory": "networks"
    },
    {
      "id": 160,
      "code": "security-scan",
      "name": "Scanning",
      "access": "none",
      "subCategory": "security"
    },
    {
      "id": 53,
      "code": "scheduling-execute",
      "name": "Scheduling - Execute",
      "access": "none",
      "subCategory": "library"
    },
    {
      "id": 52,
      "code": "scheduling-power",
      "name": "Scheduling - Power",
      "access": "none",
      "subCategory": "library"
    },
    {
      "id": 144,
      "code": "infrastructure-securityGroups",
      "name": "Security Groups",
      "access": "none",
      "subCategory": "infrastructure"
    },
    {
      "id": 143,
      "code": "infrastructure-network-server-groups",
      "name": "Server Groups",
      "access": "none",
      "subCategory": "networks"
    },
    {
      "id": 114,
      "code": "services-network-registry",
      "name": "Service Mesh",
      "access": "none",
      "subCategory": "provisioning"
    },
    {
      "id": 7,
      "code": "admin-servicePlans",
      "name": "Service Plans",
      "access": "none",
      "subCategory": "admin"
    },
    {
      "id": 41,
      "code": "snapshots",
      "name": "Snapshots",
      "access": "none",
      "subCategory": "snapshots"
    },
    {
      "id": 42,
      "code": "snapshots-linked-clone",
      "name": "Snapshots: Linked Clone",
      "access": "none",
      "subCategory": "snapshots"
    },
    {
      "id": 71,
      "code": "provisioning-state",
      "name": "State",
      "access": "none",
      "subCategory": "provisioning"
    },
    {
      "id": 146,
      "code": "infrastructure-state",
      "name": "State",
      "access": "none",
      "subCategory": "infrastructure"
    },
    {
      "id": 131,
      "code": "infrastructure-network-dhcp-routes",
      "name": "Static Routes",
      "access": "none",
      "subCategory": "networks"
    },
    {
      "id": 124,
      "code": "infrastructure-storage",
      "name": "Storage",
      "access": "none",
      "subCategory": "infrastructure"
    },
    {
      "id": 125,
      "code": "infrastructure-storage-browser",
      "name": "Storage Browser",
      "access": "none",
      "subCategory": "infrastructure"
    },
    {
      "id": 104,
      "code": "tasks",
      "name": "Tasks",
      "access": "none",
      "subCategory": "library"
    },
    {
      "id": 108,
      "code": "task-scripts",
      "name": "Tasks - Script Engines",
      "access": "none",
      "subCategory": "library"
    },
    {
      "id": 11,
      "code": "library-templates",
      "name": "Templates",
      "access": "none",
      "subCategory": "library"
    },
    {
      "id": 4,
      "code": "admin-accounts",
      "name": "Tenant",
      "access": "none",
      "subCategory": "admin"
    },
    {
      "id": 5,
      "code": "admin-accounts-users",
      "name": "Tenant - Impersonate Users",
      "access": "none",
      "subCategory": "admin"
    },
    {
      "id": 82,
      "code": "thresholds",
      "name": "Thresholds",
      "access": "none",
      "subCategory": "library"
    },
    {
      "id": 110,
      "code": "trust-services",
      "name": "Trust Integrations",
      "access": "none",
      "subCategory": "infrastructure"
    },
    {
      "id": 56,
      "code": "account-usage",
      "name": "Usage",
      "access": "none",
      "subCategory": "operations"
    },
    {
      "id": 13,
      "code": "admin-users",
      "name": "Users",
      "access": "none",
      "subCategory": "admin"
    },
    {
      "id": 119,
      "code": "services-vdi-pools",
      "name": "VDI Pools",
      "access": "none",
      "subCategory": "virtual-desktop"
    },
    {
      "id": 113,
      "code": "virtual-images",
      "name": "Virtual Images",
      "access": "none",
      "subCategory": "library"
    },
    {
      "id": 15,
      "code": "admin-whitelabel",
      "name": "Whitelabel Settings",
      "access": "none",
      "subCategory": "admin"
    },
    {
      "id": 46,
      "code": "operations-wiki",
      "name": "Wiki",
      "access": "none",
      "subCategory": "operations"
    }
  ],
  "globalSiteAccess": "none",
  "sites": [],
  "globalZoneAccess": "full",
  "zones": [],
  "globalInstanceTypeAccess": "none",
  "instanceTypePermissions": [],
  "globalAppTemplateAccess": "none",
  "appTemplatePermissions": [],
  "globalCatalogItemTypeAccess": "none",
  "catalogItemTypePermissions": [],
  "globalPersonaAccess": "full",
  "personaPermissions": [
    {
      "id": 2,
      "code": "serviceCatalog",
      "name": "Service Catalog",
      "access": "full"
    }
  ],
  "globalVdiPoolAccess": "none",
  "vdiPoolPermissions": [],
  "globalReportTypeAccess": "none",
  "reportTypePermissions": [],
  "globalTaskAccess": "none",
  "taskPermissions": [],
  "globalTaskSetAccess": "none",
  "taskSetPermissions": [],
  "globalClusterTypeAccess": "none",
  "clusterTypePermissions": []
}
`
	// permissions stored in the state file which were obtained from API
	// after not specifying any permissions, i.e. all computed defaults.
	// note how the ordering differs to permissionsTestAPIComputedFull
	permissionsTestAPIComputedFullStatefile = `
{
  "appTemplatePermissions": [],
  "catalogItemTypePermissions": [],
  "clusterTypePermissions": [],
  "featurePermissions": [
    {
      "access": "none",
      "code": "activity",
      "id": 45,
      "name": "Activity",
      "subCategory": "operations"
    },
    {
      "access": "none",
      "code": "provisioning-admin",
      "id": 72,
      "name": "Administrator",
      "subCategory": "provisioning"
    },
    {
      "access": "none",
      "code": "library-advanced-node-type-options",
      "id": 83,
      "name": "Advanced Node Type Options",
      "subCategory": "provisioning"
    },
    {
      "access": "none",
      "code": "operations-alarms",
      "id": 48,
      "name": "Alarms",
      "subCategory": "operations"
    },
    {
      "access": "none",
      "code": "reports-analytics",
      "id": 123,
      "name": "Analytics",
      "subCategory": "operations"
    },
    {
      "access": "none",
      "code": "integrations-ansible",
      "id": 159,
      "name": "Ansible",
      "subCategory": "admin"
    },
    {
      "access": "none",
      "code": "app-templates",
      "id": 34,
      "name": "App Blueprints",
      "subCategory": "library"
    },
    {
      "access": "none",
      "code": "admin-appliance",
      "id": 14,
      "name": "Appliance Settings",
      "subCategory": "admin"
    },
    {
      "access": "none",
      "code": "operations-approvals",
      "id": 50,
      "name": "Approvals",
      "subCategory": "operations"
    },
    {
      "access": "none",
      "code": "apps",
      "id": 33,
      "name": "Apps",
      "subCategory": "provisioning"
    },
    {
      "access": "none",
      "code": "services-archives",
      "id": 116,
      "name": "Archives",
      "subCategory": "tools"
    },
    {
      "access": "none",
      "code": "admin-backupSettings",
      "id": 17,
      "name": "Backup Settings",
      "subCategory": "admin"
    },
    {
      "access": "none",
      "code": "backups",
      "id": 40,
      "name": "Backups",
      "subCategory": "backups"
    },
    {
      "access": "none",
      "code": "billing",
      "id": 55,
      "name": "Billing",
      "subCategory": "api"
    },
    {
      "access": "none",
      "code": "arm-template",
      "id": 35,
      "name": "Blueprints - ARM",
      "subCategory": "library"
    },
    {
      "access": "none",
      "code": "cloudFormation-template",
      "id": 37,
      "name": "Blueprints - CloudFormation",
      "subCategory": "library"
    },
    {
      "access": "none",
      "code": "helm-template",
      "id": 39,
      "name": "Blueprints - Helm",
      "subCategory": "library"
    },
    {
      "access": "none",
      "code": "kubernetes-template",
      "id": 38,
      "name": "Blueprints - Kubernetes",
      "subCategory": "library"
    },
    {
      "access": "none",
      "code": "terraform-template",
      "id": 36,
      "name": "Blueprints - Terraform",
      "subCategory": "library"
    },
    {
      "access": "none",
      "code": "infrastructure-boot",
      "id": 145,
      "name": "Boot",
      "subCategory": "infrastructure"
    },
    {
      "access": "none",
      "code": "operations-budgets",
      "id": 51,
      "name": "Budgets",
      "subCategory": "operations"
    },
    {
      "access": "none",
      "code": "service-catalog",
      "id": 97,
      "name": "Catalog",
      "subCategory": "catalog"
    },
    {
      "access": "none",
      "code": "catalog",
      "id": 96,
      "name": "Catalog Items",
      "subCategory": "library"
    },
    {
      "access": "none",
      "code": "admin-certificates",
      "id": 111,
      "name": "Certificates",
      "subCategory": "infrastructure"
    },
    {
      "access": "none",
      "code": "admin-clients",
      "id": 31,
      "name": "Clients",
      "subCategory": "admin"
    },
    {
      "access": "none",
      "code": "admin-zones",
      "id": 3,
      "name": "Clouds",
      "subCategory": "infrastructure"
    },
    {
      "access": "none",
      "code": "infrastructure-cluster",
      "id": 153,
      "name": "Clusters",
      "subCategory": "infrastructure"
    },
    {
      "access": "none",
      "code": "deployments",
      "id": 100,
      "name": "Code Deployments",
      "subCategory": "provisioning"
    },
    {
      "access": "none",
      "code": "deployment-services",
      "id": 101,
      "name": "Code Integrations",
      "subCategory": "provisioning"
    },
    {
      "access": "none",
      "code": "code-repositories",
      "id": 102,
      "name": "Code Repositories",
      "subCategory": "provisioning"
    },
    {
      "access": "none",
      "code": "admin-servers",
      "id": 1,
      "name": "Compute",
      "subCategory": "infrastructure"
    },
    {
      "access": "none",
      "code": "services-vdi-copy",
      "id": 120,
      "name": "Copy/Paste",
      "subCategory": "virtual-desktop"
    },
    {
      "access": "none",
      "code": "credentials",
      "id": 109,
      "name": "Credentials",
      "subCategory": "infrastructure"
    },
    {
      "access": "none",
      "code": "services-cypher",
      "id": 115,
      "name": "Cypher",
      "subCategory": "tools"
    },
    {
      "access": "none",
      "code": "dashboard",
      "id": 44,
      "name": "Dashboard",
      "subCategory": "operations"
    },
    {
      "access": "none",
      "code": "service-catalog-dashboard",
      "id": 98,
      "name": "Dashboard",
      "subCategory": "catalog"
    },
    {
      "access": "none",
      "code": "infrastructure-network-dhcp-relay",
      "id": 130,
      "name": "DHCP Relays",
      "subCategory": "networks"
    },
    {
      "access": "none",
      "code": "infrastructure-network-dhcp-server",
      "id": 129,
      "name": "DHCP Servers",
      "subCategory": "networks"
    },
    {
      "access": "none",
      "code": "admin-distributed-workers",
      "id": 25,
      "name": "Distributed Workers",
      "subCategory": "admin"
    },
    {
      "access": "none",
      "code": "infrastructure-domains",
      "id": 133,
      "name": "Domains",
      "subCategory": "networks"
    },
    {
      "access": "none",
      "code": "admin-environments",
      "id": 20,
      "name": "Environment Settings",
      "subCategory": "admin"
    },
    {
      "access": "none",
      "code": "provisioning-environment",
      "id": 65,
      "name": "Environment Variables",
      "subCategory": "lifecycle"
    },
    {
      "access": "none",
      "code": "provisioning-execute-script",
      "id": 68,
      "name": "Execute Script",
      "subCategory": "provisioning"
    },
    {
      "access": "none",
      "code": "provisioning-execute-task",
      "id": 69,
      "name": "Execute Task",
      "subCategory": "provisioning"
    },
    {
      "access": "none",
      "code": "provisioning-execute-workflow",
      "id": 70,
      "name": "Execute Workflow",
      "subCategory": "provisioning"
    },
    {
      "access": "none",
      "code": "execution-request",
      "id": 157,
      "name": "Execution Request",
      "subCategory": "api"
    },
    {
      "access": "none",
      "code": "executions",
      "id": 103,
      "name": "Executions",
      "subCategory": "provisioning"
    },
    {
      "access": "none",
      "code": "admin-export-import",
      "id": 6,
      "name": "Export/Import",
      "subCategory": "admin"
    },
    {
      "access": "none",
      "code": "lifecycle-extend",
      "id": 95,
      "name": "Extend Expirations",
      "subCategory": "lifecycle"
    },
    {
      "access": "none",
      "code": "infrastructure-network-firewalls",
      "id": 137,
      "name": "Firewalls",
      "subCategory": "networks"
    },
    {
      "access": "none",
      "code": "infrastructure-floating-ips",
      "id": 136,
      "name": "Floating IPs",
      "subCategory": "networks"
    },
    {
      "access": "none",
      "code": "admin-groups",
      "id": 2,
      "name": "Groups",
      "subCategory": "infrastructure"
    },
    {
      "access": "none",
      "code": "guidance",
      "id": 54,
      "name": "Guidance",
      "subCategory": "operations"
    },
    {
      "access": "none",
      "code": "admin-guidanceSettings",
      "id": 16,
      "name": "Guidance Settings",
      "subCategory": "admin"
    },
    {
      "access": "none",
      "code": "admin-health",
      "id": 47,
      "name": "Health",
      "subCategory": "admin"
    },
    {
      "access": "none",
      "code": "admin-identity-sources",
      "id": 12,
      "name": "Identity Source",
      "subCategory": "admin"
    },
    {
      "access": "none",
      "code": "services-image-builder",
      "id": 118,
      "name": "Image Builder",
      "subCategory": "tools"
    },
    {
      "access": "none",
      "code": "provisioning-import-image",
      "id": 67,
      "name": "Import Image",
      "subCategory": "provisioning"
    },
    {
      "access": "none",
      "code": "admin-containers",
      "id": 9,
      "name": "Instance Types",
      "subCategory": "library"
    },
    {
      "access": "none",
      "code": "provisioning-add",
      "id": 59,
      "name": "Instances: Add",
      "subCategory": "provisioning"
    },
    {
      "access": "none",
      "code": "provisioning-clone",
      "id": 66,
      "name": "Instances: Clone",
      "subCategory": "provisioning"
    },
    {
      "access": "none",
      "code": "provisioning-delete",
      "id": 61,
      "name": "Instances: Delete",
      "subCategory": "provisioning"
    },
    {
      "access": "none",
      "code": "provisioning-edit",
      "id": 60,
      "name": "Instances: Edit",
      "subCategory": "provisioning"
    },
    {
      "access": "none",
      "code": "provisioning-force-delete",
      "id": 73,
      "name": "Instances: Force Delete",
      "subCategory": "provisioning"
    },
    {
      "access": "none",
      "code": "provisioning",
      "id": 58,
      "name": "Instances: List",
      "subCategory": "provisioning"
    },
    {
      "access": "none",
      "code": "provisioning-lock",
      "id": 62,
      "name": "Instances: Lock/Unlock",
      "subCategory": "provisioning"
    },
    {
      "access": "none",
      "code": "provisioning-remove-control",
      "id": 74,
      "name": "Instances: Remove From Control",
      "subCategory": "provisioning"
    },
    {
      "access": "none",
      "code": "provisioning-scale",
      "id": 63,
      "name": "Instances: Scale",
      "subCategory": "provisioning"
    },
    {
      "access": "none",
      "code": "provisioning-settings",
      "id": 64,
      "name": "Instances: Settings",
      "subCategory": "provisioning"
    },
    {
      "access": "none",
      "code": "infrastructure-network-integrations",
      "id": 128,
      "name": "Integration",
      "subCategory": "networks"
    },
    {
      "access": "none",
      "code": "admin-cm",
      "id": 23,
      "name": "Integrations",
      "subCategory": "admin"
    },
    {
      "access": "none",
      "code": "backup-services",
      "id": 43,
      "name": "Integrations",
      "subCategory": "backups"
    },
    {
      "access": "none",
      "code": "automation-services",
      "id": 105,
      "name": "Integrations",
      "subCategory": "library"
    },
    {
      "access": "none",
      "code": "operations-invoices",
      "id": 158,
      "name": "Invoices",
      "subCategory": "operations"
    },
    {
      "access": "none",
      "code": "infrastructure-ippools",
      "id": 135,
      "name": "IP Pools",
      "subCategory": "networks"
    },
    {
      "access": "none",
      "code": "job-executions",
      "id": 106,
      "name": "Job Executions",
      "subCategory": "provisioning"
    },
    {
      "access": "none",
      "code": "job-templates",
      "id": 107,
      "name": "Jobs",
      "subCategory": "provisioning"
    },
    {
      "access": "none",
      "code": "admin-keypairs",
      "id": 112,
      "name": "Keypairs",
      "subCategory": "infrastructure"
    },
    {
      "access": "none",
      "code": "services-kubernetes",
      "id": 117,
      "name": "Kubernetes",
      "subCategory": "tools"
    },
    {
      "access": "none",
      "code": "infrastructure-kube-cntl",
      "id": 154,
      "name": "Kubernetes Control",
      "subCategory": "infrastructure"
    },
    {
      "access": "none",
      "code": "library-packages",
      "id": 161,
      "name": "Library Packages",
      "subCategory": "library"
    },
    {
      "access": "none",
      "code": "admin-licenses",
      "id": 21,
      "name": "License Settings",
      "subCategory": "admin"
    },
    {
      "access": "none",
      "code": "infrastructure-loadbalancer",
      "id": 126,
      "name": "Load Balancers",
      "subCategory": "infrastructure"
    },
    {
      "access": "none",
      "code": "services-vdi-printer",
      "id": 121,
      "name": "Local Printer",
      "subCategory": "virtual-desktop"
    },
    {
      "access": "none",
      "code": "admin-logSettings",
      "id": 22,
      "name": "Log Settings",
      "subCategory": "admin"
    },
    {
      "access": "none",
      "code": "logs",
      "id": 49,
      "name": "Logs",
      "subCategory": "monitoring"
    },
    {
      "access": "none",
      "code": "admin-motd",
      "id": 29,
      "name": "Message of the day",
      "subCategory": "admin"
    },
    {
      "access": "none",
      "code": "monitoring",
      "id": 57,
      "name": "Monitoring",
      "subCategory": "monitoring"
    },
    {
      "access": "none",
      "code": "admin-monitorSettings",
      "id": 18,
      "name": "Monitoring Settings",
      "subCategory": "admin"
    },
    {
      "access": "none",
      "code": "infrastructure-move-server",
      "id": 155,
      "name": "Move Servers",
      "subCategory": "infrastructure"
    },
    {
      "access": "none",
      "code": "infrastructure-networks",
      "id": 127,
      "name": "Networks",
      "subCategory": "infrastructure"
    },
    {
      "access": "none",
      "code": "library-options",
      "id": 10,
      "name": "Options",
      "subCategory": "library"
    },
    {
      "access": "none",
      "code": "service-catalog-inventory",
      "id": 99,
      "name": "Order History",
      "subCategory": "catalog"
    },
    {
      "access": "none",
      "code": "admin-packages",
      "id": 26,
      "name": "Packages",
      "subCategory": "admin"
    },
    {
      "access": "none",
      "code": "admin-plugins",
      "id": 24,
      "name": "Plugins",
      "subCategory": "admin"
    },
    {
      "access": "none",
      "code": "admin-policies",
      "id": 27,
      "name": "Policies",
      "subCategory": "infrastructure"
    },
    {
      "access": "none",
      "code": "admin-global-policies",
      "id": 28,
      "name": "Policies",
      "subCategory": "admin"
    },
    {
      "access": "none",
      "code": "provisioning-power",
      "id": 93,
      "name": "Power Control",
      "subCategory": "lifecycle"
    },
    {
      "access": "none",
      "code": "admin-profiles",
      "id": 30,
      "name": "Profiles",
      "subCategory": "admin"
    },
    {
      "access": "none",
      "code": "projects",
      "id": 156,
      "name": "Projects",
      "subCategory": "projects"
    },
    {
      "access": "none",
      "code": "admin-provisioningSettings",
      "id": 19,
      "name": "Provisioning Settings",
      "subCategory": "admin"
    },
    {
      "access": "none",
      "code": "infrastructure-proxies",
      "id": 134,
      "name": "Proxies",
      "subCategory": "networks"
    },
    {
      "access": "none",
      "code": "provisioning-reconfigure",
      "id": 84,
      "name": "Reconfigure",
      "subCategory": "lifecycle"
    },
    {
      "access": "none",
      "code": "provisioning-reconfigure-change-plan",
      "id": 85,
      "name": "Reconfigure: Change Plan",
      "subCategory": "lifecycle"
    },
    {
      "access": "none",
      "code": "provisioning-reconfigure-add-disk",
      "id": 90,
      "name": "Reconfigure: Disk Add",
      "subCategory": "lifecycle"
    },
    {
      "access": "none",
      "code": "provisioning-reconfigure-disk-type",
      "id": 92,
      "name": "Reconfigure: Disk Change Type",
      "subCategory": "lifecycle"
    },
    {
      "access": "none",
      "code": "provisioning-reconfigure-modify-disk",
      "id": 89,
      "name": "Reconfigure: Disk Modify",
      "subCategory": "lifecycle"
    },
    {
      "access": "none",
      "code": "provisioning-reconfigure-remove-disk",
      "id": 91,
      "name": "Reconfigure: Disk Remove",
      "subCategory": "lifecycle"
    },
    {
      "access": "none",
      "code": "provisioning-reconfigure-add-network",
      "id": 87,
      "name": "Reconfigure: Network Add",
      "subCategory": "lifecycle"
    },
    {
      "access": "none",
      "code": "provisioning-reconfigure-modify-network",
      "id": 86,
      "name": "Reconfigure: Network Modify",
      "subCategory": "lifecycle"
    },
    {
      "access": "none",
      "code": "provisioning-reconfigure-remove-network",
      "id": 88,
      "name": "Reconfigure: Network Remove",
      "subCategory": "lifecycle"
    },
    {
      "access": "none",
      "code": "terminal",
      "id": 76,
      "name": "Remote Console",
      "subCategory": "provisioning"
    },
    {
      "access": "none",
      "code": "terminal-access",
      "id": 77,
      "name": "Remote Console Auto Login",
      "subCategory": "provisioning"
    },
    {
      "access": "none",
      "code": "reports",
      "id": 122,
      "name": "Reports",
      "subCategory": "operations"
    },
    {
      "access": "none",
      "code": "lifecycle-retry-cancel",
      "id": 94,
      "name": "Retry/Cancel",
      "subCategory": "lifecycle"
    },
    {
      "access": "none",
      "code": "admin-roles",
      "id": 8,
      "name": "Roles",
      "subCategory": "admin"
    },
    {
      "access": "none",
      "code": "infrastructure-network-router-firewalls",
      "id": 141,
      "name": "Router Firewalls",
      "subCategory": "networks"
    },
    {
      "access": "none",
      "code": "infrastructure-network-router-interfaces",
      "id": 140,
      "name": "Router Interfaces",
      "subCategory": "networks"
    },
    {
      "access": "none",
      "code": "infrastructure-nat",
      "id": 138,
      "name": "Router NAT",
      "subCategory": "networks"
    },
    {
      "access": "none",
      "code": "infrastructure-network-router-redistribution",
      "id": 142,
      "name": "Router Redistribution",
      "subCategory": "networks"
    },
    {
      "access": "none",
      "code": "infrastructure-network-router-routes",
      "id": 139,
      "name": "Router Routes",
      "subCategory": "networks"
    },
    {
      "access": "none",
      "code": "infrastructure-routers",
      "id": 132,
      "name": "Routers",
      "subCategory": "networks"
    },
    {
      "access": "none",
      "code": "security-scan",
      "id": 160,
      "name": "Scanning",
      "subCategory": "security"
    },
    {
      "access": "none",
      "code": "scheduling-execute",
      "id": 53,
      "name": "Scheduling - Execute",
      "subCategory": "library"
    },
    {
      "access": "none",
      "code": "scheduling-power",
      "id": 52,
      "name": "Scheduling - Power",
      "subCategory": "library"
    },
    {
      "access": "none",
      "code": "infrastructure-securityGroups",
      "id": 144,
      "name": "Security Groups",
      "subCategory": "infrastructure"
    },
    {
      "access": "none",
      "code": "infrastructure-network-server-groups",
      "id": 143,
      "name": "Server Groups",
      "subCategory": "networks"
    },
    {
      "access": "none",
      "code": "services-network-registry",
      "id": 114,
      "name": "Service Mesh",
      "subCategory": "provisioning"
    },
    {
      "access": "none",
      "code": "admin-servicePlans",
      "id": 7,
      "name": "Service Plans",
      "subCategory": "admin"
    },
    {
      "access": "none",
      "code": "snapshots",
      "id": 41,
      "name": "Snapshots",
      "subCategory": "snapshots"
    },
    {
      "access": "none",
      "code": "snapshots-linked-clone",
      "id": 42,
      "name": "Snapshots: Linked Clone",
      "subCategory": "snapshots"
    },
    {
      "access": "none",
      "code": "provisioning-state",
      "id": 71,
      "name": "State",
      "subCategory": "provisioning"
    },
    {
      "access": "none",
      "code": "infrastructure-state",
      "id": 146,
      "name": "State",
      "subCategory": "infrastructure"
    },
    {
      "access": "none",
      "code": "infrastructure-network-dhcp-routes",
      "id": 131,
      "name": "Static Routes",
      "subCategory": "networks"
    },
    {
      "access": "none",
      "code": "infrastructure-storage",
      "id": 124,
      "name": "Storage",
      "subCategory": "infrastructure"
    },
    {
      "access": "none",
      "code": "infrastructure-storage-browser",
      "id": 125,
      "name": "Storage Browser",
      "subCategory": "infrastructure"
    },
    {
      "access": "none",
      "code": "tasks",
      "id": 104,
      "name": "Tasks",
      "subCategory": "library"
    },
    {
      "access": "none",
      "code": "task-scripts",
      "id": 108,
      "name": "Tasks - Script Engines",
      "subCategory": "library"
    },
    {
      "access": "none",
      "code": "library-templates",
      "id": 11,
      "name": "Templates",
      "subCategory": "library"
    },
    {
      "access": "none",
      "code": "admin-accounts",
      "id": 4,
      "name": "Tenant",
      "subCategory": "admin"
    },
    {
      "access": "none",
      "code": "admin-accounts-users",
      "id": 5,
      "name": "Tenant - Impersonate Users",
      "subCategory": "admin"
    },
    {
      "access": "none",
      "code": "thresholds",
      "id": 82,
      "name": "Thresholds",
      "subCategory": "library"
    },
    {
      "access": "none",
      "code": "trust-services",
      "id": 110,
      "name": "Trust Integrations",
      "subCategory": "infrastructure"
    },
    {
      "access": "none",
      "code": "account-usage",
      "id": 56,
      "name": "Usage",
      "subCategory": "operations"
    },
    {
      "access": "none",
      "code": "admin-users",
      "id": 13,
      "name": "Users",
      "subCategory": "admin"
    },
    {
      "access": "none",
      "code": "services-vdi-pools",
      "id": 119,
      "name": "VDI Pools",
      "subCategory": "virtual-desktop"
    },
    {
      "access": "none",
      "code": "virtual-images",
      "id": 113,
      "name": "Virtual Images",
      "subCategory": "library"
    },
    {
      "access": "none",
      "code": "admin-whitelabel",
      "id": 15,
      "name": "Whitelabel Settings",
      "subCategory": "admin"
    },
    {
      "access": "none",
      "code": "operations-wiki",
      "id": 46,
      "name": "Wiki",
      "subCategory": "operations"
    }
  ],
  "globalAppTemplateAccess": "none",
  "globalCatalogItemTypeAccess": "none",
  "globalClusterTypeAccess": "none",
  "globalInstanceTypeAccess": "none",
  "globalPersonaAccess": "none",
  "globalReportTypeAccess": "none",
  "globalSiteAccess": "none",
  "globalTaskAccess": "none",
  "globalTaskSetAccess": "none",
  "globalVdiPoolAccess": "none",
  "globalZoneAccess": "none",
  "instanceTypePermissions": [],
  "personaPermissions": [],
  "reportTypePermissions": [],
  "sites": [],
  "taskPermissions": [],
  "taskSetPermissions": [],
  "vdiPoolPermissions": [],
  "zones": []
}
`
	// permissions obtained from the API after creating a role without setting permissions
	// i.e. all computed default values
	permissionsTestAPIComputedFull = `
{
  "featurePermissions": [
    {
      "id": 45,
      "code": "activity",
      "name": "Activity",
      "access": "none",
      "subCategory": "operations"
    },
    {
      "id": 72,
      "code": "provisioning-admin",
      "name": "Administrator",
      "access": "none",
      "subCategory": "provisioning"
    },
    {
      "id": 83,
      "code": "library-advanced-node-type-options",
      "name": "Advanced Node Type Options",
      "access": "none",
      "subCategory": "provisioning"
    },
    {
      "id": 48,
      "code": "operations-alarms",
      "name": "Alarms",
      "access": "none",
      "subCategory": "operations"
    },
    {
      "id": 123,
      "code": "reports-analytics",
      "name": "Analytics",
      "access": "none",
      "subCategory": "operations"
    },
    {
      "id": 159,
      "code": "integrations-ansible",
      "name": "Ansible",
      "access": "none",
      "subCategory": "admin"
    },
    {
      "id": 34,
      "code": "app-templates",
      "name": "App Blueprints",
      "access": "none",
      "subCategory": "library"
    },
    {
      "id": 14,
      "code": "admin-appliance",
      "name": "Appliance Settings",
      "access": "none",
      "subCategory": "admin"
    },
    {
      "id": 50,
      "code": "operations-approvals",
      "name": "Approvals",
      "access": "none",
      "subCategory": "operations"
    },
    {
      "id": 33,
      "code": "apps",
      "name": "Apps",
      "access": "none",
      "subCategory": "provisioning"
    },
    {
      "id": 116,
      "code": "services-archives",
      "name": "Archives",
      "access": "none",
      "subCategory": "tools"
    },
    {
      "id": 17,
      "code": "admin-backupSettings",
      "name": "Backup Settings",
      "access": "none",
      "subCategory": "admin"
    },
    {
      "id": 40,
      "code": "backups",
      "name": "Backups",
      "access": "none",
      "subCategory": "backups"
    },
    {
      "id": 55,
      "code": "billing",
      "name": "Billing",
      "access": "none",
      "subCategory": "api"
    },
    {
      "id": 35,
      "code": "arm-template",
      "name": "Blueprints - ARM",
      "access": "none",
      "subCategory": "library"
    },
    {
      "id": 37,
      "code": "cloudFormation-template",
      "name": "Blueprints - CloudFormation",
      "access": "none",
      "subCategory": "library"
    },
    {
      "id": 39,
      "code": "helm-template",
      "name": "Blueprints - Helm",
      "access": "none",
      "subCategory": "library"
    },
    {
      "id": 38,
      "code": "kubernetes-template",
      "name": "Blueprints - Kubernetes",
      "access": "none",
      "subCategory": "library"
    },
    {
      "id": 36,
      "code": "terraform-template",
      "name": "Blueprints - Terraform",
      "access": "none",
      "subCategory": "library"
    },
    {
      "id": 145,
      "code": "infrastructure-boot",
      "name": "Boot",
      "access": "none",
      "subCategory": "infrastructure"
    },
    {
      "id": 51,
      "code": "operations-budgets",
      "name": "Budgets",
      "access": "none",
      "subCategory": "operations"
    },
    {
      "id": 97,
      "code": "service-catalog",
      "name": "Catalog",
      "access": "none",
      "subCategory": "catalog"
    },
    {
      "id": 96,
      "code": "catalog",
      "name": "Catalog Items",
      "access": "none",
      "subCategory": "library"
    },
    {
      "id": 111,
      "code": "admin-certificates",
      "name": "Certificates",
      "access": "none",
      "subCategory": "infrastructure"
    },
    {
      "id": 31,
      "code": "admin-clients",
      "name": "Clients",
      "access": "none",
      "subCategory": "admin"
    },
    {
      "id": 3,
      "code": "admin-zones",
      "name": "Clouds",
      "access": "none",
      "subCategory": "infrastructure"
    },
    {
      "id": 153,
      "code": "infrastructure-cluster",
      "name": "Clusters",
      "access": "none",
      "subCategory": "infrastructure"
    },
    {
      "id": 100,
      "code": "deployments",
      "name": "Code Deployments",
      "access": "none",
      "subCategory": "provisioning"
    },
    {
      "id": 101,
      "code": "deployment-services",
      "name": "Code Integrations",
      "access": "none",
      "subCategory": "provisioning"
    },
    {
      "id": 102,
      "code": "code-repositories",
      "name": "Code Repositories",
      "access": "none",
      "subCategory": "provisioning"
    },
    {
      "id": 1,
      "code": "admin-servers",
      "name": "Compute",
      "access": "none",
      "subCategory": "infrastructure"
    },
    {
      "id": 120,
      "code": "services-vdi-copy",
      "name": "Copy/Paste",
      "access": "none",
      "subCategory": "virtual-desktop"
    },
    {
      "id": 109,
      "code": "credentials",
      "name": "Credentials",
      "access": "none",
      "subCategory": "infrastructure"
    },
    {
      "id": 115,
      "code": "services-cypher",
      "name": "Cypher",
      "access": "none",
      "subCategory": "tools"
    },
    {
      "id": 44,
      "code": "dashboard",
      "name": "Dashboard",
      "access": "none",
      "subCategory": "operations"
    },
    {
      "id": 98,
      "code": "service-catalog-dashboard",
      "name": "Dashboard",
      "access": "none",
      "subCategory": "catalog"
    },
    {
      "id": 130,
      "code": "infrastructure-network-dhcp-relay",
      "name": "DHCP Relays",
      "access": "none",
      "subCategory": "networks"
    },
    {
      "id": 129,
      "code": "infrastructure-network-dhcp-server",
      "name": "DHCP Servers",
      "access": "none",
      "subCategory": "networks"
    },
    {
      "id": 25,
      "code": "admin-distributed-workers",
      "name": "Distributed Workers",
      "access": "none",
      "subCategory": "admin"
    },
    {
      "id": 133,
      "code": "infrastructure-domains",
      "name": "Domains",
      "access": "none",
      "subCategory": "networks"
    },
    {
      "id": 20,
      "code": "admin-environments",
      "name": "Environment Settings",
      "access": "none",
      "subCategory": "admin"
    },
    {
      "id": 65,
      "code": "provisioning-environment",
      "name": "Environment Variables",
      "access": "none",
      "subCategory": "lifecycle"
    },
    {
      "id": 68,
      "code": "provisioning-execute-script",
      "name": "Execute Script",
      "access": "none",
      "subCategory": "provisioning"
    },
    {
      "id": 69,
      "code": "provisioning-execute-task",
      "name": "Execute Task",
      "access": "none",
      "subCategory": "provisioning"
    },
    {
      "id": 70,
      "code": "provisioning-execute-workflow",
      "name": "Execute Workflow",
      "access": "none",
      "subCategory": "provisioning"
    },
    {
      "id": 157,
      "code": "execution-request",
      "name": "Execution Request",
      "access": "none",
      "subCategory": "api"
    },
    {
      "id": 103,
      "code": "executions",
      "name": "Executions",
      "access": "none",
      "subCategory": "provisioning"
    },
    {
      "id": 6,
      "code": "admin-export-import",
      "name": "Export/Import",
      "access": "none",
      "subCategory": "admin"
    },
    {
      "id": 95,
      "code": "lifecycle-extend",
      "name": "Extend Expirations",
      "access": "none",
      "subCategory": "lifecycle"
    },
    {
      "id": 137,
      "code": "infrastructure-network-firewalls",
      "name": "Firewalls",
      "access": "none",
      "subCategory": "networks"
    },
    {
      "id": 136,
      "code": "infrastructure-floating-ips",
      "name": "Floating IPs",
      "access": "none",
      "subCategory": "networks"
    },
    {
      "id": 2,
      "code": "admin-groups",
      "name": "Groups",
      "access": "none",
      "subCategory": "infrastructure"
    },
    {
      "id": 54,
      "code": "guidance",
      "name": "Guidance",
      "access": "none",
      "subCategory": "operations"
    },
    {
      "id": 16,
      "code": "admin-guidanceSettings",
      "name": "Guidance Settings",
      "access": "none",
      "subCategory": "admin"
    },
    {
      "id": 47,
      "code": "admin-health",
      "name": "Health",
      "access": "none",
      "subCategory": "admin"
    },
    {
      "id": 12,
      "code": "admin-identity-sources",
      "name": "Identity Source",
      "access": "none",
      "subCategory": "admin"
    },
    {
      "id": 118,
      "code": "services-image-builder",
      "name": "Image Builder",
      "access": "none",
      "subCategory": "tools"
    },
    {
      "id": 67,
      "code": "provisioning-import-image",
      "name": "Import Image",
      "access": "none",
      "subCategory": "provisioning"
    },
    {
      "id": 9,
      "code": "admin-containers",
      "name": "Instance Types",
      "access": "none",
      "subCategory": "library"
    },
    {
      "id": 59,
      "code": "provisioning-add",
      "name": "Instances: Add",
      "access": "none",
      "subCategory": "provisioning"
    },
    {
      "id": 66,
      "code": "provisioning-clone",
      "name": "Instances: Clone",
      "access": "none",
      "subCategory": "provisioning"
    },
    {
      "id": 61,
      "code": "provisioning-delete",
      "name": "Instances: Delete",
      "access": "none",
      "subCategory": "provisioning"
    },
    {
      "id": 60,
      "code": "provisioning-edit",
      "name": "Instances: Edit",
      "access": "none",
      "subCategory": "provisioning"
    },
    {
      "id": 73,
      "code": "provisioning-force-delete",
      "name": "Instances: Force Delete",
      "access": "none",
      "subCategory": "provisioning"
    },
    {
      "id": 58,
      "code": "provisioning",
      "name": "Instances: List",
      "access": "none",
      "subCategory": "provisioning"
    },
    {
      "id": 62,
      "code": "provisioning-lock",
      "name": "Instances: Lock/Unlock",
      "access": "none",
      "subCategory": "provisioning"
    },
    {
      "id": 74,
      "code": "provisioning-remove-control",
      "name": "Instances: Remove From Control",
      "access": "none",
      "subCategory": "provisioning"
    },
    {
      "id": 63,
      "code": "provisioning-scale",
      "name": "Instances: Scale",
      "access": "none",
      "subCategory": "provisioning"
    },
    {
      "id": 64,
      "code": "provisioning-settings",
      "name": "Instances: Settings",
      "access": "none",
      "subCategory": "provisioning"
    },
    {
      "id": 128,
      "code": "infrastructure-network-integrations",
      "name": "Integration",
      "access": "none",
      "subCategory": "networks"
    },
    {
      "id": 23,
      "code": "admin-cm",
      "name": "Integrations",
      "access": "none",
      "subCategory": "admin"
    },
    {
      "id": 43,
      "code": "backup-services",
      "name": "Integrations",
      "access": "none",
      "subCategory": "backups"
    },
    {
      "id": 105,
      "code": "automation-services",
      "name": "Integrations",
      "access": "none",
      "subCategory": "library"
    },
    {
      "id": 158,
      "code": "operations-invoices",
      "name": "Invoices",
      "access": "none",
      "subCategory": "operations"
    },
    {
      "id": 135,
      "code": "infrastructure-ippools",
      "name": "IP Pools",
      "access": "none",
      "subCategory": "networks"
    },
    {
      "id": 106,
      "code": "job-executions",
      "name": "Job Executions",
      "access": "none",
      "subCategory": "provisioning"
    },
    {
      "id": 107,
      "code": "job-templates",
      "name": "Jobs",
      "access": "none",
      "subCategory": "provisioning"
    },
    {
      "id": 112,
      "code": "admin-keypairs",
      "name": "Keypairs",
      "access": "none",
      "subCategory": "infrastructure"
    },
    {
      "id": 117,
      "code": "services-kubernetes",
      "name": "Kubernetes",
      "access": "none",
      "subCategory": "tools"
    },
    {
      "id": 154,
      "code": "infrastructure-kube-cntl",
      "name": "Kubernetes Control",
      "access": "none",
      "subCategory": "infrastructure"
    },
    {
      "id": 161,
      "code": "library-packages",
      "name": "Library Packages",
      "access": "none",
      "subCategory": "library"
    },
    {
      "id": 21,
      "code": "admin-licenses",
      "name": "License Settings",
      "access": "none",
      "subCategory": "admin"
    },
    {
      "id": 126,
      "code": "infrastructure-loadbalancer",
      "name": "Load Balancers",
      "access": "none",
      "subCategory": "infrastructure"
    },
    {
      "id": 121,
      "code": "services-vdi-printer",
      "name": "Local Printer",
      "access": "none",
      "subCategory": "virtual-desktop"
    },
    {
      "id": 22,
      "code": "admin-logSettings",
      "name": "Log Settings",
      "access": "none",
      "subCategory": "admin"
    },
    {
      "id": 49,
      "code": "logs",
      "name": "Logs",
      "access": "none",
      "subCategory": "monitoring"
    },
    {
      "id": 29,
      "code": "admin-motd",
      "name": "Message of the day",
      "access": "none",
      "subCategory": "admin"
    },
    {
      "id": 57,
      "code": "monitoring",
      "name": "Monitoring",
      "access": "none",
      "subCategory": "monitoring"
    },
    {
      "id": 18,
      "code": "admin-monitorSettings",
      "name": "Monitoring Settings",
      "access": "none",
      "subCategory": "admin"
    },
    {
      "id": 155,
      "code": "infrastructure-move-server",
      "name": "Move Servers",
      "access": "none",
      "subCategory": "infrastructure"
    },
    {
      "id": 127,
      "code": "infrastructure-networks",
      "name": "Networks",
      "access": "none",
      "subCategory": "infrastructure"
    },
    {
      "id": 10,
      "code": "library-options",
      "name": "Options",
      "access": "none",
      "subCategory": "library"
    },
    {
      "id": 99,
      "code": "service-catalog-inventory",
      "name": "Order History",
      "access": "none",
      "subCategory": "catalog"
    },
    {
      "id": 26,
      "code": "admin-packages",
      "name": "Packages",
      "access": "none",
      "subCategory": "admin"
    },
    {
      "id": 24,
      "code": "admin-plugins",
      "name": "Plugins",
      "access": "none",
      "subCategory": "admin"
    },
    {
      "id": 27,
      "code": "admin-policies",
      "name": "Policies",
      "access": "none",
      "subCategory": "infrastructure"
    },
    {
      "id": 28,
      "code": "admin-global-policies",
      "name": "Policies",
      "access": "none",
      "subCategory": "admin"
    },
    {
      "id": 93,
      "code": "provisioning-power",
      "name": "Power Control",
      "access": "none",
      "subCategory": "lifecycle"
    },
    {
      "id": 30,
      "code": "admin-profiles",
      "name": "Profiles",
      "access": "none",
      "subCategory": "admin"
    },
    {
      "id": 156,
      "code": "projects",
      "name": "Projects",
      "access": "none",
      "subCategory": "projects"
    },
    {
      "id": 19,
      "code": "admin-provisioningSettings",
      "name": "Provisioning Settings",
      "access": "none",
      "subCategory": "admin"
    },
    {
      "id": 134,
      "code": "infrastructure-proxies",
      "name": "Proxies",
      "access": "none",
      "subCategory": "networks"
    },
    {
      "id": 84,
      "code": "provisioning-reconfigure",
      "name": "Reconfigure",
      "access": "none",
      "subCategory": "lifecycle"
    },
    {
      "id": 85,
      "code": "provisioning-reconfigure-change-plan",
      "name": "Reconfigure: Change Plan",
      "access": "none",
      "subCategory": "lifecycle"
    },
    {
      "id": 90,
      "code": "provisioning-reconfigure-add-disk",
      "name": "Reconfigure: Disk Add",
      "access": "none",
      "subCategory": "lifecycle"
    },
    {
      "id": 92,
      "code": "provisioning-reconfigure-disk-type",
      "name": "Reconfigure: Disk Change Type",
      "access": "none",
      "subCategory": "lifecycle"
    },
    {
      "id": 89,
      "code": "provisioning-reconfigure-modify-disk",
      "name": "Reconfigure: Disk Modify",
      "access": "none",
      "subCategory": "lifecycle"
    },
    {
      "id": 91,
      "code": "provisioning-reconfigure-remove-disk",
      "name": "Reconfigure: Disk Remove",
      "access": "none",
      "subCategory": "lifecycle"
    },
    {
      "id": 87,
      "code": "provisioning-reconfigure-add-network",
      "name": "Reconfigure: Network Add",
      "access": "none",
      "subCategory": "lifecycle"
    },
    {
      "id": 86,
      "code": "provisioning-reconfigure-modify-network",
      "name": "Reconfigure: Network Modify",
      "access": "none",
      "subCategory": "lifecycle"
    },
    {
      "id": 88,
      "code": "provisioning-reconfigure-remove-network",
      "name": "Reconfigure: Network Remove",
      "access": "none",
      "subCategory": "lifecycle"
    },
    {
      "id": 76,
      "code": "terminal",
      "name": "Remote Console",
      "access": "none",
      "subCategory": "provisioning"
    },
    {
      "id": 77,
      "code": "terminal-access",
      "name": "Remote Console Auto Login",
      "access": "none",
      "subCategory": "provisioning"
    },
    {
      "id": 122,
      "code": "reports",
      "name": "Reports",
      "access": "none",
      "subCategory": "operations"
    },
    {
      "id": 94,
      "code": "lifecycle-retry-cancel",
      "name": "Retry/Cancel",
      "access": "none",
      "subCategory": "lifecycle"
    },
    {
      "id": 8,
      "code": "admin-roles",
      "name": "Roles",
      "access": "none",
      "subCategory": "admin"
    },
    {
      "id": 141,
      "code": "infrastructure-network-router-firewalls",
      "name": "Router Firewalls",
      "access": "none",
      "subCategory": "networks"
    },
    {
      "id": 140,
      "code": "infrastructure-network-router-interfaces",
      "name": "Router Interfaces",
      "access": "none",
      "subCategory": "networks"
    },
    {
      "id": 138,
      "code": "infrastructure-nat",
      "name": "Router NAT",
      "access": "none",
      "subCategory": "networks"
    },
    {
      "id": 142,
      "code": "infrastructure-network-router-redistribution",
      "name": "Router Redistribution",
      "access": "none",
      "subCategory": "networks"
    },
    {
      "id": 139,
      "code": "infrastructure-network-router-routes",
      "name": "Router Routes",
      "access": "none",
      "subCategory": "networks"
    },
    {
      "id": 132,
      "code": "infrastructure-routers",
      "name": "Routers",
      "access": "none",
      "subCategory": "networks"
    },
    {
      "id": 160,
      "code": "security-scan",
      "name": "Scanning",
      "access": "none",
      "subCategory": "security"
    },
    {
      "id": 53,
      "code": "scheduling-execute",
      "name": "Scheduling - Execute",
      "access": "none",
      "subCategory": "library"
    },
    {
      "id": 52,
      "code": "scheduling-power",
      "name": "Scheduling - Power",
      "access": "none",
      "subCategory": "library"
    },
    {
      "id": 144,
      "code": "infrastructure-securityGroups",
      "name": "Security Groups",
      "access": "none",
      "subCategory": "infrastructure"
    },
    {
      "id": 143,
      "code": "infrastructure-network-server-groups",
      "name": "Server Groups",
      "access": "none",
      "subCategory": "networks"
    },
    {
      "id": 114,
      "code": "services-network-registry",
      "name": "Service Mesh",
      "access": "none",
      "subCategory": "provisioning"
    },
    {
      "id": 7,
      "code": "admin-servicePlans",
      "name": "Service Plans",
      "access": "none",
      "subCategory": "admin"
    },
    {
      "id": 41,
      "code": "snapshots",
      "name": "Snapshots",
      "access": "none",
      "subCategory": "snapshots"
    },
    {
      "id": 42,
      "code": "snapshots-linked-clone",
      "name": "Snapshots: Linked Clone",
      "access": "none",
      "subCategory": "snapshots"
    },
    {
      "id": 71,
      "code": "provisioning-state",
      "name": "State",
      "access": "none",
      "subCategory": "provisioning"
    },
    {
      "id": 146,
      "code": "infrastructure-state",
      "name": "State",
      "access": "none",
      "subCategory": "infrastructure"
    },
    {
      "id": 131,
      "code": "infrastructure-network-dhcp-routes",
      "name": "Static Routes",
      "access": "none",
      "subCategory": "networks"
    },
    {
      "id": 124,
      "code": "infrastructure-storage",
      "name": "Storage",
      "access": "none",
      "subCategory": "infrastructure"
    },
    {
      "id": 125,
      "code": "infrastructure-storage-browser",
      "name": "Storage Browser",
      "access": "none",
      "subCategory": "infrastructure"
    },
    {
      "id": 104,
      "code": "tasks",
      "name": "Tasks",
      "access": "none",
      "subCategory": "library"
    },
    {
      "id": 108,
      "code": "task-scripts",
      "name": "Tasks - Script Engines",
      "access": "none",
      "subCategory": "library"
    },
    {
      "id": 11,
      "code": "library-templates",
      "name": "Templates",
      "access": "none",
      "subCategory": "library"
    },
    {
      "id": 4,
      "code": "admin-accounts",
      "name": "Tenant",
      "access": "none",
      "subCategory": "admin"
    },
    {
      "id": 5,
      "code": "admin-accounts-users",
      "name": "Tenant - Impersonate Users",
      "access": "none",
      "subCategory": "admin"
    },
    {
      "id": 82,
      "code": "thresholds",
      "name": "Thresholds",
      "access": "none",
      "subCategory": "library"
    },
    {
      "id": 110,
      "code": "trust-services",
      "name": "Trust Integrations",
      "access": "none",
      "subCategory": "infrastructure"
    },
    {
      "id": 56,
      "code": "account-usage",
      "name": "Usage",
      "access": "none",
      "subCategory": "operations"
    },
    {
      "id": 13,
      "code": "admin-users",
      "name": "Users",
      "access": "none",
      "subCategory": "admin"
    },
    {
      "id": 119,
      "code": "services-vdi-pools",
      "name": "VDI Pools",
      "access": "none",
      "subCategory": "virtual-desktop"
    },
    {
      "id": 113,
      "code": "virtual-images",
      "name": "Virtual Images",
      "access": "none",
      "subCategory": "library"
    },
    {
      "id": 15,
      "code": "admin-whitelabel",
      "name": "Whitelabel Settings",
      "access": "none",
      "subCategory": "admin"
    },
    {
      "id": 46,
      "code": "operations-wiki",
      "name": "Wiki",
      "access": "none",
      "subCategory": "operations"
    }
  ],
  "globalSiteAccess": "none",
  "sites": [],
  "globalZoneAccess": "none",
  "zones": [],
  "globalInstanceTypeAccess": "none",
  "instanceTypePermissions": [],
  "globalAppTemplateAccess": "none",
  "appTemplatePermissions": [],
  "globalCatalogItemTypeAccess": "none",
  "catalogItemTypePermissions": [],
  "globalPersonaAccess": "none",
  "personaPermissions": [],
  "globalVdiPoolAccess": "none",
  "vdiPoolPermissions": [],
  "globalReportTypeAccess": "none",
  "reportTypePermissions": [],
  "globalTaskAccess": "none",
  "taskPermissions": [],
  "globalTaskSetAccess": "none",
  "taskSetPermissions": [],
  "globalClusterTypeAccess": "none",
  "clusterTypePermissions": []
}
`
)
