// (C) Copyright 2021-2023 Hewlett Packard Enterprise Development LP

package models

type CloudRespBody struct {
	ID         int         `json:"id"`
	UUID       string      `json:"uuid"`
	Externalid interface{} `json:"externalId"`
	Name       string      `json:"name"`
	Code       string      `json:"code"`
	Location   string      `json:"location"`
	Owner      struct {
		ID   int    `json:"id"`
		Name string `json:"name"`
	} `json:"owner"`
	Accountid int `json:"accountId"`
	Account   struct {
		ID   int    `json:"id"`
		Name string `json:"name"`
	} `json:"account"`
	Visibility    string      `json:"visibility"`
	Enabled       bool        `json:"enabled"`
	Status        string      `json:"status"`
	Statusmessage interface{} `json:"statusMessage"`
	Zonetype      struct {
		ID   int    `json:"id"`
		Code string `json:"code"`
		Name string `json:"name"`
	} `json:"zoneType"`
	Zonetypeid      int         `json:"zoneTypeId"`
	Guidancemode    string      `json:"guidanceMode"`
	Storagemode     string      `json:"storageMode"`
	Agentmode       string      `json:"agentMode"`
	Userdatalinux   interface{} `json:"userDataLinux"`
	Userdatawindows interface{} `json:"userDataWindows"`
	Consolekeymap   string      `json:"consoleKeymap"`
	Containermode   string      `json:"containerMode"`
	Serviceversion  interface{} `json:"serviceVersion"`
	Costingmode     interface{} `json:"costingMode"`
	Inventorylevel  string      `json:"inventoryLevel"`
	Timezone        interface{} `json:"timezone"`
	Apiproxy        struct {
		ID   int    `json:"id"`
		Name string `json:"name"`
	} `json:"apiProxy"`
	Provisioningproxy struct {
		ID   int    `json:"id"`
		Name string `json:"name"`
	} `json:"provisioningProxy"`
	Networkdomain         interface{} `json:"networkDomain"`
	Domainname            string      `json:"domainName"`
	Regioncode            string      `json:"regionCode"`
	Autorecoverpowerstate bool        `json:"autoRecoverPowerState"`
	Scalepriority         int         `json:"scalePriority"`
	Config                struct {
		Apiurl            string      `json:"apiUrl"`
		Username          string      `json:"username"`
		Password          string      `json:"password"`
		Datacenter        string      `json:"datacenter"`
		Cluster           string      `json:"cluster"`
		Resourcepoolid    interface{} `json:"resourcePoolId"`
		Resourcepool      string      `json:"resourcePool"`
		Rpcmode           string      `json:"rpcMode"`
		Hidehostselection string      `json:"hideHostSelection"`
		Enablevnc         string      `json:"enableVnc"`
		Diskstoragetype   string      `json:"diskStorageType"`
		Applianceurl      string      `json:"applianceUrl"`
		Datacentername    string      `json:"datacenterName"`
		NetworkserverID   string      `json:"networkServer.id"`
		Networkserver     struct {
			ID string `json:"id"`
		} `json:"networkServer"`
		Securitymode               string `json:"securityMode"`
		Securityserver             string `json:"securityServer"`
		Certificateprovider        string `json:"certificateProvider"`
		Backupmode                 string `json:"backupMode"`
		Replicationmode            string `json:"replicationMode"`
		Datacenterid               string `json:"datacenterId"`
		Apiversion                 string `json:"apiVersion"`
		Serviceregistryid          string `json:"serviceRegistryId"`
		Configmanagementid         string `json:"configManagementId"`
		Importexisting             string `json:"importExisting"`
		Enabledisktypeselection    string `json:"enableDiskTypeSelection"`
		Enablenetworktypeselection string `json:"enableNetworkTypeSelection"`
	} `json:"config"`
	Datecreated string `json:"dateCreated"`
	Lastupdated string `json:"lastUpdated"`
	Groups      []struct {
		ID        int    `json:"id"`
		Name      string `json:"name"`
		Accountid int    `json:"accountId"`
	} `json:"groups"`
	Securityserver interface{} `json:"securityServer"`
	Stats          struct {
		Servercounts struct {
			All           int `json:"all"`
			Host          int `json:"host"`
			Hypervisor    int `json:"hypervisor"`
			Containerhost int `json:"containerHost"`
			VM            int `json:"vm"`
			Baremetal     int `json:"baremetal"`
			Unmanaged     int `json:"unmanaged"`
		} `json:"serverCounts"`
	} `json:"stats"`
	Servercount int `json:"serverCount"`
}

