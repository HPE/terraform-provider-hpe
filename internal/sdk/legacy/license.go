package morpheus

import (
	"fmt"
	"time"
)

var (
	// LicensePath is the API endpoint for license configuration
	LicensePath = "/api/license"
)

// ApplianceLicense is one licence as /api/license reports it: the active one
// under "license" and every installed one under "installedLicenses". KeyId
// (the first eight characters of the key) and Hash (the SHA-256 of the key)
// identify a licence without the key itself being readable.
type ApplianceLicense struct {
	ID             int64     `json:"id"`
	KeyID          string    `json:"keyId"`
	Hash           string    `json:"hash"`
	LicenseVersion int       `json:"licenseVersion"`
	Producttier    string    `json:"productTier"`
	Startdate      time.Time `json:"startDate"`
	Enddate        time.Time `json:"endDate"`
	Maxinstances   int       `json:"maxInstances"`
	Maxmemory      int       `json:"maxMemory"`
	Maxstorage     int       `json:"maxStorage"`
	Hardlimit      bool      `json:"hardLimit"`
	Freetrial      bool      `json:"freeTrial"`
	Multitenant    bool      `json:"multiTenant"`
	Whitelabel     bool      `json:"whitelabel"`
	Reportstatus   bool      `json:"reportStatus"`
	Supportlevel   string    `json:"supportLevel"`
	Accountname    string    `json:"accountName"`
	Config         struct {
	} `json:"config"`
	Amazonproductcodes interface{} `json:"amazonProductCodes"`
	Features           struct {
		Dashboard                bool `json:"dashboard"`
		Guidance                 bool `json:"guidance"`
		Discovery                bool `json:"discovery"`
		Analytics                bool `json:"analytics"`
		Scheduling               bool `json:"scheduling"`
		Approvals                bool `json:"approvals"`
		Usage                    bool `json:"usage"`
		Activity                 bool `json:"activity"`
		Instances                bool `json:"instances"`
		Apps                     bool `json:"apps"`
		Templates                bool `json:"templates"`
		Automation               bool `json:"automation"`
		Virtualimages            bool `json:"virtualImages"`
		Library                  bool `json:"library"`
		Migrations               bool `json:"migrations"`
		Deployments              bool `json:"deployments"`
		Groups                   bool `json:"groups"`
		Clouds                   bool `json:"clouds"`
		Hosts                    bool `json:"hosts"`
		Network                  bool `json:"network"`
		Loadbalancers            bool `json:"loadBalancers"`
		Storage                  bool `json:"storage"`
		Keypairs                 bool `json:"keyPairs"`
		Sslcertificates          bool `json:"sslCertificates"`
		Boot                     bool `json:"boot"`
		Backups                  bool `json:"backups"`
		Cypher                   bool `json:"cypher"`
		Archives                 bool `json:"archives"`
		Imagebuilder             bool `json:"imageBuilder"`
		Tenants                  bool `json:"tenants"`
		Plans                    bool `json:"plans"`
		Pricing                  bool `json:"pricing"`
		Users                    bool `json:"users"`
		Usergroups               bool `json:"userGroups"`
		Monitoring               bool `json:"monitoring"`
		Logging                  bool `json:"logging"`
		Monitoringservices       bool `json:"monitoringServices"`
		Loggingservices          bool `json:"loggingServices"`
		Backupservices           bool `json:"backupServices"`
		Dnsservices              bool `json:"dnsServices"`
		Codeservice              bool `json:"codeService"`
		Buildservices            bool `json:"buildServices"`
		Loadbalancerservices     bool `json:"loadBalancerServices"`
		Ipamservices             bool `json:"ipamServices"`
		Approvalservices         bool `json:"approvalServices"`
		Cmdbservices             bool `json:"cmdbServices"`
		Deploymentservices       bool `json:"deploymentServices"`
		Automationservices       bool `json:"automationServices"`
		Servicediscoveryservices bool `json:"serviceDiscoveryServices"`
		Identityservices         bool `json:"identityServices"`
		Trustservices            bool `json:"trustServices"`
		Securityservices         bool `json:"securityServices"`
	} `json:"features"`
	Zonetypes   interface{} `json:"zoneTypes"`
	Lastupdated time.Time   `json:"lastUpdated"`
	Datecreated time.Time   `json:"dateCreated"`
}

// GetLicenseResult structures for use in request and response payloads.
// License is nil when no licence is active on the appliance.
type GetLicenseResult struct {
	License           *ApplianceLicense  `json:"license"`
	InstalledLicenses []ApplianceLicense `json:"installedLicenses"`
	Currentusage      struct {
		Memory    int64 `json:"memory"`
		Storage   int64 `json:"storage"`
		Workloads int   `json:"workloads"`
	} `json:"currentUsage"`
}

type UninstallLicenseResult struct {
	DeleteResult
}

// Client request methods
func (client *Client) GetLicense(req *Request) (*Response, error) {
	return client.Execute(&Request{
		Method:      "GET",
		Path:        LicensePath,
		QueryParams: req.QueryParams,
		Result:      &GetLicenseResult{},
	})
}

func (client *Client) InstallLicense(req *Request) (*Response, error) {
	return client.Execute(&Request{
		Method:      "POST",
		Path:        LicensePath,
		QueryParams: req.QueryParams,
		Body:        req.Body,
		Result:      &GetLicenseResult{},
	})
}

func (client *Client) TestLicense(req *Request) (*Response, error) {
	return client.Execute(&Request{
		Method:      "POST",
		Path:        fmt.Sprintf("%s/test", LicensePath),
		QueryParams: req.QueryParams,
		Body:        req.Body,
		Result:      &GetLicenseResult{},
	})
}

// UninstallLicense deactivates the installed licence whose key id (the first
// eight characters of its key, as GET /api/license reports under keyId) is
// keyID. An empty keyID deactivates every installed licence, which is what the
// API does when the parameter is absent. There is no per-licence path: the
// endpoint is DELETE /api/license, selected by the keyId query parameter.
func (client *Client) UninstallLicense(keyID string, req *Request) (*Response, error) {
	queryParams := map[string]string{}
	for k, v := range req.QueryParams {
		queryParams[k] = v
	}

	if keyID != "" {
		queryParams["keyId"] = keyID
	}

	return client.Execute(&Request{
		Method:      "DELETE",
		Path:        LicensePath,
		QueryParams: queryParams,
		Body:        req.Body,
		Result:      &UninstallLicenseResult{},
	})
}