type CloudsResp struct {
	Clouds []CloudRespBody `json:"zones"`
}

type GetSpecificCloud struct {
	Cloud CloudRespBody `json:"zone"`
}

type ResourcePoolRespBody struct {
	ID          int         `json:"id"`
	Name        string      `json:"name"`
	Description interface{} `json:"description"`
	Zone        struct {
		ID   int    `json:"id"`
		Name string `json:"name"`
	} `json:"zone"`
	Parent struct {
		ID   int    `json:"id"`
		Name string `json:"name"`
	} `json:"parent"`
	Type        string      `json:"type"`
	Externalid  string      `json:"externalId"`
	Regioncode  interface{} `json:"regionCode"`
	Iacid       interface{} `json:"iacId"`
	Visibility  string      `json:"visibility"`
	Readonly    bool        `json:"readOnly"`
	Defaultpool bool        `json:"defaultPool"`
	Active      bool        `json:"active"`
	Status      string      `json:"status"`
	Config      struct{}    `json:"config"`
	Tenants     []struct {
		ID   int    `json:"id"`
		Name string `json:"name"`
	} `json:"tenants"`
	Resourcepermission struct {
		All      bool          `json:"all"`
		Allplans bool          `json:"allPlans"`
		Sites    []interface{} `json:"sites"`
		Plans    []interface{} `json:"plans"`
	} `json:"resourcePermission"`
	Depth int `json:"depth"`
}

type ResourcePoolsResp struct {
	ResourcePools []ResourcePoolRespBody `json:"resourcePools"`
}

type GetSpecificCloudResourcePool struct {
	ResourcePools ResourcePoolRespBody `json:"resourcePool"`
}

type DataStoresRespBody struct {
	ID   int    `json:"id"`
	Name string `json:"name"`
	Zone struct {
		ID   int    `json:"id"`
		Name string `json:"name"`
	} `json:"zone"`
	Type       string `json:"type"`
	Freespace  int64  `json:"freeSpace"`
	Online     bool   `json:"online"`
	Active     bool   `json:"active"`
	Visibility string `json:"visibility"`
	Tenants    []struct {
		ID            int    `json:"id"`
		Name          string `json:"name"`
		Defaultstore  bool   `json:"defaultStore"`
		Defaulttarget bool   `json:"defaultTarget"`
	} `json:"tenants"`
	Resourcepermission struct {
		All      bool          `json:"all"`
		Allplans bool          `json:"allPlans"`
		Sites    []interface{} `json:"sites"`
		Plans    []interface{} `json:"plans"`
	} `json:"resourcePermission"`
}

type DataStoresResp struct {
	Datastores []DataStoresRespBody `json:"datastores"`
}

type GetSpecificCloudDataStores struct {
	Datastores DataStoresRespBody `json:"datastore"`
}

type GetAllCloudNetworks struct {
	Data DataGetNetworkInterface `json:"data"`
}

type DataGetNetworkInterface struct {
	NetworkTypes []GetNetworkInterfaceNetworkTypes `json:"networkTypes"`
}

type GetNetworkInterfaceNetworkTypes struct {
	ID   int    `json:"id"`
	Name string `json:"name"`
	Code string `json:"code"`
}

type GetAllCloudFolders struct {
	Folders []GetCloudFolder `json:"folders"`
}

type GetCloudFolder struct {
	ID            int         `json:"id"`
	Name          string      `json:"name"`
	Parent        interface{} `json:"parent"`
	Type          string      `json:"type"`
	ExternalID    string      `json:"externalId"`
	Visibility    string      `json:"visibility"`
	ReadOnly      bool        `json:"readOnly"`
	DefaultFolder bool        `json:"defaultFolder"`
	DefaultStore  bool        `json:"defaultStore"`
	Active        bool        `json:"active"`
}

type GetSpecificCloudFolder struct {
	Folder GetCloudFolder `json:"folder"`
}
